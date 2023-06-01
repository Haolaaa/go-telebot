package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"telebot_v2/canal"
	"telebot_v2/global"
	"telebot_v2/model"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type ProcessingState struct {
	sync.Mutex
	isProcessing bool
}

var startProcessing = &ProcessingState{}
var allProcessing = &ProcessingState{}

func StartHandler(bot *tele.Bot, reader *kafka.Reader) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		startProcessing.Lock()
		defer startProcessing.Unlock()

		if startProcessing.isProcessing {
			bot.Send(ctx.Chat(), "已经在运行了")
			return nil
		}

		startProcessing.isProcessing = true
		go func() {
			defer func() {
				startProcessing.Lock()
				startProcessing.isProcessing = false
				startProcessing.Unlock()
			}()

			processKafkaMessages(bot, ctx.Chat(), reader)
		}()

		return nil
	}
}

func processKafkaMessages(bot *tele.Bot, chat *tele.Chat, reader *kafka.Reader) {
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			if err == io.EOF {
				log.Panicln("reached end of partition")
				time.Sleep(10 * time.Second)
				continue
			}
			log.Printf("error while reading message: %v", err)
			continue
		}

		log.Printf("read message at offfse %v from partition %v", message.Offset, message.Partition)

		var text model.VideoReleaseKafkaMessage
		err = json.Unmarshal(message.Value, &text)
		if err != nil {
			log.Printf("error while unmarshalling message: %v", err)
			continue
		}

		sendMessage := formatMessage(text)

		log.Println(sendMessage)

		_, err = bot.Send(chat, sendMessage, &tele.SendOptions{
			ParseMode: tele.ModeMarkdownV2,
		})
		if err != nil {
			log.Printf("error while sending message: %v", err)
			continue
		}
	}
}

func formatMessage(text model.VideoReleaseKafkaMessage) string {
	text.CreatedAt = strings.Replace(text.CreatedAt, "T", " ", 1)
	text.CreatedAt = strings.Split(text.CreatedAt, "+")[0]

	sendText := "***任务名称***: ``%s`` \n**发布站点**: `%v`\n**视频标题**: `%v`\n**视频ID**: `%v`\n**播放链接**: `%v`\n**直连状态**: %v\n**直连地址**: %v\n**CDN状态**: %v\n**CDN地址**: %v\n**CF状态**: %v\n**CF地址**: %v\n**上传时间**: %v\n"

	sendMessage := fmt.Sprintf(sendText, text.TaskName, text.PublishedSiteName, text.Title, text.VideoId, text.PlayUrl, text.DirectPlayUrlStatus, text.DirectPlayUrl, text.CDNPlayUrlStatus, text.CDNPlayUrl, text.DirectPlayUrlStatus, text.CFPlayUrl, text.CreatedAt)

	sendMessage = strings.Replace(sendMessage, "-", "\\-", -1)
	sendMessage = strings.Replace(sendMessage, ".", "\\.", -1)

	return sendMessage
}

func UserHandler(bot *tele.Bot) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		text := fmt.Sprintf("user: %v", ctx.Sender().ID)
		_, err := bot.Send(ctx.Chat(), text)
		return err
	}
}

func ChatHandler(bot *tele.Bot) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		text := fmt.Sprintf("chat: %v", ctx.Chat().ID)
		_, err := bot.Send(ctx.Chat(), text)
		return err
	}
}

func HealthHandler(bot *tele.Bot) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		text := fmt.Sprintln("💊Server Health: OK✅")
		_, err := bot.Send(ctx.Chat(), text)
		return err
	}
}

func AllVideosHandler(bot *tele.Bot) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		allProcessing.Lock()
		defer allProcessing.Unlock()

		if allProcessing.isProcessing {
			bot.Send(ctx.Chat(), "已经在运行了")
			return nil
		}

		allProcessing.isProcessing = true
		errChan := make(chan error, 1)
		go func() {
			defer func() {
				allProcessing.Lock()
				allProcessing.isProcessing = false
				allProcessing.Unlock()
			}()
			videos, err := VideosService.GetVideos()
			if err != nil {
				global.LOG.Error("GetVideos failed", zap.Error(err))
				errChan <- err
				return
			}

			var releasedVideos []model.VideoReleaseMessage

			releasedVideos, err = filterReleasedVideo(videos)
			if err != nil {
				global.LOG.Error("filter videos failed", zap.Error(err))
				errChan <- err
				return
			}

			videosCount := len(releasedVideos)
			pinMsg := fmt.Sprintf("📺当前共有 %v 个视频", videosCount)

			var msg tele.Editable
			msg, err = bot.Send(ctx.Chat(), pinMsg)
			if err != nil {
				global.LOG.Error("send message failed", zap.Error(err))
				errChan <- err
			}

			err = bot.Pin(msg)
			if err != nil {
				global.LOG.Error("pin message failed", zap.Error(err))
				errChan <- err
			}

			for _, releasedVideo := range releasedVideos {

				messageBytes, err := json.Marshal(releasedVideo)
				if err != nil {
					global.LOG.Error("decode message failed", zap.Error(err))
					errChan <- err
					return
				}

				err = global.Writer.WriteMessages(
					context.Background(),
					kafka.Message{
						Topic: "video_read_all",
						Key:   []byte("video_read_all"),
						Value: messageBytes,
					},
				)
			}
		}()

		select {
		case err := <-errChan:
			return err
		default:
			return nil
		}
	}
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
