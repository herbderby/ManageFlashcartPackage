package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// flashcartSkill contains the embedded domain knowledge that teaches
// Claude how to compose the MCP primitive tools into flashcart
// management workflows. This text is returned by the
// flashcart_knowledge MCP prompt.
const flashcartSkill = `# Flashcart SD Card Management

You have MCP tools for managing Nintendo DS flashcart SD cards. The
card runs a dual-layer firmware: Wood R4 1.62 kernel (base) with
TWiLight Menu++ (menu system) on top. The tools are low-level
primitives; this document teaches you how to compose them.

## SD Card Directory Structure

A fully set up Ace3DS+ card looks like this:

` + "```" + `
/
+-- _DSMENU.dat                           # TWiLight autoboot wrapper
+-- _DS_MENU.dat                          # TWiLight autoboot wrapper
+-- BOOT.NDS                              # TWiLight Menu++ binary
+-- Wfwd.dat                              # Ace3DS+ flashcart forwarder
+-- snemul.cfg                            # SNES emulator config
+-- __rpg/                                # Wood R4 1.62 kernel
|   +-- globalsettings.ini                #   (showHiddenFiles = 0!)
|   +-- fonts/  language/  ui/ACE/
|   +-- game.dldi  savesize.bin  backlight.ini
+-- _wfwd/                                # TWiLight's kernel loader
+-- _nds/
|   +-- nds-bootstrap-release.nds
|   +-- GBARunner2_*.nds
|   +-- TWiLightMenu/
|       +-- emulators/                    # 22 emulator binaries
|       +-- boxart/                       # Cover art PNGs
|       +-- gamesettings/
|       +-- extras/
+-- roms/                                 # System-specific ROM dirs
|   +-- nds/ gb/ gbc/ gba/ nes/ snes/ gen/ sms/ gg/ pce/
|   +-- a26/ a52/ a78/ xegs/ col/ int/ ngp/ ws/ dsi/ mini/ sg/
+-- Games/                                # User's ROM directory
` + "```" + `

System files: __rpg/, _wfwd/, _nds/, _DSMENU.dat, _DS_MENU.dat,
BOOT.NDS, Wfwd.dat. User ROMs go in Games/ or roms/<system>/.
Box art goes in _nds/TWiLightMenu/boxart/.

## ROM Identification

**NDS ROMs (.nds):** Use read_bytes to extract identification fields
from the ROM header:

- **Game title:** offset 0x000, length 12 bytes. Uppercase ASCII.
- **Title ID (TID):** offset 0x00C, length 4 bytes. Uppercase ASCII.
  This is the key used for GameTDB box art lookups.
- **Makercode:** offset 0x010, length 2 bytes.
- **Unitcode:** offset 0x012, length 1 byte. 00h=NDS, 02h=NDS+DSi,
  03h=DSi only.

Homebrew ROMs have TID "####" or all zeros -- no box art on GameTDB.

The fourth character of the TID indicates region:
- J = Japan, E = North America, P = Europe, K = Korea

**Non-NDS ROMs:** Identify by file extension. Use the full filename
(No-Intro naming convention) for box art lookups.

## Box Art Workflow

All box art goes to _nds/TWiLightMenu/boxart/ on the card.
Recommended size: 128x115 pixels. Maximum: 208x143. Cached mode
limit: 44 KiB per image.

### NDS Games (GameTDB)

1. Read 4 bytes at offset 0x0C from the ROM to get the TID.
2. Download from GameTDB:
   https://art.gametdb.com/ds/coverS/US/{TID}.png
3. Fallback regions: US -> EN -> JA.
4. GameTDB coverS images are typically already 128x115.
5. Save as _nds/TWiLightMenu/boxart/{TID}.png.

### Non-NDS Games (libretro thumbnails)

1. Determine the system from the file extension.
2. Map extension to libretro system name (see table below).
3. URL-encode the ROM filename (without extension) for the URL.
4. Download from libretro:
   https://thumbnails.libretro.com/{System}/Named_Boxarts/{Name}.png
5. Resize to 128x115 with resize_image (libretro images are oversized).
6. Save as _nds/TWiLightMenu/boxart/{full_rom_filename}.png.

File extension to libretro system name:
  .gb  -> Nintendo - Game Boy
  .gbc -> Nintendo - Game Boy Color
  .gba -> Nintendo - Game Boy Advance
  .nes -> Nintendo - Nintendo Entertainment System
  .sfc -> Nintendo - Super Nintendo Entertainment System
  .sms -> Sega - Master System - Mark III
  .gg  -> Sega - Game Gear
  .gen -> Sega - Mega Drive - Genesis
  .pce -> NEC - PC Engine - TurboGrafx 16
  .a26 -> Atari - 2600
  .a52 -> Atari - 5200
  .a78 -> Atari - 7800
  .col -> Coleco - ColecoVision
  .int -> Mattel - Intellivision
  .ngp -> SNK - Neo Geo Pocket
  .ws  -> Bandai - WonderSwan

Libretro filenames follow No-Intro naming conventions exactly.
Mismatches return 404. Try URL-encoded filename first.

IMPORTANT: After writing box art to the card, run clean_dot_files.

## Firmware Installation

Setup is three phases. Always extract archives to /tmp/ first, then
copy to the card. CDN downloads can fail -- retry or use mirrors.

### Phase 1: Wood R4 Kernel

1. Download the kernel:
   https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip
   (CDN redirects can fail -- retry or use a SourceForge mirror)
2. Extract to /tmp/, then copy contents to card root.
3. Create a Games/ directory at card root for ROMs.
4. Fix __rpg/globalsettings.ini: use read_file to read it, change
   "showHiddenFiles = 1" to "showHiddenFiles = 0", write_file back.
5. Run clean_dot_files on the entire volume.

### Phase 2: TWiLight Menu++

1. Download TWiLightMenu-Flashcard.7z from GitHub releases latest:
   https://github.com/DS-Homebrew/TWiLightMenu/releases/latest
2. Extract to /tmp/.
3. Copy to card in this order (order matters!):
   a. _nds/ to card root (merges with existing)
   b. BOOT.NDS to card root
   c. roms/ to card root
   d. snemul.cfg to card root
   e. Contents of Autoboot/Ace3DS+/ to card root
      (overwrites _DSMENU.dat and _DS_MENU.dat)
   f. Contents of Flashcart Loader/Ace3DS+/ to card root
      (creates Wfwd.dat and _wfwd/)
4. Fix _wfwd/globalsettings.ini: set showHiddenFiles = 0.
5. IMPORTANT: Do NOT overwrite _nds/TWiLightMenu/settings.ini if it
   already exists -- it contains the user's preferences.
6. Run clean_dot_files on the entire volume.

### Phase 3: Virtual Console Emulators

1. Download AddOn-VirtualConsole.7z from the same GitHub release.
2. Extract to /tmp/.
3. Copy _nds/ to card root (merges emulators into
   _nds/TWiLightMenu/emulators/).
4. Run clean_dot_files on the entire volume.

## Working Directory Management

The user maintains a persistent folder on their computer (not on the
SD card) that survives card reformats and swaps. Structure:

` + "```" + `
~/DS_Flashcart/
+-- backups/
|   +-- saves/                           # .sav files from SD cards
|   +-- cards/                           # Card state snapshots
+-- staging/                             # ROMs waiting to be organized
+-- cache/
|   +-- boxart/                          # Downloaded box art PNGs
|   +-- twilight_menu/                   # Cached firmware releases
+-- config.json                          # Region, preferences
` + "```" + `

On first session, if the working directory is empty, create this
structure using create_directory.

## macOS FAT32 Hygiene

macOS creates AppleDouble resource fork files (._*) on FAT32 volumes
because FAT32 has no extended attribute support. These are invisible
in Finder but visible to the flashcart menu system, causing clutter.

CRITICAL: clean_dot_files MUST be the final step of every operation
that writes to the SD card. Treat it as a post-condition, like
flushing a write buffer. During field testing, a single kernel
installation generated 66 AppleDouble files.

Both INI files need showHiddenFiles = 0:
- __rpg/globalsettings.ini (Wood R4 kernel)
- _wfwd/globalsettings.ini (TWiLight's kernel loader)

## Safety Rules

- ALWAYS back up saves before modifying a card: copy saves/*.sav
  from the card to the working directory backups/saves/ folder.
- ALWAYS scan the card with list_directory before making changes
  so the user can see the current state.
- NEVER delete ROMs without explicit user approval.
- ALWAYS confirm before overwriting files on the card.
- Use file_exists to check before overwriting settings files.
- ALWAYS run clean_dot_files after writing to a FAT32 volume.
- ALWAYS extract archives to /tmp/ first, then copy to the card.
`

// handleFlashcartKnowledge returns the embedded domain knowledge as
// an MCP prompt message. Claude receives this text at session start
// and uses it to compose the primitive tools into flashcart workflows.
func handleFlashcartKnowledge(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Domain knowledge for managing DS flashcart SD cards",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: flashcartSkill},
			},
		},
	}, nil
}
