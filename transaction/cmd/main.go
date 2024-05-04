package main

import (
	"log"

	"github.com/ShmelJUJ/software-engineering/transaction/config"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to create new config: ", err)
	}

	app.Run(cfg)
}
