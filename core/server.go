package core

import (
	"net/http"
	"telebot_v2/global"
	"time"

	"github.com/gin-gonic/gin"
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

	NewKafka()
	global.Writer = Writer()
	Reader()

	err = RunCanal(true)
	if err != nil {
		global.LOG.Error("Run canal failed", zap.Error(err))
	}

	s := initServer(":8082", Router)
	time.Sleep(10 * time.Microsecond)
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
