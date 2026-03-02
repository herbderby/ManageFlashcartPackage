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
textfile.go             # read_file, write_file tools
network.go              # download_file, fetch_url tools
archive.go              # extract_archive tool (7z + zip)
image.go                # resize_image, image_info tools
json_tools.go           # read_json, write_json tools
dotclean.go             # clean_dot_files tool (AppleDouble removal)
skill.go                # flashcart_knowledge prompt (embedded)
Makefile                # Build and pack targets
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

Next: commit, reinstall .mcpb, test three-phase setup in Claude
Desktop, then decide on skill sub-task implementation.
