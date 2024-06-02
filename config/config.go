package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

func Init(filename string) {
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

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}
