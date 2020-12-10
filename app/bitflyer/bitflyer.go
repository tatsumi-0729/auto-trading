package bitflyer

import (
	"auto-trading/app/config"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var baseURL = config.Config.BaseUrl

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

// コンストラクタ
func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

// ヘッダーを作成するメソッド(hmacで認証)
func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	// Unix()を使ってタイムスタンプ型にする
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)

	// api.secretのハッシュを作成
	mac := hmac.New(sha256.New, []byte(api.secret))
	// ハッシュに加えたいmessageを追加
	mac.Write([]byte(message))
	// ハッシュの終わりを意味するnilを末尾に加える
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

// httpリクエストを送信するメソッド
func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	// urlが正しいものか、パースで検証する
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}
	// リクエスト送信先のエンドポイントを作成する
	endpoint := baseURL.ResolveReference(apiURL).String()
	log.Printf("action=doRequest endpoint=%s", endpoint)
	// リクエストを作成する
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	// クエリ(リクエストパラメータ)を取得する
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	// クエリを足したURLをエンコードする
	req.URL.RawQuery = q.Encode()

	// ヘッダーをリクエストに追加する
	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}
	// リクエストを送信する
	res, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// レスポンスボディを読み込む
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

// 自分の資産情報を取得する
func (api *APIClient) GetBalance() ([]Balance, error) {
	url := config.Config.GetBalanceUrl
	res, err := api.doRequest("GET", url, map[string]string{}, nil)
	log.Printf("url=%s res=%s", url, string(res))
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	// レスポンスのjsonとbalanceをmarshalで紐付ける
	var balance []Balance
	err = json.Unmarshal(res, &balance)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	return balance, nil
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// 売値と買値の中間を取得する
func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

// DB保存用の時間を変換する
func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=DateTime, err=%s", err.Error())
	}
	return dateTime
}

// 分と秒をトランケイト(0にする)
func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

// Bitcoinのデータを取得する
func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := config.Config.GetTickerUrl
	res, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	if err != nil {
		return nil, err
	}
	// レスポンスのjsonとbalanceをmarshalで紐付ける
	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil {
		return nil, err
	}
	return &ticker, nil
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

// tickerをリアルタイムで取得する
func (api *APIClient) GetRealTimeTicker(symbol string, ch chan<- Ticker) {
	u := url.URL{Scheme: config.Config.GetRealTimeTickerSchema, Host: config.Config.GetRealTimeTickerHost, Path: config.Config.GetRealTimeTickerPath}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshaTic, err := json.Marshal(binary)
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

// 注文を入れるメソッド
func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	// orderをMarshalでjsonにする
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	url := config.Config.SendOrderUrl
	res, err := api.doRequest("POST", url, map[string]string{}, data)
	if err != nil {
		return nil, err
	}
	// レスポンスとresponseをUnmarshalで紐付ける
	var response ResponseSendChildOrder
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// 注文の履歴・状態を取得する
func (api *APIClient) ListOrder(query map[string]string) ([]Order, error) {
	res, err := api.doRequest("GET", config.Config.ListOrderUrl, query, nil)
	if err != nil {
		return nil, err
	}
	var responseListOrder []Order
	// レスポンスとOrderをUnmarshalで紐付ける
	err = json.Unmarshal(res, &responseListOrder)
	if err != nil {
		return nil, err
	}
	return responseListOrder, nil
}
