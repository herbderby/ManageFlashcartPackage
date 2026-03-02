package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	utls "github.com/refraction-networking/utls"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// browserClient is an HTTP client whose TLS fingerprint mimics
// Chrome. CDNs like blobfrii (used by archive.flashcarts.net)
// perform JA3/JA4 TLS fingerprinting and reject Go's standard
// crypto/tls ClientHello with a TCP RST.
var browserClient = &http.Client{
	Transport: &http.Transport{
		DialTLS: func(network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}
			conn, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			// Get the Chrome ClientHello spec and force HTTP/1.1
			// ALPN. The preset advertises h2, but Go's
			// http.Transport can't handle HTTP/2 over a custom
			// DialTLS connection, so we strip it.
			spec, err := utls.UTLSIdToSpec(utls.HelloChrome_Auto)
			if err != nil {
				conn.Close()
				return nil, err
			}
			for _, ext := range spec.Extensions {
				if alpn, ok := ext.(*utls.ALPNExtension); ok {
					alpn.AlpnProtocols = []string{"http/1.1"}
				}
			}
			uconn := utls.UClient(conn, &utls.Config{
				ServerName: host,
			}, utls.HelloCustom)
			if err := uconn.ApplyPreset(&spec); err != nil {
				conn.Close()
				return nil, err
			}
			if err := uconn.Handshake(); err != nil {
				conn.Close()
				return nil, err
			}
			return uconn, nil
		},
		DialContext: (&net.Dialer{}).DialContext,
	},
}

// newRequest creates an HTTP GET request with a browser-like
// User-Agent header.
func newRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	return req, nil
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
	httpReq, err := newRequest(ctx, in.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := browserClient.Do(httpReq)
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
	httpReq, err := newRequest(ctx, in.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := browserClient.Do(httpReq)
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
