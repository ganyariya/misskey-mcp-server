package note

import (
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

const (
	toolName    = "post_misskey_note"
	description = "Post a note to Misskey"
)

func NewPostNoteTool() *postNoteTool {
	return &postNoteTool{
		Name:        toolName,
		Description: description,
	}
}

type postNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

type postNoteTool struct {
	Name        string
	Description string
}

func (p *postNoteTool) Register(server *mcp_golang.Server, misskeyClient *misskey.Client) error {
	err := server.RegisterTool(
		p.GetName(),
		p.GetDescription(),
		func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error) {
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
	return err
}

func (p *postNoteTool) GetName() string {
	return p.Name
}
func (p *postNoteTool) GetDescription() string {
	return p.Description
}
