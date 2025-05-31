package tools

import (
	"github.com/ganyariya/misskey-mcp-server/internal/misskey_tools/note"
	"github.com/ganyariya/misskey-mcp-server/internal/misskey_tools/user"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/yitsushi/go-misskey"
)

type MisskeyTool interface {
	GetName() string
	GetDescription() string
	Register(*mcp_golang.Server, *misskey.Client) error
}

func RegisterMisskeyTools(server *mcp_golang.Server, misskeyClient *misskey.Client) error {
	var tools []MisskeyTool

	// Register Misskey tools here
	// TODO: dynamically load tools from a directory by reflection or plugin
	tools = append(tools, note.NewPostNoteTool())
	tools = append(tools, user.NewGetUserNotesTool())

	for _, tool := range tools {
		err := tool.Register(server, misskeyClient)
		if err != nil {
			return err
		}
	}

	return nil
}
