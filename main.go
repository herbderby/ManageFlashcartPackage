// Package main implements a minimal MCP server for the HelloSDCard
// plugin. It registers a single tool, list_volumes, and serves
// requests over stdio using the official Go MCP SDK.
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hello-sdcard",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_volumes",
		Description: "List mounted volumes with filesystem type, total size, and free space.",
	}, handleListVolumes)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
