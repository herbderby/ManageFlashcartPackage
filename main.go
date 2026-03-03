// Package main implements an MCP server that provides low-level
// tools for managing Nintendo DS flashcart SD cards. The tools are
// primitives (filesystem, network, image, archive); all domain
// knowledge lives in the embedded flashcart_knowledge prompt.
package main

import (
	"context"
	_ "image/png"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newServer constructs and configures the MCP server with all tools
// and prompts registered. It is called by main and by tests.
func newServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "flashcart-tools",
		Version: "0.3.0",
	}, nil)

	// Prompts: identification, domain knowledge, and step-by-step
	// workflow recipes.
	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_identify",
		Description: "Identify a flashcart model from photographs of the cartridge",
	}, handleFlashcartIdentify)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_knowledge",
		Description: "Domain knowledge for managing DS flashcart SD cards",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "flashcart_model",
				Description: "Model ID from the registry (e.g. ace3ds_plus, r4ils)",
				Required:    true,
			},
		},
	}, handleFlashcartKnowledge)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_init",
		Description: "Step-by-step Wood R4 kernel installation",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "flashcart_model",
				Description: "Model ID from the registry (e.g. ace3ds_plus, r4ils)",
				Required:    true,
			},
		},
	}, handleFlashcartInit)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_twilight_install",
		Description: "Step-by-step TWiLight Menu++ installation",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "flashcart_model",
				Description: "Model ID from the registry (e.g. ace3ds_plus, r4ils)",
				Required:    true,
			},
		},
	}, handleFlashcartTwilight)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_emulators",
		Description: "Step-by-step Virtual Console emulator add-on installation",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "flashcart_model",
				Description: "Model ID from the registry (e.g. ace3ds_plus, r4ils)",
				Required:    true,
			},
		},
	}, handleFlashcartEmulators)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_boxart",
		Description: "Scan all ROMs and download box art for each",
	}, handleFlashcartBoxart)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_add_game",
		Description: "Add a single ROM to the flashcart with box art",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "source_path",
				Description: "Absolute path to the ROM file to add",
				Required:    true,
			},
		},
	}, handleFlashcartAddGame)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_cleanup",
		Description: "Clean AppleDouble files and verify card structure",
	}, handleFlashcartCleanup)

	server.AddPrompt(&mcp.Prompt{
		Name:        "flashcart_manual",
		Description: "User manual: getting started, workflows, and troubleshooting",
	}, handleFlashcartManual)

	// Identification tool -- Chat sees "flashcart_identify" in
	// the tool list and calls it when asked to identify a cart.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "flashcart_identify",
		Description: "Return the flashcart visual identification guide. Call this before identifying a DS flashcart from photos. Covers PCB color, shell indent patterns, label analysis, and kernel selection.",
	}, handleFlashcartIdentifyTool)

	// Help tool -- always visible in the tool list so Chat finds
	// it reliably when someone asks for help or a manual.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "flashcart_help",
		Description: "Show the Flashcart Tools user manual: getting started, available workflows, troubleshooting.",
	}, handleFlashcartHelp)

	// Volume tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_volumes",
		Description: "List mounted volumes with filesystem type, total size, and free space.",
	}, handleListVolumes)

	// Filesystem tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_directory",
		Description: "List files and subdirectories at a path with names, sizes, and modification times.",
	}, handleListDirectory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_directory",
		Description: "Create a directory and any necessary parent directories.",
	}, handleCreateDirectory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "move_file",
		Description: "Move or rename a file or directory.",
	}, handleMoveFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "copy_file",
		Description: "Copy a file or directory from one path to another. Set recursive to true for directories.",
	}, handleCopyFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_file",
		Description: "Delete a single file (not recursive).",
	}, handleDeleteFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "file_exists",
		Description: "Check whether a path exists and whether it is a file or directory.",
	}, handleFileExists)

	// Byte reading tool.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_bytes",
		Description: "Read N bytes at a byte offset from a file. Returns hex and ASCII representations.",
	}, handleReadBytes)

	// Hash tool.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "compute_sha1",
		Description: "Compute the SHA1 hash of a file.",
	}, handleComputeSHA1)

	// No-Intro database lookup tool.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lookup_nointro",
		Description: "Look up a ROM's No-Intro canonical name and console system from its SHA1 hash.",
	}, handleLookupNoIntro)

	// Network tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "download_file",
		Description: "Download a URL to a local file path.",
	}, handleDownloadFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch_url",
		Description: "Fetch a URL and return the response body as text (truncated at 1 MiB).",
	}, handleFetchURL)

	// Archive tool.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "extract_archive",
		Description: "Extract a .7z or .zip archive to a directory.",
	}, handleExtractArchive)

	// Image tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "resize_image",
		Description: "Resize a PNG image to given dimensions using high-quality interpolation.",
	}, handleResizeImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "image_info",
		Description: "Return the dimensions and file size of an image.",
	}, handleImageInfo)

	// Text file tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_file",
		Description: "Read a text file and return its contents. For INI files, configs, and other text.",
	}, handleReadFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "write_file",
		Description: "Write text content to a file, creating parent directories as needed.",
	}, handleWriteFile)

	// JSON tools.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_json",
		Description: "Read and parse a JSON file.",
	}, handleReadJSON)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "write_json",
		Description: "Write a JSON object to a file with indentation.",
	}, handleWriteJSON)

	// macOS volume hygiene.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "clean_dot_files",
		Description: "Remove AppleDouble ._* resource fork files from a directory tree. Run this after writing files to a FAT32 volume on macOS.",
	}, handleCleanDotFiles)

	return server
}

func main() {
	server := newServer()
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
