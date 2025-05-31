package user

import (
	"encoding/json"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/services/notes"
)

const (
	toolName    = "get_misskey_user_notes"
	description = "Get user notes from Misskey"
)

func NewGetUserNotesTool() *getUserNotesTool {
	return &getUserNotesTool{
		Name:        toolName,
		Description: description,
	}
}

type getUserNotesArguments struct {
	UserID string `json:"userId" jsonschema:"required,description=The ID of the user to get notes for"`
}

type getUserNotesTool struct {
	Name        string
	Description string
}

func (t *getUserNotesTool) Register(server *mcp_golang.Server, misskeyClient *misskey.Client) error {
	err := server.RegisterTool(
		t.GetName(),
		t.GetDescription(),
		func(arguments getUserNotesArguments) (*mcp_golang.ToolResponse, error) {
			userID := arguments.UserID
			response, err := misskeyClient.Notes().Search(notes.SearchRequest{
				Query:  "*", // Changed Query to "*" to pass validation
				UserID: core.NewString(userID),
				Limit:  100,
			})
			if err != nil {
				return nil, err
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				return nil, err
			}

			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(responseJSON))), nil
		},
	)
	return err
}

func (t *getUserNotesTool) GetName() string {
	return t.Name
}

func (t *getUserNotesTool) GetDescription() string {
	return t.Description
}
