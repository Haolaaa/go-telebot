package core

import (
	"fmt"
	"telebot_v2/core/internal"
	"telebot_v2/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	fmt.Println("Viper init...")
	v := viper.New()
	v.SetConfigFile(internal.ConfigDefaultFile)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err = v.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	if err = v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
