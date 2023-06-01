package system_api

import (
	"telebot_v2/canal"
	"telebot_v2/global"
	"telebot_v2/model"
	"telebot_v2/model/response"
	"telebot_v2/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysVideosApi struct{}

func (a *SysVideosApi) GetVideos(c *gin.Context) {
	videos, err := services.VideosService.GetVideos()
	if err != nil {
		global.LOG.Error("GetVideos failed", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	var releasedVideos []model.VideoReleaseMessage

	if videos == nil {
		releasedVideos = []model.VideoReleaseMessage{}
	} else {
		releasedVideos, err = filterReleasedVideo(videos)
		if err != nil {
			global.LOG.Error("GetVideos failed", zap.Error(err))
			response.FailWithMessage("获取失败", c)
			return
		}
	}

	response.OkWithDetailed(releasedVideos, "获取成功", c)
}

func filterReleasedVideo(videos []model.Video) (filteredVideos []model.VideoReleaseMessage, err error) {
	var releasedVideos []model.VideoRelease

	err = global.DB.Table("video_release").Find(&releasedVideos).Error
	if err != nil {
		global.LOG.Error("GetVideos failed", zap.Error(err))
		return nil, err
	}

	for _, video := range videos {
		for _, releasedVideo := range releasedVideos {
			if releasedVideo.Status == 1 && video.ID == uint(releasedVideo.VideoId) {
				publishedSite, err := getPublishedSiteName(int(releasedVideo.SiteID))
				if err != nil {
					global.LOG.Error("GetVideos failed", zap.Error(err))
					return nil, err
				}

				playUrl := canal.FormatM3u8Url(video.PlayUrl, publishedSite.SiteKey)

				filteredVideos = append(filteredVideos, model.VideoReleaseMessage{
					PublishedSiteName: publishedSite.SiteName,
					VideoId:           int(video.ID),
					Title:             video.Title,
					PlayUrl:           video.PlayUrl,
					DirectPlayUrl:     publishedSite.DirectPlayUrl + playUrl,
					CFPlayUrl:         publishedSite.CFPlayUrl + playUrl,
					CDNPlayUrl:        publishedSite.CDNPlayUrl + playUrl,
					CreatedAt:         video.CreatedAt,
				})
			}
		}
	}

	return
}

func getPublishedSiteName(publishedSiteId int) (publishedSite model.SiteVideoUrls, err error) {
	err = global.DB.Where("site_id = ?", publishedSiteId).Find(&publishedSite).Error
	if err != nil {
		global.LOG.Error("GetVideos failed", zap.Error(err))
		return
	}

	return
}
