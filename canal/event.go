package canal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"telebot_v2/global"
	"telebot_v2/model"
	"telebot_v2/utils"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	StatusPublished = 1
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func GetParseValue(event *canal.RowsEvent) (*canal.RowsEvent, map[string]interface{}) {
	var rows = map[string]interface{}{}

	for colIndex, currCol := range event.Table.Columns {
		colValue := event.Rows[len(event.Rows)-1][colIndex]
		switch colValue.(type) {
		case []uint8:
			colValue = utils.Uint8TosString(colValue.([]uint8)...)
		case decimal.Decimal:
			v, _ := colValue.(decimal.Decimal).Float64()
			colValue = float32(v)
		case string:
			switch currCol.Name {
			case "created_at", "updated_at":
				if colValue.(string) != "" && event.Table.Name == "video_changes" {
					var err error
					colValue, err = time.Parse("2006-01-02 15:04:05", colValue.(string))
					if err != nil {
						logError("time.Parse failed", err)
					}
				}
			}
		}

		rows[currCol.Name] = colValue
	}

	return event, rows
}

func TableEventDispatcher(event *canal.RowsEvent, row map[string]interface{}) {
	if global.Writer == nil {
		logError("global.Writer is nil", nil)
		return
	}

	modelStruct := model.GetModelStruct(event.Table.Name)
	rowModel, ok := MapStructureRow(event, modelStruct, row)
	if !ok {
		return
	}
	switch event.Action {
	case canal.InsertAction, canal.UpdateAction:
		if event.Table.Name == "video_release" {
			rowModel := rowModel.(model.VideoRelease)
			if rowModel.Status == StatusPublished {
				video, err := GetVideo(rowModel.VideoId)
				if err != nil {
					logError("GetVideo failed", err)
					return
				}

				siteVideoUrls, err := utils.GetSitePlayUrls(int(rowModel.SiteID))
				if err != nil {
					global.LOG.Error("GetSitePlayUrls failed", zap.Error(err))

					return
				}

				playUrl := FormatM3u8Url(video.PlayUrl, siteVideoUrls.SiteKey)
				downUrl := utils.FormatUrl(video.DownUrl)
				coverUrl := utils.FormatUrl(video.Cover)

				message := model.VideoReleaseMessage{
					PublishedSiteName: siteVideoUrls.SiteName,
					VideoId:           int(video.ID),
					Title:             video.Title,
					CreatedAt:         video.CreatedAt,
				}

				if siteVideoUrls.CDNPlayUrl != "" {
					message.CDNPlayUrl = siteVideoUrls.CDNPlayUrl + playUrl
				}
				if siteVideoUrls.CFPlayUrl != "" {
					message.CFPlayUrl = siteVideoUrls.CFPlayUrl + playUrl
				}
				if siteVideoUrls.DirectPlayUrl != "" {
					message.DirectPlayUrl = siteVideoUrls.DirectPlayUrl + playUrl
				}
				if downUrl != "" {
					message.DownUrl = downUrl
				}
				if coverUrl != "" {
					message.CoverUrl = coverUrl
				}

				fmt.Println(message)

				messageBytes, err := json.Marshal(message)
				if err != nil {
					logError("json.Marshal failed", err)
					return
				}

				err = global.Writer.WriteMessages(
					context.Background(),
					kafka.Message{
						Topic: "video_release",
						Key:   []byte("video_release"),
						Value: messageBytes,
					},
				)
				if err != nil {
					logError("kafka.Writer.WriteMessages failed", err)
					return
				}
			}
		}
	}
}

func MapStructureRow(event *canal.RowsEvent, model interface{}, row map[string]interface{}) (interface{}, bool) {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &model,
	})
	err := decoder.Decode(row)
	if err != nil {
		return nil, false
	}

	return model, true
}

func FormatM3u8Url(url string, siteKey string) string {
	url = strings.Replace(url, "/data/org", "", -1)
	if strings.Contains(url, "index.m3u8") {
		url = strings.Replace(url, "index.m3u8", siteKey+".m3u8", -1)
	}

	return url
}

func GetVideo(videoId int) (video model.Video, err error) {
	err = global.DB.Table("video").Where("id = ?", videoId).First(&video).Error
	return
}

func logError(message string, err error) {
	global.LOG.Error(message, zap.Error(err))
}
