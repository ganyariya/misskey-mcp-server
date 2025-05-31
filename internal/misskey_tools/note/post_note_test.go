package note

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ganyariya/misskey-mcp-server/internal/mocks"
	mcp_golang "github.com/metoro-io/mcp-golang"
	// "github.com/yitsushi/go-misskey/core" // Keep for core.NewString if used, but not directly in this corrected test logic for Create call args
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
)

func TestNewPostNoteTool(t *testing.T) {
	tool := NewPostNoteTool()
	if tool.Name != toolName {
		t.Errorf("expected tool name %s, got %s", toolName, tool.Name)
	}
	if tool.Description != description {
		t.Errorf("expected tool description %s, got %s", description, tool.Description)
	}
}

func TestPostNoteTool_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMCPServer := mocks.NewMockServer(ctrl)
	mockMisskeyClient := mocks.NewMockClient(ctrl)
	mockNotesService := mocks.NewMockClientInterface(ctrl)

	tool := NewPostNoteTool()

	expectedNoteText := "Test note from MCP"
	expectedAPIResponseNoteID := "testNoteId123"

	var registeredFunc func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error)

	// Expect RegisterTool to be called on the MCP Server and capture the function
	mockMCPServer.EXPECT().
		RegisterTool(
			tool.GetName(),
			tool.GetDescription(),
			gomock.Any(),
		).
		DoAndReturn(func(name, desc string, f interface{}) error {
			// Type assertion for the function
			cb, ok := f.(func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error))
			if !ok {
				t.Fatalf("RegisterTool was called with an unexpected function signature.")
			}
			registeredFunc = cb
			return nil
		}).
		Return(nil). // Ensure this is present if the mocked function is expected to return a value.
		Times(1)

	// Setup expectations for the Misskey client for when the registeredFunc is called
	// This setup needs to be here because registeredFunc will call these.
	mockMisskeyClient.EXPECT().Notes().Return(mockNotesService).Times(1)
	mockNotesService.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(req notes.CreateRequest) (*notes.CreateResponse, error) {
			if req.Text == nil || *req.Text != expectedNoteText {
				t.Errorf("expected text '%s', got '%s'", expectedNoteText, *req.Text)
			}
			if req.Visibility != models.VisibilityPublic {
				t.Errorf("expected visibility '%s', got '%s'", models.VisibilityPublic, req.Visibility)
			}
			return &notes.CreateResponse{
				CreatedNote: models.Note{ID: models.ID(expectedAPIResponseNoteID)},
			}, nil
		}).Times(1)

	// Call Register to capture the function
	err := tool.Register(mockMCPServer, mockMisskeyClient)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Ensure the function was captured
	if registeredFunc == nil {
		t.Fatal("Registered function was not captured by the mock")
	}

	// Now, call the captured function
	args := postNoteArguments{Text: expectedNoteText}
	response, err := registeredFunc(args)

	if err != nil {
		t.Fatalf("Registered function execution error = %v", err)
	}
	if response == nil {
		t.Fatal("Registered function response was nil")
	}

	expectedToolResponseContent := "Note posted successfully: " + expectedNoteText + expectedAPIResponseNoteID
	if len(response.Content) != 1 || response.Content[0].Text == nil || *response.Content[0].Text != expectedToolResponseContent {
		t.Errorf("Expected tool response content '%s', got '%v'", expectedToolResponseContent, response.Content)
	}
}

func TestPostNoteTool_Register_MisskeyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMCPServer := mocks.NewMockServer(ctrl)
	mockMisskeyClient := mocks.NewMockClient(ctrl)
	mockNotesService := mocks.NewMockClientInterface(ctrl)

	tool := NewPostNoteTool()
	expectedError := errors.New("misskey API error")

	var registeredFunc func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error)

	mockMCPServer.EXPECT().
		RegisterTool(tool.GetName(), tool.GetDescription(), gomock.Any()).
		DoAndReturn(func(name, desc string, f interface{}) error {
			cb, ok := f.(func(arguments postNoteArguments) (*mcp_golang.ToolResponse, error))
			if !ok {
				t.Fatalf("RegisterTool was called with an unexpected function signature.")
			}
			registeredFunc = cb
			return nil
		}).
		Return(nil).
		Times(1)

	// Setup expectations for the Misskey client for when the registeredFunc is called
	mockMisskeyClient.EXPECT().Notes().Return(mockNotesService).Times(1)
	mockNotesService.EXPECT().
		Create(gomock.Any()).
		Return(nil, expectedError). // Simulate an error from Misskey
		Times(1)

	err := tool.Register(mockMCPServer, mockMisskeyClient)
	if err != nil {
		t.Fatalf("Register() error during setup = %v", err)
	}

	if registeredFunc == nil {
		t.Fatal("Registered function was not captured")
	}

	args := postNoteArguments{Text: "test"}
	response, err := registeredFunc(args)

	if err == nil {
		t.Errorf("Expected an error from registered function, but got nil")
	}
	if err.Error() != expectedError.Error() { // Compare error messages as errors.Is/As might be more involved with mock errors
		t.Errorf("Expected error message '%s', got '%s'", expectedError.Error(), err.Error())
	}
	if response != nil {
		t.Errorf("Expected nil response on error, got %v", response)
	}
}

func TestPostNoteTool_GetName(t *testing.T) {
	tool := NewPostNoteTool()
	if name := tool.GetName(); name != toolName {
		t.Errorf("GetName() = %s; want %s", name, toolName)
	}
}

func TestPostNoteTool_GetDescription(t *testing.T) {
	tool := NewPostNoteTool()
	if desc := tool.GetDescription(); desc != description {
		t.Errorf("GetDescription() = %s; want %s", desc, description)
	}
}
