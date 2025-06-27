package main

import (
	"chatrelaybot/internal/backend"
	"chatrelaybot/internal/slackbot"
	"chatrelaybot/internal/utils"
	"context"
	"log"
)

func main() {
	// Initialize OpenTelemetry
	shutdown := utils.InitTracer("chatrelaybot")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal("Error shutting down tracer:", err)
		}
	}()

	// Start app
	backend.StartMockServer()
	slackbot.StartBot()
}
