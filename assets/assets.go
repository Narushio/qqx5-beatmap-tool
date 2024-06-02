package assets

import (
	"embed"
	"io/fs"
	"os"

	"github.com/Narushio/qqx5-beatmap-tool/config"
)

var (
	//go:embed all:*
	embedFS embed.FS

	FS fs.FS
)

func New() {
	FS = os.DirFS("./assets")
	if config.Viper.GetString("environment") != "development" {
		FS = embedFS
	}
}
