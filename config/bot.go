package config

type Bot struct {
	Token  string `mapstructure:"token" json:"token" yaml:"token"`
	ChatID int64  `mapstructure:"chat_id" json:"chat_id" yaml:"chat_id"`
}
