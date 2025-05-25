package main

import (
	"flag"
	"os" // Keep os for os.Getenv

	"github.com/ganyariya/misskey-mcp-server/internal/actions"
	"github.com/ganyariya/misskey-mcp-server/internal/actions/notes"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
)

const SERVER_NAME = "misskey-mcp-server"

func run(
	// transport string, // This parameter is not used
	// logLevel string,  // Log level is configured in main and available globally via logrus
	misskeyClient *misskey.Client,
) error {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	// Initialize and use the new action registry
	actionRegistry := actions.NewRegistry()

	// Create and register actions
	postNoteAction := notes.NewPostNoteAction()
	if err := actionRegistry.Register(postNoteAction); err != nil {
	postNoteAction := notes.NewPostNoteAction()
	if err := actionRegistry.Register(postNoteAction); err != nil {
		logrus.WithError(err).Errorf("Error registering action %s", postNoteAction.Name())
		return err
	}

	// Create and register the search action
	searchNotesAction := notes.NewSearchNotesAction()
	if err := actionRegistry.Register(searchNotesAction); err != nil {
		logrus.WithError(err).Errorf("Error registering action %s", searchNotesAction.Name())
		return err
	}

	// Register all actions with the MCP server
	// Pass the standard logrus logger to the registry, it can create child loggers from it.
	if err := actionRegistry.RegisterActionsWithServer(server, misskeyClient, logrus.StandardLogger()); err != nil {
		logrus.WithError(err).Error("Error registering actions with MCP server")
		return err
	}

	logrus.Infof("MCP Server '%s' starting with registered actions...", SERVER_NAME)

	if err := server.Serve(); err != nil {
		logrus.WithError(err).Error("Server failed to start or exited unexpectedly")
		return err
	}

	<-done // This will keep the server running until 'done' is closed.

	return nil
}

func main() {
	// TODO: Add support http transport
	var transport string
	// var transport string // transport flag is not used
	// flag.StringVar(&transport, "t", "stdio", "transport type (stdio only for now)")
	// flag.StringVar(&transport, "transport", "stdio", "transport type (stdio only for now)")

	var logLevelFlag string
	flag.StringVar(&logLevelFlag, "log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	flag.StringVar(&logLevelFlag, "l", logLevelFlag, "Log level (trace, debug, info, warn, error, fatal, panic) (shorthand)")

	flag.Parse()

	// Configure logrus
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	level, err := logrus.ParseLevel(logLevelFlag)
	if err != nil {
		logrus.WithError(err).Warnf("Invalid log level '%s', defaulting to 'info'", logLevelFlag)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.Infof("Log level set to: %s", level.String())

	misskeyApiToken := os.Getenv("MISSKEY_API_TOKEN")
	misskeyProtocol := os.Getenv("MISSKEY_PROTOCOL")
	misskeyDomain := os.Getenv("MISSKEY_DOMAIN")
	misskeyPath := os.Getenv("MISSKEY_PATH")

	if misskeyApiToken == "" || misskeyProtocol == "" || misskeyDomain == "" {
		logrus.Fatal("Missing one or more Misskey environment variables: MISSKEY_API_TOKEN, MISSKEY_PROTOCOL, MISSKEY_DOMAIN")
	}

	misskeyClient, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyApiToken),
		misskey.WithBaseURL(misskeyProtocol, misskeyDomain, misskeyPath),
		// go-misskey uses its own logrus instance. We can set its level.
		// Let's make it consistent with the application's log level.
		misskey.WithLogLevel(level),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create misskey client")
	}
	// The misskey client's SetLogLevel method was an assumption, WithLogLevel during creation is the correct way.
	// The previous logic for misskeyClient.SetLogLevel(parsedLogLevel) has been removed.

	// The run function no longer needs logLevel string, as logrus is globally configured.
	// The transport string was also removed as it's not used.
	if err := run(misskeyClient); err != nil {
		logrus.WithError(err).Fatalf("Server run failed")
	}
}
