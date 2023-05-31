package utils

import (
	"net/http"
	"strings"
	"sync"
	"telebot_v2/global"
	"telebot_v2/model"
	"time"

	"go.uber.org/zap"
)

var wg sync.WaitGroup

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func ValidateUrls(videoReleaseMessage *model.VideoReleaseMessage, releasedVideo *model.VideoReleaseStatus) {
	wg.Add(3)
	go validateUrl(videoReleaseMessage.PlayUrl, &releasedVideo.CDNPlayUrlStatus)
	go validateUrl(videoReleaseMessage.PlayUrl, &releasedVideo.CFPlayUrlStatus)
	go validateUrl(videoReleaseMessage.PlayUrl, &releasedVideo.DirectPlayUrlStatus)
	wg.Wait()
}

func validateUrl(url string, status *string) {
	defer wg.Done()
	if url != "" {

	}
}

func ValidateM3u8Urls(playUrl string, status *string) {
	resp, err := httpClient.Get(playUrl)
	if err != nil {
		global.LOG.Error("http get failed", zap.Error(err))
		if strings.Contains(err.Error(), "Client.Timeout") {
			*status = "超时❌"
		} else {
			*status = "无法访问❌"
		}
		return
	}

	defer resp.Body.Close()

	if strings.Contains(resp.Status, "20") {
		*status = "正常✅"
	} else {
		*status = "无法访问❌"
	}
}
