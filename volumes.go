package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Volume holds information about a single mounted filesystem volume.
type Volume struct {
	Name       string `json:"name"`
	MountPoint string `json:"mountPoint"`
	FSType     string `json:"fsType"`
	TotalBytes uint64 `json:"totalBytes"`
	FreeBytes  uint64 `json:"freeBytes"`
	TotalHuman string `json:"totalHuman"`
	FreeHuman  string `json:"freeHuman"`
}

// handleListVolumes enumerates mounted volumes and returns filesystem
// metadata for each. The platform-specific listVolumes function is
// defined in volumes_darwin.go and volumes_linux.go.
func handleListVolumes(
	ctx context.Context,
	req *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, any, error) {
	volumes, err := listVolumes()
	if err != nil {
		return nil, nil, fmt.Errorf("list volumes: %w", err)
	}

	data, err := json.Marshal(volumes)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal volumes: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// humanBytes formats a byte count as a human-readable string with
// one decimal place (e.g., "29.1 GB").
func humanBytes(b uint64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
	)
	switch {
	case b >= tb:
		return fmt.Sprintf("%.1f TB", float64(b)/float64(tb))
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
