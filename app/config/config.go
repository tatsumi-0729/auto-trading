package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey                  string
	ApiSecret               string
	LogFile                 string
	ProductCode             string
	BaseUrl                 string
	GetBalanceUrl           string
	GetTickerUrl            string
	GetRealTimeTickerHost   string
	GetRealTimeTickerSchema string
	GetRealTimeTickerPath   string
	SendOrderUrl            string
	ListOrderUrl            string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("fail to load Config", err)
	}
	Config = ConfigList{
		ApiKey:                  cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:               cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:                 cfg.Section("auto-trading").Key("log_file").String(),
		ListOrderUrl:            cfg.Section("auto-trading").Key("list_order_url").String(),
		BaseUrl:                 cfg.Section("auto-trading").Key("base_url").String(),
		GetBalanceUrl:           cfg.Section("auto-trading").Key("get_balance_url").String(),
		GetTickerUrl:            cfg.Section("auto-trading").Key("get_ticker_url").String(),
		GetRealTimeTickerHost:   cfg.Section("auto-trading").Key("get_realtime_ticker_host").String(),
		GetRealTimeTickerSchema: cfg.Section("auto-trading").Key("get_realtime_ticker_schema").String(),
		GetRealTimeTickerPath:   cfg.Section("auto-trading").Key("get_realtime_ticker_path").String(),
		SendOrderUrl:            cfg.Section("auto-trading").Key("send_order_url").String(),
		ProductCode:             cfg.Section("auto-trading").Key("product_code").String(),
	}
}
