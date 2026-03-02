package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// connect creates an in-process server/client pair for testing.
// The returned ClientSession is ready for ListPrompts, GetPrompt,
// and CallTool calls. The caller should defer session.Close().
func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()

	server := newServer()
	t1, t2 := mcp.NewInMemoryTransports()

	_, err := server.Connect(ctx, t1, nil)
	if err != nil {
		t.Fatalf("server.Connect: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "0.1.0",
	}, nil)

	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	return session
}

// TestPromptsList verifies that all expected prompts are registered.
func TestPromptsList(t *testing.T) {
	session := connect(t)
	defer session.Close()

	result, err := session.ListPrompts(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListPrompts: %v", err)
	}

	want := map[string]bool{
		"flashcart_knowledge":        false,
		"flashcart_init":             false,
		"flashcart_twilight_install": false,
		"flashcart_emulators":        false,
		"flashcart_boxart":           false,
		"flashcart_add_game":         false,
		"flashcart_cleanup":          false,
		"flashcart_manual":           false,
	}

	for _, p := range result.Prompts {
		if _, ok := want[p.Name]; ok {
			want[p.Name] = true
		}
	}

	for name, found := range want {
		if !found {
			t.Errorf("prompt %q not found in prompts/list", name)
		}
	}

	if len(result.Prompts) != len(want) {
		t.Errorf("ListPrompts returned %d prompts, want %d",
			len(result.Prompts), len(want))
	}
}

// TestPromptsGet verifies that each prompt returns non-empty content.
func TestPromptsGet(t *testing.T) {
	session := connect(t)
	defer session.Close()

	prompts := []struct {
		name    string
		args    map[string]string
		wantSub string // substring expected in the response text
	}{
		{"flashcart_knowledge", nil, "Flashcart SD Card Management"},
		{"flashcart_init", nil, "Wood R4 Kernel Installation"},
		{"flashcart_twilight_install", nil, "TWiLight Menu++ Installation"},
		{"flashcart_emulators", nil, "Virtual Console Emulators"},
		{"flashcart_boxart", nil, "Box Art Download"},
		{"flashcart_add_game", map[string]string{"source_path": "/tmp/test.nds"}, "Add a ROM"},
		{"flashcart_cleanup", nil, "Flashcart Volume Cleanup"},
		{"flashcart_manual", nil, "User Manual"},
	}

	for _, tt := range prompts {
		t.Run(tt.name, func(t *testing.T) {
			result, err := session.GetPrompt(context.Background(), &mcp.GetPromptParams{
				Name:      tt.name,
				Arguments: tt.args,
			})
			if err != nil {
				t.Fatalf("GetPrompt(%q): %v", tt.name, err)
			}
			if len(result.Messages) == 0 {
				t.Fatalf("GetPrompt(%q) returned 0 messages", tt.name)
			}
			text, ok := result.Messages[0].Content.(*mcp.TextContent)
			if !ok {
				t.Fatalf("GetPrompt(%q) message is not TextContent", tt.name)
			}
			if !strings.Contains(text.Text, tt.wantSub) {
				t.Errorf("GetPrompt(%q) text does not contain %q\ngot: %.100s...",
					tt.name, tt.wantSub, text.Text)
			}
		})
	}
}

// TestAddGameSubstitution verifies that the source_path argument is
// substituted into the flashcart_add_game prompt.
func TestAddGameSubstitution(t *testing.T) {
	session := connect(t)
	defer session.Close()

	path := "/Users/herb/Games/MarioKart.nds"
	result, err := session.GetPrompt(context.Background(), &mcp.GetPromptParams{
		Name:      "flashcart_add_game",
		Arguments: map[string]string{"source_path": path},
	})
	if err != nil {
		t.Fatalf("GetPrompt: %v", err)
	}

	text := result.Messages[0].Content.(*mcp.TextContent).Text
	if !strings.Contains(text, path) {
		t.Errorf("source_path %q not found in response", path)
	}
	if strings.Contains(text, "{{source_path}}") {
		t.Error("unreplaced {{source_path}} placeholder in response")
	}
}

// TestToolsList verifies that all 18 tools are registered.
func TestToolsList(t *testing.T) {
	session := connect(t)
	defer session.Close()

	result, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}

	want := []string{
		"list_volumes", "list_directory", "create_directory",
		"move_file", "copy_file", "delete_file", "file_exists",
		"read_bytes", "download_file", "fetch_url",
		"extract_archive", "resize_image", "image_info",
		"read_file", "write_file", "read_json", "write_json",
		"clean_dot_files",
	}

	got := make(map[string]bool)
	for _, tool := range result.Tools {
		got[tool.Name] = true
	}

	for _, name := range want {
		if !got[name] {
			t.Errorf("tool %q not found in tools/list", name)
		}
	}

	if len(result.Tools) != len(want) {
		t.Errorf("ListTools returned %d tools, want %d",
			len(result.Tools), len(want))
	}
}
