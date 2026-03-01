---
title: "NDS ROM Header Format and Game Identification"
sources:
  - url: "https://problemkaputt.de/gbatek-ds-cartridge-header.htm"
    description: "GBATEK DS Cartridge Header reference"
  - url: "https://github.com/Roughsketch/mdnds/wiki/NDS-Format"
    description: "NDS ROM format wiki"
  - url: "https://ndspy.readthedocs.io/en/latest/api/rom.html"
    description: "ndspy ROM documentation with header details"
  - url: "https://www.romhacking.net/documents/%5B469%5Dnds_formats.htm"
    description: "Nintendo DS File Formats documentation"
  - url: "https://dsibrew.org/wiki/DSi_cartridge_header"
    description: "DSiBrew cartridge header reference"
date: "2026-03-01"
---

# NDS ROM Header Format and Game Identification

## Overview

Nintendo DS ROM files contain a standardized header structure at the beginning of each cartridge image. This header provides critical metadata about the game, including its title, game code, maker code, and technical specifications. The header format is essential for emulators, flashcarts, and ROM management tools to identify and properly load NDS games.

An NDS ROM file typically begins with the header information followed by the ARM9 and ARM7 code sections, file system data (NitroROM or NitroARC), and additional game data.

## NDS ROM Header Structure (First 0x20 Bytes)

The NDS ROM header contains vital identification and structural information in the first 32 bytes (0x00 to 0x1F). The following table shows the key fields:

| Offset | Size (bytes) | Field Name | Description |
|--------|-------------|-----------|-------------|
| 0x000 | 12 | Game Title | Game name in uppercase ASCII, padded with 0x00 |
| 0x00C | 4 | Gamecode | 4-character uppercase ASCII game identifier |
| 0x010 | 2 | Makercode | 2-character uppercase ASCII developer identifier |
| 0x012 | 1 | Unitcode | Hardware type: 0x00=NDS, 0x02=NDS+DSi, 0x03=DSi Only |
| 0x013 | 1 | Encryption Seed Select | Encryption seed for secure area (0-3) |
| 0x014 | 1 | Device Capacity | Cartridge capacity: 2^(17+X) bytes |
| 0x015 | 7 | Reserved | Reserved, typically 0x00 |
| 0x01C | 4 | Reserved | Region/revision info (reserved for future use) |

## Critical Header Fields for Game Identification

### Game Title (0x000 - 12 bytes)

The game title is stored as uppercase ASCII characters, with unused bytes padded to 0x00. This field contains the name of the game as it would be displayed in the ROM information.

**Example:**
```
"POKEMON DIAMOND" would be stored as:
50 4F 4B 45 4D 4F 4E 20 44 49 41 4D (hex)
```

### Gamecode (0x00C - 4 bytes)

The gamecode uniquely identifies the specific game and is the most critical field for game identification. It follows the NTR (Nitro) format: a letter followed by three alphanumeric characters (e.g., "APDE" for Pokémon Diamond European version).

**Structure:**
- Character 1: Game type (usually 'A', 'B', 'C', etc.)
- Characters 2-4: Game identifier

**Format Example:**
```
APDE = Pokémon Diamond (Europe)
AXPE = Pokémon Pearl (Europe)
ASME = Super Mario 64 DS (Europe)
```

**Homebrew/Unlicensed Identification:**
- **All zeros (0x00000000):** Indicates a homebrew or unlicensed ROM without an official gamecode
- **Four # characters (####):** Alternative homebrew indicator
- **Missing or invalid codes:** Typically indicates custom ROMs or modified versions

### Makercode (0x010 - 2 bytes)

The makercode identifies the publisher or developer of the game in two uppercase ASCII characters.

**Common Makercodes:**
- **01** - Nintendo
- **02** - Hudson Soft
- **08** - Capcom
- **09** - Hot B
- **0A** - (I-Max)
- **0L** - Acclaim
- **P0** - Thq
- **00** - Unlicensed/Homebrew

### Unitcode (0x012 - 1 byte)

The unitcode specifies the hardware compatibility of the ROM.

| Value | Hardware | Notes |
|-------|----------|-------|
| 0x00 | NDS | Original Nintendo DS |
| 0x02 | NDS+DSi | Compatible with both NDS and DSi |
| 0x03 | DSi | DSi-exclusive title |

## Extended Header Structure

Beyond the first 0x20 bytes, the NDS ROM header continues with additional technical specifications:

| Offset | Size (bytes) | Field Name | Description |
|--------|-------------|-----------|-------------|
| 0x020 | 4 | ARM9 ROM Offset | Offset to ARM9 executable code (minimum 0x4000, aligned to 0x1000) |
| 0x024 | 4 | ARM9 Entry Address | Entry point address for ARM9 code execution |
| 0x028 | 4 | ARM9 RAM Address | RAM destination for ARM9 code |
| 0x02C | 4 | ARM9 Size | Size of ARM9 code in bytes |
| 0x030 | 4 | ARM7 ROM Offset | Offset to ARM7 executable code |
| 0x034 | 4 | ARM7 Entry Address | Entry point address for ARM7 code execution |
| 0x038 | 4 | ARM7 RAM Address | RAM destination for ARM7 code (typically 0x02800000) |
| 0x03C | 4 | ARM7 Size | Size of ARM7 code in bytes |
| 0x040 | 4 | File System ROM Offset | Offset to file allocation table (FAT) / NitroROM header |
| 0x044 | 4 | File System Size | Size of file system data |

## Programmatic Reading of Header Fields

### Pseudo-code for Reading NDS Header

```
function readNDSHeader(romData):
    gameTitle = readASCII(romData, 0x000, 12)
    gamecode = readASCII(romData, 0x00C, 4)
    makercode = readASCII(romData, 0x010, 2)
    unitcode = readByte(romData, 0x012)
    encryptionSeed = readByte(romData, 0x013)
    deviceCapacity = readByte(romData, 0x014)

    // Determine cartridge size
    cartridgeSize = 2^(17 + deviceCapacity)

    // Determine hardware type
    if unitcode == 0x00:
        hardware = "NDS"
    else if unitcode == 0x02:
        hardware = "NDS+DSi"
    else if unitcode == 0x03:
        hardware = "DSi Only"

    // Check if homebrew
    if gamecode == "0000" or gamecode == "####":
        isHomebrew = true
    else:
        isHomebrew = false

    return {
        title: gameTitle,
        code: gamecode,
        maker: makercode,
        hardware: hardware,
        isHomebrew: isHomebrew,
        cartridgeSize: cartridgeSize
    }
```

### Python Example

```python
import struct

def read_nds_header(rom_path):
    with open(rom_path, 'rb') as f:
        header = f.read(0x50)

    # Extract header fields
    game_title = header[0x000:0x00C].decode('ascii').rstrip('\x00')
    gamecode = header[0x00C:0x010].decode('ascii')
    makercode = header[0x010:0x012].decode('ascii')
    unitcode = header[0x012]
    encryption_seed = header[0x013]
    device_capacity = header[0x014]

    # Calculate cartridge size
    cartridge_size = 2 ** (17 + device_capacity)

    # Map unitcode to hardware
    hardware_map = {
        0x00: "NDS",
        0x02: "NDS+DSi",
        0x03: "DSi Only"
    }
    hardware = hardware_map.get(unitcode, "Unknown")

    # Check if homebrew
    is_homebrew = (gamecode == "0000" or gamecode == "####")

    # Read ARM9 and ARM7 info
    arm9_offset = struct.unpack('<I', header[0x020:0x024])[0]
    arm7_offset = struct.unpack('<I', header[0x030:0x034])[0]

    return {
        'title': game_title,
        'gamecode': gamecode,
        'makercode': makercode,
        'unitcode': unitcode,
        'hardware': hardware,
        'is_homebrew': is_homebrew,
        'cartridge_size': cartridge_size,
        'arm9_offset': arm9_offset,
        'arm7_offset': arm7_offset
    }
```

## Homebrew ROM Identification

Homebrew and unlicensed NDS ROMs can be identified by examining key header fields:

### Identification Methods

1. **Gamecode Check**
   - Homebrew ROMs typically have gamecode of `0000` or `####`
   - Official Nintendo ROMs have valid four-character codes (e.g., `APDE`, `ASME`)

2. **Makercode Check**
   - Homebrew ROMs usually have makercode of `00` or `??`
   - Known publishers use specific codes (e.g., `01` for Nintendo)

3. **Title Format**
   - Homebrew titles may be short or use unusual naming conventions
   - Official releases follow Nintendo's naming guidelines

4. **Unitcode**
   - Homebrew can have any unitcode value
   - Official releases are carefully categorized

### Homebrew Identification Example

```python
def is_likely_homebrew(nds_header):
    gamecode = nds_header['gamecode'].strip()
    makercode = nds_header['makercode'].strip()

    homebrew_indicators = [
        gamecode == "0000" or gamecode == "####",
        makercode == "00" or makercode == "??",
        len(gamecode) < 4,
        len(makercode) < 2,
    ]

    return any(homebrew_indicators)
```

## Comparison with Other Consoles

### GBA (Game Boy Advance) ROM Header

GBA cartridges use a 192-byte header (0x00-0xBF) with the following key fields:

| Offset | Size | Field |
|--------|------|-------|
| 0x00 | 4 | Entry Point (ARM opcode) |
| 0x04 | 156 | Nintendo Logo (required for bootability) |
| 0xA0 | 12 | Game Title (uppercase ASCII, padded with 0x00) |
| 0xAC | 4 | Game Code (e.g., "AXVE" format) |
| 0xB0 | 2 | Maker Code |
| 0xB2 | 1 | Fixed Value (0x96) |
| 0xB3 | 1 | Main Unit Code (0x00 for standard GBA) |

GBA headers are located in ROM address space and include a required Nintendo logo that bootlegs lack.

### SNES (Super Famicom) ROM Header

SNES cartridges have an internal header located at CPU address range $00FFC0-$00FFDF (at the end of the ROM address space), with an optional extended header at $00FFB0-$00FFBF:

**Key Fields:**
- **Game Title:** 16 bytes of SNES-formatted text
- **Cartridge Type:** Specifies memory mapping hardware (LoROM, HiROM, etc.)
- **ROM/RAM Size:** Bytes specifying total ROM and RAM capacity
- **Maker Code:** 2-byte identifier
- **Checksum:** Used to verify ROM integrity

Some SNES ROM images also include a separate 512-byte (0x200) SMC header at the file's beginning, which is not part of the on-cartridge ROM.

### Game Boy / Game Boy Color Header

Game Boy and GBC ROMs contain a header in the memory range $0100-$014F, with the actual header content at $0104-$014F:

| Offset | Size | Field |
|--------|------|-------|
| 0x134 | 16 | Game Title (uppercase ASCII, padded) |
| 0x143 | 1 | CGB Compatibility Flag |
| 0x144 | 2 | License Code (New licensee ID) |
| 0x146 | 1 | Super Game Boy Flag |
| 0x147 | 1 | Cartridge Type (MBC, RAM, Battery) |
| 0x148 | 1 | ROM Size |
| 0x149 | 1 | RAM Size |
| 0x14A | 1 | Destination Code (Region: 0=Japan, 1=International) |
| 0x14B | 1 | License Code (Old licensee ID) |
| 0x14C | 1 | ROM Version Number |
| 0x14D | 1 | Header Checksum |
| 0x14E | 2 | Global Checksum |

The Game Boy header includes region codes and cartridge type information critical for determining MBC (Memory Bank Controller) type.

### NES (iNES Format) Header

NES ROM files use the iNES format with a simple 16-byte header:

| Offset | Size | Content | Description |
|--------|------|---------|-------------|
| 0x00 | 3 | "NES" | Magic number (ASCII) |
| 0x03 | 1 | 0x1A | DOS EOF marker |
| 0x04 | 1 | PRG ROM | Number of 16 KB PRG ROM banks |
| 0x05 | 1 | CHR ROM | Number of 8 KB CHR ROM banks |
| 0x06 | 1 | Flags 6 | Mapper lower nibble, mirroring, battery, trainer |
| 0x07 | 1 | Flags 7 | Mapper upper nibble, console type |
| 0x08 | 1 | PRG RAM | Number of 8 KB PRG RAM banks |
| 0x09 | 1 | Flags 9 | TV system (NTSC/PAL) |
| 0x0A | 1 | Flags 10 | TV system, PRG RAM presence |
| 0x0B-0x0F | 5 | Reserved | Unused/reserved for NES 2.0 |

The NES iNES header is much simpler than cartridge-based systems, using a file format convention rather than on-cartridge ROM data. The magic bytes "NES\x1A" allow instant identification of NES ROM files.

**Note:** Older iNES implementations ignored bytes 7-15, leading to variable interpretations. The NES 2.0 format extends this header with better mapper and console support.

## Summary

The NDS ROM header provides essential identification information through standardized fields:

- **Game Title (0x000):** Human-readable game name
- **Gamecode (0x00C):** Unique game identifier crucial for emulation and flashcart compatibility
- **Makercode (0x010):** Publisher/developer identification
- **Unitcode (0x012):** Hardware compatibility specification
- **Homebrew Identification:** Detectable through gamecode and makercode values

Understanding these header fields is essential for ROM management, emulation, flashcart functionality, and game identification tools. The header structure has been consistent across NDS and DSi cartridges, making it highly reliable for programmatic analysis.

Compared to other console ROM formats, the NDS header is well-standardized and comprehensive, providing more metadata than NES iNES but less extensive than SNES internal headers. The Game Boy Color header is similarly detailed but uses different offset locations and formatting conventions.
