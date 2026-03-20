package olamtl

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/siuyin/dflt"
)

func ExampleFromMCP() {
	svr := dflt.EnvString("SERVER", "serverV2")
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	transport := &mcp.CommandTransport{Command: exec.Command(svr)}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer session.Close()

	lt, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		log.Fatal("list tools: ", err)
	}

	olamaTools, err := FromMCP(lt.Tools)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", olamaTools)
	// // Output:
}
