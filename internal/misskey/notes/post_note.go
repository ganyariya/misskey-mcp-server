package notes

import (
	"fmt" // Added for error wrapping
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/sirupsen/logrus" // Added for logger
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

const (
	toolName    = "post_misskey_note"
	description = "Post a note to Misskey"
)

// NewPostNoteTool creates a new instance of the postNoteTool.
// It now accepts a logger instance.
func NewPostNoteTool(logger *logrus.Logger) *postNoteTool {
	return &postNoteTool{
		Name:        toolName,
		Description: description,
		logger:      logger,
	}
}

type postNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

type postNoteTool struct {
	Name        string
	Description string
	logger      *logrus.Logger // Logger instance for the tool
}

func (p *postNoteTool) Register(server *mcp_golang.Server, misskeyClient *misskey.Client) error {
	err := server.RegisterTool(
		p.GetName(),
		p.GetDescription(),
		func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error) {
			p.logger.Infof("Attempting to post note with text: %s", arguments.Text)
			text := arguments.Text
			response, err := misskeyClient.Notes().Create(notes.CreateRequest{
				Text:       core.NewString(text),
				Visibility: models.VisibilityPublic,
			})
			if err != nil {
				p.logger.Errorf("Error creating note via Misskey API: %v", err)
				return nil, fmt.Errorf("failed to create note via misskey api: %w", err)
			}

			p.logger.Infof("Note posted successfully: %s (ID: %s)", text, response.CreatedNote.ID)
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Note posted successfully: " + text + " ID: " + response.CreatedNote.ID)), nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to register tool %s: %w", p.GetName(), err)
	}
	return nil
}

func (p *postNoteTool) GetName() string {
	return p.Name
}
func (p *postNoteTool) GetDescription() string {
	return p.Description
}
