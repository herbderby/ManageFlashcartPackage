# ManageFlashcartPackage

## Overview

An MCPB desktop extension that bundles a compiled Go MCP server
binary for managing Nintendo DS flashcart SD cards. Claude uses
the MCP tools plus embedded domain knowledge to prepare and
maintain cards running TWiLight Menu++ or Wood R4.

## Language and Build

- **Language:** Go
- **Module:** `github.com/herbderby/ManageFlashcartPackage`
- **Build:** `make build` (strips debug info, builds darwin/arm64 + linux/arm64)
- **Test:** `go test ./...` (in-process MCP API tests)
- **Pack:** `make pack` (creates `flashcart-tools.mcpb`)
- **MCP SDK:** `github.com/modelcontextprotocol/go-sdk` v1.4.0

## Project Structure

```
main.go                 # MCP server entry point, tool registration
volumes.go              # list_volumes handler + shared types
volumes_darwin.go       # macOS volume enumeration (syscall.Statfs_t)
volumes_linux.go        # Linux stub (returns empty list)
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
skill.go                # 9 prompts + 3 "docs as tools" (read_me_first,
                        #   flashcart_identify, flashcart_help)
mcp_test.go             # In-process MCP API tests
Makefile                # Build and pack targets
tools/
  gen_nointro.go        # Generator: merges NoIntro.db + Myrient DATs
ext/
  manifest.json         # MCPB manifest (dxt_version 0.1)
  server/
    flashcart-tools               # Shell launcher (dispatches by OS/arch)
    flashcart-tools-darwin-arm64  # macOS Apple Silicon binary
    flashcart-tools-linux-arm64   # Linux arm64 binary (Cowork VM)
PRD.md                  # Product requirements document
research/               # Background research
.mcp.json               # Go language server MCP for this workspace
```

## Key References

- `PRD.md` -- product requirements, tool surface, architecture
- `research/flashcart-setup-report.md` -- field test report (authoritative)
- `research/flashcart_identification_guide.md` -- visual ID guide (fetched at runtime)
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

Phase 2.3 (complete): Chat discoverability. Chat cannot find
MCP prompts, only tools. Added three "docs as tools":
`read_me_first` (routing guide with workflows), `flashcart_identify`
(visual ID guide fetched from GitHub), `flashcart_help` (user
manual). Restructured identification guide to prioritize PCB
color and shell indents over label text.
23 tools + 9 prompts.

Next: field test read_me_first in Chat, then evaluate whether
prompts are still needed (see NEXT.md for details).

## Chat Discoverability

Chat (Claude Desktop) cannot reliably find MCP prompts. It
sees tools in the tool list but not prompts. Pattern: create a
zero-parameter tool that returns the prompt content as text.
The tool name triggers reflexive invocation (e.g., Chat sees
`flashcart_identify` and calls it when asked to identify a
cart). Three instances: `read_me_first`, `flashcart_identify`,
`flashcart_help`.

## Dual-Binary Build (Cowork Support)

Cowork runs MCP servers inside a Linux/aarch64 VM, not
natively on macOS. A macOS-only binary will silently fail to
start, causing all tools to disappear from the Cowork session.

### Requirements for Cowork compatibility

1. **Two binaries.** The Linux/arm64 binary runs inside the
   Cowork VM. The Darwin/arm64 binary runs natively for Chat
   and Code. Both are cross-compiled with `CGO_ENABLED=0`.

2. **Shell launcher.** `ext/server/flashcart-tools` is a
   POSIX shell script that dispatches by `uname -s` /
   `uname -m` to the correct platform binary. It must NOT
   be a compiled binary itself.

3. **Platform declaration.** `manifest.json` must include
   `"platforms": ["darwin", "linux"]` in the compatibility
   section. If only `"darwin"` is listed, Cowork will see
   the plugin in the connectors menu but refuse to load it.

4. **Platform-specific Go code.** Any code using
   Darwin-only syscalls (e.g., `Statfs_t.Fstypename`) must
   be guarded with build tags (`volumes_darwin.go` /
   `volumes_linux.go`). The Linux build will fail otherwise.

5. **Version bump on every `.mcpb` rebuild.** Cowork caches
   plugins by version. If the version is unchanged, Cowork
   may serve stale binaries. Always bump the version in both
   `ext/manifest.json` and `main.go` before `make pack`.

### Binary layout

- `flashcart-tools` -- shell launcher (dispatches by OS/arch)
- `flashcart-tools-darwin-arm64` -- Chat, Code (macOS Apple Silicon)
- `flashcart-tools-linux-arm64` -- Cowork (Linux VM)

`make build` compiles both. `make pack` includes the launcher
and both binaries in the `.mcpb`. The launcher also handles
`darwin/x86_64` and `linux/x86_64` cases for future use.
