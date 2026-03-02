package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// setStandardHeaders adds User-Agent and Accept-Encoding headers to
// an HTTP request. Some CDNs (notably blobfrii, used by
// archive.flashcarts.net) drop connections from bare clients that
// send no User-Agent.
func setStandardHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Flashcart-Tools/0.2.3")
}

// DownloadFileInput holds the parameters for the download_file tool.
type DownloadFileInput struct {
	URL         string `json:"url" jsonschema:"URL to download"`
	Destination string `json:"destination" jsonschema:"Absolute path to save the downloaded file"`
}

// DownloadFileResult reports the outcome of a file download.
type DownloadFileResult struct {
	BytesWritten int64 `json:"bytesWritten"`
	StatusCode   int   `json:"statusCode"`
}

// handleDownloadFile downloads a URL to a local file path. Parent
// directories of the destination are created automatically.
func handleDownloadFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in DownloadFileInput,
) (*mcp.CallToolResult, any, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, in.URL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}
	setStandardHeaders(httpReq)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("download %q: %w", in.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		result := DownloadFileResult{StatusCode: resp.StatusCode}
		data, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
			IsError: true,
		}, nil, nil
	}

	if err := resolvedMkdirAll(filepath.Dir(in.Destination), 0o755); err != nil {
		return nil, nil, fmt.Errorf("create destination directory for %q: %w", in.Destination, err)
	}

	f, err := os.Create(in.Destination)
	if err != nil {
		return nil, nil, fmt.Errorf("create file %q: %w", in.Destination, err)
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("write to file: %w", err)
	}

	result := DownloadFileResult{
		BytesWritten: n,
		StatusCode:   resp.StatusCode,
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

// maxFetchBody is the maximum number of bytes returned by fetch_url.
const maxFetchBody = 1 << 20 // 1 MiB

// FetchURLInput holds the parameters for the fetch_url tool.
type FetchURLInput struct {
	URL string `json:"url" jsonschema:"URL to fetch"`
}

// handleFetchURL fetches a URL and returns the response body as text,
// truncated to 1 MiB.
func handleFetchURL(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in FetchURLInput,
) (*mcp.CallToolResult, any, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, in.URL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}
	setStandardHeaders(httpReq)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch %q: %w", in.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFetchBody))
	if err != nil {
		return nil, nil, fmt.Errorf("read response body: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(body)}},
	}, nil, nil
}
