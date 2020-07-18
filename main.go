package main

import (
	"fmt"
	"log"
	"os"
	"whatthecard/pkg/game"
	"whatthecard/pkg/logger"
	"whatthecard/pkg/server"
)

func main() {
	logLevel := os.Getenv("LOGLEVEL")
	logger := logger.NewLogger(logLevel)
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	hub := server.NewHub(logger)
	gameService := game.NewService(logger)
	server := server.New(hub, gameService, logger)

	if err := server.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
