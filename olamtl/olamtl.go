// Package olamtl provide tools to convert mcp tools to ollama tools.
package olamtl

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type ToolParams struct {
	Type       string                      `json:"type"`
	Defs       any                         `json:"$defs,omitempty"`
	Items      any                         `json:"items,omitempty"`
	Required   []string                    `json:"required"`
	Properties map[string]api.ToolProperty `json:"properties"`
}

// FromMCP converts MCP tools to ollama tools.
func FromMCP(mcpTools []*mcp.Tool) ([]api.Tool, error) {
	ollamaTools := make([]api.Tool, len(mcpTools))
	for i, tool := range mcpTools {
		// Convert mcp.ToolProperty to the required JSON schema format for Ollama.
		properties := make(map[string]api.ToolProperty)
		required := []string{}
		required = append(required, tool.InputSchema.Required...)
		for propName, prop := range tool.InputSchema.Properties {
			properties[propName] = api.ToolProperty{
				Type:        api.PropertyType{prop.Type},
				Description: prop.Description,
			}
		}

		// Marshal the parameters into the format expected by the Ollama API
		ollamaTools[i] = api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: ToolParams{
					Type:       "object",
					Properties: properties,
					Required:   required,
				},
			},
		}
	}
	return ollamaTools, nil
}
