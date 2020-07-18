package main

import (
	"log"
	"os"
	"whatthecard/pkg/game"
	"whatthecard/pkg/logger"
	"whatthecard/pkg/server"
)

func main() {
	logLevel := os.Getenv("LOGLEVEL")
	logger := logger.NewLogger(logLevel)
	hub := server.NewHub(logger)
	gameService := game.NewService(logger)
	server := server.New(hub, gameService, logger)

	if err := server.Start(":4000"); err != nil {
		log.Fatal(err)
	}
}
