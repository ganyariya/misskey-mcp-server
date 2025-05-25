package actions

import (
	"fmt"
	// "log" // Replaced with logrus

	"github.com/arrow2nd/mcp-golang"
	"github.com/misskey-dev/misskey-go/misskey"
	"github.com/sirupsen/logrus"
)

// Registry holds registered actions.
type Registry struct {
	actions map[string]Action
	logger  logrus.FieldLogger
}

// NewRegistry creates a new Registry.
// It now accepts a logger. If nil, it uses the standard logger.
func NewRegistry(logger logrus.FieldLogger) *Registry {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return &Registry{
		actions: make(map[string]Action),
		logger:  logger.WithField("component", "action_registry"),
	}
}

// Register adds an action to the registry.
// It returns an error if an action with the same name is already registered.
func (r *Registry) Register(action Action) error {
	name := action.Name()
	if _, exists := r.actions[name]; exists {
		r.logger.Errorf("Action with name '%s' already registered", name)
		return fmt.Errorf("action with name '%s' already registered", name)
	}
	r.actions[name] = action
	r.logger.Infof("Action '%s' registered successfully.", name)
	return nil
}

// RegisterActionsWithServer registers all actions in the registry with the MCP server.
// It now requires a baseLogger to create context-specific loggers.
func (r *Registry) RegisterActionsWithServer(server *mcp_golang.Server, misskeyClient *misskey.Client, baseLogger logrus.FieldLogger) error {
	if baseLogger == nil {
		baseLogger = logrus.StandardLogger() // Fallback, though main.go should provide one
	}
	actionCtxLogger := baseLogger.WithField("component", "action_executor")
	actionCtx := NewActionContext(misskeyClient, actionCtxLogger)

	if len(r.actions) == 0 {
		r.logger.Info("No actions registered in the registry to register with MCP server.")
		return nil
	}

	r.logger.Infof("Registering %d actions with MCP server...", len(r.actions))

	for actionName, action := range r.actions {
		// Capture loop variable for the closure
		currentAction := action // Essential for correct operation in closure
		currentActionName := actionName

		r.logger.WithFields(logrus.Fields{
			"action_name": currentActionName,
			"description": currentAction.Description(),
		}).Info("Attempting to register tool with MCP server")

		err := server.RegisterTool(
			currentActionName,
			currentAction.Description(),
			func(args interface{}) (*mcp_golang.ToolResponse, error) {
				// This logger will have the "action_executor" field from actionCtx.Logger()
				actionSpecificLogger := actionCtx.Logger().WithField("action_name", currentActionName)
				actionSpecificLogger.WithField("params", args).Info("Executing action")
				
				// Create a new ActionContext for this specific execution, with a more specific logger
				execCtx := NewActionContext(misskeyClient, actionSpecificLogger)

				response, err := currentAction.Execute(execCtx, args) // Pass the new context
				if err != nil {
					actionSpecificLogger.WithError(err).Error("Action execution failed")
				} else {
					actionSpecificLogger.Info("Action executed successfully")
				}
				return response, err
			},
			currentAction.Params(),
		)
		if err != nil {
			r.logger.WithError(err).Errorf("Failed to register tool %s with MCP server", currentActionName)
			return fmt.Errorf("failed to register tool %s: %w", currentActionName, err)
		}
		r.logger.Infof("Tool '%s' registered successfully with MCP server.", currentActionName)
	}
	r.logger.Info("All actions registered with MCP server.")
	return nil
}
