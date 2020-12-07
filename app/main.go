package main

import (
	// go mod initで、プロジェクトのパスを登録しないとエラーになる。
	"auto-trading/app/bitflyer"
	"auto-trading/app/config"
	"auto-trading/app/util"
	"fmt"
)

func main() {
	util.Logging(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	fmt.Println(apiClient.GetBalance())
	fmt.Println(apiClient.GetTicker("FX_BTC_JPY"))
}
