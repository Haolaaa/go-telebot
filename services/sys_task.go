package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"telebot_v2/global"
	"telebot_v2/model"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func AllVideosHandlerTaskV2() error {
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

	chat, err := global.Bot.ChatByID(-1001954537168)
	_, err = global.Bot.Send(chat, pinMsg)
	if err != nil {
		zap.L().Error("send message failed", zap.Error(err))
		return err
	}

	if err != nil {
		zap.L().Error("pin message failed", zap.Error(err))
		return err
	}

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

func SystemHealth() {
	text := "程序状态: OK✅ "

	chat, err := global.Bot.ChatByID(-1001954537168)
	_, err = global.Bot.Send(chat, text)
	if err != nil {
		zap.L().Error("send message failed", zap.Error(err))
		return
	}
	return
}

func HW() {
	// Calculate the next run time
	nextRunTime := time.Now().Add(4 * time.Hour)

	// Send a message to the Telegram chat
	msg := fmt.Sprintf("Task completed. The next task will run at %s.", nextRunTime.Format("15:04"))

	chat, err := global.Bot.ChatByID(-1001954537168)
	if err != nil {
		log.Fatal(err)
	}

	sentMsg, err := global.Bot.Send(chat, msg)

	global.Bot.UnpinAll(chat)
	global.Bot.Pin(sentMsg)
}
