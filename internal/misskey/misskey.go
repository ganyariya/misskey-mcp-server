package misskey

import (
	"github.com/ganyariya/misskey-mcp-server/internal/misskey/notes" // New import path
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
)

// Tool defines the interface for a Misskey tool.
// Renamed from MisskeyTool for brevity.
type Tool interface {
	GetName() string
	GetDescription() string
	// Register registers the tool with the MCP server.
	// The misskey.Client is provided for interacting with Misskey.
	Register(*mcp_golang.Server, *misskey.Client) error
}

// RegisterTools initializes and registers all available Misskey tools.
// Renamed from RegisterMisskeyTools.
// It now accepts a logger instance.
func RegisterTools(server *mcp_golang.Server, client *misskey.Client, logger *logrus.Logger) error {
	// TODO: Consider whether the logger should be passed to each tool's Register method
	// or if tools should have an Init(logger) method, or access a global logger.
	// For now, the logger is here if RegisterTools needs to log something directly,
	// or if a tool's New...() function needs it.

	logger.Debugf("Registering Misskey tools...")

	var tools []Tool

	// Explicitly register tools.
	// TODO: Explore dynamic tool registration later if needed.
	postNoteTool := notes.NewPostNoteTool(logger) // Pass logger to the tool constructor
	tools = append(tools, postNoteTool)

	for _, tool := range tools {
		logger.Infof("Registering tool: %s", tool.GetName())
		if err := tool.Register(server, client); err != nil {
			logger.Errorf("Failed to register tool %s: %v", tool.GetName(), err)
			return err // Stop registration if one tool fails
		}
	}

	logger.Info("All Misskey tools registered successfully.")
	return nil
}
