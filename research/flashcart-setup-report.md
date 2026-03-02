# Preparing an Ace3DS+ Flashcart SD Card with Cowork

## A Field Report for Building Claude Code Skills

**Date:** March 2026
**Hardware:** Ace3DS+/R4iLS clone flashcart, 128 GB micro SD card, MacBook Pro M1 Max
**Software:** Cowork with Flashcart Tools MCP server

---

## 1. Overview

This report describes an end-to-end session using Anthropic's Cowork desktop application to prepare a micro SD card for a Nintendo DS flashcart. The work covered four major phases: installing the Wood R4 1.62 base kernel, installing TWiLight Menu++ v27.22.2 with Virtual Console emulators, downloading and resizing box art for installed games, and dealing with a persistent class of macOS-specific filesystem problems throughout.

The goal is to document the procedures, the tooling that worked, the bugs we hit, and the domain knowledge required, so that a comprehensive Claude Code skill can be written to automate all of this for any DS flashcart.


## 2. The Flashcart Tools MCP Server

All SD card operations were performed through a custom MCP (Model Context Protocol) server called "Flashcart Tools" running locally on macOS. This server provides the following operations relevant to flashcart management:

**Filesystem operations:** `list_volumes`, `list_directory`, `file_exists`, `create_directory`, `copy_file`, `move_file`, `delete_file`, `read_file`, `write_file`, `read_bytes`, `read_json`, `write_json`

**Archive and download:** `download_file`, `extract_archive` (supports .7z and .zip), `fetch_url`

**Image processing:** `image_info` (returns width, height, file size), `resize_image` (PNG only, with high-quality interpolation)

**macOS cleanup:** `clean_dot_files` (recursively removes AppleDouble `._*` resource fork files)

The MCP server operates directly on mounted volumes. The SD card appeared at `/Volumes/NDS` after formatting as FAT32.

### 2.1 Tool Evolution During the Session

The Flashcart Tools server was updated twice during the session to fix bugs we discovered in real time:

**First version:** `extract_archive` silently failed when targeting the FAT32 volume -- it returned success but wrote no files. `copy_file` had no recursive mode, so directory trees could not be copied. `create_directory` failed on FAT32 with a misleading "not a directory" error.

**Second version:** `extract_archive` was fixed to work on FAT32. `create_directory` was fixed. `write_file` was added for text content. `clean_dot_files` was added, which turned out to be essential.

**Third version:** `read_file` for text files (not just JSON) was added, enabling us to read and modify INI configuration files.

This iterative development pattern is important context for skill design: the tools need to handle FAT32 quirks, and the `clean_dot_files` operation must be called after every batch of write operations to the SD card.


## 3. Phase 1: Wood R4 1.62 Kernel Installation

### 3.1 SD Card Preparation

The SD card was pre-formatted as FAT32 and labeled "NDS" by the user. It appeared at `/Volumes/NDS` with `fsType: "msdos"` (macOS's identifier for FAT32). The card was 128 GB, which is within the FAT32 maximum volume size for 64 KB clusters.

### 3.2 Kernel Download

The kernel source is the DS-Homebrew flashcard archive. The canonical URL is:

```
https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip
```

This URL redirects through a CDN (`us.dl.blobfrii.com`). The first download attempt failed with "unexpected EOF" from the CDN. A fallback download from a SourceForge mirror succeeded. A skill should implement retry logic with mirror fallback.

### 3.3 Extraction and File Layout

The zip was extracted to `/tmp/` first, then its contents were placed on the SD card root. The kernel consists of:

- `_DSMENU.DAT` (427 KB) -- boot file for Ace3DS+ variant
- `_DS_MENU.DAT` (427 KB) -- boot file for R4iLS variant
- `__rpg/` -- kernel data directory containing fonts (.pcf files), UI themes (BMP images), language files, `globalsettings.ini`, `savesize.bin`, `game.dldi`, and `backlight.ini`

After extraction, a `Games/` directory was created at the SD card root for ROM storage.

### 3.4 The globalsettings.ini Fix

The default `__rpg/globalsettings.ini` ships with `showHiddenFiles = 1`. This is a problem on macOS because every write operation to a FAT32 volume causes macOS to create AppleDouble `._*` resource fork files. With hidden files visible, the Wood R4 file browser displays these junk files alongside actual game ROMs.

The fix: read the INI file with `read_file`, change `showHiddenFiles = 1` to `showHiddenFiles = 0`, write it back with `write_file`. Then run `clean_dot_files` on the entire volume.

### 3.5 AppleDouble Cleanup

After every batch of file writes to the SD card, `clean_dot_files` must be called on `/Volumes/NDS`. In the initial kernel installation, 66 AppleDouble files were created. After cleaning and writing the fixed `globalsettings.ini`, 2 more appeared. This is unavoidable on macOS when writing to FAT32, and the cleanup must be the final step of any write sequence.


## 4. Phase 2: TWiLight Menu++ Installation

### 4.1 What TWiLight Menu++ Is

TWiLight Menu++ is an open-source replacement for the DS system menu, developed by Rocket Robz and the DS-Homebrew community. On a flashcart, it replaces the Wood R4 file browser as the primary interface and adds support for launching emulators, displaying box art, and running homebrew. The Wood R4 kernel is preserved underneath as a "flashcart loader" that TWiLight can invoke to actually boot NDS games.

### 4.2 Downloads Required

Two archives were downloaded from the GitHub releases page at `https://github.com/DS-Homebrew/TWiLightMenu/releases/latest`:

1. `TWiLightMenu-Flashcard.7z` (v27.22.2, ~40.8 MB) -- the main flashcard distribution
2. `AddOn-VirtualConsole.7z` (~3.7 MB) -- emulator binaries for retro consoles

### 4.3 Flashcard Installation Structure

The flashcard .7z archive contains a complex structure with files for many different flashcart models. For the Ace3DS+, the relevant pieces are:

**From the archive root:**
- `_nds/` folder -- TWiLight core files including nds-bootstrap, GBARunner2, color LUTs, theme/skin data, fonts, widescreen patches, and various `.srldr` loader modules
- `BOOT.NDS` (213 KB) -- TWiLight's main binary
- `roms/` folder -- pre-created subdirectories for 17 different console systems (nds, gba, dsi, gb, nes, snes, etc.)
- `snemul.cfg` (9.7 KB) -- SNES emulator configuration

**From `Autoboot/Ace3DS+/`:**
- `_DSMENU.dat` and `_DS_MENU.dat` -- these replace the Wood R4 originals to make TWiLight auto-boot when the flashcart starts. They are smaller wrapper files that chainload TWiLight's `BOOT.NDS`.

**From `Flashcart Loader/Ace3DS+/`:**
- `Wfwd.dat` -- the flashcart forwarder binary
- `_wfwd/` folder (29 files) -- a stripped-down Wood R4 kernel that TWiLight uses when Game Loader is set to "Kernel" mode. Contains its own `globalsettings.ini`, DLDI driver, UI theme, and language files.

### 4.4 Installation Sequence

The installation order matters because the autoboot files overwrite the original Wood R4 boot files:

1. Extract the archive to a temp directory
2. Copy `_nds/` to SD root (merging with any existing `_nds/` content)
3. Copy `BOOT.NDS` to SD root
4. Copy `roms/` to SD root
5. Copy `snemul.cfg` to SD root
6. Copy contents of `Autoboot/Ace3DS+/` to SD root (overwrites `_DSMENU.dat` and `_DS_MENU.dat`)
7. Copy contents of `Flashcart Loader/Ace3DS+/` to SD root (creates `Wfwd.dat` and `_wfwd/`)
8. Run `clean_dot_files` on entire volume

The original Wood R4 `__rpg/` directory is preserved. TWiLight's kernel loader uses `_wfwd/` instead.

### 4.5 Virtual Console Add-on

The Virtual Console add-on extracts its `_nds/` folder to the SD root, which merges into the existing `_nds/TWiLightMenu/emulators/` directory. After installation, 22 emulator binaries were present:

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

ROMs for these systems go in the corresponding `roms/<system>/` subdirectories.


## 5. Phase 3: Box Art

### 5.1 How TWiLight Menu++ Box Art Works

TWiLight displays small box art thumbnails in its game browser. The images are stored as PNG files in `_nds/TWiLightMenu/boxart/`. The naming convention differs by system:

**NDS games:** The file is named with the game's 4-character Title ID (TID), e.g., `ADME.png`. The TID is embedded in the NDS ROM header at byte offset `0x0C` (4 bytes, ASCII).

**Other systems (GBC, GBA, A78, etc.):** The file is named after the full ROM filename with `.png` appended, e.g., `Galaga - Destination Earth (USA).gbc.png`.

**Size constraints:** Recommended size is 128x115 pixels. Maximum is 208x143. For cached mode, file size should be 44 KB or less.

### 5.2 Extracting NDS Title IDs

For each NDS ROM, we used `read_bytes` to read 4 bytes at offset 12 (0x0C) from the ROM file. This extracts the ASCII TID directly from the ROM header. Example:

```
read_bytes(path="/Volumes/NDS/Games/Chrono Trigger (USA) (En,Fr).nds", offset=12, length=4)
-> ASCII: "YQUE"
```

The five NDS games yielded these TIDs:

| Game | TID |
|------|-----|
| Ace Attorney Investigations - Miles Edgeworth (USA) | C32E |
| Advance Wars - Dual Strike (USA, Australia) | AWRE |
| Animal Crossing - Wild World (USA) (Rev 1) | ADME |
| Castlevania - Dawn of Sorrow (USA) | ACVE |
| Chrono Trigger (USA) (En,Fr) | YQUE |

### 5.3 NDS Cover Art Sources

NDS cover art was downloaded from GameTDB, a community database of game metadata and art:

```
https://art.gametdb.com/ds/coverS/US/{TID}.png
```

The `coverS` path returns small covers. All five downloaded as 128x115 PNGs between 27-34 KB -- already the correct size for TWiLight. No resizing was needed.

### 5.4 Non-NDS Cover Art Sources

For non-NDS systems, the primary source is the libretro thumbnail database on GitHub:

```
https://thumbnails.libretro.com/{System}/Named_Boxarts/{Filename}.png
```

Where `{System}` is the libretro system name (e.g., "Nintendo - Game Boy Color" for GBC, "Atari - 7800" for A78). The filename must match the ROM name as catalogued by No-Intro or similar databases.

For the GBC game "Galaga - Destination Earth (USA).gbc", the libretro thumbnail was 256x251 and needed resizing to 128x115.

For homebrew titles not in any database (e.g., "Crystal Depths" by quackdev), we fell back to sourcing a banner image from the developer's itch.io page, downloading it, and resizing to 128x115.

For the Atari 7800 game "Centipede (USA).a78", the libretro thumbnail was 512x714 (tall Atari box format) and was resized to 128x115.

### 5.5 Image Resizing

The `resize_image` tool handles resizing with high-quality interpolation. It only works with PNG files. The operation is:

```
resize_image(source="/tmp/original.png", destination="/Volumes/NDS/_nds/TWiLightMenu/boxart/output.png", width=128, height=115)
```

Note that aspect ratio is not preserved automatically -- the image is stretched to fit the target dimensions. For most box art this is acceptable at thumbnail scale, but it means tall cover art (like Atari 7800 boxes at ~1:1.4) gets noticeably squished. A skill could implement letterboxing for extreme aspect ratios.

### 5.6 Final Box Art Cleanup

After writing all box art files, `clean_dot_files` was run on the boxart directory to remove the 9 AppleDouble files macOS had created during the process.


## 6. Key Lessons for Skill Design

### 6.1 The macOS FAT32 Problem

This was the single most persistent issue throughout the entire session. Every time any file is written to a FAT32 volume from macOS, the operating system creates an invisible `._*` companion file containing extended attributes and resource fork data. On the DS, these files are visible and create clutter.

A flashcart management skill **must** call `clean_dot_files` as the final step of every operation that writes to the SD card. This is not optional. It should be treated as a post-condition, analogous to flushing a write buffer.

### 6.2 Download Reliability

Both major downloads (Wood kernel and TWiLight Menu++) involved CDN redirects that sometimes failed. The Flashcart Tools `download_file` function follows redirects, but the CDNs themselves can return truncated responses.

A skill should:
- Implement a list of mirror URLs for each known download
- Verify file integrity after download (check file size against expected size, or check a known magic number at byte offset 0)
- Cache downloads in `/tmp/` so they survive retries within a session

### 6.3 Archive Extraction

The `.7z` format is standard in the DS homebrew community because it compresses better than zip. The `extract_archive` tool handles both. However, the TWiLight Menu++ archive contains files for dozens of different flashcart models, so a skill needs to know which subdirectories to extract for a given cart.

The mapping from flashcart model to subdirectory names is:

```
Ace3DS+ -> Autoboot/Ace3DS+, Flashcart Loader/Ace3DS+
R4iLS   -> Autoboot/R4iLS (if no Ace3DS+ folder), Flashcart Loader/R4iLS
DSTT    -> (uses YSMenu instead)
...
```

This mapping would need to be maintained as a lookup table in the skill.

### 6.4 Configuration File Handling

The Wood R4 kernel and TWiLight both use INI-format configuration files. The Flashcart Tools server has `read_file` and `write_file` for text content, which is sufficient for modifying these files. A skill should know the default values that need changing:

- `__rpg/globalsettings.ini`: set `showHiddenFiles = 0`
- `_wfwd/globalsettings.ini` (TWiLight's kernel loader): same fix may be needed

### 6.5 ROM Header Parsing

NDS ROM headers follow the official Nintendo format. The Title ID at offset 0x0C is the key field for box art lookup. A skill for box art automation should:

1. Scan the Games directory for ROM files
2. Identify the file extension to determine the system
3. For `.nds` files: read 4 bytes at offset 0x0C to get the TID
4. For other extensions: use the filename directly
5. Construct the appropriate download URL based on system and identifier
6. Download, resize if needed, save with the correct naming convention

### 6.6 Box Art URL Construction

The URL patterns for different systems:

**NDS (GameTDB):**
```
https://art.gametdb.com/ds/coverS/US/{TID}.png
```

**Other systems (libretro thumbnails):**
```
https://thumbnails.libretro.com/{LibretroSystem}/Named_Boxarts/{ROMFilenameWithoutExtension}.png
```

The libretro system name mapping:
```
.gb   -> Nintendo - Game Boy
.gbc  -> Nintendo - Game Boy Color
.gba  -> Nintendo - Game Boy Advance
.nes  -> Nintendo - Nintendo Entertainment System
.sfc  -> Nintendo - Super Nintendo Entertainment System
.sms  -> Sega - Master System - Mark III
.gg   -> Sega - Game Gear
.gen  -> Sega - Mega Drive - Genesis
.pce  -> NEC - PC Engine - TurboGrafx 16
.a26  -> Atari - 2600
.a52  -> Atari - 5200
.a78  -> Atari - 7800
.col  -> Coleco - ColecoVision
.int  -> Mattel - Intellivision
.ngp  -> SNK - Neo Geo Pocket
.ws   -> Bandai - WonderSwan
```

The libretro filenames follow No-Intro naming conventions. ROM filenames that don't match the database exactly will return 404. A skill could use URL-encoded filenames and implement fuzzy matching as a fallback.


## 7. Final SD Card Layout

```
/Volumes/NDS/
|-- _DSMENU.dat          TWiLight autoboot (was Wood R4 boot file)
|-- _DS_MENU.dat         TWiLight autoboot (was Wood R4 boot file)
|-- BOOT.NDS             TWiLight Menu++ main binary (213 KB)
|-- Wfwd.dat             Ace3DS+ flashcart kernel forwarder
|-- snemul.cfg           SNES emulator config (9.7 KB)
|-- __rpg/               Wood R4 1.62 kernel (preserved)
|   |-- globalsettings.ini   (showHiddenFiles = 0)
|   |-- fonts/
|   |-- language/
|   |-- ui/ACE/
|   |-- game.dldi
|   |-- savesize.bin
|   `-- backlight.ini
|-- _wfwd/               TWiLight's Ace3DS+ kernel loader (29 files)
|-- _nds/
|   |-- nds-bootstrap-release.nds
|   |-- GBARunner2_*.nds         (6 variants)
|   |-- colorLut/                (color filter LUTs)
|   `-- TWiLightMenu/
|       |-- main.srldr, settings.srldr, etc.
|       |-- emulators/           (22 emulator binaries)
|       |-- boxart/              (8 PNG files)
|       |-- extras/fonts/, extras/widescreen.pck
|       `-- {theme directories}
|-- roms/
|   |-- nds/    gb/     gbc/    gba/    nes/
|   |-- snes/   gen/    sms/    gg/     pce/
|   |-- a26/    a52/    a78/    xegs/   col/
|   |-- int/    ngp/    ws/     dsi/    mini/
|   `-- sg/
`-- Games/               NDS and other ROMs
    |-- Ace Attorney Investigations - Miles Edgeworth (USA).nds
    |-- Advance Wars - Dual Strike (USA, Australia).nds
    |-- Animal Crossing - Wild World (USA) (Rev 1).nds
    |-- Castlevania - Dawn of Sorrow (USA).nds
    |-- Chrono Trigger (USA) (En,Fr).nds
    |-- Galaga - Destination Earth (USA).gbc
    |-- Crystal Depths (v1.4).gbc
    `-- Centipede (USA).a78
```


## 8. Proposed Skill Architecture

A complete flashcart management skill for Claude Code would consist of several sub-tasks, each of which could be invoked independently or chained together:

### `flashcart-init`
Format verification, kernel download and installation, configuration patching, directory creation. Input: flashcart model. Output: bootable SD card with kernel.

### `flashcart-twilight-install`
Download and install TWiLight Menu++ with appropriate autoboot and loader files for the detected flashcart. Input: flashcart model (or auto-detect from existing kernel files). Output: TWiLight Menu++ overlaid on existing kernel.

### `flashcart-emulators`
Download and install the Virtual Console add-on. Optionally download BIOS files if the user provides them.

### `flashcart-boxart`
Scan all ROM files on the card, determine system and identifier for each, download cover art from the appropriate source, resize to 128x115, save with the correct naming convention. Report any games where art could not be found.

### `flashcart-add-game`
Copy a ROM file to the appropriate location, automatically download box art for it.

### `flashcart-cleanup`
Run `clean_dot_files` on the entire volume. Verify directory structure integrity. Report disk usage.

Each of these sub-tasks should end with `clean_dot_files` as a mandatory final step.


## 9. References

- DS-Homebrew Flashcard Archive: `https://github.com/DS-Homebrew/flashcard-archive`
- TWiLight Menu++ GitHub: `https://github.com/DS-Homebrew/TWiLightMenu`
- TWiLight Menu++ Wiki (flashcard installation): `https://wiki.ds-homebrew.com/twilightmenu/installing-flashcard`
- TWiLight Menu++ Wiki (add-ons): `https://wiki.ds-homebrew.com/twilightmenu/installing-addons`
- TWiLight Menu++ Wiki (box art): `https://wiki.ds-homebrew.com/twilightmenu/how-to-get-box-art`
- GameTDB cover art API: `https://art.gametdb.com/ds/coverS/US/{TID}.png`
- Libretro thumbnail database: `https://thumbnails.libretro.com/`
- Flashcart Guides (Ace3DS+): `https://flashcart-guides.github.io/wiki/cart-guides/ace3ds/`
- Ace3DS+ Wood kernel NAND save fix (davidmorom): `https://gbatemp.net/threads/wood-r4-1-62-kernel-fix-for-ace3ds-r4ils.639921/`
- NDS ROM header format: offset 0x00 = game title (12 bytes), offset 0x0C = game code / TID (4 bytes)
