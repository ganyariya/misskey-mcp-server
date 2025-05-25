package notes_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ganyariya/misskey-mcp-server/internal/actions"
	"github.com/ganyariya/misskey-mcp-server/internal/actions/notes"
	mcp_golang_metoro "github.com/metoro-io/mcp-golang" // As used in post.go
	"github.com/misskey-dev/misskey-go/misskey"
	misskey_models "github.com/misskey-dev/misskey-go/misskey/models"
	misskey_notes "github.com/misskey-dev/misskey-go/misskey/services/notes"
	"github.com/sirupsen/logrus"
)

// mockMisskeyNotesService implements misskey.NotesServiceIface for testing.
type mockMisskeyNotesService struct {
	CreateFunc func(req misskey_notes.CreateRequest) (*misskey_notes.CreateResponse, error)
	// Add other methods if needed by other actions or tests
}

func (m *mockMisskeyNotesService) Create(req misskey_notes.CreateRequest) (*misskey_notes.CreateResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return nil, errors.New("CreateFunc not implemented in mock")
}

// Unused methods to satisfy the interface
func (m *mockMisskeyNotesService) Children(req misskey_notes.ChildrenRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Conversation(req misskey_notes.ConversationRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Delete(req misskey_notes.DeleteRequest) error { return nil }
func (m *mockMisskeyNotesService) Favorites(req misskey_notes.FavoritesRequest) ([]*misskey_models.NoteFavorite, error) { return nil, nil }
func (m *mockMisskeyNotesService) Featured(req misskey_notes.FeaturedRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) GlobalTimeline(req misskey_notes.GlobalTimelineRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) HybridTimeline(req misskey_notes.HybridTimelineRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) LocalTimeline(req misskey_notes.LocalTimelineRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Mentions(req misskey_notes.MentionsRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) PollVote(req misskey_notes.PollVoteRequest) error { return nil }
func (m *mockMisskeyNotesService) Reactions(req misskey_notes.ReactionsRequest) ([]*misskey_models.NoteReaction, error) { return nil, nil }
func (m *mockMisskeyNotesService) ReactionsCreate(req misskey_notes.ReactionsCreateRequest) error { return nil, nil }
func (m *mockMisskeyNotesService) ReactionsDelete(req misskey_notes.ReactionsDeleteRequest) error { return nil, nil }
func (m *mockMisskeyNotesService) Renote(req misskey_notes.RenoteRequest) (*misskey_notes.CreateResponse, error) { return nil, nil }
func (m *mockMisskeyNotesService) Replies(req misskey_notes.RepliesRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) SearchByTag(req misskey_notes.SearchByTagRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Search(req misskey_notes.SearchRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Show(req misskey_notes.ShowRequest) (*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) State(req misskey_notes.StateRequest) (*misskey_models.NoteState, error) { return nil, nil }
func (m *mockMisskeyNotesService) Timeline(req misskey_notes.TimelineRequest) ([]*misskey_models.Note, error) { return nil, nil }
func (m *mockMisskeyNotesService) Translate(req misskey_notes.TranslateRequest) (*misskey_notes.TranslateResponse, error) { return nil, nil }
func (m *mockMisskeyNotesService) Unrenote(req misskey_notes.UnrenoteRequest) error { return nil }
func (m *mockMisskeyNotesService) UserListTimeline(req misskey_notes.UserListTimelineRequest) ([]*misskey_models.Note, error) { return nil, nil }


// mockMisskeyClient implements parts of Misskey client for testing.
type mockMisskeyClient struct {
	notesService misskey.NotesServiceIface
}

func (m *mockMisskeyClient) Notes() misskey.NotesServiceIface {
	return m.notesService
}
// Add other service getters if needed by other actions

func TestPostNoteAction_Name(t *testing.T) {
	action := notes.NewPostNoteAction()
	if action.Name() != "post_misskey_note" {
		t.Errorf("Expected Name 'post_misskey_note', got '%s'", action.Name())
	}
}

func TestPostNoteAction_Description(t *testing.T) {
	action := notes.NewPostNoteAction()
	if action.Description() != "Post a note to Misskey" {
		t.Errorf("Expected Description 'Post a note to Misskey', got '%s'", action.Description())
	}
}

func TestPostNoteAction_Params(t *testing.T) {
	action := notes.NewPostNoteAction()
	params := action.Params()
	expectedType := &notes.PostNoteArguments{}
	if reflect.TypeOf(params) != reflect.TypeOf(expectedType) {
		t.Errorf("Expected Params type %T, got %T", expectedType, params)
	}
}

func TestPostNoteAction_Execute(t *testing.T) {
	logger := logrus.NewTestLogger()
	action := notes.NewPostNoteAction()

	t.Run("Success Case", func(t *testing.T) {
		mockNotesSvc := &mockMisskeyNotesService{
			CreateFunc: func(req misskey_notes.CreateRequest) (*misskey_notes.CreateResponse, error) {
				if req.Text == nil || *req.Text != "Hello Test" {
					return nil, fmt.Errorf("unexpected text in CreateRequest: %v", req.Text)
				}
				return &misskey_notes.CreateResponse{
					CreatedNote: misskey_models.Note{ID: "testnote123", CreatedAt: time.Now(), Text: req.Text},
				}, nil
			},
		}
		mockClient := &mockMisskeyClient{notesService: mockNotesSvc}
		ctx := actions.NewActionContext(mockClient, logger)
		params := &notes.PostNoteArguments{Text: "Hello Test"}

		response, err := action.Execute(ctx, params)

		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if response == nil {
			t.Fatal("Execute response is nil")
		}
		expectedContent := mcp_golang_metoro.NewTextContent("Note posted successfully. ID: testnote123, Text: Hello Test")
		if !reflect.DeepEqual(response.Content, expectedContent) {
			t.Errorf("Expected response content %+v, got %+v", expectedContent, response.Content)
		}
	})

	t.Run("Misskey API Error Case", func(t *testing.T) {
		apiError := errors.New("misskey API error")
		mockNotesSvc := &mockMisskeyNotesService{
			CreateFunc: func(req misskey_notes.CreateRequest) (*misskey_notes.CreateResponse, error) {
				return nil, apiError
			},
		}
		mockClient := &mockMisskeyClient{notesService: mockNotesSvc}
		ctx := actions.NewActionContext(mockClient, logger)
		params := &notes.PostNoteArguments{Text: "Error Test"}

		response, err := action.Execute(ctx, params)

		if err == nil {
			t.Fatal("Execute should have returned an error, but it was nil")
		}
		if !errors.Is(err, apiError) { // Check if the original error is wrapped
		    // Check if the error message contains the apiError message
		    if ! (err.Error() == fmt.Sprintf("failed to post note for action %s: %s", action.Name(), apiError.Error())) {
			    t.Errorf("Expected error to wrap '%v', got '%v'", apiError, err)
            }
		}
		if response != nil {
			t.Errorf("Execute response should be nil on error, got %+v", response)
		}
	})

	t.Run("Invalid Argument Type Case", func(t *testing.T) {
		mockClient := &mockMisskeyClient{notesService: &mockMisskeyNotesService{}} // Service won't be called
		ctx := actions.NewActionContext(mockClient, logger)
		params := "not a postNoteArguments struct" // Invalid type

		response, err := action.Execute(ctx, params)

		if err == nil {
			t.Fatal("Execute should have returned an error for invalid argument type, but it was nil")
		}
		expectedErrorMsg := fmt.Sprintf("invalid arguments type for %s, expected *notes.PostNoteArguments, got string", action.Name())
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
		if response != nil {
			t.Errorf("Execute response should be nil on invalid argument error, got %+v", response)
		}
	})

	t.Run("Nil MisskeyClient Case", func(t *testing.T) {
		ctx := actions.NewActionContext(nil, logger) // Nil MisskeyClient
		params := &notes.PostNoteArguments{Text: "Nil Client Test"}

		response, err := action.Execute(ctx, params)

		if err == nil {
			t.Fatal("Execute should have returned an error for nil MisskeyClient, but it was nil")
		}
		expectedErrorMsg := fmt.Sprintf("misskey client is not initialized in ActionContext for action %s", action.Name())
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
		if response != nil {
			t.Errorf("Execute response should be nil on nil client error, got %+v", response)
		}
	})
}
