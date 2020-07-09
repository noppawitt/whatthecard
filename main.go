package main

import (
	"log"
	"whatthecard/game"
	"whatthecard/logger"
	"whatthecard/server"
)

func main() {
	logger := logger.NewLogger("debug")

	hub := server.NewHub(logger)
	gameService := game.NewService(logger)
	server := server.New(hub, gameService, logger)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
