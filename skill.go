package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// flashcartSkill contains the embedded domain knowledge that teaches
// Claude how to compose the MCP primitive tools into flashcart
// management workflows. Placeholders like {{model_name}} are replaced
// by [handleFlashcartKnowledge] with model-specific values.
const flashcartSkill = `# Flashcart SD Card Management

You have MCP tools for managing Nintendo DS flashcart SD cards. The
card runs a dual-layer firmware: Wood R4 1.62 kernel (base) with
TWiLight Menu++ (menu system) on top. The tools are low-level
primitives; this document teaches you how to compose them.

This knowledge is configured for the **{{model_name}}** flashcart.

## SD Card Directory Structure

A fully set up {{model_name}} card looks like this:

` + "```" + `
/
+-- _DSMENU.dat                           # TWiLight autoboot wrapper
+-- _DS_MENU.dat                          # TWiLight autoboot wrapper
+-- BOOT.NDS                              # TWiLight Menu++ binary
+-- {{forwarder_file}}                    # {{model_name}} flashcart forwarder
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
BOOT.NDS, {{forwarder_file}}. User ROMs go in Games/ or roms/<system>/.
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

### Non-NDS Games (SHA1 lookup + libretro thumbnails)

Always hash the ROM first to resolve the canonical No-Intro name.
This handles renamed ROMs correctly. Fall back to the filename only
for homebrew, hacks, and bad dumps not in the database.

1. Call compute_sha1 on the ROM file to get its SHA1 hash.
2. Call lookup_nointro with the SHA1 hash.
3. If found (found=true): use the returned name and system for the URL.
4. If not found (found=false): fall back to the ROM filename without
   extension as the name, and map the extension to the libretro system
   name using the table below.
5. URL-encode the name for the URL. Download from libretro:
   https://thumbnails.libretro.com/{System}/Named_Boxarts/{Name}.png
6. Resize to 128x115 with resize_image (libretro images are oversized).
7. Save as _nds/TWiLightMenu/boxart/{full_rom_filename}.png.

Fallback extension to libretro system name (used only when SHA1 lookup
returns found=false):
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

IMPORTANT: After writing box art to the card, run clean_dot_files.

## Firmware Installation

Setup is three phases. Always extract archives to /tmp/ first, then
copy to the card. CDN downloads can fail -- retry or use mirrors.

### Phase 1: Wood R4 Kernel

1. Download the kernel:
   {{kernel_url}}
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
   e. Contents of Autoboot/{{autoboot_dir}}/ to card root
      (overwrites _DSMENU.dat and _DS_MENU.dat)
   f. Contents of Flashcart Loader/{{loader_dir}}/ to card root
      (creates {{forwarder_file}} and _wfwd/)
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
// an MCP prompt message, parameterized for the requested flashcart
// model. Claude receives this text at session start and uses it to
// compose the primitive tools into flashcart workflows.
func handleFlashcartKnowledge(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	text, err := substituteModel(flashcartSkill, req.Params.Arguments)
	if err != nil {
		return promptError(err)
	}
	return &mcp.GetPromptResult{
		Description: "Domain knowledge for managing DS flashcart SD cards",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartInitPrompt contains the step-by-step procedure for
// installing the Wood R4 1.62 kernel on a flashcart SD card.
// Placeholders are replaced by [handleFlashcartInit] with
// model-specific values.
const flashcartInitPrompt = `# Wood R4 Kernel Installation

Step-by-step procedure for installing the Wood R4 1.62 kernel on a
{{model_name}} flashcart SD card. Follow each step in order, calling
the named tool at each step.

## Step 1: Detect the SD Card

Call list_volumes. Look for a FAT32 volume (fsType "msdos") that is
NOT the system disk. Confirm the volume path with the user before
proceeding. Let CARD = the confirmed volume path (e.g., /Volumes/NDS).

## Step 2: Scan Existing Contents

Call list_directory on CARD. If __rpg/ already exists, warn the user
that a kernel is already installed and ask whether to overwrite.

## Step 3: Download the Kernel

Call download_file with:
  url: {{kernel_url}}
  path: /tmp/{{kernel_archive}}

The CDN sometimes drops connections. If the download fails, retry once.

## Step 4: Extract the Archive

Call extract_archive with:
  path: /tmp/{{kernel_archive}}
  destination: /tmp/wood_r4_kernel

## Step 5: Copy Kernel Files to Card

Call list_directory on /tmp/wood_r4_kernel to see extracted contents.
Then copy each item to CARD with copy_file (use recursive=true for
directories):

  copy_file: /tmp/wood_r4_kernel/__rpg -> CARD/__rpg (recursive)
  copy_file: /tmp/wood_r4_kernel/_DSMENU.DAT -> CARD/_DSMENU.DAT
  copy_file: /tmp/wood_r4_kernel/_DS_MENU.DAT -> CARD/_DS_MENU.DAT

## Step 6: Create the Games Directory

Call create_directory with path: CARD/Games

## Step 7: Fix globalsettings.ini

Call read_file with path: CARD/__rpg/globalsettings.ini
Change "showHiddenFiles = 1" to "showHiddenFiles = 0" in the text.
Call write_file with the modified text to: CARD/__rpg/globalsettings.ini

This prevents macOS AppleDouble junk files from appearing in the
Wood R4 file browser.

## Step 8: Clean AppleDouble Files

Call clean_dot_files with path: CARD
This removes all ._* resource fork files macOS created during copying.

## Step 9: Verify

Call list_directory on CARD. Confirm these items exist:
  __rpg/  Games/  _DSMENU.DAT  _DS_MENU.DAT

Report success to the user. The card is now bootable with the Wood R4
kernel. Next step: install TWiLight Menu++ (flashcart_twilight_install).
`

// handleFlashcartInit returns the Wood R4 kernel installation
// procedure as an MCP prompt message, parameterized for the requested
// flashcart model.
func handleFlashcartInit(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	text, err := substituteModel(flashcartInitPrompt, req.Params.Arguments)
	if err != nil {
		return promptError(err)
	}
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for installing Wood R4 1.62 kernel",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartTwilightPrompt contains the step-by-step procedure for
// installing TWiLight Menu++ on a flashcart card with Wood R4.
// Placeholders are replaced by [handleFlashcartTwilight] with
// model-specific values.
const flashcartTwilightPrompt = `# TWiLight Menu++ Installation

Step-by-step procedure for installing TWiLight Menu++ on a
{{model_name}} flashcart SD card that already has the Wood R4 1.62
kernel installed.

## Step 1: Detect the SD Card

Call list_volumes. Find the FAT32 volume (fsType "msdos"). Confirm
the path with the user. Let CARD = the confirmed volume path.

## Step 2: Verify Wood R4 Is Present

Call file_exists with path: CARD/__rpg
If __rpg/ does not exist, stop and tell the user to run
flashcart_init first.

## Step 3: Back Up Saves

Call list_directory on CARD/Games. If any .sav files exist, copy
them to a safe location (e.g., ~/DS_Flashcart/backups/saves/) using
copy_file before proceeding.

## Step 4: Download TWiLight Menu++

Call download_file with:
  url: https://github.com/DS-Homebrew/TWiLightMenu/releases/latest/download/TWiLightMenu-Flashcard.7z
  path: /tmp/TWiLightMenu-Flashcard.7z

## Step 5: Extract the Archive

Call extract_archive with:
  path: /tmp/TWiLightMenu-Flashcard.7z
  destination: /tmp/twilight_menu

## Step 6: Copy Files to Card (Order Matters!)

Copy these items from /tmp/twilight_menu to CARD in this exact order:

a. copy_file: /tmp/twilight_menu/_nds -> CARD/_nds (recursive)
   Creates the TWiLight core directory with nds-bootstrap,
   GBARunner2, themes, and settings.

b. copy_file: /tmp/twilight_menu/BOOT.NDS -> CARD/BOOT.NDS

c. copy_file: /tmp/twilight_menu/roms -> CARD/roms (recursive)
   Creates subdirectories for 17+ console systems.

d. copy_file: /tmp/twilight_menu/snemul.cfg -> CARD/snemul.cfg

e. Copy contents of Autoboot/{{autoboot_dir}}/ to CARD root:
   Call list_directory on "/tmp/twilight_menu/Autoboot/{{autoboot_dir}}/"
   Copy each file to CARD root. These overwrite the Wood R4 boot
   files with TWiLight autoboot wrappers (_DSMENU.dat, _DS_MENU.dat).

f. Copy contents of Flashcart Loader/{{loader_dir}}/ to CARD root:
   Call list_directory on "/tmp/twilight_menu/Flashcart Loader/{{loader_dir}}/"
   Copy each item to CARD root (recursive for directories).
   This creates {{forwarder_file}} and _wfwd/ (TWiLight's kernel loader).

## Step 7: Fix _wfwd/globalsettings.ini

Call read_file with path: CARD/_wfwd/globalsettings.ini
Change "showHiddenFiles = 1" to "showHiddenFiles = 0".
Call write_file with the modified text.

## Step 8: Clean AppleDouble Files

Call clean_dot_files with path: CARD

## Step 9: Verify

Call list_directory on CARD. Confirm these items exist:
  __rpg/  _nds/  _wfwd/  roms/  Games/
  _DSMENU.dat  _DS_MENU.dat  BOOT.NDS  {{forwarder_file}}  snemul.cfg

Report success. Next step: install emulators (flashcart_emulators).
`

// handleFlashcartTwilight returns the TWiLight Menu++ installation
// procedure as an MCP prompt message, parameterized for the requested
// flashcart model.
func handleFlashcartTwilight(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	text, err := substituteModel(flashcartTwilightPrompt, req.Params.Arguments)
	if err != nil {
		return promptError(err)
	}
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for installing TWiLight Menu++",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartEmulatorsPrompt contains the step-by-step procedure for
// installing the Virtual Console emulator add-on. The {{model_name}}
// placeholder is replaced by [handleFlashcartEmulators].
const flashcartEmulatorsPrompt = `# Virtual Console Emulators Installation

Step-by-step procedure for installing the Virtual Console emulator
add-on for TWiLight Menu++ on a {{model_name}} flashcart SD card.

## Step 1: Detect the SD Card

Call list_volumes. Find the FAT32 volume (fsType "msdos"). Confirm
the path with the user. Let CARD = the confirmed volume path.

## Step 2: Verify TWiLight Menu++ Is Present

Call file_exists with path: CARD/_nds/TWiLightMenu
If TWiLightMenu/ does not exist, stop and tell the user to run
flashcart_twilight_install first.

## Step 3: Download the Virtual Console Add-On

Call download_file with:
  url: https://github.com/DS-Homebrew/TWiLightMenu/releases/latest/download/AddOn-VirtualConsole.7z
  path: /tmp/AddOn-VirtualConsole.7z

## Step 4: Extract the Archive

Call extract_archive with:
  path: /tmp/AddOn-VirtualConsole.7z
  destination: /tmp/virtual_console

## Step 5: Copy Emulators to Card

Call copy_file with:
  source: /tmp/virtual_console/_nds
  destination: CARD/_nds
  recursive: true

This merges the emulator binaries into
CARD/_nds/TWiLightMenu/emulators/ (22 emulators covering Atari
2600/5200/7800, ColecoVision, Game Boy, Genesis, NES, SNES, PC
Engine, Neo Geo Pocket, WonderSwan, Intellivision, and more).

## Step 6: Clean AppleDouble Files

Call clean_dot_files with path: CARD

## Step 7: Verify

Call list_directory on CARD/_nds/TWiLightMenu/emulators/ to confirm
emulator binaries are present.

Report success. The card now supports retro console emulation via
TWiLight Menu++. Next step: add box art (flashcart_boxart).
`

// handleFlashcartEmulators returns the Virtual Console installation
// procedure as an MCP prompt message, parameterized for the requested
// flashcart model.
func handleFlashcartEmulators(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	text, err := substituteModel(flashcartEmulatorsPrompt, req.Params.Arguments)
	if err != nil {
		return promptError(err)
	}
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for installing Virtual Console emulators",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartBoxartPrompt contains the step-by-step procedure for
// scanning ROMs and downloading box art for each one.
const flashcartBoxartPrompt = `# Box Art Download and Installation

Step-by-step procedure for scanning all ROMs on a flashcart SD card
and downloading cover art for each one.

## Step 1: Detect the SD Card

Call list_volumes. Find the FAT32 volume (fsType "msdos"). Confirm
the path with the user. Let CARD = the confirmed volume path.

## Step 2: Create Box Art Directory

Call create_directory with path: CARD/_nds/TWiLightMenu/boxart

## Step 3: Scan for NDS ROMs

Call list_directory on CARD/Games. Collect all files ending in .nds.

## Step 4: Scan for Non-NDS ROMs

Call list_directory on CARD/roms. For each subdirectory, call
list_directory to find ROM files. Collect files with these extensions:
.gb .gbc .gba .nes .sfc .sms .gg .gen .pce .a26 .a52 .a78 .col .int
.ngp .ws

## Step 5: Download NDS Box Art

For each .nds file found:

a. Call read_bytes with path=<ROM path>, offset=12, length=4
   This reads the 4-byte Title ID (TID) from the ROM header.

b. If the TID is "####" or all zeros, skip it (homebrew, no art).

c. Call download_file with:
     url: https://art.gametdb.com/ds/coverS/US/<TID>.png
     path: CARD/_nds/TWiLightMenu/boxart/<TID>.png

   If US region fails (404), try fallback regions in order:
     https://art.gametdb.com/ds/coverS/EN/<TID>.png
     https://art.gametdb.com/ds/coverS/JA/<TID>.png

d. GameTDB coverS images are already 128x115 -- no resizing needed.

## Step 6: Download Non-NDS Box Art

For each non-NDS ROM file found:

a. Call compute_sha1 on the ROM file to get its SHA1 hash.

b. Call lookup_nointro with the SHA1 hash.

c. If found (found=true): use the returned name and system from the
   lookup result. If not found (found=false): use the ROM filename
   without extension as the name, and map the extension to the libretro
   system name (see flashcart_knowledge for the mapping table).

d. URL-encode the name (spaces become %20, parentheses stay as-is).
   Call download_file with:
     url: https://thumbnails.libretro.com/<System>/Named_Boxarts/<Name>.png
     path: /tmp/boxart_<filename>.png

e. Call resize_image with:
     source: /tmp/boxart_<filename>.png
     destination: CARD/_nds/TWiLightMenu/boxart/<full_rom_filename>.png
     width: 128
     height: 115

   The boxart filename for non-NDS games is the full ROM filename
   with .png appended (e.g., "Galaga - Destination Earth (USA).gbc.png").

## Step 7: Clean AppleDouble Files

Call clean_dot_files with path: CARD

## Step 8: Report Results

List all games and whether box art was found for each one. Note any
games where art could not be downloaded (homebrew, obscure titles).
Report the total number of box art images installed.
`

// handleFlashcartBoxart returns the box art download procedure as
// an MCP prompt message.
func handleFlashcartBoxart(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for downloading box art for all ROMs",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: flashcartBoxartPrompt},
			},
		},
	}, nil
}

// flashcartAddGamePrompt contains the step-by-step procedure for
// adding a single ROM to the flashcart with box art. The
// {{source_path}} placeholder is replaced with the actual path by
// [handleFlashcartAddGame].
const flashcartAddGamePrompt = `# Add a ROM to the Flashcart

Step-by-step procedure for copying a single ROM file to the flashcart
SD card and downloading its box art.

Source ROM: {{source_path}}

## Step 1: Detect the SD Card

Call list_volumes. Find the FAT32 volume (fsType "msdos"). Confirm
the path with the user. Let CARD = the confirmed volume path.

## Step 2: Identify the ROM Type

Determine the system from the file extension of the source ROM:

  .nds = Nintendo DS -> copy to CARD/Games/
  .gb .gbc .gba .nes .sfc .sms .gg .gen .pce .a26 .a52 .a78
  .col .int .ngp .ws = retro console -> copy to CARD/roms/<ext>/

Where <ext> is the extension without the dot (e.g., .gbc -> roms/gbc/).

## Step 3: Copy the ROM

Call copy_file with:
  source: {{source_path}}
  destination: CARD/Games/<filename> (for .nds)
           or: CARD/roms/<ext>/<filename> (for non-NDS)

## Step 4: Download Box Art

For NDS ROMs:
  a. Call read_bytes on the destination path, offset=12, length=4
     to read the 4-byte Title ID (TID) from the ROM header.
  b. If TID is "####" or all zeros, skip (homebrew, no art available).
  c. Call download_file with:
       url: https://art.gametdb.com/ds/coverS/US/<TID>.png
       path: CARD/_nds/TWiLightMenu/boxart/<TID>.png
     Fallback regions if US fails: EN, then JA.

For non-NDS ROMs:
  a. Call compute_sha1 on the ROM file to get its SHA1 hash.
  b. Call lookup_nointro with the SHA1 hash.
  c. If found (found=true): use the returned name and system.
     If not found (found=false): use the ROM filename without
     extension as the name, and map the extension to the libretro
     system name (see flashcart_knowledge for the mapping table).
  d. URL-encode the name. Call download_file with:
       url: https://thumbnails.libretro.com/<System>/Named_Boxarts/<Name>.png
       path: /tmp/boxart_temp.png
  e. Call resize_image with:
       source: /tmp/boxart_temp.png
       destination: CARD/_nds/TWiLightMenu/boxart/<full_rom_filename>.png
       width: 128
       height: 115

## Step 5: Clean AppleDouble Files

Call clean_dot_files with path: CARD

## Step 6: Confirm

Report the ROM name, where it was placed on the card, and whether
box art was successfully installed.
`

// handleFlashcartAddGame returns the add-game procedure as an MCP
// prompt message, with the source_path argument substituted in.
func handleFlashcartAddGame(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	sourcePath := req.Params.Arguments["source_path"]
	text := strings.ReplaceAll(flashcartAddGamePrompt, "{{source_path}}", sourcePath)
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for adding a ROM with box art",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartCleanupPrompt contains the step-by-step procedure for
// cleaning and verifying a flashcart SD card.
const flashcartCleanupPrompt = `# Flashcart Volume Cleanup

Step-by-step procedure for cleaning and verifying a flashcart SD card.

## Step 1: Detect the SD Card

Call list_volumes. Find the FAT32 volume (fsType "msdos"). Confirm
the path with the user. Let CARD = the confirmed volume path.

## Step 2: Clean AppleDouble Files

Call clean_dot_files with path: CARD
Report how many ._* files were removed.

## Step 3: Verify Directory Structure

Check that these key directories and files exist using file_exists:
  CARD/__rpg/                   (Wood R4 kernel)
  CARD/__rpg/globalsettings.ini
  CARD/_nds/TWiLightMenu/       (TWiLight Menu++)
  CARD/_wfwd/                   (Flashcart loader)
  CARD/Games/                   (NDS ROMs)
  CARD/roms/                    (Non-NDS ROMs)
  CARD/BOOT.NDS                 (TWiLight binary)
  CARD/_DSMENU.dat              (Autoboot wrapper)
  CARD/_DS_MENU.dat             (Autoboot wrapper)

Report which items are present and which are missing.

## Step 4: Verify INI Settings

Call read_file with path: CARD/__rpg/globalsettings.ini
Verify showHiddenFiles = 0. If it is 1, fix it with write_file.

Call read_file with path: CARD/_wfwd/globalsettings.ini
Verify showHiddenFiles = 0. If it is 1, fix it with write_file.

If either file was modified, call clean_dot_files on CARD again.

## Step 5: Report Disk Usage

Call list_volumes to get the volume's total and free space.
Calculate used space = total - free.
Report total, used, and free space in human-readable format.

## Step 6: Scan ROMs

Call list_directory on CARD/Games to count NDS ROMs.
Call list_directory on CARD/roms, then each subdirectory, to count
non-NDS ROMs by system.

Report ROM counts by system and total.
`

// handleFlashcartCleanup returns the volume cleanup procedure as an
// MCP prompt message.
func handleFlashcartCleanup(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Step-by-step procedure for cleaning and verifying the card",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: flashcartCleanupPrompt},
			},
		},
	}, nil
}

// flashcartManualPrompt contains a short user-facing manual that
// explains what the flashcart tools extension does, how to get
// started, and what workflows are available.
const flashcartManualPrompt = `# Flashcart Tools -- User Manual

Flashcart Tools is a Claude Desktop extension that helps you set up
and maintain Nintendo DS flashcart SD cards. Plug in your card, open
Claude, and ask -- Claude handles the firmware, ROMs, box art, and
cleanup.

## Getting Started

1. Install the extension by double-clicking flashcart-tools.mcpb.
2. Format your micro SD card as FAT32 (use Disk Utility or the
   Nintendo DS Flashcart Tool).
3. Insert the card into your computer.
4. Open Claude Desktop and say: "I want to set up my flashcart
   SD card." If Claude does not know your flashcart model, it will
   ask you to photograph the front and back of the cart so it can
   identify the model from the sticker text.

Claude will detect the card and walk you through the rest.

## Available Workflows

Run these in order for a fresh card, or individually as needed.
Ask Claude by name (e.g., "Run flashcart_init") or just describe
what you want and Claude will pick the right one.

| Prompt                     | What It Does                             |
|----------------------------|------------------------------------------|
| flashcart_identify         | Identify your flashcart model from photos |
| flashcart_init             | Install the Wood R4 1.62 base kernel     |
| flashcart_twilight_install | Install TWiLight Menu++ over the kernel   |
| flashcart_emulators        | Add Virtual Console emulators (22 cores)  |
| flashcart_boxart           | Scan all ROMs and download cover art      |
| flashcart_add_game         | Copy one ROM to the card with box art     |
| flashcart_cleanup          | Clean junk files and verify card layout   |

## Adding Games

- **NDS games:** Drop .nds files into the Games/ folder on the card
  (or ask Claude to copy them for you with flashcart_add_game).
- **Retro games:** Drop ROMs into roms/<system>/ on the card. Claude
  knows which folder matches which file extension (.gb, .gbc, .gba,
  .nes, .sfc, .gen, .sms, .gg, .pce, .a26, .a52, .a78, .col, .int,
  .ngp, .ws).
- **Box art:** Ask Claude to fetch box art after adding games. NDS
  art comes from GameTDB; retro art comes from libretro thumbnails.

## Maintenance

- Run flashcart_cleanup periodically to remove macOS junk files
  (AppleDouble ._* files) and verify the card's directory structure.
- Claude backs up your save files before making changes to the card.
- If you reformat or swap cards, your box art cache and saves are
  preserved in ~/DS_Flashcart/ on your computer.

## Supported Systems

The card runs NDS games natively. With the Virtual Console emulators
installed, it also plays: Game Boy, Game Boy Color, Game Boy Advance,
NES, SNES, Genesis/Mega Drive, Master System, Game Gear, PC Engine,
Atari 2600/5200/7800, ColecoVision, Intellivision, Neo Geo Pocket,
WonderSwan, and Pokemon Mini.

## Troubleshooting

- **Download fails ("unexpected EOF"):** The CDN sometimes drops
  connections. Ask Claude to retry -- it handles mirrors and fallback
  URLs automatically.
- **Junk files in the game menu:** Run flashcart_cleanup. This clears
  AppleDouble files and fixes the showHiddenFiles setting.
- **Card not detected:** Make sure it is formatted as FAT32 and
  mounted in Finder. Claude looks for FAT32 volumes that are not
  the system disk.
- **Box art missing for a game:** Homebrew ROMs and obscure titles
  may not have cover art in GameTDB or libretro. Claude will tell
  you which games it could not find art for.
`

// identifyGuideURL is the raw GitHub URL for the comprehensive
// flashcart identification guide. Fetched at runtime so guide
// updates take effect without rebuilding the binary.
const identifyGuideURL = "https://raw.githubusercontent.com/herbderby/ManageFlashcartPackage/main/research/flashcart_identification_guide.md"

// manualURL is the raw GitHub URL for MANUAL.md. The manual is
// fetched at runtime so edits to the file take effect without
// rebuilding the binary.
const manualURL = "https://raw.githubusercontent.com/herbderby/ManageFlashcartPackage/main/MANUAL.md"

// maxManualBody is the maximum number of bytes read from the
// fetched manual (256 KiB -- the actual file is ~2 KB).
const maxManualBody = 256 << 10

// fetchIdentifyGuide downloads the comprehensive flashcart
// identification guide from GitHub. On any error it returns the
// embedded flashcartIdentifyPrompt as a fallback.
func fetchIdentifyGuide(ctx context.Context) string {
	req, err := newRequest(ctx, identifyGuideURL)
	if err != nil {
		return flashcartIdentifyPrompt
	}
	resp, err := browserClient.Do(req)
	if err != nil {
		return flashcartIdentifyPrompt
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return flashcartIdentifyPrompt
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManualBody))
	if err != nil || len(body) == 0 {
		return flashcartIdentifyPrompt
	}
	return string(body)
}

// fetchManual tries to download MANUAL.md from GitHub. On any
// error (network, non-200 status, read failure) it returns the
// embedded fallback constant.
func fetchManual(ctx context.Context) string {
	req, err := newRequest(ctx, manualURL)
	if err != nil {
		return flashcartManualPrompt
	}
	resp, err := browserClient.Do(req)
	if err != nil {
		return flashcartManualPrompt
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return flashcartManualPrompt
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManualBody))
	if err != nil || len(body) == 0 {
		return flashcartManualPrompt
	}
	return string(body)
}

// handleFlashcartManual returns the user-facing manual as an MCP
// prompt message. It fetches the latest MANUAL.md from GitHub,
// falling back to the embedded constant if the fetch fails.
func handleFlashcartManual(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	text := fetchManual(ctx)
	return &mcp.GetPromptResult{
		Description: "User manual for Flashcart Tools",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// flashcartIdentifyPrompt teaches Claude how to identify a flashcart
// model from photographs of the front and back of the cartridge.
// The {{model_list}} placeholder is replaced at runtime with the
// current registry contents.
const flashcartIdentifyPrompt = `# Flashcart Identification

Help the user identify their Nintendo DS flashcart model by examining
photographs of the cartridge.

## Instructions

1. Ask the user to take two photos of their flashcart cartridge:
   - **Front:** the label/sticker side
   - **Back:** the circuit board / connector side

2. Examine the photos for these identifying features:
   - **Printed URL** on the sticker (most reliable identifier)
   - **Brand name** or model text on the label
   - **Color and design** of the sticker artwork
   - **PCB color** and chip layout on the back

3. Match the observed features against the known models listed below.

4. Report your findings:
   - The identified model ID (use this with other flashcart prompts)
   - The display name
   - Whether the model is fully supported for automated setup
   - If not supported, explain what the user can do instead

5. If you cannot identify the model:
   - List the features you observed
   - Suggest the user post photos to r/flashcarts on Reddit for
     community identification
   - Note that Wood R4-based carts are the most likely to work with
     this tool

## Known Models

{{model_list}}

## Common Visual Patterns

- **"ace3ds.com" on sticker** -> Ace3DS+ (ace3ds_plus)
- **"r4ils.com" on sticker** -> R4iLS (r4ils)
- **"r4ds.com" or "r4ds.cn" on sticker** -> Original R4 (r4_original)
- **"r4sdhc.com" on sticker** -> R4 SDHC (r4sdhc)
- **"dstt.net" or "ndstt.com" on sticker** -> DSTT (dstt)
- **"r4isdhc.com" with year on sticker** -> R4i-SDHC DEMON (r4i_sdhc_demon)
- **"r4ids.cn" on sticker** -> R4i Gold (r4i_gold)
- **"gateway-3ds.com" blue card** -> Gateway Blue (gateway_blue)
- **"acekard.com" on sticker** -> Acekard 2i (acekard_2i)
- **"supercard.sc" on sticker** -> SuperCard DSONE (supercard_dsone)

## Important Notes

- Many R4 clones look identical but use different firmware. The
  printed URL on the sticker is the most reliable way to tell them
  apart.
- If the sticker is missing or unreadable, the PCB markings and
  chip layout may help narrow it down.
- Carts labeled "R4" without further identification are often Ace3DS+
  or R4 SDHC clones.
`

// handleFlashcartIdentify returns the flashcart identification
// guide as an MCP prompt message. It fetches the comprehensive
// guide from GitHub, falling back to the embedded constant if the
// fetch fails. The model registry is always appended so Claude
// knows which models support automated setup.
func handleFlashcartIdentify(
	ctx context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	guide := fetchIdentifyGuide(ctx)
	text := guide + "\n\n" + modelListText()
	return &mcp.GetPromptResult{
		Description: "Identify a flashcart model from photographs",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

// handleFlashcartHelp is a tool (not a prompt) that returns the user
// manual. Tools are always visible in Chat's tool list, making this
// reliably discoverable when someone asks for help.
func handleFlashcartHelp(
	ctx context.Context,
	req *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, any, error) {
	text := fetchManual(ctx)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, nil, nil
}

// substituteModel looks up the flashcart model from the prompt
// arguments and replaces all model-specific placeholders in text.
// Returns an error if the model argument is missing or unknown, or
// if the model is recognized but not yet supported.
func substituteModel(text string, args map[string]string) (string, error) {
	id, ok := args["flashcart_model"]
	if !ok || id == "" {
		return "", fmt.Errorf("missing required argument: flashcart_model")
	}
	m, err := lookupModel(id)
	if err != nil {
		return "", err
	}
	if !m.Supported {
		return "", fmt.Errorf("flashcart model %q (%s) is recognized but not yet supported for automated setup", m.ID, m.DisplayName)
	}

	r := strings.NewReplacer(
		"{{model_name}}", m.DisplayName,
		"{{kernel_url}}", m.KernelURL,
		"{{kernel_archive}}", m.KernelArchive,
		"{{autoboot_dir}}", m.AutobootDir,
		"{{loader_dir}}", m.LoaderDir,
		"{{forwarder_file}}", m.ForwarderFile,
	)
	return r.Replace(text), nil
}

// promptError returns a GetPromptResult containing an error message.
// This is used when model lookup fails so the caller still gets a
// readable prompt response rather than a protocol-level error.
func promptError(err error) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Error",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: "Error: " + err.Error()},
			},
		},
	}, nil
}
