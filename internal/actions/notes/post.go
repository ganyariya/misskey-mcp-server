package notes

import (
	"fmt"

	"github.com/ganyariya/misskey-mcp-server/internal/actions"
	mcp_golang "github.com/metoro-io/mcp-golang" // Standardizing to metoro-io
	"github.com/misskey-dev/misskey-go/misskey"
	"github.com/misskey-dev/misskey-go/misskey/core"
	"github.com/misskey-dev/misskey-go/misskey/models"
	"github.com/misskey-dev/misskey-go/misskey/services/notes"
	// "github.com/sirupsen/logrus" // Not needed if using ctx.Logger() which is already a logrus.FieldLogger
)

// postNoteArguments struct remains the same
type postNoteArguments struct {
	Text string `json:"text" jsonschema:"required,description=The text of the note to post"`
}

type PostNoteAction struct{}

func NewPostNoteAction() *PostNoteAction {
	return &PostNoteAction{}
}

func (a *PostNoteAction) Name() string {
	return "post_misskey_note"
}

func (a *PostNoteAction) Description() string {
	return "Post a note to Misskey"
}

func (a *PostNoteAction) Params() interface{} {
	// Return a new instance, not the zero value of the type itself,
	// as mcp-golang will use reflection on this instance.
	return &postNoteArguments{}
}

func (a *PostNoteAction) Execute(ctx actions.ActionContext, params interface{}) (*mcp_golang.ToolResponse, error) {
	logger := ctx.Logger() // Get the logger from context

	args, ok := params.(*postNoteArguments)
	if !ok {
		// It's important to log this error before returning, as the caller (registry) might only log the error it receives.
		logger.Errorf("Invalid arguments type for %s, expected *postNoteArguments, got %T", a.Name(), params)
		return nil, fmt.Errorf("invalid arguments type for %s, expected *postNoteArguments, got %T", a.Name(), params)
	}

	if ctx.MisskeyClient == nil {
		logger.Error("Misskey client is not initialized in ActionContext")
		return nil, fmt.Errorf("misskey client is not initialized in ActionContext for action %s", a.Name())
	}

	logger.WithField("text_length", len(args.Text)).Info("Attempting to post note")

	response, err := ctx.MisskeyClient.Notes().Create(notes.CreateRequest{
		Text:       core.NewString(args.Text),
		Visibility: models.VisibilityPublic, // Or make this configurable via params
	})
	if err != nil {
		logger.WithError(err).Error("Failed to post note to Misskey")
		return nil, fmt.Errorf("failed to post note for action %s: %w", a.Name(), err)
	}

	logger.Infof("Note posted successfully. ID: %s, Text: %s", response.CreatedNote.ID, args.Text)
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Note posted successfully. ID: %s, Text: %s", response.CreatedNote.ID, args.Text))), nil
}
