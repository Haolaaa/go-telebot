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
		// core.RegisterTables()
		defer db.Close()
	}

	core.RunServer()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		global.Bot.Stop()
		os.Exit(1)
	}()

	select {}
}
