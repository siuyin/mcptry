package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HiParams struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

func SayHi(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[HiParams]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Hi " + req.Arguments.Name}},
	}, nil
}

func main() {
	log.Println("myserver running")
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Println("run: ", err)
	}

}
