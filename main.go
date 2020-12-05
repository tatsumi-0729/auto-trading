package main

import (
	// go mod initで、プロジェクトのパスを登録しないとエラーになる。
	"auto-trading/config"
	"fmt"
)

func main() {
	fmt.Println(config.Config.ApiKey, config.Config.ApiSecret)
}
