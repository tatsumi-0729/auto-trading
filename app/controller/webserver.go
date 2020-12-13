package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"auto-trading/app/config"
)

// 表示したいhtmlをキャッシュ(読み込み)しておく
var templates = template.Must(template.ParseFiles("view/chart.html"))

// ハンドラを定義する
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ハンドラを登録する
func StartWebServer() error {
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
