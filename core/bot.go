package core

import (
	"telebot_v2/global"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func NewBot(token string) (*tele.Bot, error) {
	pref := tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		global.LOG.Error("NewBot failed", zap.Error(err))
		return nil, err
	}

	return b, nil
}
