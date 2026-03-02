package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReadFileInput holds the parameters for the read_file tool.
type ReadFileInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the text file to read"`
}

// handleReadFile reads a text file and returns its contents as a
// string. Use this for configuration files, INI files, and other
// human-readable text. For binary data use read_bytes instead.
func handleReadFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ReadFileInput,
) (*mcp.CallToolResult, any, error) {
	data, err := os.ReadFile(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("read %q: %w", in.Path, err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

// WriteFileInput holds the parameters for the write_file tool.
type WriteFileInput struct {
	Path    string `json:"path" jsonschema:"Absolute path to the file to write"`
	Content string `json:"content" jsonschema:"Text content to write to the file"`
}

// handleWriteFile writes text content to a file, creating parent
// directories as needed. The file is created with mode 0644. If the
// file already exists it is overwritten.
func handleWriteFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in WriteFileInput,
) (*mcp.CallToolResult, any, error) {
	if err := resolvedMkdirAll(filepath.Dir(in.Path), 0o755); err != nil {
		return nil, nil, fmt.Errorf("create directory for %q: %w", in.Path, err)
	}

	if err := os.WriteFile(in.Path, []byte(in.Content), 0o644); err != nil {
		return nil, nil, fmt.Errorf("write %q: %w", in.Path, err)
	}

	msg := fmt.Sprintf("wrote %d bytes to %s", len(in.Content), in.Path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}
