package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"telebot_v2/global"
	"telebot_v2/model"

	"go.uber.org/zap"
)

type Response struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type IntermediateResponse struct {
	SiteConfig model.SiteVideoUrls `json:"siteConfig"`
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
