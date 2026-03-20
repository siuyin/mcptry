// Package olamtl provide tools to convert mcp tools to ollama tools.
package olamtl

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

type mySchema struct {
	Properties map[string]map[string]string `json:"properties"`
	Required   []string                     `json:"required"`
}

// FromMCP converts MCP tools to ollama tools.
func FromMCP(mcpTools []*mcp.Tool) ([]api.Tool, error) {
	debugStr := dflt.EnvString("DEBUG", "false")
	debug := false
	if debugStr != "false" {
		debug = true
	}

	ollamaTools := make([]api.Tool, len(mcpTools))
	for i, tool := range mcpTools {
		// Convert mcp.ToolProperty to the required JSON schema format for Ollama.
		toolPropsMap := api.NewToolPropertiesMap()

		if debug {
			fmt.Printf("Tool: %s\n", tool.Name)
		}

		b, err := json.Marshal(tool.InputSchema)
		if err != nil {
			log.Fatal(err)
		}
		sch := &mySchema{}
		if err := json.Unmarshal(b, sch); err != nil {
			log.Fatal(err)
		}
		for k, v := range sch.Properties {
			if debug {
				fmt.Printf("%q:", k)
			}
			for kk, vv := range v {
				if debug {
					fmt.Printf(" %q:%q", kk, vv)
				}
			}
			toolPropsMap.Set(k, api.ToolProperty{Type: api.PropertyType{v["type"]}, Description: v["description"]})
			if debug {
				fmt.Println()
			}
		}
		if debug {
			fmt.Printf("required: %v\n", sch.Required)
			fmt.Println("-------")
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
					Required:   sch.Required,
				},
			},
		}
	}
	return ollamaTools, nil
}
