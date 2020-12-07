package bitflyer

import (
	"auto-trading/app/config"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var baseUrl = config.Config.BaseUrl

type ApiClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

// コンストラクタ
func New(key, secret string) *ApiClient {
	apiClient := &ApiClient{key, secret, &http.Client{}}
	return apiClient
}

// ヘッダーを作成するメソッド(hmacで認証)
func (api ApiClient) header(method, endpoint string, body []byte) map[string]string {
	// Unix()でタイムスタンプ型にする
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(timestamp)
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
func (api *ApiClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	// urlが正しいものか、パースで検証する
	baseUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("URL:%sが不正です。", baseUrl)
	}
	apiUrl, err := url.Parse(urlPath)
	if err != nil {
		log.Fatalf("URLのパス:%sが不正です。", apiUrl)
	}
	// リクエスト送信先のエンドポイントを作成する
	endpoint := baseUrl.ResolveReference(apiUrl).String()
	log.Printf("endpoint=%s", endpoint)
	// リクエストを作成する
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("リクエストの作成に失敗しました。リクエスト=%s", req, err.Error())
	}
	// クエリ(リクエストパラメータ)を取得する
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	// クエリを足したURLをエンコードする
	req.URL.RawQuery = q.Encode()

	// ヘッダーがあればリクエストに追加する
	for k, v := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(k, v)
	}

	// リクエストを送信する
	res, err := api.httpClient.Do(req)
	if err != nil {
		log.Fatalln("リクエストの送信に失敗しました。", err.Error())
	}

	defer res.Body.Close()
	// レスポンスボディを読み込む
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("レスポンスボディの読み込みに失敗しました", err.Error())
	}
	return body, err
}

type Balance struct {
	CurrentCode string  `json:"current_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

// 自分の資産情報を取得する
func (api *ApiClient) GetBalance() ([]Balance, error) {
	url := config.Config.GetBalanceUrl
	res, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		log.Fatalln("リクエストの作成に失敗しました。", err.Error())
	}
	// レスポンスのjsonとbalanceをmarshalで紐付ける
	var balance []Balance
	if err = json.Unmarshal(res, &balance); err != nil {
		log.Fatalln("レスポンスのjsonの紐付けに失敗しました", err.Error())
	}
	return balance, err
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          float64 `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// Bitcoinのデーターを取得する
func (api *ApiClient) GetTicker(product_code string) (*Ticker, error) {
	url := config.Config.GetTickerUrl
	res, err := api.doRequest("GET", url, map[string]string{"product_code": product_code}, nil)
	if err != nil {
		log.Fatalln("リクエストの作成に失敗しました。", err.Error())
	}
	// レスポンスのjsonとbalanceをmarshalで紐付ける
	var ticker Ticker
	if err = json.Unmarshal(res, &ticker); err != nil {
		log.Fatalln("レスポンスのjsonの紐付けに失敗しました", err.Error())
	}
	return &ticker, err
}
