package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReadJSONInput holds the parameters for the read_json tool.
type ReadJSONInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the JSON file to read"`
}

// handleReadJSON reads a JSON file and returns the parsed contents.
func handleReadJSON(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ReadJSONInput,
) (*mcp.CallToolResult, any, error) {
	data, err := os.ReadFile(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("read %q: %w", in.Path, err)
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, nil, fmt.Errorf("parse JSON from %q: %w", in.Path, err)
	}

	out, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(out)}},
	}, nil, nil
}

// WriteJSONInput holds the parameters for the write_json tool.
type WriteJSONInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the JSON file to write"`
	Data any    `json:"data" jsonschema:"JSON data to write to the file"`
}

// handleWriteJSON marshals data as indented JSON and writes it to a
// file. Parent directories are created automatically.
func handleWriteJSON(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in WriteJSONInput,
) (*mcp.CallToolResult, any, error) {
	out, err := json.MarshalIndent(in.Data, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal JSON: %w", err)
	}

	if err := resolvedMkdirAll(filepath.Dir(in.Path), 0o755); err != nil {
		return nil, nil, fmt.Errorf("create directory for %q: %w", in.Path, err)
	}

	if err := os.WriteFile(in.Path, append(out, '\n'), 0o644); err != nil {
		return nil, nil, fmt.Errorf("write %q: %w", in.Path, err)
	}

	msg := fmt.Sprintf("wrote JSON to %s", in.Path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}
