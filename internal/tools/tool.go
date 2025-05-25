package tools

import (
	"github.com/ganyariya/misskey-mcp-server/internal/misskey" // Changed import
	"github.com/ganyariya/misskey-mcp-server/internal/tools/note"
	mcp_golang "github.com/metoro-io/mcp-golang"
)

type MisskeyTool interface {
	GetName() string
	GetDescription() string
	Register(*mcp_golang.Server, misskey.Client) error // Changed parameter type
}

func RegisterMisskeyTools(server *mcp_golang.Server, misskeyClient misskey.Client) error { // Changed parameter type
	var tools []MisskeyTool

	// Register Misskey tools here
	// TODO: dynamically load tools from a directory by reflection or plugin
	// Pass misskeyClient to NewPostNoteTool constructor
	tools = append(tools, note.NewPostNoteTool(misskeyClient))

	for _, tool := range tools {
		err := tool.Register(server, misskeyClient)
		if err != nil {
			return err
		}
	}

	return nil
}
