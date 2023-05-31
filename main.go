package main

import (
	"os"
	"os/signal"
	"syscall"
	"telebot_v2/core"
	"telebot_v2/global"

	"go.uber.org/zap"
)

func main() {
	global.VP = core.Viper()

	global.LOG = core.Zap()
	zap.ReplaceGlobals(global.LOG)

	global.DB = core.Gorm()
	if global.DB != nil {
		db, _ := global.DB.DB()
		core.RegisterTables()
		defer db.Close()
	}

	var err error
	global.Bot, err = core.NewBot("6026225894:AAHksPMok37YLUMJobEPbeoRkYPk6Itxs_Q")
	if err != nil {
		global.LOG.Error("NewBot failed", zap.Error(err))
	}

	go global.Bot.Start()

	go core.RunServer()

	core.NewKafka()

	core.Reader()

	global.Writer = core.Writer()

	err = core.RunCanal(true)
	if err != nil {
		global.LOG.Error("Run canal failed", zap.Error(err))
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		global.Bot.Stop()
		os.Exit(1)
	}()

	select {}
}
