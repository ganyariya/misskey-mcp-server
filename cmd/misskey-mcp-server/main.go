package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ganyariya/misskey-mcp-server/internal/misskey" // Updated import path
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
)

var log = logrus.New()

const SERVER_NAME = "misskey-mcp-server"

func run(
	transport string,
	logLevel string,
	misskeyClient *misskey.Client,
	log *logrus.Logger,
) error {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	// TODO: Add logger
	// Updated to call misskey.RegisterTools and pass the logger
	err := misskey.RegisterTools(server, misskeyClient, log)
	if err != nil {
		// Error handling for tool registration
		// The logger is already passed to RegisterTools, so it can log specifics.
		return fmt.Errorf("failed to register misskey tools: %w", err)
	}

	if err := server.Serve(); err != nil {
		return fmt.Errorf("mcp server error: %w", err)
	}

	<-done

	return nil
}

func main() {
	// TODO: Add support http transport
	var transport string
	flag.StringVar(&transport, "t", "stdio", "transport type (stdio only for now)")
	flag.StringVar(&transport, "transport", "stdio", "transport type (stdio only for now)")

	// TODO: Add LogLevel
	var logLevel string
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&logLevel, "l", "info", "Log level (debug, info, warn, error)")

	flag.Parse()

	// Initialize logger
	log.SetOutput(os.Stderr)
	log.SetFormatter(&logrus.TextFormatter{})
	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel) // Default to info
	}

	misskeyApiToken := os.Getenv("MISSKEY_API_TOKEN")
	misskeyProtocol := os.Getenv("MISSKEY_PROTOCOL")
	misskeyDomain := os.Getenv("MISSKEY_DOMAIN")
	misskeyPath := os.Getenv("MISSKEY_PATH")

	misskeyClient, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyApiToken),
		misskey.WithBaseURL(
			misskeyProtocol,
			misskeyDomain,
			misskeyPath,
		),
		misskey.WithLogLevel(logrus.DebugLevel),
	)
	if err != nil {
		// Use %w to wrap the error, allowing further unwrapping if needed.
		// However, since we log it directly and return, %v is also acceptable here
		// if we consider this the final point for this specific error.
		// For consistency in wrapping, %w is good practice if this were a library.
		// Given it's main, and we log.Errorf then return, the immediate next step is exit.
		// Let's use a more descriptive message and keep %v as log.Errorf doesn't wrap.
		log.Errorf("Initialization failed: Unable to create misskey client: %v", err)
		// If we wanted to return it for further wrapping before a log.Fatalf:
		// return fmt.Errorf("failed to create misskey client: %w", err)
		return // Exiting because client creation failed.
	}

	if err := run(transport, logLevel, misskeyClient, log); err != nil {
		log.Fatalf("Startup error: %v", err) // Use Fatalf to log and exit
	}
}
