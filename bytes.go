package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReadBytesInput holds the parameters for the read_bytes tool.
type ReadBytesInput struct {
	Path   string `json:"path" jsonschema:"Absolute path to the file to read"`
	Offset int64  `json:"offset" jsonschema:"Byte offset to start reading from"`
	Length int    `json:"length" jsonschema:"Number of bytes to read"`
}

// ReadBytesResult contains the hex and ASCII representations of the
// bytes read from a file.
type ReadBytesResult struct {
	Hex   string `json:"hex"`
	ASCII string `json:"ascii"`
}

// handleReadBytes reads length bytes at the given offset from a file
// and returns both hex and ASCII representations. Non-printable ASCII
// bytes are replaced with '.'.
func handleReadBytes(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ReadBytesInput,
) (*mcp.CallToolResult, any, error) {
	f, err := os.Open(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %q: %w", in.Path, err)
	}
	defer f.Close()

	if _, err := f.Seek(in.Offset, 0); err != nil {
		return nil, nil, fmt.Errorf("seek to offset %d: %w", in.Offset, err)
	}

	buf := make([]byte, in.Length)
	n, err := f.Read(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("read %d bytes at offset %d: %w", in.Length, in.Offset, err)
	}
	buf = buf[:n]

	ascii := make([]byte, n)
	for i, b := range buf {
		if b >= 0x20 && b <= 0x7E {
			ascii[i] = b
		} else {
			ascii[i] = '.'
		}
	}

	result := ReadBytesResult{
		Hex:   hex.EncodeToString(buf),
		ASCII: string(ascii),
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
