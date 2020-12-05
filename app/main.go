package main

import (
	// go mod initで、プロジェクトのパスを登録しないとエラーになる。
	"auto-trading/app/bitflyer"
	"auto-trading/app/config"
	"auto-trading/app/util"
	"fmt"
	"log"
)

func main() {
	util.Logging(config.Config.LogFile)
	log.Println(bitflyer.Bitflyer())
	fmt.Println(config.Config.ApiKey, config.Config.ApiSecret)

}
