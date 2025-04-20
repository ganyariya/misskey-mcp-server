package main

import (
	"flag"
	"fmt"
	"os"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

const SERVER_NAME = "misskey-mcp-server"

type PostNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

func run(
	transport string,
	logLevel string,
	misskeyClient *misskey.Client,
) error {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName(SERVER_NAME))

	err := server.RegisterTool(
		"misskey-note-post",
		"Post a note to Misskey",
		func(arguments PostNoteArguments) (*mcp_golang.ToolResponse, error) {
			text := arguments.Text
			response, err := misskeyClient.Notes().Create(notes.CreateRequest{
				Text:       core.NewString(text),
				Visibility: models.VisibilityPublic,
			})
			if err != nil {
				return nil, err
			}

			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Note posted successfully: " + text + response.CreatedNote.ID)), nil
		},
	)
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

	var logLevel string
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

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
