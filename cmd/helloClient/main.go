package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	for _, t := range lt.Tools {
		fmt.Printf("tool name: %s, type: %v, required params: %v, types: %v \n",
			t.Name, t.InputSchema.Type, t.InputSchema.Required, requiredPropertyTypes(t))
	}

	// Call a tool on the server.
	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": name},
	}
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

func requiredPropertyTypes(tl *mcp.Tool) []string {
	pairs := []string{}
	for _, r := range tl.InputSchema.Required {
		typ := tl.InputSchema.Properties[r].Type
		pairs = append(pairs, fmt.Sprintf("%s: %s", r, typ))
	}

	return pairs
}
