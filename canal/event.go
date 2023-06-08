package canal

import (
	"bytes"
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

type Response struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type IntermediateResponse struct {
	SiteConfig model.SiteVideoUrls `json:"siteConfig"`
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

				siteVideoUrls, err := GetSitePlayUrls(int(rowModel.SiteID))
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

func GetSitePlayUrls(siteId int) (sitePlayUrls model.SiteVideoUrls, err error) {
	reqBody := map[string]interface{}{
		"siteID":     siteId,
		"parentName": "集团",
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		global.LOG.Error("json.Marshal failed", zap.Error(err))
		return
	}

	req, err := http.NewRequest("POST", "https://3yzt.com/siteConfig/getByID", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		global.LOG.Error("http.NewRequest failed", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		global.LOG.Error("httpClient.Do failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected HTTP status: %v", resp.StatusCode)
		global.LOG.Error("resp.StatusCode != http.StatusOK", zap.Error(err))
		return
	}

	var body Response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		global.LOG.Error("json.NewDecoder.Decode failed", zap.Error(err))
		return
	}

	if body.Code != 0 {
		err = fmt.Errorf("unexpected body code: %v", body.Code)
		global.LOG.Error("body.Code != 0", zap.Error(err))
		return
	}

	var interResp IntermediateResponse
	err = json.Unmarshal(body.Data, &interResp)
	if err != nil {
		global.LOG.Error("json.Unmarshal failed", zap.Error(err))
		return
	}

	sitePlayUrls = interResp.SiteConfig

	return
}

func logError(message string, err error) {
	global.LOG.Error(message, zap.Error(err))
}
