package main

import (
	// go mod initで、プロジェクトのパスを登録しないとエラーになる。
	// "auto-trading/app/bitflyer"
	"auto-trading/app/config"
	"auto-trading/app/util"
	"fmt"
)

func main() {
	util.Logging(config.Config.LogFile)
	fmt.Println(config.Config.ApiKey, config.Config.ApiSecret)

}
