package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

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

	s := server.NewMCPServer(
		"misskey-mcp-server",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	notePostTool := mcp.NewTool(
		"misskey-note-post",
		mcp.WithDescription("Post a note to Misskey"),
		mcp.WithString("text", mcp.Required(), mcp.Description("Note text")),
	)

	s.AddTool(notePostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		text := request.Params.Arguments["text"].(string)
		response, err := misskeyClient.Notes().Create(notes.CreateRequest{
			Text:       core.NewString(text),
			Visibility: models.VisibilityPublic,
		})
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText("Note posted successfully: " + text + response.CreatedNote.ID), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}
}
