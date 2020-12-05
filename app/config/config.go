package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey    string
	ApiSecret string
	LogFile   string
	BaseUrl   string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("fail to load Config", err)
	}
	Config = ConfigList{
		ApiKey:    cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:   cfg.Section("auto-trading").Key("log_file").String(),
		BaseUrl:   cfg.Section("auto-trading").Key("base_url").String(),
	}
}
