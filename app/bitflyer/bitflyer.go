package bitflyer

import (
	"auto-trading/app/config"
)

func Bitflyer() string {
	baseUrl := config.Config.BaseUrl
	return baseUrl
}
