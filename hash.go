package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ComputeSHA1Input holds the parameters for the compute_sha1 tool.
type ComputeSHA1Input struct {
	Path string `json:"path" jsonschema:"Absolute path to the file to hash"`
}

// ComputeSHA1Result contains the SHA1 hash and file size of the
// hashed file.
type ComputeSHA1Result struct {
	SHA1     string `json:"sha1"`
	FileSize int64  `json:"fileSize"`
}

// handleComputeSHA1 computes the SHA1 hash of a file using streaming
// reads, so it handles files of any size without loading the entire
// contents into memory.
func handleComputeSHA1(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ComputeSHA1Input,
) (*mcp.CallToolResult, any, error) {
	f, err := os.Open(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %q: %w", in.Path, err)
	}
	defer f.Close()

	h := sha1.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return nil, nil, fmt.Errorf("read %q: %w", in.Path, err)
	}

	result := ComputeSHA1Result{
		SHA1:     hex.EncodeToString(h.Sum(nil)),
		FileSize: n,
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
