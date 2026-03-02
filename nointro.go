package main

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed nointro.json.gz
var nointroDB []byte

// nointroEntry represents a single ROM in the embedded No-Intro
// database. The Name is the canonical No-Intro title (without file
// extension), and Extension is the ROM extension without the leading
// dot (e.g. "gb", "nes", "a26").
type nointroEntry struct {
	Name      string
	Extension string
}

var (
	nointroMap  map[string]*nointroEntry
	nointroOnce sync.Once
	nointroErr  error
)

// consoleToLibretro maps ROM file extensions to libretro system names
// used in thumbnail download URLs. This is the single source of truth
// for extension-to-system mapping in the binary.
var consoleToLibretro = map[string]string{
	"gb":   "Nintendo - Game Boy",
	"gbc":  "Nintendo - Game Boy Color",
	"gba":  "Nintendo - Game Boy Advance",
	"nes":  "Nintendo - Nintendo Entertainment System",
	"sfc":  "Nintendo - Super Nintendo Entertainment System",
	"sms":  "Sega - Master System - Mark III",
	"gg":   "Sega - Game Gear",
	"gen":  "Sega - Mega Drive - Genesis",
	"pce":  "NEC - PC Engine - TurboGrafx 16",
	"a26":  "Atari - 2600",
	"a52":  "Atari - 5200",
	"a78":  "Atari - 7800",
	"col":  "Coleco - ColecoVision",
	"int":  "Mattel - Intellivision",
	"ngp":  "SNK - Neo Geo Pocket",
	"ws":   "Bandai - WonderSwan",
	"mini": "Nintendo - Pokemon Mini",
}

// loadNoIntro decompresses and parses the embedded nointro.json.gz
// database on first call. Subsequent calls are no-ops.
func loadNoIntro() {
	nointroOnce.Do(func() {
		r, err := gzip.NewReader(bytes.NewReader(nointroDB))
		if err != nil {
			nointroErr = fmt.Errorf("decompress nointro db: %w", err)
			return
		}
		defer r.Close()

		// The JSON encodes map[string][2]string: sha1 -> [ext, name].
		var raw map[string][2]string
		if err := json.NewDecoder(r).Decode(&raw); err != nil {
			nointroErr = fmt.Errorf("parse nointro db: %w", err)
			return
		}

		nointroMap = make(map[string]*nointroEntry, len(raw))
		for sha1, pair := range raw {
			nointroMap[sha1] = &nointroEntry{
				Extension: pair[0],
				Name:      pair[1],
			}
		}
	})
}

// lookupNoIntro returns the No-Intro entry for the given SHA1 hex
// string, or nil if the hash is not in the database.
func lookupNoIntro(sha1hex string) (*nointroEntry, error) {
	loadNoIntro()
	if nointroErr != nil {
		return nil, nointroErr
	}
	return nointroMap[sha1hex], nil
}

// LookupNoIntroInput holds the parameters for the lookup_nointro tool.
type LookupNoIntroInput struct {
	SHA1 string `json:"sha1" jsonschema:"SHA1 hash of the ROM file (40-character hex string)"`
}

// LookupNoIntroResult contains the result of a No-Intro database
// lookup. Found is true when the SHA1 matched an entry.
type LookupNoIntroResult struct {
	Found     bool   `json:"found"`
	Name      string `json:"name,omitempty"`
	Extension string `json:"extension,omitempty"`
	System    string `json:"system,omitempty"`
}

// handleLookupNoIntro looks up a ROM's No-Intro canonical name and
// console system from its SHA1 hash. Returns found=false for hashes
// not in the database (homebrew, hacks, bad dumps).
func handleLookupNoIntro(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in LookupNoIntroInput,
) (*mcp.CallToolResult, any, error) {
	entry, err := lookupNoIntro(in.SHA1)
	if err != nil {
		return nil, nil, err
	}

	result := LookupNoIntroResult{Found: entry != nil}
	if entry != nil {
		result.Name = entry.Name
		result.Extension = entry.Extension
		result.System = consoleToLibretro[entry.Extension]
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
