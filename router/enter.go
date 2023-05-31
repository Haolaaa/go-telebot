package router

import "telebot_v2/router/system_router"

type RouterGroup struct {
	SystemRouter system_router.SystemRouter
}

var RouterGroupApp = new(RouterGroup)
