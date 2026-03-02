package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExtractArchiveInput holds the parameters for the extract_archive
// tool.
type ExtractArchiveInput struct {
	ArchivePath string `json:"archivePath" jsonschema:"Absolute path to the archive file (.7z or .zip)"`
	Destination string `json:"destination" jsonschema:"Absolute path to the directory to extract into"`
}

// handleExtractArchive extracts a .7z or .zip archive to the given
// destination directory. The archive format is detected by file
// extension. Returns the list of extracted file paths.
func handleExtractArchive(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ExtractArchiveInput,
) (*mcp.CallToolResult, any, error) {
	if err := resolvedMkdirAll(in.Destination, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create destination %q: %w", in.Destination, err)
	}

	ext := strings.ToLower(filepath.Ext(in.ArchivePath))
	var extracted []string
	var err error

	switch ext {
	case ".7z":
		extracted, err = extract7z(in.ArchivePath, in.Destination)
	case ".zip":
		extracted, err = extractZip(in.ArchivePath, in.Destination)
	default:
		return nil, nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
	if err != nil {
		return nil, nil, err
	}

	data, err := json.Marshal(extracted)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

// extract7z extracts a 7z archive to dst and returns the list of
// extracted file paths relative to dst.
func extract7z(archivePath, dst string) ([]string, error) {
	r, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open 7z %q: %w", archivePath, err)
	}
	defer r.Close()

	var extracted []string
	for _, f := range r.File {
		target := filepath.Join(dst, f.Name)

		// Guard against zip slip.
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dst)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := resolvedMkdirAll(target, 0o755); err != nil {
				return nil, fmt.Errorf("create dir %q in %q: %w", f.Name, dst, err)
			}
			continue
		}

		if err := resolvedMkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, fmt.Errorf("create parent dir for %q in %q: %w", f.Name, dst, err)
		}

		if err := extractFile(f, target); err != nil {
			return nil, err
		}
		extracted = append(extracted, f.Name)
	}
	return extracted, nil
}

// extractFile writes a single file from a 7z archive entry to disk.
func extractFile(f *sevenzip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open entry %q: %w", f.Name, err)
	}
	defer rc.Close()

	out, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("create %q: %w", target, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, rc); err != nil {
		return fmt.Errorf("extract %q: %w", f.Name, err)
	}
	return nil
}

// extractZip extracts a zip archive to dst and returns the list of
// extracted file paths relative to dst.
func extractZip(archivePath, dst string) ([]string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open zip %q: %w", archivePath, err)
	}
	defer r.Close()

	var extracted []string
	for _, f := range r.File {
		target := filepath.Join(dst, f.Name)

		// Guard against zip slip.
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dst)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := resolvedMkdirAll(target, 0o755); err != nil {
				return nil, fmt.Errorf("create dir %q in %q: %w", f.Name, dst, err)
			}
			continue
		}

		if err := resolvedMkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, fmt.Errorf("create parent dir for %q in %q: %w", f.Name, dst, err)
		}

		if err := extractZipFile(f, target); err != nil {
			return nil, err
		}
		extracted = append(extracted, f.Name)
	}
	return extracted, nil
}

// extractZipFile writes a single file from a zip archive entry to
// disk.
func extractZipFile(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open entry %q: %w", f.Name, err)
	}
	defer rc.Close()

	out, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("create %q: %w", target, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, rc); err != nil {
		return fmt.Errorf("extract %q: %w", f.Name, err)
	}
	return nil
}
