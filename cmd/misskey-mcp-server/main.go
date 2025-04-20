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

func run(
	transport string,
	logLevel string,
	misskeyClient *misskey.Client,
) error {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	// TODO: Add logger
	err := tools.RegisterMisskeyTools(server, misskeyClient)
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
		fmt.Printf("failed to create misskey client: %v\n", err)
		return
	}

	if err := run(transport, logLevel, misskeyClient); err != nil {
		panic(err)
	}
}
