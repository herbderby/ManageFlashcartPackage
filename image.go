package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/image/draw"
)

// ResizeImageInput holds the parameters for the resize_image tool.
type ResizeImageInput struct {
	Source      string `json:"source" jsonschema:"Absolute path to the source PNG image"`
	Destination string `json:"destination" jsonschema:"Absolute path for the resized output PNG"`
	Width       int    `json:"width" jsonschema:"Target width in pixels"`
	Height      int    `json:"height" jsonschema:"Target height in pixels"`
}

// ResizeImageResult reports the dimensions and file size of the
// resized image.
type ResizeImageResult struct {
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	FileSize int64 `json:"fileSize"`
}

// handleResizeImage scales a PNG image to the given dimensions using
// CatmullRom interpolation and writes the result as PNG.
func handleResizeImage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ResizeImageInput,
) (*mcp.CallToolResult, any, error) {
	srcFile, err := os.Open(in.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("open source %q: %w", in.Source, err)
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, nil, fmt.Errorf("decode image %q: %w", in.Source, err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, in.Width, in.Height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	outFile, err := os.Create(in.Destination)
	if err != nil {
		return nil, nil, fmt.Errorf("create output %q: %w", in.Destination, err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, dst); err != nil {
		return nil, nil, fmt.Errorf("encode png: %w", err)
	}

	info, err := outFile.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("stat output: %w", err)
	}

	result := ResizeImageResult{
		Width:    in.Width,
		Height:   in.Height,
		FileSize: info.Size(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}

// ImageInfoInput holds the parameters for the image_info tool.
type ImageInfoInput struct {
	Path string `json:"path" jsonschema:"Absolute path to the image file"`
}

// ImageInfoResult reports the dimensions and file size of an image.
type ImageInfoResult struct {
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	FileSize int64 `json:"fileSize"`
}

// handleImageInfo returns the dimensions and file size of an image
// without fully decoding it.
func handleImageInfo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ImageInfoInput,
) (*mcp.CallToolResult, any, error) {
	f, err := os.Open(in.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %q: %w", in.Path, err)
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, nil, fmt.Errorf("decode config %q: %w", in.Path, err)
	}

	info, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("stat %q: %w", in.Path, err)
	}

	result := ImageInfoResult{
		Width:    cfg.Width,
		Height:   cfg.Height,
		FileSize: info.Size(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
