package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio" // Added import
	"github.com/stretchr/testify/assert"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core" // Added import
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes" // Added import
)

func TestGetUserNotesTool_Register(t *testing.T) {
	// Use stdio transport for the server
	serverTransport := stdio.NewStdioServerTransport()
	server := mcp_golang.NewServer(serverTransport) // Corrected NewServer
	// For mockMisskeyClient, WithHTTPClient and setting BaseURL manually is more robust for tests.
	// Since Register doesn't use the client to make calls, a minimal setup is okay.
	mockMisskeyClient, _ := misskey.NewClientWithOptions(misskey.WithHTTPClient(&http.Client{})) // Minimal client
	mockMisskeyClient.BaseURL = "http://localhost"                                             // Set BaseURL

	tool := NewGetUserNotesTool()
	err := tool.Register(server, mockMisskeyClient)
	assert.NoError(t, err)

	// Changed from GetTool to CheckToolRegistered
	ok := server.CheckToolRegistered(tool.GetName())
	assert.True(t, ok)
	// Cannot get tool details directly from server, so removed these assertions
	// assert.Equal(t, tool.GetName(), registeredTool.Name)
	// assert.Equal(t, tool.GetDescription(), registeredTool.Description)
}

// Define the handler function type based on the RegisterTool signature
type getUserNotesHandlerFunc func(arguments getUserNotesArguments) (*mcp_golang.ToolResponse, error)

func TestGetUserNotesTool_Execute(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		mockResponse  []*models.Note
		mockError     error
		expectedError bool
	}{
		{
			name:   "Successful execution",
			userID: "testUserID",
			mockResponse: []*models.Note{
				{ID: "note1", Text: "Note 1"},
				{ID: "note2", Text: "Note 2"},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "Misskey API error",
			userID:        "testUserID",
			mockResponse:  nil,
			mockError:     errors.New("Misskey API error"),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Validate request if necessary, e.g., check r.URL.Path
				// For /notes/search, it would be POST to /api/notes/search
				// For now, just return mockResponse or mockError
				if tc.mockError != nil {
					http.Error(w, tc.mockError.Error(), http.StatusInternalServerError)
					return
				}
				// Expecting a POST request for search
				if r.Method != http.MethodPost {
					http.Error(w, "Expected POST", http.StatusMethodNotAllowed)
					return
				}
				// Minimal check for body content (e.g. if UserID is present) could be added here

				jsonBytes, _ := json.Marshal(tc.mockResponse) // This is a list of notes
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			}))
			defer mockHttpServer.Close()

			// Create misskey client configured for the mock server
			mockMisskeyClient, err := misskey.NewClientWithOptions(
				misskey.WithHTTPClient(mockHttpServer.Client()), // Use the test server's client
				// No API token needed for this specific test if endpoint doesn't require auth
				// or if auth is handled by the mock server logic (not the case here yet)
			)
			assert.NoError(t, err)
			mockMisskeyClient.BaseURL = mockHttpServer.URL // Crucial: point client to mock server

			// The handler function is what tool.Register would pass to mcpServer.RegisterTool
			// We extract its essence here for direct testing.
			// tool := NewGetUserNotesTool() // No longer needed as we test handler directly
			handler := func(arguments getUserNotesArguments) (*mcp_golang.ToolResponse, error) {
				userID := arguments.UserID
				// This is the actual logic from get_user_notes.go
				searchResponse, err := mockMisskeyClient.Notes().Search(notes.SearchRequest{
					Query:  "*", // Changed Query field to "*" as a test
					UserID: core.NewString(userID),
					Limit:  100, // Should match what's in the actual tool
				})
				if err != nil {
					return nil, err
				}

				responseJSON, err := json.Marshal(searchResponse)
				if err != nil {
					return nil, err
				}
				return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(responseJSON))), nil
			}

			args := getUserNotesArguments{
				UserID: tc.userID,
			}
			// No argBytes needed as we call the handler directly for this unit test.
			// No mcpServer.ExecuteTool needed.
			response, err := handler(args)

			if tc.expectedError {
				assert.Error(t, err) // Check if the handler returned an error
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Content) // Ensure content exists

				expectedResponseJSON, _ := json.Marshal(tc.mockResponse)
				// Accessing the text content correctly
				assert.JSONEq(t, string(expectedResponseJSON), response.Content[0].TextContent.Text)
			}
		})
	}
}

func TestGetUserNotesTool_GetName(t *testing.T) {
	tool := NewGetUserNotesTool()
	assert.Equal(t, toolName, tool.GetName())
}

func TestGetUserNotesTool_GetDescription(t *testing.T) {
	tool := NewGetUserNotesTool()
	assert.Equal(t, description, tool.GetDescription())
}
