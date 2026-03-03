# ManageFlashcartPackage

## Overview

An MCPB desktop extension that bundles a compiled Go MCP server
binary for managing Nintendo DS flashcart SD cards. Claude uses
the MCP tools plus embedded domain knowledge to prepare and
maintain cards running TWiLight Menu++ or Wood R4.

## Language and Build

- **Language:** Go
- **Module:** `github.com/herbderby/ManageFlashcartPackage`
- **Build:** `make build` (strips debug info, cross-compiles)
- **Test:** `go test ./...` (in-process MCP API tests)
- **Pack:** `make pack` (creates `flashcart-tools.mcpb`)
- **MCP SDK:** `github.com/modelcontextprotocol/go-sdk` v1.4.0

## Project Structure

```
main.go                 # MCP server entry point, tool registration
volumes.go              # list_volumes tool
filesystem.go           # list_directory, create_directory, move_file,
                        #   copy_file (recursive), delete_file, file_exists
pathutil.go             # resolvedMkdirAll (firmlink-safe mkdir)
bytes.go                # read_bytes tool
hash.go                 # compute_sha1 tool
nointro.go              # lookup_nointro tool, embedded No-Intro database
nointro.json.gz         # Embedded DB: SHA1 -> No-Intro name (~21K entries)
textfile.go             # read_file, write_file tools
network.go              # download_file, fetch_url tools
archive.go              # extract_archive tool (7z + zip)
image.go                # resize_image, image_info tools
json_tools.go           # read_json, write_json tools
dotclean.go             # clean_dot_files tool (AppleDouble removal)
models.go               # Flashcart model registry and lookup
skill.go                # All 9 prompts (identify, knowledge, workflows, manual)
mcp_test.go             # In-process MCP API tests
Makefile                # Build and pack targets
tools/
  gen_nointro.go        # Generator: merges NoIntro.db + Myrient DATs
ext/
  manifest.json         # MCPB manifest (dxt_version 0.1)
  server/
    flashcart-tools     # Compiled binary (darwin/arm64)
PRD.md                  # Product requirements document
research/               # Background research
.mcp.json               # Go language server MCP for this workspace
```

## Key References

- `PRD.md` -- product requirements, tool surface, architecture
- `research/flashcart-setup-report.md` -- field test report (authoritative)
- `research/nds-rom-header-format.md` -- ROM header parsing
- `research/twilight-menu-plus-plus.md` -- firmware details
- `research/gametdb-box-art.md` -- box art sourcing
- `research/r4-flashcart-setup.md` -- R4 flashcart setup

## Go Dependencies

- `github.com/modelcontextprotocol/go-sdk` -- MCP server SDK
- `github.com/bodgit/sevenzip` -- pure Go 7z extraction
- `golang.org/x/image` -- image resizing (CatmullRom)

## MCP SDK Notes

- `jsonschema` struct tags are the description directly, not
  `description=...` prefixed
- StdioTransport uses newline-delimited JSON (not
  Content-Length framed)
- Prompt handler signature:
  `func(ctx, *GetPromptRequest) (*GetPromptResult, error)`
- Tool handler signature:
  `func(ctx, *CallToolRequest, In) (*CallToolResult, Out, error)`
- No-input tools use `struct{}` as `In` type

## Development Phase

Phase 0 (complete): `list_volumes` PoC validated via CLI and
MCPB install in Claude Desktop.

Phase 1 (complete): All 18 tools + `flashcart_knowledge` prompt
implemented. Binary compiles, `.mcpb` packs at 3.4 MB.

Phase 1.1 (complete): FAT32 volume fixes (firmlink-safe mkdir,
recursive copy, AppleDouble cleanup), text file tools, Wood R4
domain knowledge added.

Phase 2 (complete): Field test with real Ace3DS+ SD card via
Chat. PRD and embedded prompt rewritten with field-tested
knowledge. See `research/flashcart-setup-report.md`.

Phase 2.1 (complete): Flashcart model registry. Hardcoded
Ace3DS+ references replaced with parameterized prompts.
`flashcart_identify` prompt added for photo-based cart ID.
Models: `models.go`, 18 tools + 9 prompts.

Phase 2.2 (complete): Non-NDS box art SHA1 lookup. Embedded
No-Intro database (~21K entries) for hash-based ROM
identification. `compute_sha1` and `lookup_nointro` tools.
Box art prompts rewritten to hash first, filename fallback.
22 tools + 9 prompts.

Next: improve flashcart identification from photos
(see NEXT.md for details).
