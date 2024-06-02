package view

import (
	"embed"
	"io/fs"
	"os"

	"github.com/Narushio/qqx5-beatmap-tool/config"
)

var (
	//go:embed *.html
	embedFS embed.FS

	FS fs.FS
)

func New() {
	FS = os.DirFS("./view")
	if config.Viper.GetString("environment") != "development" {
		FS = embedFS
	}
}
