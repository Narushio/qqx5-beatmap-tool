package main

import (
	"github.com/Narushio/qqx5-beatmap-tool/config"
	"github.com/Narushio/qqx5-beatmap-tool/server"
)

func main() {
	config.Init("config")
	server.Init()
}
