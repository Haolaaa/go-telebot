package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"telebot_v2/model"
	"time"

	"github.com/segmentio/kafka-go"
	tele "gopkg.in/telebot.v3"
)

var wg sync.WaitGroup

func StartHandler(bot *tele.Bot, chat *tele.Chat, reader *kafka.Reader) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		go processKafkaMessages(bot, chat, reader)

		return nil
	}
}

func processKafkaMessages(bot *tele.Bot, chat *tele.Chat, reader *kafka.Reader) {
	// var message string

	for {
		message, err := reader.ReadMessage(context.Background())
		log.Printf("read message at offfse %v from partition %v", message.Offset, message.Partition)

		if err != nil {
			if err == io.EOF {
				log.Panicln("reached end of partition")
				time.Sleep(10 * time.Second)
				continue
			}
			log.Printf("error while reading message: %v", err)
			continue
		}

		var text model.VideoReleaseKafkaMessage

		err = json.Unmarshal(message.Value, &text)
		if err != nil {
			log.Printf("error while unmarshalling message: %v", err)
			continue
		}

		sendText := "***任务名称***: ``%s`` \n**发布站点**: `%v`\n**视频标题**: `%v`\n**视频ID**: `%v`\n**播放链接**: `%v`\n**直连状态**: %v\n**直连地址**: %v\n**CDN状态**: %v\n**CDN地址**: %v\n**CF状态**: %v\n**CF地址**: %v\n**上传时间**: %v\n"

		text.CreatedAt = strings.Replace(text.CreatedAt, "T", " ", 1)
		text.CreatedAt = strings.Split(text.CreatedAt, "+")[0]

		sendMessage := fmt.Sprintf(sendText, text.TaskName, text.PublishedSiteName, text.Title, text.VideoId, text.PlayUrl, text.DirectPlayUrlStatus, text.DirectPlayUrl, text.CDNPlayUrlStatus, text.CDNPlayUrl, text.DirectPlayUrlStatus, text.CFPlayUrl, text.CreatedAt)

		sendMessage = strings.Replace(sendMessage, "-", "\\-", -1)
		sendMessage = strings.Replace(sendMessage, ".", "\\.", -1)
		fmt.Println(sendMessage)

		_, err = bot.Send(chat, sendMessage, &tele.SendOptions{
			ParseMode: tele.ModeMarkdownV2,
		})
		if err != nil {
			log.Printf("error while sending message: %v", err)
			continue
		}
	}
}
