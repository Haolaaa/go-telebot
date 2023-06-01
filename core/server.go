package core

import (
	"fmt"
	"net/http"
	"telebot_v2/global"
	"telebot_v2/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func RunServer() {
	var err error
	Router := Routers()

	global.Bot, err = NewBot("6026225894:AAHksPMok37YLUMJobEPbeoRkYPk6Itxs_Q")
	if err != nil {
		global.LOG.Error("NewBot failed", zap.Error(err))
	}
	go global.Bot.Start()

	kafka, err := NewKafkaService()
	if err != nil {
		global.LOG.Error("NewKafkaService failed", zap.Error(err))
	}

	global.Writer = kafka.Writer
	kafka.Reader()

	cron := cron.New(cron.WithSeconds())
	_, err = cron.AddFunc("0 0 */4 * * *", func() {
		fmt.Println("AllVideosHandlerTaskV2 started")
		err := services.AllVideosHandlerTaskV2()
		if err != nil {
			global.LOG.Error("AllVideosHandlerTaskV2 failed", zap.Error(err))
		}
	})
	if err != nil {
		global.LOG.Error("AddFunc failed", zap.Error(err))
	}
	_, err = cron.AddFunc("0 0 * * * *", services.SystemHealth)
	if err != nil {
		global.LOG.Error("AddFunc failed", zap.Error(err))
	}
	cron.Start()

	err = RunCanal(true)
	if err != nil {
		global.LOG.Error("Run canal failed", zap.Error(err))
	}

	s := initServer(":8082", Router)
	global.LOG.Info("server run success on ", zap.String("address", "8082"))

	global.LOG.Error(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
}
