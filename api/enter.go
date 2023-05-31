package api

import "telebot_v2/api/system_api"

type ApiGroup struct {
	SystemApiGroup system_api.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
