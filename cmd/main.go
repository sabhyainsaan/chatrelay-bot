package main

import (
	"chatrelaybot/internal/backend"
	"chatrelaybot/internal/slackbot"
	"chatrelaybot/internal/utils"
	"context"
	"log"
)

func main() {
	//To Initialize the OpenTelemetry
	shutdown := utils.InitTracer("chatrelaybot")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal("Error shutting down tracer:", err)
		}
	}()

	// For Starting the app
	backend.StartMockServer()
	slackbot.StartBot()
}
