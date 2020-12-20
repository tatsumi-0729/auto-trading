package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"auto-trading/app/config"
	"auto-trading/app/model"
)

// 表示したいhtmlをキャッシュ(読み込み)しておく
var templates = template.Must(template.ParseFiles("view/chart.html"))

// ハンドラを定義する
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	// テンプレート側にdfCandleを渡す
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// jsonエラーを取得する
func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	// http.ResponseWriterでjsonエラーを表示
	w.Write(jsonError)
}

var apiValidPath = regexp.MustCompile("^/api/candle/$")

// こいつを通してjsonをブラウザに表示する
func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// urlがマッチングするかを確認
		m := apiValidPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			APIError(w, "Not found", http.StatusNotFound)
		}
		fn(w, r)
	}
}

// jsonにしてhttp.ResponseWriterに値を返す
func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	//　ブラウザのurlから、ajaxでデータを取得する
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		APIError(w, "No product_code param", http.StatusBadRequest)
		return
	}
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}
	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Config.Durations[duration]

	df, _ := model.GetAllCandle(productCode, durationTime, limit)

	// dfをjsonに変換
	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	// http.ResponseWriterでjsonを表示
	w.Write(js)
}

// ハンドラを登録して、サーバ通信開始
func StartWebServer() error {
	// ハンドラを登録して、パスを通す
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/chart/", viewChartHandler)
	// デフォルトパス　/
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
