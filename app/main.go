package main

import (
	// go mod initで、プロジェクトのパスを登録しないとエラーになる。
	"auto-trading/app/config"
	"auto-trading/app/controller"
	"auto-trading/app/model"
	"auto-trading/app/util"
	"fmt"
)

func main() {
	// log出力の設定
	util.Logging(config.Config.LogFile)
	fmt.Println(model.DbConnection)
	controller.StreamIngestionData()
	controller.StartWebServer()
}

// ticker, _ := apiClient.GetTicker("ETH_BTC")
// fmt.Println(ticker.GetMidPrice())
// fmt.Println(ticker.DateTime())
// fmt.Println(ticker.TruncateDateTime(time.Second))

// apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
// order := &bitflyer.Order{
// 	ProductCode:     config.Config.ProductCode,
// 	ChildOrderType:  "LIMIT",
// 	Side:            "BUY",
// 	Price:           7000,
// 	Size:            0.01,
// 	MinuteToExpires: 1,
// 	TimeInForce:     "GTC",
// }
// res, _ := apiClient.SendOrder(order)
// fmt.Println(res.ChildOrderAcceptanceID)

//i := "JRF20181012-144016-140584"
//params := map[string]string{
//	"product_code": config.Config.ProductCode,
//	"child_order_acceptance_id": i,
//}
//r, _ := apiClient.ListOrder(params)
//fmt.Println(r)

// limit := 100
// duration := "1m"
// durationTime := config.Config.Durations[duration]
// df, _ := model.GetAllCandle(config.Config.ProductCode, durationTime, limit)
