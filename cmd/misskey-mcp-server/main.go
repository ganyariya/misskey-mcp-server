package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ganyariya/misskey-mcp-server/internal/tools"
	misskeyclient "github.com/ganyariya/misskey-mcp-server/internal/misskey" // Added import
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/sirupsen/logrus"
	misskeysdk "github.com/yitsushi/go-misskey" // Aliased import
)

const SERVER_NAME = "misskey-mcp-server"

func run(
	transport string,
	logLevel string,
	client misskeyclient.Client, // Changed parameter type
) error {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	// TODO: Add logger
	err := tools.RegisterMisskeyTools(server, client) // Use the new parameter name
	if err != nil {
		return nil
	}

	err = server.Serve()
	if err != nil {
		return err
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

	misskeyApiToken := os.Getenv("MISSKEY_API_TOKEN")
	misskeyProtocol := os.Getenv("MISSKEY_PROTOCOL")
	misskeyDomain := os.Getenv("MISSKEY_DOMAIN")
	misskeyPath := os.Getenv("MISSKEY_PATH")

	sdkClient, err := misskeysdk.NewClientWithOptions( // Renamed original client
		misskeysdk.WithAPIToken(misskeyApiToken),
		misskeysdk.WithBaseURL(
			misskeyProtocol,
			misskeyDomain,
			misskeyPath,
		),
		misskeysdk.WithLogLevel(logrus.DebugLevel),
	)
	if err != nil {
		fmt.Printf("failed to create misskey sdk client: %v\n", err) // Updated log message
		return
	}

	wrappedClient := misskeyclient.NewMisskeyGoClientWrapper(sdkClient) // Create wrapped client

	if err := run(transport, logLevel, wrappedClient); err != nil { // Pass wrapped client
		panic(err)
	}
}
