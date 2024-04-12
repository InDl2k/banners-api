package main

import (
	"banners/app"
	"banners/internal/config"
)

func main() {

	cfg := config.MustLoad()

	a := app.App{}

	a.Initialize()
	a.Run(cfg.HTTPServer.Address)

}
