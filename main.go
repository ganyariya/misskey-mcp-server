package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"ganyariya's sample Calculator Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	calculator := mcp.NewTool("ganyariya-calculator tool",
		mcp.WithDescription("ganyariya basic calculator description"),
		mcp.WithString("operation", mcp.Required(), mcp.Description("operation to perform"), mcp.Enum("add", "subtract", "multiply", "divide")),
		mcp.WithNumber("x", mcp.Required()),
		mcp.WithNumber("y", mcp.Required(), mcp.Description("second operand")),
	)

	s.AddTool(calculator, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			return mcp.NewToolResultError("まだ割り算はできません"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}
}
