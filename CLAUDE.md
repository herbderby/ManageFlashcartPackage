# ManageFlashcartPackage

## Overview

A Cowork plugin that bundles a compiled Go MCP server binary for
managing Nintendo DS flashcart SD cards. Claude uses the MCP tools
plus domain knowledge from plugin skill files to prepare and
maintain cards running TWiLight Menu++.

## Language and Build

- **Language:** Go
- **Module:** `github.com/herbderby/ManageFlashcartPackage`
- **Build:** `go build -ldflags="-s -w"` for release binaries
- **Cross-compile:** `GOOS=darwin GOARCH=arm64`, etc.
- **MCP SDK:** Official Go SDK (`github.com/modelcontextprotocol/go-sdk`)
  recommended per research; alternatives documented in
  `research/go-mcp-sdk-options.md`

## Project Structure

```
PRD.md                  # Full product requirements document
research/               # Background research (libraries, protocols, ROM format)
.mcp.json               # Go language server MCP for this workspace
```

## Key References

- `PRD.md` -- complete product requirements, Phase 0 PoC spec,
  MCP tool surface, plugin architecture
- `research/go-libraries-and-build.md` -- 7z, image resizing,
  cross-compilation
- `research/go-mcp-sdk-options.md` -- SDK comparison and examples
- `research/nds-rom-header-format.md` -- ROM header parsing
- `research/twilight-menu-plus-plus.md` -- firmware details
- `research/gametdb-box-art.md` -- box art sourcing
- `research/r4-flashcart-setup.md` -- R4 flashcart setup
- `research/cowork-plugin-architecture.md` -- plugin packaging

## Go Dependencies

- `github.com/modelcontextprotocol/go-sdk` -- MCP server SDK
- `github.com/bodgit/sevenzip` -- pure Go 7z extraction
- `golang.org/x/image` -- image resizing (CatmullRom or
  BiLinear for quality downscaling)

## Development Phase

Phase 0: Proof of Concept -- build `HelloSDCard` plugin with a
single `list_volumes` MCP tool to validate the end-to-end chain
(plugin install, binary discovery, MCP launch, host filesystem
access).
