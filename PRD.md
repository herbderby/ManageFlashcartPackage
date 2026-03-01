# ManageFlashcartPackage -- Product Requirements Document

## Problem

Nintendo DS and DS Lite owners use R4 flashcarts with micro SD cards to play
game libraries. A properly set up SD card requires: firmware/kernel files in
specific directories, ROM files organized by system, box art images sized and
named to match each ROM, cheat databases, and emulator cores. Getting all of
this right is technically challenging for most users -- it involves downloading
files from multiple sources, understanding directory conventions, parsing ROM
headers to find title IDs, and resizing images.

## Product Vision

A single installable Cowork plugin that gives Claude the tools and knowledge
to prepare and maintain a flashcart SD card on the user's behalf. The user
installs the plugin, plugs in their SD card, opens a Cowork session, and
Claude handles the rest: installing firmware, organizing ROMs, fetching box
art, and validating the card layout.

## Target Firmware: TWiLight Menu++

TWiLight Menu++ is the actively maintained, open-source menu system for DS
flashcarts. It supports NDS games natively, includes emulators for GBA, SNES,
GB, GBC, and NES, and displays box art in its menu. Releases are published on
GitHub at [DS-Homebrew/TWiLightMenu](https://github.com/DS-Homebrew/TWiLightMenu/releases).

### SD Card Directory Structure

```
/
+-- BOOT.NDS                              # TWiLight Menu++ bootloader
+-- _nds/
|   +-- TWiLightMenu/
|       +-- boxart/                       # Box art PNGs
|       +-- gamesettings/                 # Per-game settings
|       +-- extras/
|           +-- usrcheat.dat              # Cheat database
+-- roms/
|   +-- nds/                              # NDS ROMs (.nds)
|   +-- gba/                              # GBA ROMs (.gba)
|   +-- snes/                             # SNES ROMs (.smc, .sfc)
|   +-- gb/                               # Game Boy ROMs (.gb, .gbc)
|   +-- nes/                              # NES ROMs (.nes)
+-- saves/                                # Save files (.sav)
```

### Box Art Specifications

- Format: PNG
- Naming: `{GAMECODE}.png` (preferred) or `{romfilename}.nds.png`
- Recommended size: 128x115, maximum: 208x143
- Cached mode limit: 44 KiB per image
- Location: `_nds/TWiLightMenu/boxart/`
- Source: GameTDB at `https://art.gametdb.com/ds/coverS/{REGION}/{GAMECODE}.png`

### NDS ROM Header (identification fields)

| Offset | Size     | Field     | Notes                              |
|--------|----------|-----------|------------------------------------|
| 0x000  | 12 bytes | Game Title| Uppercase ASCII, null-padded       |
| 0x00C  | 4 bytes  | Gamecode  | Uppercase ASCII, used as title ID  |
| 0x010  | 2 bytes  | Makercode | Uppercase ASCII publisher ID       |
| 0x012  | 1 byte   | Unitcode  | 00h=NDS, 02h=NDS+DSi, 03h=DSi     |

The 4-byte gamecode at offset 0x00C is the title ID used for box art lookups
and game identification throughout the system.

## Precondition: Persistent Working Directory

The user must maintain a folder on their computer that persists between SD
card maintenance sessions. This is the folder the user opens in Cowork. It
is NOT on the SD card -- it lives on the computer's internal storage and
survives card reformats, card swaps, and card loss.

### Purpose

The working directory serves as the user's flashcart home base:

- **Save backups.** Copies of `.sav` files pulled from the SD card before
  making changes. If a card dies or gets corrupted, saves are not lost.
- **ROM staging.** A place to drop `.nds` files before the tool organizes
  them onto the card. ROMs can also be downloaded here from repositories.
- **Box art cache.** Downloaded box art PNGs, so re-setting up a card
  (or setting up a second card) does not re-download from GameTDB.
- **Configuration.** ROM repository URLs, preferred GameTDB region,
  per-game settings, cheat preferences -- anything the user has
  customized across sessions.
- **Session history.** A record of what was done to which card and when,
  so Claude can pick up where the last session left off.

### Directory Structure

```
~/DS_Flashcart/                          # User chooses the name and location
+-- backups/
|   +-- saves/                           # .sav files backed up from SD cards
|   +-- cards/                           # Snapshots of card state before changes
+-- staging/                             # ROMs waiting to be organized onto a card
+-- cache/
|   +-- boxart/                          # Downloaded box art PNGs (by gamecode)
|   +-- twilight_menu/                   # Cached TWiLight Menu++ releases
+-- config.json                          # Repository URLs, region, preferences
```

### Workflow Implications

- **First session:** If the working directory is empty or missing expected
  structure, the plugin skill guides the user through initial setup --
  creating the folder structure and setting preferences (e.g., GameTDB
  region).
- **Card operations:** Before modifying an SD card, the tool backs up saves
  to `backups/saves/`. Before installing or updating firmware, it snapshots
  the card's `_nds/` state to `backups/cards/`.
- **Card swaps:** The user can maintain multiple SD cards from the same
  working directory. The box art cache and ROM staging area are shared
  across all cards.
- **Cowork sessions:** The user opens this folder in Cowork each time.
  Claude reads `config.json` and session history to resume context.

## Delivery: Cowork Plugin with Bundled MCP Server

### Why a Plugin

Claude Cowork supports plugins -- self-contained packages that bundle skills
(domain knowledge), slash commands (user workflows), and MCP server
configurations (tools). A plugin is the natural delivery vehicle because it
combines all three layers into a single install.

Cowork's sandboxed VM cannot access the SD card directly (the VM cannot reach
external volumes). However, MCP servers configured by a plugin run on the
**host machine**, not inside the VM. An MCP server running on the host has
full access to mounted SD cards.

### Plugin Architecture

The plugin bundles a compiled Go binary as the MCP server, referenced via
`${CLAUDE_PLUGIN_ROOT}`:

```
ManageFlashcartPlugin/
  .claude-plugin/
    plugin.json                          # Plugin manifest
  .mcp.json                             # MCP server configuration
  bin/
    manage-flashcart                     # Launcher script (detects OS/arch)
    manage-flashcart.cmd                 # Windows launcher
    manage-flashcart-darwin-arm64        # macOS Apple Silicon binary
    manage-flashcart-darwin-amd64        # macOS Intel binary
    manage-flashcart-windows-amd64.exe   # Windows binary
  skills/
    flashcart-management/
      SKILL.md                           # Domain knowledge and workflow guidance
  commands/
    setup-card/
      SKILL.md                           # /setup-card slash command
    scan-card/
      SKILL.md                           # /scan-card slash command
```

### .mcp.json

```json
{
  "mcpServers": {
    "flashcart": {
      "command": "${CLAUDE_PLUGIN_ROOT}/bin/manage-flashcart",
      "args": ["--serve"],
      "cwd": "${CLAUDE_PLUGIN_ROOT}"
    }
  }
}
```

On macOS, the `manage-flashcart` launcher script detects the architecture
via `uname -m` and execs the correct binary. On Windows, the `.cmd` launcher
runs the `.exe`.

### Installation

User installs the plugin from the Cowork plugin directory (one click) or
uploads it manually. The plugin includes all binaries -- no separate
downloads, no runtime dependencies, no restarts required beyond what plugin
installation already entails.

### Size Budget

Go binaries are typically 5-15 MB each. Three platform binaries total
roughly 20-40 MB. The Cowork plugin upload limit is 50 MB. This is tight
but feasible, especially with `go build -ldflags="-s -w"` to strip debug
info (reduces binary size by ~30%).

## Design Philosophy: Low-Level Tools, Smart Agent

The MCP server provides **primitive operations** -- filesystem access,
HTTP downloads, binary reads, and image manipulation. It does not encode
flashcart-specific workflows. All domain knowledge (directory conventions,
ROM identification, box art naming, firmware installation steps) lives in
the **plugin skills**, which are markdown files.

This split has three advantages:

1. **The Go binary is simple.** Each tool does one thing. Easy to test,
   easy to audit, easy to maintain.
2. **Workflows are editable.** Updating how ROMs are organized or adding
   support for a new firmware variant means editing a markdown skill file,
   not recompiling and redistributing a binary.
3. **Claude can adapt.** When something unexpected happens (unusual card
   layout, download failure, unrecognized file), Claude reasons about it
   using domain knowledge rather than hitting a hardcoded error path.

## MCP Tool Surface (Primitives)

The server exposes low-level tools. Claude composes them using domain
knowledge from the skill files.

**`list_volumes`:**
List mounted removable volumes with mount point, filesystem type, total
size, and free space.

**`list_directory`:**
List files and subdirectories at a path with sizes and modification times.

**`create_directory`:**
Create a directory (and parents) at a path.

**`move_file`:**
Move or rename a file or directory.

**`copy_file`:**
Copy a file from one path to another.

**`delete_file`:**
Delete a single file (not recursive).

**`file_exists`:**
Check whether a path exists and whether it is a file or directory.

**`read_bytes`:**
Read N bytes at a byte offset from a file. Returns hex and ASCII. Claude
uses this to parse ROM headers.

**`download_file`:**
Download a URL to a local file path. Reports progress for large files.
Supports resume.

**`fetch_url`:**
Fetch a URL and return the body as text. For HTML directory listings, API
responses, etc.

**`extract_archive`:**
Extract a .7z or .zip archive to a directory.

**`resize_image`:**
Resize a PNG to given dimensions and save. Reports output file size (for
the 44 KiB cache limit).

**`image_info`:**
Return dimensions and file size of an image.

**`read_json`:**
Read and parse a JSON file.

**`write_json`:**
Write a JSON object to a file.

Each tool returns structured JSON. Each tool does one thing.

## Plugin Skills (Domain Knowledge)

The skills teach Claude everything it needs to compose the primitives into
flashcart workflows. The Go binary knows nothing about flashcarts.

### Skill: flashcart-management (SKILL.md)

This is the core skill, always active when the plugin is enabled. It
teaches Claude:

**TWiLight Menu++ directory conventions:**
- Exact directory tree expected on the card (the structure from the
  "SD Card Directory Structure" section above)
- Which files are system files vs. user files
- Where ROMs go by system type (nds, gba, snes, gb, nes)
- Where saves, box art, cheat databases, and settings live

**ROM identification:**
- How to use `read_bytes` at offset 0x00C (4 bytes) to extract the
  gamecode from an NDS ROM header
- How to use `read_bytes` at offset 0x000 (12 bytes) to extract the
  game title
- The gamecode format: 4 uppercase ASCII characters
- Homebrew ROMs have gamecode `####` or all zeros

**Box art workflow:**
- Gamecode maps to box art URL: `https://art.gametdb.com/ds/coverS/{REGION}/{GAMECODE}.png`
- Default region is US; fall back to EN, then JA
- Download to cache first (`working_dir/cache/boxart/`), then copy to card
- Check dimensions with `image_info`; resize with `resize_image` if needed
- Naming convention: `{GAMECODE}.png` in `_nds/TWiLightMenu/boxart/`

**Firmware installation:**
- Latest release URL: `https://github.com/DS-Homebrew/TWiLightMenu/releases/latest`
- Download `TWiLightMenu-Flashcard.7z` to `working_dir/cache/twilight_menu/`
- Extract with `extract_archive`
- Copy `BOOT.NDS` and `_nds/` to card root using `copy_file`
- Preserve user settings: do not overwrite `_nds/TWiLightMenu/settings.ini`
  if it already exists

**Working directory management:**
- On first session, create the working directory structure
- Back up saves before modifying a card: copy `saves/*.sav` from card to
  `working_dir/backups/saves/`
- Cache downloads to avoid re-fetching

**Safety rules:**
- Always confirm before deleting or overwriting files on the card
- Never delete ROMs without explicit user approval
- Back up saves before any card modification
- Scan the card before making changes so the user sees current state first

### Slash Command: /setup-card

Guided workflow that composes the primitives:
1. `list_volumes` to find the card
2. `list_directory` to assess current state
3. Confirm plan with user
4. Download and install TWiLight Menu++ (download_file, extract_archive,
   copy_file, create_directory)
5. Create ROM directories (create_directory)
6. Scan for existing ROMs (list_directory, read_bytes to identify each)
7. Organize ROMs into correct folders (move_file)
8. Fetch box art (download_file, resize_image, copy_file)
9. Report final state

### Slash Command: /scan-card

Quick assessment that composes the primitives:
1. `list_volumes` to find the card
2. `list_directory` on card root and key subdirectories
3. `read_bytes` on each ROM to extract title and gamecode
4. `list_directory` on boxart folder to check coverage
5. `file_exists` for `usrcheat.dat`, `BOOT.NDS`, `settings.ini`
6. Report to user

## Implementation Language and Build

- **Language:** Go
- **Build system:** Go modules
- **Rationale:** Go produces self-contained static binaries for macOS
  (ARM64 + x86_64) and Windows (x86_64) with trivial cross-compilation.
  A single `GOOS=... GOARCH=... go build` per target. No runtime
  dependencies on the user's machine.
- **Module name:** `github.com/herbderby/ManageFlashcartPackage`

### Go Dependencies (vendored, not user-installed)

- MCP protocol: Go MCP SDK or hand-rolled JSON-RPC over stdio
- 7z extraction: `github.com/bodgit/sevenzip` (pure Go, no CGo)
- Image resizing: `golang.org/x/image` or `image/png` from stdlib
- HTTP client: `net/http` from stdlib
- Filesystem: `os` and `path/filepath` from stdlib
- No flashcart-specific logic -- all domain knowledge is in the skills

## Platform Support

| Platform           | Binary                              | SD Card Mount       |
|--------------------|-------------------------------------|---------------------|
| macOS Apple Silicon| `manage-flashcart-darwin-arm64`     | `/Volumes/{NAME}`   |
| macOS Intel        | `manage-flashcart-darwin-amd64`     | `/Volumes/{NAME}`   |
| Windows x86_64     | `manage-flashcart-windows-amd64.exe`| `E:\`, `F:\`, etc.  |

The `detect_cards` tool uses platform-specific logic to find removable
volumes:
- macOS: enumerate `/Volumes/` and check for removable media
- Windows: enumerate drive letters and check for removable drives

## Non-Goals

- **Formatting the SD card.** That is handled by the Nintendo DS Flashcart
  Tool, a separate project (macOS only).
- **Running on the DS itself.** This tool runs on the user's computer.
- **Supporting non-TWiLight Menu++ firmware.** TWiLight Menu++ only for now.
- **ROM dumping or ripping.** The user provides ROMs or points to
  repositories.
- **Linux support in Phase 1.** Can be added later trivially (Go
  cross-compiles to Linux with one command).

## User Experience Walkthrough

```
User installs the ManageFlashcart plugin from the Cowork plugin
directory (one click). Inserts their SD card. Opens a Cowork session.

User: I want to set up my DS flashcart SD card.

Claude: Let me check for SD cards...

I found one removable volume:

  SDCARD (32 GB FAT32, 29.1 GB free)
  Mounted at /Volumes/SDCARD
  Status: Empty -- no TWiLight Menu++ detected.

Would you like me to set this card up from scratch? I'll install
TWiLight Menu++, create the ROM folder structure, and install the
cheat database.

User: Yes.

Claude: Setting up your card...

  [x] Downloaded TWiLight Menu++ v27.22.2 from GitHub
  [x] Extracted and installed BOOT.NDS and _nds/ to card root
  [x] Created roms/nds/, roms/gba/, roms/snes/, roms/gb/, roms/nes/
  [x] Created saves/ directory
  [x] Downloaded and installed usrcheat.dat

Your card is ready for games. You can:
  - Give me a URL to a ROM repository and I'll help you browse it
  - Copy .nds files to the card yourself and I'll organize them
    and fetch box art

User: I copied a bunch of nds files to the card.

Claude: Let me scan the card...

I found 47 .nds files in the card root. I read each ROM's header
and identified them:

  Mario Kart DS (AMCE)
  New Super Mario Bros. (A2DE)
  Pokemon Diamond (ADAE)
  Phoenix Wright: Ace Attorney (ABCE)
  ... and 43 more.

I'll move them all into roms/nds/ and fetch box art for each from
GameTDB. Go ahead?

User: Yes please.

Claude: Done.

  [x] Moved 47 ROMs to roms/nds/
  [x] Fetched box art for 45 of 47 games
  [ ] No box art found for 2 titles:
      - Homebrew_Menu.nds (homebrew, not in GameTDB)
      - ROM_Hack_v2.nds (ROM hack, not in GameTDB)

Your card has 47 games with box art. Eject the card, pop it in
your R4, and try it out.
```

## Phase 0: Proof of Concept

Before building the full flashcart tool, validate the core architecture with
a trivial plugin that bundles a compiled Go MCP server binary. This proves
(or disproves) the end-to-end chain: plugin install, binary discovery,
MCP server launch, tool execution on the host filesystem.

### What to Build

A plugin called `HelloSDCard` with one MCP tool: `list_volumes`. The tool
lists mounted volumes on the host machine and returns their names, mount
points, and free space. No flashcart logic, no ROM parsing -- just proof
that a bundled Go binary can talk to Cowork through the MCP protocol.

### Plugin Structure

```
HelloSDCard/
  .claude-plugin/
    plugin.json
  .mcp.json
  bin/
    hello-sdcard                     # Launcher script (macOS)
    hello-sdcard-darwin-arm64        # Compiled Go binary
  skills/
    hello/
      SKILL.md                       # Minimal: "Use list_volumes to show drives"
```

### plugin.json

```json
{
  "name": "HelloSDCard",
  "version": "0.1.0",
  "description": "Proof of concept: list mounted volumes from a bundled Go MCP server"
}
```

### .mcp.json

```json
{
  "mcpServers": {
    "hello-sdcard": {
      "command": "${CLAUDE_PLUGIN_ROOT}/bin/hello-sdcard",
      "cwd": "${CLAUDE_PLUGIN_ROOT}"
    }
  }
}
```

### Go MCP Server (hello-sdcard)

A single `main.go` (~100 lines) that:
1. Reads JSON-RPC requests from stdin
2. Responds to `tools/list` with one tool: `list_volumes`
3. Responds to `tools/call` for `list_volumes` by enumerating `/Volumes/`
   (macOS) and returning JSON with volume name, mount point, filesystem
   type, and free space
4. Handles the MCP `initialize` handshake

### Launcher Script (bin/hello-sdcard)

```sh
#!/bin/sh
ARCH=$(uname -m)
case "$ARCH" in
  arm64)  exec "${0}-darwin-arm64" "$@" ;;
  x86_64) exec "${0}-darwin-amd64" "$@" ;;
  *)      echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac
```

### Build and Test Steps

1. Write `main.go` with MCP stdio transport and `list_volumes` tool
2. Build: `GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/hello-sdcard-darwin-arm64`
3. Create the plugin directory structure with launcher script and skill
4. Test locally first:
   - Run the binary directly and send MCP JSON-RPC on stdin to verify
     it responds correctly
   - Confirm `list_volumes` returns real volume data
5. Install in Cowork:
   - Upload the plugin directory as a zip, or use `--plugin-dir` for
     development mode
   - Open a Cowork session with the plugin enabled
   - Ask Claude: "What drives are connected to my computer?"
   - Verify Claude calls `list_volumes` and reports the mounted volumes

### What This Validates

| Risk                                           | Validated by                              |
|------------------------------------------------|-------------------------------------------|
| Can a plugin bundle a native binary?           | Plugin installs without error             |
| Are execute permissions preserved?             | Binary launches successfully              |
| Does `${CLAUDE_PLUGIN_ROOT}` resolve correctly?| MCP server starts from the right path     |
| Does the launcher script work?                 | Correct arch binary is selected and exec'd|
| Can MCP tools access the host filesystem?      | `list_volumes` returns real `/Volumes/` data |
| Can Claude call tools from a bundled binary?   | Claude uses `list_volumes` in conversation|
| Is an SD card visible to the MCP server?       | Insert SD card, verify it appears in results |

### What This Does NOT Validate

- Windows support (macOS-only for PoC)
- Plugin size limits (PoC binary will be ~5 MB, well under 50 MB)
- Go MCP SDK maturity (PoC hand-rolls the protocol)
- 7z extraction, image resizing, HTTP downloads

### Success Criteria

The PoC is successful if: a user installs the plugin in Cowork, asks
"what volumes are mounted?", and Claude calls the bundled Go binary's
`list_volumes` tool and correctly reports the mounted drives, including
an inserted SD card.

If this works, every assumption in the full plugin architecture is proven.
If it fails, we learn exactly which part of the chain breaks and can
adjust before investing in the full implementation.

## Open Questions

1. **7z extraction.** TWiLight Menu++ releases are distributed as `.7z`
   archives. The Go library `github.com/bodgit/sevenzip` handles this
   in pure Go (no CGo), but needs validation that it works with the
   specific archives TWiLight Menu++ publishes.

2. **Go MCP SDK.** Is there a mature Go MCP SDK, or should we implement
   the JSON-RPC stdio transport directly? The protocol is simple enough
   to hand-roll.

3. **GameTDB rate limiting.** Does GameTDB have rate limits or terms of
   service for programmatic cover downloads? A card with 500+ ROMs would
   generate many requests.

4. **Non-NDS ROM identification.** GBA, SNES, GB, and NES ROMs have their
   own header formats for title extraction. Should Phase 1 support these,
   or focus on NDS only?

5. **Plugin size limit.** Three Go binaries may approach the 50 MB Cowork
   upload limit. If too large, options include: ship macOS-only initially,
   use UPX compression, or distribute via git-based marketplace (no size
   limit) instead of web upload.

6. **Plugin launcher scripts.** The `manage-flashcart` shell script and
   `.cmd` file need to detect OS/architecture and exec the right binary.
   Need to verify this works when launched by the MCP server spawner in
   Claude Desktop.

## References

- [TWiLight Menu++ Flashcard Installation](https://wiki.ds-homebrew.com/twilightmenu/installing-flashcard)
- [TWiLight Menu++ Box Art Guide](https://wiki.ds-homebrew.com/twilightmenu/how-to-get-box-art)
- [TWiLight Menu++ GitHub Releases](https://github.com/DS-Homebrew/TWiLightMenu/releases)
- [GameTDB DS Downloads](https://www.gametdb.com/DS/Downloads)
- [GBATEK DS Cartridge Header](https://problemkaputt.de/gbatek-ds-cartridge-header.htm)
- [NDS ROM Format](https://github.com/Roughsketch/mdnds/wiki/NDS-Format)
- [Cowork Plugins](https://claude.com/blog/cowork-plugins)
- [Plugin Reference](https://code.claude.com/docs/en/plugins-reference)
- [Plugin Marketplaces](https://code.claude.com/docs/en/plugin-marketplaces)
- [Get Started with Cowork](https://support.claude.com/en/articles/13345190-get-started-with-cowork)
