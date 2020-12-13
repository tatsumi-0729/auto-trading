package controller

import (
	"auto-trading/app/bitflyer"
	"auto-trading/app/config"
	"auto-trading/app/model"
	"log"
)

func StreamIngestionData() {
	// tickerを扱うチャネルを作成
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	// goroutinでリアルタイム取得処理を走らせる
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	go func() {
		for ticker := range tickerChannel {
			log.Printf("action=StreamIngestionData, %v", ticker)
			// tickerからデータを取得する度に、1s、1m、1hのそれぞれのテーブルにデータを追加する
			for _, duration := range config.Config.Durations {
				isCreated := model.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreated == true && duration == config.Config.TradeDuration {
					// TODO
				}
			}
		}
	}()
}
