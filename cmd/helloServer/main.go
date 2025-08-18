package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	log.Println("myserver running")
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "welcome a person by name by saying hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "bye", Description: "send off a person by name by saying goodbye"}, SayBye)
	mcp.AddTool(server, &mcp.Tool{Name: "utcTime", Description: "get the current time in UTC."}, utcTime)
	mcp.AddTool(server, &mcp.Tool{Name: "weather", Description: "get the weather forecast (temperature, humidity) for a given location"}, weather)
	mcp.AddTool(server, &mcp.Tool{Name: "stocks", Description: "gets current stock price for a stock ticker. eg. AAPL"}, stock)

	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Println("run: ", err)
	}

}

type NameParams struct {
	Name string `json:"name" jsonschema:"the name of the person to address"`
}

func SayHi(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[NameParams]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Hi " + req.Arguments.Name}},
	}, nil
}

func SayBye(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[NameParams]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Goodbye " + req.Arguments.Name}},
	}, nil
}

type timeInput struct {
	Location string `json:"location,omitempty"`
}

func utcTime(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[timeInput]) (*mcp.CallToolResultFor[any], error) {
	ret := ""
	t := time.Now().UTC().Format("15:04:05.000")

	if req.Arguments.Location != "" {
		ret = fmt.Sprintf("I can't provide the time for %s. But the time in UTC is %s", req.Arguments.Location, t)
		return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: ret}}}, nil
	}

	ret = fmt.Sprintf("The time in UTC is %s", t)
	return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: ret}}}, nil
}

type weatherInput struct {
	Location string `json:"location"`
}

func weather(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[weatherInput]) (*mcp.CallToolResultFor[any], error) {
	ret := fmt.Sprintf("The weather in %s is a Sunny 30°C. Rain is expected later.", req.Arguments.Location)
	return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: ret}}}, nil
}

type stockInput struct {
	Code string `json:"code" jsonschema:"the stock code to retrieve eg. AAPL"`
}

func stock(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[stockInput]) (*mcp.CallToolResultFor[any], error) {
	ret := fmt.Sprintf("The current price for  %s is USD300.", req.Arguments.Code)
	return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: ret}}}, nil
}
