package global

import (
	"telebot_v2/config"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	DB2    *gorm.DB
	VP     *viper.Viper
	LOG    *zap.Logger
	CONFIG config.Config
	Writer *kafka.Writer
	Bot    *tele.Bot
)
