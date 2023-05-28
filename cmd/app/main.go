package main

import (
	"github.com/DenisTaztdinov/PhoneBook/config"
	"github.com/DenisTaztdinov/PhoneBook/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
