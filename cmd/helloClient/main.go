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
	svr := dflt.EnvString("SERVER", "serverV2")
	log.Printf("SERVER=%s NAME=%s", svr, name)

	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout
	transport := &mcp.CommandTransport{Command: exec.Command(svr)}
	session, err := client.Connect(ctx, transport, nil)
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
		fmt.Printf("tool name: %s %q\nInput Schema: %v\n", t.Name, t.Description, t.InputSchema)
	}

	callTool(ctx, session, "greet", map[string]any{"name": name})
	callTool(ctx, session, "utcTime", map[string]any{})
	callTool(ctx, session, "weather", map[string]any{"location": "Bukit Batok,Singapore"})
	callTool(ctx, session, "stocks", map[string]any{"code": "GOOGL"})
	callTool(ctx, session, "bye", map[string]any{"name": name})

}
func callTool(ctx context.Context, session *mcp.ClientSession, name string, args map[string]any) {
	params := &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool for %s failed: %v", name, err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}

//func requiredPropertyTypes(tl *mcp.Tool) []string {
//	pairs := []string{}
//	for _, r := range tl.InputSchema.Required {
//		typ := tl.InputSchema.Properties[r].Type
//		pairs = append(pairs, fmt.Sprintf("%s: %s", r, typ))
//	}
//
//	return pairs
//}
