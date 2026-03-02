package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CleanDotFilesInput holds the parameters for the clean_dot_files
// tool.
type CleanDotFilesInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the directory to clean recursively"`
}

// handleCleanDotFiles removes AppleDouble resource fork files
// (names starting with "._") from a directory tree. macOS creates
// these automatically on FAT32 volumes because FAT32 has no
// extended attribute support. Flashcart menu systems like Wood R4
// display them as clutter when showHiddenFiles is enabled.
func handleCleanDotFiles(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in CleanDotFilesInput,
) (*mcp.CallToolResult, any, error) {
	var removed []string
	err := filepath.WalkDir(in.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Also skip .Trashes, .Spotlight-V100, .fseventsd
			// directories that macOS creates on external volumes.
			name := d.Name()
			if name == ".Trashes" || name == ".Spotlight-V100" || name == ".fseventsd" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), "._") {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove %q: %w", path, err)
			}
			rel, _ := filepath.Rel(in.Path, path)
			removed = append(removed, rel)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("clean dot files in %q: %w", in.Path, err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "removed %d AppleDouble files", len(removed))
	show := removed
	if len(show) > 20 {
		b.WriteString(" (showing first 20)")
		show = show[:20]
	}
	if len(show) > 0 {
		b.WriteString(":\n")
		for _, r := range show {
			b.WriteString("  ")
			b.WriteString(r)
			b.WriteByte('\n')
		}
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
	}, nil, nil
}
