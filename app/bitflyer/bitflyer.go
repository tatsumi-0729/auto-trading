package bitflyer

import (
	"auto-trading/app/config"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
func (api *ApiClient) sendRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
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

// 自分の資産情報を取得する
func (api *ApiClient) GetBalance() (string, error) {
	url := config.Config.BaseUrl + config.Config.GetBalanceUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("リクエストの作成に失敗しました。", err.Error())
	}
	res, err := api.httpClient.Do(req)
	if err != nil {
		log.Fatalln("リクエストの送信に失敗しました。", err.Error())
	}
	// レスポンスボディを読み込む
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("レスポンスボディの読み込みに失敗しました", err.Error())
	}
	return string(body), err
}
