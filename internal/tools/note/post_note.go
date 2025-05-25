package note

import (
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/ganyariya/misskey-mcp-server/internal/misskey" // New import
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

const (
	toolName    = "post_misskey_note"
	description = "Post a note to Misskey"
)

// NewPostNoteTool creates a new postNoteTool.
// It now requires a misskey.Client.
func NewPostNoteTool(client misskey.Client) *postNoteTool {
	return &postNoteTool{
		Name:        toolName,
		Description: description,
		client:      client, // Store the client
	}
}

type postNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

type postNoteTool struct {
	Name        string
	Description string
	client      misskey.Client // Added misskey client field
}

// handleRequest is the actual logic for handling the tool's execution.
func (p *postNoteTool) handleRequest(arguments postNoteArguments) (*mcp_golang.ToolResponse, error) {
	text := arguments.Text
	response, err := p.client.Notes().Create(notes.CreateRequest{
		Text:       core.NewString(text),
		Visibility: models.VisibilityPublic,
	})
	if err != nil {
		return nil, err
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Note posted successfully: " + text + response.CreatedNote.ID)), nil
}

func (p *postNoteTool) Register(server *mcp_golang.Server, misskeyClient misskey.Client) error { // misskeyClient param is still here due to interface MisskeyTool
	// The misskeyClient passed to Register might be different from p.client if Register is called directly
	// For consistency, ensure p.client is set, or this method should exclusively use the client it was constructed with.
	// Current design: NewPostNoteTool sets the client. Register will use the client from the struct.
	// The misskeyClient parameter for Register might become redundant if all tools follow this pattern.
	if p.client == nil {
		// This case should ideally not happen if constructed with NewPostNoteTool
		// Or, if the tool can be registered with a *different* client later, this logic needs thought.
		// For now, assume p.client is the one to use.
		// If misskeyClient param in Register is to override, then use it.
		// Let's stick to using the client set during construction (p.client).
	}

	err := server.RegisterTool(
		p.GetName(),
		p.GetDescription(),
		p.handleRequest, // Pass the method reference
	)
	return err
}

func (p *postNoteTool) GetName() string {
	return p.Name
}
func (p *postNoteTool) GetDescription() string {
	return p.Description
}
