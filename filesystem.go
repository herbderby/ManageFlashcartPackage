package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDirectoryInput holds the parameters for the list_directory tool.
type ListDirectoryInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the directory to list"`
}

// DirEntry represents a single file or directory entry returned by
// list_directory.
type DirEntry struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

// handleListDirectory lists files and subdirectories at a path with
// sizes and modification times.
func handleListDirectory(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ListDirectoryInput,
) (*mcp.CallToolResult, any, error) {
	entries, err := os.ReadDir(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("list directory %q: %w", in.Path, err)
	}

	var result []DirEntry
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		result = append(result, DirEntry{
			Name:    e.Name(),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal entries: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

// CreateDirectoryInput holds the parameters for the create_directory tool.
type CreateDirectoryInput struct {
	Path string `json:"path" jsonschema:"Absolute path of the directory to create (parents created automatically)"`
}

// handleCreateDirectory creates a directory and any necessary parents.
func handleCreateDirectory(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in CreateDirectoryInput,
) (*mcp.CallToolResult, any, error) {
	if err := resolvedMkdirAll(in.Path, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create directory %q: %w", in.Path, err)
	}
	msg := fmt.Sprintf("created directory %s", in.Path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// MoveFileInput holds the parameters for the move_file tool.
type MoveFileInput struct {
	Source      string `json:"source" jsonschema:"Absolute path of the file or directory to move"`
	Destination string `json:"destination" jsonschema:"Absolute path of the destination"`
}

// handleMoveFile moves or renames a file or directory.
func handleMoveFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in MoveFileInput,
) (*mcp.CallToolResult, any, error) {
	if err := os.Rename(in.Source, in.Destination); err != nil {
		return nil, nil, fmt.Errorf("move %q to %q: %w", in.Source, in.Destination, err)
	}
	msg := fmt.Sprintf("moved %s to %s", in.Source, in.Destination)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// CopyFileInput holds the parameters for the copy_file tool.
type CopyFileInput struct {
	Source      string `json:"source" jsonschema:"Absolute path of the source file or directory"`
	Destination string `json:"destination" jsonschema:"Absolute path of the destination file or directory"`
	Recursive   bool   `json:"recursive,omitempty" jsonschema:"Set true to copy a directory tree recursively"`
}

// handleCopyFile copies a file or directory from source to
// destination. Parent directories of the destination are created
// automatically. When Recursive is true and the source is a
// directory, the entire tree is copied.
func handleCopyFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in CopyFileInput,
) (*mcp.CallToolResult, any, error) {
	info, err := os.Stat(in.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("stat source %q: %w", in.Source, err)
	}

	if info.IsDir() {
		if !in.Recursive {
			return nil, nil, fmt.Errorf("source %q is a directory; set recursive to true", in.Source)
		}
		count, err := copyDirRecursive(in.Source, in.Destination)
		if err != nil {
			return nil, nil, err
		}
		msg := fmt.Sprintf("copied %d files from %s to %s", count, in.Source, in.Destination)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		}, nil, nil
	}

	n, err := copySingleFile(in.Source, in.Destination)
	if err != nil {
		return nil, nil, err
	}
	msg := fmt.Sprintf("copied %d bytes from %s to %s", n, in.Source, in.Destination)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// copySingleFile copies a single file from src to dst, creating
// parent directories as needed. Returns the number of bytes copied.
func copySingleFile(src, dst string) (int64, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("open source %q: %w", src, err)
	}
	defer in.Close()

	if err := resolvedMkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, fmt.Errorf("create destination directory for %q: %w", dst, err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("create destination %q: %w", dst, err)
	}
	defer out.Close()

	n, err := io.Copy(out, in)
	if err != nil {
		return 0, fmt.Errorf("copy %q to %q: %w", src, dst, err)
	}
	return n, nil
}

// copyDirRecursive walks src and copies every file into dst,
// preserving the directory structure. Returns the number of files
// copied.
func copyDirRecursive(src, dst string) (int, error) {
	var count int
	err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("compute relative path for %q: %w", path, err)
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return resolvedMkdirAll(target, 0o755)
		}

		if _, err := copySingleFile(path, target); err != nil {
			return err
		}
		count++
		return nil
	})
	return count, err
}

// DeleteFileInput holds the parameters for the delete_file tool.
type DeleteFileInput struct {
	Path string `json:"path" jsonschema:"Absolute path of the file to delete (not recursive)"`
}

// handleDeleteFile deletes a single file. It does not delete
// directories or operate recursively.
func handleDeleteFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in DeleteFileInput,
) (*mcp.CallToolResult, any, error) {
	if err := os.Remove(in.Path); err != nil {
		return nil, nil, fmt.Errorf("delete %q: %w", in.Path, err)
	}
	msg := fmt.Sprintf("deleted %s", in.Path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// FileExistsInput holds the parameters for the file_exists tool.
type FileExistsInput struct {
	Path string `json:"path" jsonschema:"Absolute path to check"`
}

// FileExistsResult contains existence and metadata for a path.
type FileExistsResult struct {
	Exists bool  `json:"exists"`
	IsDir  bool  `json:"isDir"`
	Size   int64 `json:"size"`
}

// handleFileExists checks whether a path exists and reports whether
// it is a file or directory along with its size.
func handleFileExists(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in FileExistsInput,
) (*mcp.CallToolResult, any, error) {
	var result FileExistsResult
	info, err := os.Stat(in.Path)
	if err == nil {
		result.Exists = true
		result.IsDir = info.IsDir()
		result.Size = info.Size()
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
