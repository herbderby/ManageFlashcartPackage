package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

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

// handleListVolumes enumerates /Volumes/ and returns filesystem
// metadata for each mounted volume. It uses syscall.Statfs to
// retrieve filesystem type, total size, and free space.
func handleListVolumes(
	ctx context.Context,
	req *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, any, error) {
	entries, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil, nil, fmt.Errorf("read /Volumes: %w", err)
	}

	var volumes []Volume
	for _, entry := range entries {
		mountPoint := "/Volumes/" + entry.Name()
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPoint, &stat); err != nil {
			continue
		}

		totalBytes := uint64(stat.Bsize) * stat.Blocks
		freeBytes := uint64(stat.Bsize) * stat.Bavail

		volumes = append(volumes, Volume{
			Name:       entry.Name(),
			MountPoint: mountPoint,
			FSType:     int8SliceToString(stat.Fstypename[:]),
			TotalBytes: totalBytes,
			FreeBytes:  freeBytes,
			TotalHuman: humanBytes(totalBytes),
			FreeHuman:  humanBytes(freeBytes),
		})
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

// int8SliceToString converts a null-terminated int8 slice (as
// returned by Darwin's Statfs_t) to a Go string.
func int8SliceToString(s []int8) string {
	buf := make([]byte, 0, len(s))
	for _, b := range s {
		if b == 0 {
			break
		}
		buf = append(buf, byte(b))
	}
	return string(buf)
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
