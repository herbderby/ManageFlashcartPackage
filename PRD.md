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

An MCPB desktop extension that gives Claude the tools and domain knowledge
to prepare and maintain a flashcart SD card on the user's behalf. The user
installs the `.mcpb` in Claude Desktop, plugs in their SD card, opens a
session, and Claude handles the rest: installing firmware, organizing ROMs,
fetching box art, and validating the card layout.

## Target Firmware: Wood R4 + TWiLight Menu++

The SD card runs a dual-layer firmware stack. The base layer is the
**Wood R4 1.62 kernel**, which provides the flashcart-specific boot
mechanism and NDS game loading. The upper layer is **TWiLight Menu++**,
which replaces the Wood R4 file browser with a modern menu system that
supports box art, emulators, and homebrew.

This dual-layer architecture was validated in a real-world setup session
documented in `research/flashcart-setup-report.md`.

### Wood R4 1.62 Kernel (Base Layer)

The Wood R4 kernel is the base firmware for Ace3DS+, R4iLS, and many R4
clone flashcarts. It provides the DLDI driver and NDS game loader that
TWiLight delegates to when running NDS games in "Kernel" mode.

**Download source:** DS-Homebrew Flashcard Archive

```
https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip
```

This URL redirects through a CDN that can fail with truncated responses.
See "Download Reliability" below.

**Kernel contents:**

- `_DSMENU.DAT` (427 KB) -- boot file for Ace3DS+ variant
- `_DS_MENU.DAT` (427 KB) -- boot file for R4iLS variant
- `__rpg/` -- kernel data directory:
  - `globalsettings.ini` -- global settings (see macOS fix below)
  - `savesize.bin` -- save size database
  - `game.dldi` -- DLDI driver
  - `backlight.ini` -- backlight settings
  - `fonts/` -- PCF font files
  - `ui/ACE/` -- UI theme (BMP images)
  - `language/` -- language files

**Post-install fix:** The default `__rpg/globalsettings.ini` ships with
`showHiddenFiles = 1`. On macOS, this causes AppleDouble `._*` files to
appear in the game menu. Must be changed to `showHiddenFiles = 0`.

### TWiLight Menu++ (Upper Layer)

TWiLight Menu++ is the actively maintained, open-source menu system for DS
flashcarts, developed by Rocket Robz and the DS-Homebrew community. On a
flashcart, it replaces the Wood R4 file browser as the primary interface
and adds support for launching emulators, displaying box art, and running
homebrew. Releases are published on GitHub at
[DS-Homebrew/TWiLightMenu](https://github.com/DS-Homebrew/TWiLightMenu/releases).

**Downloads required (from GitHub releases `latest`):**

1. `TWiLightMenu-Flashcard.7z` (~40.8 MB) -- main flashcard distribution
2. `AddOn-VirtualConsole.7z` (~3.7 MB) -- emulator binaries for retro
   consoles

**Archive structure for Ace3DS+:**

The flashcard .7z archive contains files for many different flashcart
models. For the Ace3DS+, the relevant pieces are:

- **Root files:** `_nds/`, `BOOT.NDS` (213 KB), `roms/` (pre-created
  subdirectories for 17+ console systems), `snemul.cfg`
- **`Autoboot/Ace3DS+/`:** Replacement `_DSMENU.dat` and `_DS_MENU.dat`
  that chainload TWiLight's `BOOT.NDS` instead of Wood R4
- **`Flashcart Loader/Ace3DS+/`:** `Wfwd.dat` (forwarder binary) and
  `_wfwd/` (29 files -- a stripped-down Wood R4 kernel for "Kernel"
  game loading mode, with its own `globalsettings.ini`)

**Installation order** (matters because autoboot overwrites kernel files):

1. Extract archive to `/tmp/`
2. Copy `_nds/` to SD root (merges with existing)
3. Copy `BOOT.NDS` to SD root
4. Copy `roms/` to SD root
5. Copy `snemul.cfg` to SD root
6. Copy contents of `Autoboot/Ace3DS+/` to SD root (overwrites
   `_DSMENU.dat` and `_DS_MENU.dat`)
7. Copy contents of `Flashcart Loader/Ace3DS+/` to SD root (creates
   `Wfwd.dat` and `_wfwd/`)
8. Fix `_wfwd/globalsettings.ini`: set `showHiddenFiles = 0`
9. Run `clean_dot_files` on entire volume
10. IMPORTANT: Preserve `_nds/TWiLightMenu/settings.ini` if it exists

**Flashcart model mapping** (archive subdirectory selection):

```
Ace3DS+ -> Autoboot/Ace3DS+/, Flashcart Loader/Ace3DS+/
R4iLS   -> Autoboot/R4iLS/ (fallback), Flashcart Loader/R4iLS/
DSTT    -> (uses YSMenu instead of Wood R4)
```

### Virtual Console Add-on

The Virtual Console add-on extracts its `_nds/` folder to the SD root,
merging into `_nds/TWiLightMenu/emulators/`. It provides 22 emulator
binaries:

| Emulator | System(s) |
|----------|-----------|
| StellaDS.nds | Atari 2600 |
| A5200DS.nds, A5200DSi.nds | Atari 5200 |
| A7800DS.nds | Atari 7800 |
| A8DS.nds | Atari XEGS |
| ColecoDS.nds | ColecoVision |
| SugarDS.nds | Amstrad CPC |
| gameyob.nds | Game Boy / Game Boy Color |
| jEnesisDS.nds, PicoDriveTWL.nds | Sega Genesis / Mega Drive |
| S8DS.nds | Sega Master System / Game Gear / SG-1000 |
| nesDS.nds | NES |
| SNEmulDS.nds (+ DSi/3DS variants) | SNES |
| NitroGrafx.nds | PC Engine / TurboGrafx-16 |
| NGPDS.nds | Neo Geo Pocket |
| NitroSwan.nds | WonderSwan |
| NINTV-DS.nds | Intellivision |
| PokeMini.nds | Pokemon Mini |

ROMs for these systems go in `roms/<system>/` subdirectories.

### SD Card Directory Structure (Field-Tested)

```
/
+-- _DSMENU.dat                           # TWiLight autoboot (was Wood R4)
+-- _DS_MENU.dat                          # TWiLight autoboot (was Wood R4)
+-- BOOT.NDS                              # TWiLight Menu++ main binary
+-- Wfwd.dat                              # Ace3DS+ flashcart forwarder
+-- snemul.cfg                            # SNES emulator config
+-- __rpg/                                # Wood R4 1.62 kernel (preserved)
|   +-- globalsettings.ini                #   (showHiddenFiles = 0)
|   +-- fonts/
|   +-- language/
|   +-- ui/ACE/
|   +-- game.dldi
|   +-- savesize.bin
|   +-- backlight.ini
+-- _wfwd/                                # TWiLight's Ace3DS+ kernel loader
+-- _nds/
|   +-- nds-bootstrap-release.nds
|   +-- GBARunner2_*.nds                  # 6 variants
|   +-- colorLut/
|   +-- TWiLightMenu/
|       +-- main.srldr, settings.srldr, ...
|       +-- emulators/                    # 22 emulator binaries
|       +-- boxart/                       # Box art PNGs
|       +-- gamesettings/
|       +-- extras/
+-- roms/                                 # Pre-created by TWiLight
|   +-- nds/  gb/  gbc/  gba/  nes/  snes/
|   +-- gen/  sms/  gg/  pce/
|   +-- a26/  a52/  a78/  xegs/  col/
|   +-- int/  ngp/  ws/  dsi/  mini/  sg/
+-- Games/                                # User's ROM directory
```

### Box Art Specifications

- **Format:** PNG
- **Recommended size:** 128x115, maximum: 208x143
- **Cached mode limit:** 44 KiB per image
- **Location:** `_nds/TWiLightMenu/boxart/`

**NDS games (GameTDB):**

- **Naming:** `{TID}.png` where TID is the 4-character title ID from the
  ROM header at offset 0x0C
- **Source:** `https://art.gametdb.com/ds/coverS/US/{TID}.png`
- **Fallback regions:** US -> EN -> JA
- GameTDB `coverS` images are typically already 128x115

**Non-NDS games (libretro thumbnails):**

- **Naming:** `{full_rom_filename}.png` (e.g.,
  `Galaga - Destination Earth (USA).gbc.png`)
- **Source:**
  `https://thumbnails.libretro.com/{LibretroSystem}/Named_Boxarts/{ROMName}.png`
- Filenames follow No-Intro naming conventions; mismatches return 404
- Libretro images are often oversized and need resizing to 128x115

**File extension to libretro system name mapping:**

| Extension | Libretro System Name |
|-----------|---------------------|
| .gb | Nintendo - Game Boy |
| .gbc | Nintendo - Game Boy Color |
| .gba | Nintendo - Game Boy Advance |
| .nes | Nintendo - Nintendo Entertainment System |
| .sfc | Nintendo - Super Nintendo Entertainment System |
| .sms | Sega - Master System - Mark III |
| .gg | Sega - Game Gear |
| .gen | Sega - Mega Drive - Genesis |
| .pce | NEC - PC Engine - TurboGrafx 16 |
| .a26 | Atari - 2600 |
| .a52 | Atari - 5200 |
| .a78 | Atari - 7800 |
| .col | Coleco - ColecoVision |
| .int | Mattel - Intellivision |
| .ngp | SNK - Neo Geo Pocket |
| .ws | Bandai - WonderSwan |

### NDS ROM Header (identification fields)

| Offset | Size     | Field     | Notes                              |
|--------|----------|-----------|------------------------------------|
| 0x000  | 12 bytes | Game Title| Uppercase ASCII, null-padded       |
| 0x00C  | 4 bytes  | Gamecode  | Uppercase ASCII, used as title ID  |
| 0x010  | 2 bytes  | Makercode | Uppercase ASCII publisher ID       |
| 0x012  | 1 byte   | Unitcode  | 00h=NDS, 02h=NDS+DSi, 03h=DSi     |

The 4-byte gamecode at offset 0x00C is the title ID used for box art lookups
and game identification throughout the system.

### macOS FAT32 Hygiene

macOS automatically creates AppleDouble resource fork files (prefixed with
`._`) on FAT32 volumes because FAT32 does not support extended attributes.
These files are invisible in Finder but visible to the flashcart menu system,
causing clutter.

**This was the single most persistent issue during field testing.** Every
write operation to the SD card generates `._*` files. During initial kernel
installation alone, 66 AppleDouble files were created.

**Rules:**

- `clean_dot_files` MUST be the final step of every operation that writes
  to the SD card. Treat it as a post-condition, like flushing a write buffer.
- Both `__rpg/globalsettings.ini` and `_wfwd/globalsettings.ini` must have
  `showHiddenFiles = 0`.
- After running `clean_dot_files`, the operation itself creates 0-2 more
  `._*` files (from writing the INI fix). Run it once more or accept the
  residual.

### Download Reliability

Both major downloads (Wood R4 kernel and TWiLight Menu++) involve CDN
redirects that can fail with truncated responses ("unexpected EOF").

**Mitigations:**

- Maintain a list of mirror URLs for each known download
- Always extract to `/tmp/` first, then copy to the SD card
- Verify file integrity after download (check file size against expected
  size, or check a known magic number at byte offset 0)
- Cache downloads in `/tmp/` so they survive retries within a session

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

## Delivery: MCPB Desktop Extension

### Why MCPB

The Cowork plugin format has a known bug: bundled native binaries are not
launched (only system PATH commands like `npx` work). The MCPB desktop
extension format (`@anthropic-ai/mcpb pack`) supports `"type": "binary"`
and works reliably in Claude Desktop via double-click install.

The MCP server runs on the **host machine**, giving it full access to
mounted SD cards.

### Extension Architecture

```
ext/
  manifest.json                          # MCPB manifest (dxt_version 0.1)
  server/
    flashcart-tools                      # Compiled Go binary (darwin/arm64)
```

The manifest declares the binary entry point and MCP server configuration:

```json
{
  "dxt_version": "0.1",
  "name": "flashcart-tools",
  "version": "0.2.2",
  "server": {
    "type": "binary",
    "entry_point": "server/flashcart-tools",
    "mcp_config": {
      "command": "${__dirname}/server/flashcart-tools"
    }
  }
}
```

### Build and Install

```
make build    # CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build
make pack     # zip ext/ into flashcart-tools.mcpb
```

Install by double-clicking `flashcart-tools.mcpb` in Finder, which opens
Claude Desktop and registers the MCP server.

### Domain Knowledge

Domain knowledge (TWiLight Menu++ conventions, ROM identification, box art
workflow, safety rules) is embedded in the binary as a Go string constant
and served via the `flashcart_knowledge` MCP prompt. This avoids external
dependencies -- the binary is fully self-contained.

### Size

The stripped binary is ~8 MB; the `.mcpb` ZIP is ~3.4 MB. Single platform
(darwin/arm64) for now; cross-compilation is trivial to add later.

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

**Volume:**

- **`list_volumes`:** List mounted removable volumes with mount point,
  filesystem type, total size, and free space.

**Filesystem:**

- **`list_directory`:** List files and subdirectories at a path with sizes
  and modification times.
- **`create_directory`:** Create a directory (and parents) at a path.
  Uses firmlink-safe mkdir for FAT32 volumes.
- **`move_file`:** Move or rename a file or directory.
- **`copy_file`:** Copy a file or directory. Supports recursive directory
  copying for installing firmware trees.
- **`delete_file`:** Delete a single file (not recursive).
- **`file_exists`:** Check whether a path exists and whether it is a file
  or directory.

**Text and binary:**

- **`read_file`:** Read a text file and return its contents. Used for
  reading INI configuration files on the SD card.
- **`write_file`:** Write text content to a file, creating parent
  directories as needed. Used for modifying INI files.
- **`read_bytes`:** Read N bytes at a byte offset from a file. Returns
  hex and ASCII. Claude uses this to parse ROM headers.
- **`read_json`:** Read and parse a JSON file.
- **`write_json`:** Write a JSON object to a file with indentation.

**Network:**

- **`download_file`:** Download a URL to a local file path. Follows
  redirects. CDN failures are possible (see Download Reliability).
- **`fetch_url`:** Fetch a URL and return the body as text (truncated at
  1 MiB). For HTML directory listings, API responses, etc.

**Archive:**

- **`extract_archive`:** Extract a .7z or .zip archive to a directory.
  Uses firmlink-safe mkdir for FAT32 volumes.

**Image:**

- **`resize_image`:** Resize a PNG to given dimensions using high-quality
  interpolation. Does not preserve aspect ratio (stretches to fit).
- **`image_info`:** Return dimensions and file size of an image.

**macOS hygiene:**

- **`clean_dot_files`:** Remove AppleDouble `._*` resource fork files
  from a directory tree recursively. MUST be called after every batch of
  writes to a FAT32 volume.

Each tool returns structured JSON. Each tool does one thing.
18 tools + 9 prompts total.

## Domain Knowledge (Embedded Prompts)

The domain knowledge teaches Claude everything it needs to compose the
primitives into flashcart workflows. It is embedded in the Go binary as
string constants and served via MCP prompts.

### flashcart_identify prompt

Teaches Claude to identify a flashcart model from photographs. Lists
visual features (sticker URLs, brand text) for each known model. Claude
reads the printed URL/text on the cart sticker, matches it to a model
ID, and returns the ID for use with the workflow prompts.

### flashcart_knowledge prompt

This is the core knowledge, parameterized by `flashcart_model`. It is
returned when Claude requests the prompt with a model ID. It covers:

**Dual-layer firmware architecture:**
- Wood R4 1.62 as base layer (boot files, DLDI driver, NDS game loading)
- TWiLight Menu++ as upper layer (menu system, emulators, box art)
- How autoboot files chainload TWiLight from the Wood R4 boot sequence
- Flashcart model mapping (which archive subdirectories to use)

**SD card directory conventions:**
- Field-tested directory tree with all system directories
- Wood R4 directories (`__rpg/`, `_wfwd/`) vs TWiLight directories
  (`_nds/`, `roms/`)
- Which files are system files vs. user files
- Where ROMs go by system type (17+ systems with emulators)

**ROM identification:**
- NDS: `read_bytes` at offset 0x0C (4 bytes) for title ID
- Non-NDS: use the ROM filename directly (No-Intro naming convention)
- File extension determines the console system
- Homebrew ROMs have gamecode `####` or all zeros

**Box art workflow (multi-source):**
- NDS: GameTDB at `https://art.gametdb.com/ds/coverS/US/{TID}.png`
- Non-NDS: libretro thumbnails at
  `https://thumbnails.libretro.com/{System}/Named_Boxarts/{Name}.png`
- Complete file extension to libretro system name mapping
- NDS box art naming: `{TID}.png`
- Non-NDS box art naming: `{full_rom_filename}.png`
- Download, resize to 128x115 if needed, save to
  `_nds/TWiLightMenu/boxart/`

**Firmware installation (three-phase):**
1. Wood R4 kernel: download, extract, configure `showHiddenFiles = 0`
2. TWiLight Menu++: download, extract, copy model-specific autoboot and
   loader files, configure `_wfwd/globalsettings.ini`
3. Virtual Console: download, extract emulators to `_nds/TWiLightMenu/emulators/`

**macOS FAT32 hygiene:**
- `clean_dot_files` as mandatory post-condition for all write operations
- INI configuration patching for `showHiddenFiles`
- Why this matters (AppleDouble resource fork files)

**Download reliability:**
- CDN redirect failures, mirror fallback
- Extract to `/tmp/` first, then copy to card
- File integrity verification

**Working directory management:**
- Persistent folder on user's computer for backups, staging, cache
- Structure and first-session setup

**Safety rules:**
- Back up saves before any card modification
- Confirm before deleting or overwriting
- Scan card before making changes
- `clean_dot_files` after every write batch

### Skill Architecture (Field-Tested)

The following composable sub-tasks were validated during the field test
session documented in `research/flashcart-setup-report.md`. Each can be
invoked independently or chained together. Every sub-task that writes to
the SD card MUST end with `clean_dot_files`.

**`flashcart-identify`:**
Teach Claude to identify a flashcart model from photographs. Lists
visual features for each known model. Returns the model ID.

**`flashcart-init` (parameterized by `flashcart_model`):**
Format verification, Wood R4 kernel download and installation,
`showHiddenFiles` fix, `Games/` directory creation. Input: flashcart
model ID. Output: bootable SD card with base kernel.

**`flashcart-twilight-install` (parameterized by `flashcart_model`):**
Download and install TWiLight Menu++ with appropriate autoboot and
flashcart loader files for the specified model. Fix
`_wfwd/globalsettings.ini`. Input: flashcart model ID. Output:
TWiLight overlaid on base kernel.

**`flashcart-emulators` (parameterized by `flashcart_model`):**
Download and install the Virtual Console add-on. Merges emulator
binaries into `_nds/TWiLightMenu/emulators/`.

**`flashcart-boxart`:**
Scan all ROM files on the card, determine system and identifier for
each, download cover art from the appropriate source (GameTDB for NDS,
libretro for everything else), resize to 128x115, save with the
correct naming convention. Report any games where art could not be
found.

**`flashcart-add-game`:**
Copy a ROM file to the appropriate directory, automatically download
and install box art for it.

**`flashcart-cleanup`:**
Run `clean_dot_files` on the entire volume. Verify directory structure
integrity. Report disk usage.

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
- **ROM dumping or ripping.** The user provides ROMs or points to
  repositories.
- **Linux support in Phase 2.** Can be added later trivially (Go
  cross-compiles to Linux with one command).
- **Flashcarts beyond the Wood R4 family.** Other carts (DSTT, etc.)
  use different kernels (YSMenu). They are recognized by the model
  registry but lack automated setup. Support can be added later.

## User Experience Walkthrough

```
User installs flashcart-tools.mcpb in Claude Desktop (double-click).
Inserts their FAT32-formatted SD card. Opens a session.

User: I want to set up my Ace3DS+ flashcart SD card.

Claude: Let me check for SD cards...

I found one removable volume:

  NDS (128 GB FAT32, 119.2 GB free)
  Mounted at /Volumes/NDS
  Status: Empty -- no kernel detected.

I'll set this up in three phases:
1. Install Wood R4 1.62 kernel (base firmware)
2. Install TWiLight Menu++ (menu system with emulators)
3. Set up ROM directories and clean up macOS junk files

Shall I proceed?

User: Yes.

Claude: Phase 1: Installing Wood R4 kernel...

  [x] Downloaded Ace3DS+_R4iLS_Wood_R4_1.62.zip
  [x] Extracted kernel files to card
  [x] Fixed showHiddenFiles = 0 in __rpg/globalsettings.ini
  [x] Created Games/ directory

Phase 2: Installing TWiLight Menu++...

  [x] Downloaded TWiLightMenu-Flashcard.7z v27.22.2
  [x] Downloaded AddOn-VirtualConsole.7z
  [x] Installed _nds/, BOOT.NDS, roms/ directories
  [x] Installed Ace3DS+ autoboot files
  [x] Installed flashcart loader (Wfwd.dat, _wfwd/)
  [x] Fixed showHiddenFiles = 0 in _wfwd/globalsettings.ini
  [x] Installed 22 Virtual Console emulators

Phase 3: Cleanup...

  [x] Removed 66 AppleDouble files from volume
  [x] Verified directory structure

Your card is ready. You can copy ROM files to the Games/
directory and I'll identify them and fetch box art.

User: I copied some games -- NDS, GBC, and an Atari 7800 ROM.

Claude: Let me scan the card...

I found 8 ROM files in Games/:
  5 NDS games (read title IDs from ROM headers)
  2 GBC games
  1 Atari 7800 game

I'll fetch box art for each:
  - NDS: from GameTDB using title IDs
  - GBC: from libretro thumbnails using filenames
  - A78: from libretro thumbnails using filename

Go ahead?

User: Yes please.

Claude: Done.

  [x] Downloaded box art for 5 NDS games from GameTDB (128x115)
  [x] Downloaded and resized GBC art from libretro (256x251 -> 128x115)
  [x] Downloaded and resized A78 art from libretro (512x714 -> 128x115)
  [ ] No art found for Crystal Depths (homebrew, not in any database)
  [x] Cleaned 9 AppleDouble files from boxart directory

7 of 8 games have box art. Eject the card and try it out.
```

## Phase 0: Proof of Concept (Complete)

Phase 0 validated the end-to-end chain: Go binary compiles, serves MCP
over stdio, `list_volumes` returns real volume data, and the `.mcpb`
package installs in Claude Desktop via double-click.

Key learnings:
- Cowork plugin format cannot launch bundled binaries (known bug).
  Pivoted to MCPB desktop extension format.
- Official Go MCP SDK (`go-sdk` v1.4.0) works well: `mcp.NewServer`,
  `mcp.AddTool[In, Out]`, `server.Run` with `&mcp.StdioTransport{}`.
- `jsonschema` struct tags are the description string directly, not
  `description=...` prefixed.
- `syscall.Statfs_t.Fstypename` is `[16]int8` on Darwin, needs
  int8-to-byte conversion.

## Phase 1: Full Tool Surface (Complete)

All 18 MCP tools and 1 prompt implemented and tested:
- **Volume:** `list_volumes`
- **Filesystem:** `list_directory`, `create_directory`, `move_file`,
  `copy_file` (recursive), `delete_file`, `file_exists`
- **Text/binary:** `read_file`, `write_file`, `read_bytes`
- **Network:** `download_file`, `fetch_url`
- **Archive:** `extract_archive` (7z + zip)
- **Image:** `resize_image`, `image_info`
- **JSON:** `read_json`, `write_json`
- **macOS hygiene:** `clean_dot_files`
- **Prompt:** `flashcart_knowledge` (embedded domain knowledge)

Key fixes during Phase 1:
- FAT32 firmlink-safe `mkdir` (`resolvedMkdirAll` in `pathutil.go`)
- Recursive `copy_file` for directory trees
- `read_file` / `write_file` for INI configuration editing
- `clean_dot_files` for AppleDouble removal

## Phase 2: Field Test (Complete)

Full end-to-end setup of an Ace3DS+ flashcart SD card using Cowork with
the Flashcart Tools MCP server. Documented in
`research/flashcart-setup-report.md`.

Validated:
- Wood R4 1.62 kernel installation (download, extract, configure)
- TWiLight Menu++ v27.22.2 installation with Ace3DS+ autoboot
- Virtual Console add-on (22 emulators)
- Box art for NDS (GameTDB) and non-NDS (libretro thumbnails)
- macOS FAT32 hygiene (66 AppleDouble files cleaned in first pass alone)
- Download reliability (CDN failures, mirror fallback)

Key learnings:
- `clean_dot_files` must be called after EVERY batch of writes to FAT32
- TWiLight archive contains model-specific subdirectories; must select
  the right ones for the flashcart model
- `_wfwd/globalsettings.ini` (TWiLight's kernel loader) also needs the
  `showHiddenFiles = 0` fix
- Non-NDS box art uses the full ROM filename, not a header-derived ID
- libretro thumbnail filenames follow No-Intro conventions exactly

## Open Questions

1. **GameTDB rate limiting.** No explicit rate limits documented, but
   a card with 500+ ROMs would generate many requests. Should implement
   reasonable delays and local caching.

2. **Cross-platform.** Currently darwin/arm64 only. Adding more
   platforms is trivial with Go cross-compilation but requires testing
   volume detection on each OS.

3. **Aspect ratio in resize.** `resize_image` stretches to fit the target
   dimensions. Tall cover art (e.g., Atari 7800 at ~1:1.4) gets
   noticeably squished at 128x115. Could add letterboxing for extreme
   aspect ratios.

4. **Other flashcart models.** DSTT uses YSMenu instead of Wood R4.
   The flashcart model mapping needs expansion as more carts are tested.

## User Manual

The canonical user-facing manual lives in
[`MANUAL.md`](MANUAL.md) at the repo root. The `flashcart_manual`
MCP prompt fetches it from GitHub at runtime so edits take effect
without rebuilding the binary. The embedded constant in `skill.go`
serves as a fallback when the fetch fails (offline, rate-limited,
etc.).

## References

- [DS-Homebrew Flashcard Archive](https://github.com/DS-Homebrew/flashcard-archive)
- [TWiLight Menu++ GitHub Releases](https://github.com/DS-Homebrew/TWiLightMenu/releases)
- [TWiLight Menu++ Flashcard Installation](https://wiki.ds-homebrew.com/twilightmenu/installing-flashcard)
- [TWiLight Menu++ Add-ons](https://wiki.ds-homebrew.com/twilightmenu/installing-addons)
- [TWiLight Menu++ Box Art Guide](https://wiki.ds-homebrew.com/twilightmenu/how-to-get-box-art)
- [GameTDB DS Cover Art](https://art.gametdb.com/ds/coverS/US/)
- [Libretro Thumbnail Database](https://thumbnails.libretro.com/)
- [Flashcart Guides -- Ace3DS+](https://flashcart-guides.github.io/wiki/cart-guides/ace3ds/)
- [Ace3DS+ Wood Kernel NAND Fix](https://gbatemp.net/threads/wood-r4-1-62-kernel-fix-for-ace3ds-r4ils.639921/)
- [GBATEK DS Cartridge Header](https://problemkaputt.de/gbatek-ds-cartridge-header.htm)
- [NDS ROM Format](https://github.com/Roughsketch/mdnds/wiki/NDS-Format)
- [GameTDB DS Downloads](https://www.gametdb.com/DS/Downloads)
- [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
