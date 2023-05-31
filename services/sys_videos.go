package services

import (
	"telebot_v2/global"
	"telebot_v2/model"
	"time"

	"go.uber.org/zap"
)

type SysVideosService struct{}

var VideosService = new(SysVideosService)

func (s *SysVideosService) GetVideos() (videos []model.Video, err error) {
	now := time.Now()
	twoDaysAgo := now.Add(-24 * time.Hour * 6)

	err = global.DB.Table("video").Where("created_at BETWEEN ? AND ?", twoDaysAgo, now).Find(&videos).Error
	if err != nil {
		global.LOG.Error("GetVideos failed", zap.Error(err))
		return nil, err
	}

	return
}
