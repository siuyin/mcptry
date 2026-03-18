// Package olamtl provide tools to convert mcp tools to ollama tools.
package olamtl

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

// FromMCP converts MCP tools to ollama tools.
func FromMCP(mcpTools []*mcp.Tool) ([]api.Tool, error) {
	ollamaTools := make([]api.Tool, len(mcpTools))
	for i, tool := range mcpTools {
		// Convert mcp.ToolProperty to the required JSON schema format for Ollama.
		toolPropsMap := api.NewToolPropertiesMap()
		required := []string{}
		required = append(required, tool.InputSchema.Required...)
		for propName, prop := range tool.InputSchema.Properties {
			toolPropsMap.Set(propName, api.ToolProperty{
				Type:        api.PropertyType{prop.Type},
				Description: prop.Description,
			})
		}

		// Marshal the parameters into the format expected by the Ollama API
		ollamaTools[i] = api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: api.ToolFunctionParameters{
					Type:       "object",
					Properties: toolPropsMap,
					Required:   required,
				},
			},
		}
	}
	return ollamaTools, nil
}
