package main

import (
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

type PostNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

func main() {
	misskeyClient, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(os.Getenv("MISSKEY_API_TOKEN")),
		misskey.WithBaseURL(
			os.Getenv("MISSKEY_PROTOCOL"),
			os.Getenv("MISSKEY_DOMAIN"),
			os.Getenv("MISSKEY_PATH"),
		),
		misskey.WithLogLevel(logrus.DebugLevel),
	)
	if err != nil {
		fmt.Printf("failed to create misskey client: %v\n", err)
		return
	}

	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName("misskey-mcp-server"))

	err = server.RegisterTool(
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
		fmt.Printf("failed to register tool: %v\n", err)
		return
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
