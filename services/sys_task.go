package services

import (
	"context"
	"encoding/json"
	"fmt"
	"telebot_v2/global"
	"telebot_v2/model"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func AllVideosHandlerTaskV2(bot *tele.Bot, chat *tele.Chat) error {
	allProcessing.Lock()
	defer allProcessing.Unlock()

	if allProcessing.isProcessing {
		// Already running, return
		return nil
	}

	allProcessing.isProcessing = true
	defer func() {
		allProcessing.Lock()
		allProcessing.isProcessing = false
		allProcessing.Unlock()
	}()

	videos, err := VideosService.GetVideos()
	if err != nil {
		global.LOG.Error("GetVideos failed", zap.Error(err))
		return err
	}

	var releasedVideos []model.VideoReleaseMessage

	releasedVideos, err = filterReleasedVideo(videos)
	if err != nil {
		global.LOG.Error("filter videos failed", zap.Error(err))
		return err
	}

	videosCount := len(releasedVideos)
	pinMsg := fmt.Sprintf("正在检测过去48小时内共发布的 %v 个视频。。。", videosCount)
	msg, err := bot.Send(chat, pinMsg)
	bot.Pin(msg)

	for _, releasedVideo := range releasedVideos {
		releasedVideo.Total = totalVideos
		messageBytes, err := json.Marshal(releasedVideo)
		if err != nil {
			global.LOG.Error("decode message failed", zap.Error(err))
			return err
		}

		err = global.Writer.WriteMessages(
			context.Background(),
			kafka.Message{
				Topic: "video_read_all",
				Key:   []byte("video_read_all"),
				Value: messageBytes,
			},
		)
		if err != nil {
			return err
		}
	}

	return err
}

func SystemHealth(bot *tele.Bot, chat *tele.Chat) {
	text := "程序状态: OK✅ "

	_, err := bot.Send(chat, text)
	if err != nil {
		zap.L().Error("send message failed", zap.Error(err))
		return
	}
	return
}
