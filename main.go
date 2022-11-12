package main

import (
	"log"

	"github.com/XXena/chat/internal/app"

	"github.com/XXena/chat/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
