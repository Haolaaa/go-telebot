package core

import (
	"telebot_v2/global"

	"go.uber.org/zap"
)

func RunServer() {
	var err error

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

	err = RunCanal(true)
	if err != nil {
		global.LOG.Error("Run canal failed", zap.Error(err))
	}
}
