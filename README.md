# Model Context Protocol experiments

This uses the MCP library from: https://github.com/modelcontextprotocol/go-sdk

## hello: Client and Server
Build the executable `myserver` and place the binary in your PATH.

```
go build -o ~/bin/myserver ./cmd/helloServer/
```

Run the client:
```
go run ./cmd/helloClient/

```

## Using MCP with Ollama
MCP tools are defined in jsonschema while Ollama has its onw tool definition.

cmd/ollamaTools/main.go shows how to convert MCP tools so that they are callable from Ollama.

This uses the olamtl.FromMCP function.
