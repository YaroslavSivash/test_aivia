package main

import (
	"log"

	"test_aivia/internal/config"
	"test_aivia/internal/server"
)

func main() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	app, err := server.NewApp(conf)
	if err != nil {
		log.Fatal(err)
	}
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
