package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

func New(filename string) {
	var err error
	Viper = viper.New()
	Viper.SetConfigType("yaml")
	Viper.SetConfigName(filename)
	Viper.AddConfigPath(".")
	err = Viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file: %w", err))
	}
}
