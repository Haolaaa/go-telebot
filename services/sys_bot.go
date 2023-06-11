package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"telebot_v2/canal"
	"telebot_v2/global"
	"telebot_v2/model"
	"telebot_v2/utils"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type ProcessingState struct {
	sync.Mutex
	isProcessing bool
}

var (
	startProcessing = &ProcessingState{}
	allProcessing   = &ProcessingState{}
	totalVideos     int
	allVideosTaskID *cron.EntryID = new(cron.EntryID)
)

func RunHandler(bot *tele.Bot) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {

		fmt.Println(*allVideosTaskID)
		if *allVideosTaskID != 0 {
			bot.Send(ctx.Chat(), "已经在运行了")
			return nil
		}

		cron := cron.New(cron.WithSeconds())

		bot.Send(ctx.Chat(), "正在启动。。。")
		var err error
		*allVideosTaskID, err = cron.AddFunc("0 0 */4 * * *", func() {
			err := AllVideosHandlerTaskV2(bot, ctx.Chat())
			if err != nil {
				global.LOG.Error("AllVideosHandlerTaskV2 failed", zap.Error(err))
			}
		})

		if err != nil {
			global.LOG.Error("AddFunc failed", zap.Error(err))
		}
		_, err = cron.AddFunc("@hourly", func() {
			SystemHealth(bot, ctx.Chat())
		})
		if err != nil {
			global.LOG.Error("AddFunc failed", zap.Error(err))
		}

		cron.Start()

		bot.Send(ctx.Chat(), "启动成功")

		return nil
	}
}

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
				if r := recover(); r != nil {
					global.LOG.Error("Recovered from panic", zap.Any("panic", r))
				}

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
				global.LOG.Info("EOF, waiting for new messages")
				time.Sleep(10 * time.Second)
				continue
			}
			global.LOG.Error("error while reading message", zap.Error(err))
			continue
		}

		var text model.VideoReleaseKafkaMessage
		err = json.Unmarshal(message.Value, &text)
		if err != nil {
			global.LOG.Error("error while unmarshalling message", zap.Error(err))
			continue
		}

		sendMessage := formatMessage(text)

		if text.TaskName == "周期任务" {
			_, err = bot.Send(chat, sendMessage, &tele.SendOptions{
				ParseMode: tele.ModeMarkdownV2,
			})
			if err != nil {
				global.LOG.Error("error while sending message", zap.Error(err))
				continue
			}
		} else {
			_, err = bot.Send(chat, sendMessage, &tele.SendOptions{
				ParseMode: tele.ModeMarkdownV2,
			})
			if err != nil {
				global.LOG.Error("error while sending message", zap.Error(err))
				continue
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func formatMessage(text model.VideoReleaseKafkaMessage) string {
	text.CreatedAt = strings.Replace(text.CreatedAt, "T", " ", 1)
	text.CreatedAt = strings.Split(text.CreatedAt, "+")[0]

	sendText :=
		"***任务名称***: ``%s`` \n" +
			"**发布站点**: `%v`\n" +
			"**视频标题**: `%v`\n" +
			"**视频ID**: `%v`\n" +
			"**播放链接**: `%v`\n" +
			"**直连状态**: %v\n" +
			"**直连地址**: %v\n" +
			"**CDN状态**: %v\n" +
			"**CDN地址**: %v\n" +
			"**CF状态**: %v\n" +
			"**CF地址**: %v\n" +
			"**下载地址状态**: %v\n" +
			"**下载地址**: %v\n" +
			"**视频封面状态**: %v\n" +
			"**视频封面地址**: %v\n" +
			"**上传时间**: %v\n" +
			" @a\\_lan23"

	sendMessage := fmt.Sprintf(
		sendText,
		text.TaskName,
		text.PublishedSiteName,
		text.Title,
		text.VideoId,
		text.PlayUrl,
		text.DirectPlayUrlStatus,
		text.DirectPlayUrl,
		text.CDNPlayUrlStatus,
		text.CDNPlayUrl,
		text.DirectPlayUrlStatus,
		text.CFPlayUrl,
		text.DownUrlStatus,
		text.DownUrl,
		text.CoverUrlStatus,
		text.CoverUrl,
		text.CreatedAt,
	)

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
			bot.Send(ctx.Chat(), "程序已经在运行了")
			return nil
		}

		allProcessing.isProcessing = true
		errChan := make(chan error, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					global.LOG.Error("Recovered from panic", zap.Any("panic", r))
				}
				allProcessing.Lock()
				allProcessing.isProcessing = false
				allProcessing.Unlock()
			}()
			processVideos(bot, ctx, errChan)
		}()

		select {
		case err := <-errChan:
			return err
		default:
			return nil
		}
	}
}

func processVideos(bot *tele.Bot, ctx tele.Context, errChan chan<- error) {
	videos, err := VideosService.GetVideos()
	if err != nil {
		zap.L().Error("GetVideos failed", zap.Error(err))
		errChan <- err
		return
	}

	var releasedVideos []model.VideoReleaseMessage

	releasedVideos, err = filterReleasedVideo(videos)
	if err != nil {
		zap.L().Error("filter videos failed", zap.Error(err))
		errChan <- err
		return
	}

	for _, releasedVideo := range releasedVideos {

		releasedVideo.Total = totalVideos

		fmt.Println("released video: ", releasedVideo)

		messageBytes, err := json.Marshal(releasedVideo)
		if err != nil {
			zap.L().Error("decode message failed", zap.Error(err))
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
				publishedSite, err := utils.GetSitePlayUrls(int(releasedVideo.SiteID))
				fmt.Println("published site: ", publishedSite)
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
					DownUrl:           publishedSite.DownloadUrl + video.DownUrl,
					CoverUrl:          publishedSite.VideoCover + video.Cover,
					CreatedAt:         video.CreatedAt,
				})
			}
		}
	}

	return
}
