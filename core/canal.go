package core

import (
	myCanal "telebot_v2/canal"
	"telebot_v2/global"

	"github.com/go-mysql-org/go-mysql/canal"
	"go.uber.org/zap"
)

func getCanalConfig() *canal.Config {
	cfg := canal.NewDefaultConfig()

	cfg.Addr = global.CONFIG.Canal.Addr
	cfg.User = global.CONFIG.Canal.User
	cfg.Password = global.CONFIG.Canal.Password
	cfg.Flavor = global.CONFIG.Canal.Flavor
	cfg.Charset = global.CONFIG.Canal.Charset
	cfg.ServerID = global.CONFIG.Canal.ServerID
	cfg.ParseTime = true

	cfg.Dump.TableDB = global.CONFIG.Canal.Dump.TableDB
	cfg.Dump.DiscardErr = global.CONFIG.Canal.Dump.DiscardErr
	cfg.Dump.SkipMasterData = global.CONFIG.Canal.Dump.SkipMasterData
	cfg.Dump.Tables = global.CONFIG.Canal.Dump.Tables

	return cfg
}

func RunCanal(isPos bool) error {
	cfg := getCanalConfig()
	c, err := canal.NewCanal(cfg)

	if err != nil {
		panic(err)
	}

	c.SetEventHandler(&myCanal.EventHandler{})

	if !isPos {
		return c.Run()
	}

	masterPos, err := c.GetMasterPos()
	if err != nil {
		global.LOG.Error("GetMasterPos failed", zap.Error(err))
	}

	return c.RunFrom(masterPos)
}
