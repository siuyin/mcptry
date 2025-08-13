package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

func main() {
	name := dflt.EnvString("NAME", "Siu Yin")
	svr := dflt.EnvString("SERVER", "myserver")
	log.Printf("SERVER=%s NAME=%s", svr, name)

	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout
	transport := mcp.NewCommandTransport(exec.Command(svr))
	session, err := client.Connect(ctx, transport)
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer session.Close()

	// List Tools
	lt, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatal("list tools: ", err)
	}

	ollamaTools(lt)

	toolParam := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": name},
	}
	mcpCallTool(session, toolParam)

	toolParam = &mcp.CallToolParams{
		Name:      "bye",
		Arguments: map[string]any{"name": name},
	}
	mcpCallTool(session, toolParam)

}

func ollamaTools(lt *mcp.ListToolsResult) {
	tools, _ := ConvertMCPToolsToOllamaTools(lt.Tools)
	for _, t := range tools {
		//fmt.Printf("%v\n", t)
		b, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			log.Fatal("marshal: ", err)
		}
		fmt.Printf("%s\n", b)
	}

	toolNames := []string{}
	for _, t := range tools {
		toolNames = append(toolNames, t.Function.Name)
	}
	fmt.Printf("%v\n", toolNames)
}

type ToolParams struct {
	Type       string                      `json:"type"`
	Defs       any                         `json:"$defs,omitempty"`
	Items      any                         `json:"items,omitempty"`
	Required   []string                    `json:"required"`
	Properties map[string]api.ToolProperty `json:"properties"`
}

func ConvertMCPToolsToOllamaTools(mcpTools []*mcp.Tool) ([]api.Tool, error) {
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

func mcpCallTool(session *mcp.ClientSession, params *mcp.CallToolParams) {
	ctx := context.Background()
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
