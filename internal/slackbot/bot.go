package slackbot

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Tracer for Slack interactions
var tracer = otel.Tracer("chatrelay/slack")

// StartBot connects to Slack via Socket Mode and listens for events
func StartBot() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	backendURL := os.Getenv("BACKEND_URL")

	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))
	client := socketmode.New(api)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}
				client.Ack(*evt.Request)

				switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					go handleMention(api, ev, backendURL)
				}
			}
		}
	}()

	log.Println("Slack bot running...")
	client.Run()
}

// handleMention processes Slack app mentions and talks to the backend
func handleMention(api *slack.Client, event *slackevents.AppMentionEvent, backendURL string) {
	// Start tracing span
	_, span := tracer.Start(context.Background(), "HandleSlackMention")
	defer span.End()

	query := strings.TrimSpace(strings.Replace(event.Text, "<@"+event.User+">", "", -1))

	// Add trace metadata
	span.SetAttributes(
		attribute.String("slack.user_id", event.User),
		attribute.String("slack.channel", event.Channel),
		attribute.String("slack.query", query),
	)

	// Build request body
	reqBody := map[string]string{
		"user_id": event.User,
		"query":   query,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Call mock backend
	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Backend error:", err)
		api.PostMessage(event.Channel, slack.MsgOptionText("Error calling backend.", false))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer resp.Body.Close()

	var res struct {
		FullResponse string `json:"full_response"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	// Reply in Slack
	api.PostMessage(event.Channel, slack.MsgOptionText(res.FullResponse, false))

	// Trace the response
	span.SetAttributes(attribute.String("backend.response", res.FullResponse))
}
