package actions

import (
	"github.com/arrow2nd/mcp-golang"
	"github.com/misskey-dev/misskey-go/misskey"
	"github.com/sirupsen/logrus"
)

// ActionContext holds shared dependencies for actions.
type ActionContext struct {
	MisskeyClient *misskey.Client
	logger        logrus.FieldLogger
	// Add other shared dependencies here if needed in the future
}

// NewActionContext creates a new ActionContext.
func NewActionContext(client *misskey.Client, logger logrus.FieldLogger) *ActionContext {
	if logger == nil {
		// Fallback to standard logger if nil is provided, though it's better to ensure a logger is always passed.
		logger = logrus.StandardLogger()
	}
	return &ActionContext{
		MisskeyClient: client,
		logger:        logger,
	}
}

// Logger returns the logger associated with the ActionContext.
func (ctx *ActionContext) Logger() logrus.FieldLogger {
	return ctx.logger
}

// Action defines the interface for all actions.
type Action interface {
	// Name returns the action's name (used for MCP registration and lookup).
	Name() string
	// Description returns a user-friendly description for the action.
	Description() string
	// Params returns a new instance of the struct that defines the parameters for this action.
	// This will be used for MCP tool registration to define the expected JSON schema.
	Params() interface{}
	// Execute executes the action.
	// params will be the populated struct based on the type returned by Params().
	Execute(ctx ActionContext, params interface{}) (*mcp_golang.ToolResponse, error)
}
