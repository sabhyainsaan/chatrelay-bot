ChatRelay Bot
ChatRelay is a basic Slack bot built using Go. It listens for messages where it is mentioned, sends the message to a mock backend server, and then responds in the Slack channel. The project also includes OpenTelemetry tracing so that the entire flow is visible in the terminal.

What the Bot Does
Connects to Slack using Socket Mode

Listens for messages like @ChatRelay what is Go?

Sends that message to a mock backend server

Waits for a response

Posts the response back to the Slack channel

Logs every step using OpenTelemetry

Folder Structure
go
Copy
Edit
chatrelaybot/
├── cmd/                // Application entry point
│   └── main.go
├── internal/
│   ├── slackbot/       // Slack bot logic
│   ├── backend/        // Mock backend server
│   └── utils/          // Tracing setup
├── .env                // Environment variables
├── go.mod
└── README.md
Technologies Used
Go (Golang)

Slack API (via slack-go)

OpenTelemetry (tracing)

Standard Go HTTP server for the backend

Environment Setup
Create a .env file in the root folder with the following keys:

ini
Copy
Edit
SLACK_BOT_TOKEN=your-slack-bot-token
SLACK_APP_TOKEN=your-app-token
BACKEND_URL=http://localhost:8080/v1/chat/stream
You can get these tokens from your Slack App dashboard.

How to Run
Make sure Go is installed

Run the following commands:

bash
Copy
Edit
go mod tidy
go run cmd/main.go
In Slack, go to a channel and mention the bot like this:

sql
Copy
Edit
@ChatRelay what are goroutines?
The bot will send the message to the backend and reply with a hardcoded response.

Tracing Output
OpenTelemetry is used to track what’s happening inside the app. After running the bot and sending a message, your terminal will show output like this:

makefile
Copy
Edit
Name: HandleSlackMention
Attributes:
  - slack.query: "what are goroutines?"
  - backend.response: "Goroutines are lightweight threads managed by the Go runtime."
makefile
Copy
Edit
Name: MockBackendResponse
Attributes:
  - user_id: "U123"
  - query: "what are goroutines?"
  - response: "Goroutines are lightweight threads managed by the Go runtime."
What Can Be Improved
Add real-time streaming instead of fixed backend responses

Add tests

Use Jaeger or Zipkin for visualizing traces

Dockerize the app for easier deployment

About Me
Ananya Chaturvedi
Backend Developer (Golang)
Email: ananya516chaturvedi@gmail.com