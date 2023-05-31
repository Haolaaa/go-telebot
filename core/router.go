package core

import (
	"telebot_v2/global"
	"telebot_v2/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	systemRouter := router.RouterGroupApp.SystemRouter

	SystemGroup := Router.Group("")
	{
		SystemGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
			})
		})
	}
	systemRouter.InitVideosRouter(SystemGroup)
	global.LOG.Info("router register success")
	return Router
}
