package actions_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ganyariya/misskey-mcp-server/internal/actions"
	mcp_golang_arrow "github.com/arrow2nd/mcp-golang" // For server interface and response type as per action.go/registry.go
	"github.com/misskey-dev/misskey-go/misskey"
	"github.com/sirupsen/logrus"
)

// mockAction implements actions.Action for testing purposes.
type mockAction struct {
	name        string
	description string
	params      interface{}
	executeFunc func(ctx actions.ActionContext, params interface{}) (*mcp_golang_arrow.ToolResponse, error)
}

func (m *mockAction) Name() string {
	return m.name
}

func (m *mockAction) Description() string {
	return m.description
}

func (m *mockAction) Params() interface{} {
	// Return a new instance for safety, similar to real actions
	if m.params != nil {
		// Simple way to "clone" for testing; real params might need deep copy or specific new instances
		return reflect.New(reflect.TypeOf(m.params).Elem()).Interface()
	}
	return nil
}

func (m *mockAction) Execute(ctx actions.ActionContext, params interface{}) (*mcp_golang_arrow.ToolResponse, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, params)
	}
	return mcp_golang_arrow.NewToolResponse(mcp_golang_arrow.NewTextContent("mock execute success")), nil
}

// mockMCPTool holds registration call details.
type mockMCPToolRegistration struct {
	name        string
	description string
	handler     func(interface{}) (*mcp_golang_arrow.ToolResponse, error)
	argDef      interface{}
}

// mockMCPServer implements a minimal part of mcp_golang.Server for testing.
type mockMCPServer struct {
	registrations []mockMCPToolRegistration
}

func (m *mockMCPServer) RegisterTool(
	name string,
	description string,
	handler func(interface{}) (*mcp_golang_arrow.ToolResponse, error),
	argDef interface{},
) error {
	m.registrations = append(m.registrations, mockMCPToolRegistration{
		name:        name,
		description: description,
		handler:     handler,
		argDef:      argDef,
	})
	return nil
}

// Unused methods to satisfy potential interface needs, though not strictly required by current registry code
func (m *mockMCPServer) Serve() error { return nil }
func (m *mockMCPServer) Stop() error  { return nil }
func (m *mockMCPServer) Name() string { return "mock-server" }


func TestNewRegistry(t *testing.T) {
	logger := logrus.NewTestLogger()
	registry := actions.NewRegistry(logger)

	if registry == nil {
		t.Fatal("NewRegistry returned nil")
	}
	// Internal fields like 'logger' and 'actions' map are not exported,
	// so we can't directly test them for nil or emptiness without reflection
	// or by observing behavior (e.g., registering an action).
	// For now, we trust that if NewRegistry doesn't panic and returns non-nil, it's minimally functional.
	// We will test behavior via Register method.
}

func TestRegistry_Register(t *testing.T) {
	logger := logrus.NewTestLogger()
	registry := actions.NewRegistry(logger)

	action1 := &mockAction{name: "testAction1", description: "Test Action 1", params: &struct{}{}}

	// Test successful registration
	err := registry.Register(action1)
	if err != nil {
		t.Errorf("Registering action1 failed: %v", err)
	}

	// Test registering another action with the same name
	action2 := &mockAction{name: "testAction1", description: "Duplicate Action", params: &struct{}{}}
	err = registry.Register(action2)
	if err == nil {
		t.Errorf("Registering action2 with duplicate name should have failed, but err is nil")
	}
}

func TestRegistry_RegisterActionsWithServer(t *testing.T) {
	logger := logrus.NewTestLogger()
	registry := actions.NewRegistry(logger)
	mockServer := &mockMCPServer{}

	// 1. Test with no actions
	err := registry.RegisterActionsWithServer(mockServer, nil, logger)
	if err != nil {
		t.Errorf("RegisterActionsWithServer with no actions failed: %v", err)
	}
	if len(mockServer.registrations) != 0 {
		t.Errorf("Expected 0 registrations for empty registry, got %d", len(mockServer.registrations))
	}

	// 2. Test with one action
	type testParams struct {
		Data string `json:"data"`
	}
	action1 := &mockAction{
		name:        "testAction1",
		description: "Test Action 1",
		params:      &testParams{},
		executeFunc: func(ctx actions.ActionContext, params interface{}) (*mcp_golang_arrow.ToolResponse, error) {
			if p, ok := params.(*testParams); ok {
				return mcp_golang_arrow.NewToolResponse(mcp_golang_arrow.NewTextContent("executed " + p.Data)), nil
			}
			return nil, fmt.Errorf("bad params in mock execute")
		},
	}

	err = registry.Register(action1)
	if err != nil {
		t.Fatalf("Failed to register action1 for server test: %v", err)
	}

	err = registry.RegisterActionsWithServer(mockServer, nil, logger)
	if err != nil {
		t.Errorf("RegisterActionsWithServer with one action failed: %v", err)
	}

	if len(mockServer.registrations) != 1 {
		t.Fatalf("Expected 1 registration, got %d", len(mockServer.registrations))
	}

	reg := mockServer.registrations[0]
	if reg.name != action1.Name() {
		t.Errorf("Expected registered name '%s', got '%s'", action1.Name(), reg.name)
	}
	if reg.description != action1.Description() {
		t.Errorf("Expected registered description '%s', got '%s'", action1.Description(), reg.description)
	}
	if reflect.TypeOf(reg.argDef) != reflect.TypeOf(action1.Params()) {
		t.Errorf("Expected registered argDef type %T, got %T", action1.Params(), reg.argDef)
	}

	// 3. Test handler execution (optional but good)
	if reg.handler == nil {
		t.Fatal("Registered handler is nil")
	}

	// Create a test ActionContext
	testActionCtx := actions.NewActionContext(nil, logger.WithField("test_handler", action1.Name()))
	
	// To properly test the handler, we need to pass the correct param type.
	// The handler itself will receive this from the mcp-golang library after JSON unmarshalling.
	// Here, we simulate that the unmarshalling has happened and the handler gets the populated struct.
	handlerParam := &testParams{Data: "testdata"}
	
	response, err := reg.handler(handlerParam) // Pass the populated struct directly
	if err != nil {
		t.Errorf("Handler execution failed: %v", err)
	}
	if response == nil {
		t.Fatal("Handler response is nil")
	}
	expectedContent := mcp_golang_arrow.NewTextContent("executed testdata")
	if !reflect.DeepEqual(response.Content, expectedContent) {
		t.Errorf("Handler response content mismatch. Expected: %+v, Got: %+v", expectedContent, response.Content)
	}
}

func TestRegisterActionsWithServer_NilLogger(t *testing.T) {
	registry := actions.NewRegistry(nil) // Pass nil logger to NewRegistry
	mockServer := &mockMCPServer{}
	action1 := &mockAction{name: "testAction1", description: "Test Action 1", params: &struct{}{}}

	if err := registry.Register(action1); err != nil {
		t.Fatalf("Failed to register action: %v", err)
	}

	// Pass nil for baseLogger to RegisterActionsWithServer
	err := registry.RegisterActionsWithServer(mockServer, (*misskey.Client)(nil), nil)
	if err != nil {
		t.Errorf("RegisterActionsWithServer with nil baseLogger failed: %v", err)
	}
	if len(mockServer.registrations) != 1 {
		t.Errorf("Expected 1 registration with nil baseLogger, got %d", len(mockServer.registrations))
	}
}
