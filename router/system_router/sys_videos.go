package system_router

import (
	"telebot_v2/api"

	"github.com/gin-gonic/gin"
)

type VideosRouter struct{}

func (s *VideosRouter) InitVideosRouter(Router *gin.RouterGroup) {
	videoRouter := Router.Group("system")
	apiGroup := api.ApiGroupApp.SystemApiGroup.SysVideosApi
	{
		videoRouter.GET("videos", apiGroup.GetVideos) // 获取48小时内发布的视频列表
	}
}
