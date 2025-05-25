package notes_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ganyariya/misskey-mcp-server/internal/actions"
	"github.com/ganyariya/misskey-mcp-server/internal/actions/notes"
	mcp_golang_metoro "github.com/metoro-io/mcp-golang" // As used in search.go
	"github.com/sirupsen/logrus"
)

func TestSearchNotesAction_Name(t *testing.T) {
	action := notes.NewSearchNotesAction()
	if action.Name() != "search_misskey_notes" {
		t.Errorf("Expected Name 'search_misskey_notes', got '%s'", action.Name())
	}
}

func TestSearchNotesAction_Description(t *testing.T) {
	action := notes.NewSearchNotesAction()
	if action.Description() != "Search for notes on Misskey (stub)" {
		t.Errorf("Expected Description 'Search for notes on Misskey (stub)', got '%s'", action.Description())
	}
}

func TestSearchNotesAction_Params(t *testing.T) {
	action := notes.NewSearchNotesAction()
	params := action.Params()
	expectedType := &notes.SearchNotesArguments{}
	if reflect.TypeOf(params) != reflect.TypeOf(expectedType) {
		t.Errorf("Expected Params type %T, got %T", expectedType, params)
	}
}

func TestSearchNotesAction_Execute(t *testing.T) {
	logger := logrus.NewTestLogger()
	action := notes.NewSearchNotesAction()

	t.Run("Success Case (Stub)", func(t *testing.T) {
		// MisskeyClient can be nil as it's not used by the stub Execute method
		ctx := actions.NewActionContext(nil, logger)
		params := &notes.SearchNotesArguments{Query: "test query"}

		response, err := action.Execute(ctx, params)

		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if response == nil {
			t.Fatal("Execute response is nil")
		}
		expectedMsg := fmt.Sprintf("Search functionality for notes with query '%s' is not yet implemented.", params.Query)
		expectedContent := mcp_golang_metoro.NewTextContent(expectedMsg)
		if !reflect.DeepEqual(response.Content, expectedContent) {
			t.Errorf("Expected response content %+v, got %+v", expectedContent, response.Content)
		}
	})

	t.Run("Invalid Argument Type Case", func(t *testing.T) {
		ctx := actions.NewActionContext(nil, logger)
		invalidParams := "not a searchNotesArguments struct" // Invalid type

		response, err := action.Execute(ctx, invalidParams)

		if err == nil {
			t.Fatal("Execute should have returned an error for invalid argument type, but it was nil")
		}
		// The error message comes from fmt.Errorf in the action's Execute method.
		// Example: "invalid arguments type for search_misskey_notes, expected *notes.SearchNotesArguments, got string"
		expectedErrorMsg := fmt.Sprintf("invalid arguments type for %s, expected *notes.SearchNotesArguments, got string", action.Name())
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
		if response != nil {
			t.Errorf("Execute response should be nil on invalid argument error, got %+v", response)
		}
	})
}
