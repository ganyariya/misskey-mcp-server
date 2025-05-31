package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ganyariya/misskey-mcp-server/internal/tools"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
)

const SERVER_NAME = "misskey-mcp-server"

var log = logrus.New()

func run(
	transport string,
	misskeyClient *misskey.Client,
) error {
	log.Info("Starting MCP server...")
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	err := tools.RegisterMisskeyTools(server, misskeyClient)
	if err != nil {
		return fmt.Errorf("failed to register misskey tools: %w", err)
	}

	err = server.Serve()
	if err != nil {
		log.Errorf("Server serve error: %v", err)
		return err
	}

	<-done

	return nil
}

func main() {
	log.Out = os.Stderr
	log.Formatter = &logrus.TextFormatter{}

	// TODO: Add support http transport
	var transport string
	flag.StringVar(&transport, "t", "stdio", "transport type (stdio only for now)")
	flag.StringVar(&transport, "transport", "stdio", "transport type (stdio only for now)")

	var logLevel string
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&logLevel, "l", "info", "Log level (debug, info, warn, error)")

	flag.Parse()

	parsedLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info': %v", logLevel, err)
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(parsedLevel)
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
		misskey.WithLogLevel(logrus.DebugLevel), // This is for go-misskey's internal logger
	)
	if err != nil {
		log.Fatalf("failed to create misskey client: %v", err)
	}

	if err := run(transport, misskeyClient); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
