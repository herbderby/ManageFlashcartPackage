---
title: "GameTDB - DS Box Art Cover Database"
sources:
  - url: "https://www.gametdb.com/DS/Downloads"
    description: "GameTDB DS downloads page"
  - url: "https://www.gametdb.com/DS"
    description: "GameTDB DS games database"
  - url: "https://www.gametdb.com/Main/FAQ"
    description: "GameTDB FAQ"
  - url: "https://art.gametdb.com/"
    description: "GameTDB art server (CDN)"
  - url: "https://github.com/FM1337/box-art-downloader"
    description: "Box art downloader reference implementation"
  - url: "https://wiki.ds-homebrew.com/twilightmenu/how-to-get-box-art"
    description: "DS-Homebrew wiki: How to get box art"
  - url: "https://github.com/KirovAir/TwilightBoxart"
    description: "TwilightBoxart tool - programmatic access example"
date: "2026-03-01"
---

# GameTDB - DS Box Art Cover Database

## Overview

GameTDB is a comprehensive game information database and box art repository that serves the Nintendo DS community. The platform hosts an extensive collection of cover artwork for DS games, organized by region and available in multiple sizes and formats. No signup is required to access or download covers and game databases.

## What is GameTDB?

GameTDB functions as both a game information database and a centralized artwork repository. The project is community-driven, with contributors regularly uploading and updating cover artwork for games across multiple Nintendo platforms (DS, Wii, Wii U, 3DS, PS3, and more). The database is actively maintained with artwork packs updated on a weekly basis.

The primary purpose of GameTDB is to provide accurate game metadata and high-quality box art covers that can be used in custom menu systems, emulators, and game launchers. The database is widely integrated into homebrew projects like TWiLight Menu++, which uses GameTDB as a primary source for retrieving game cover art.

## DS Cover Art System Architecture

### Core Components

GameTDB's DS cover art system consists of two main components:

1. **Web Interface** (gametdb.com/DS) - A browsable database where users can search for games and view available artwork
2. **CDN/Art Server** (art.gametdb.com) - A content delivery network serving individual cover images with direct URL patterns

### How the System Works

The DS game database links game identifiers (gamecodes) to available artwork in multiple regions and formats. When a game is uploaded to GameTDB, contributors can submit cover artwork for different regional variants. The system automatically generates multiple size variants from submitted images and stores them on the art server with consistent URL patterns.

Users and tools can programmatically fetch covers by constructing URLs based on known gamecodes and desired region/format combinations.

## URL Patterns for Downloading Covers

### Basic URL Structure

The cover art CDN uses a consistent hierarchical URL structure:

```
https://art.gametdb.com/{platform}/{cover_type}/{region}/{gamecode}.{extension}
```

### Platform and Path Components

- **Platform**: `ds` (lowercase)
- **Cover Type**: Determines the size/quality variant (see Cover Types section)
- **Region**: Two-letter region code (see Regions section)
- **Gamecode**: 4-character game identifier from the ROM header (e.g., ASME, IPKD)
- **Extension**: Typically `.png` or `.jpg` depending on cover type

### Example URLs

```
https://art.gametdb.com/ds/coverS/EN/ASME.png
https://art.gametdb.com/ds/coverM/US/ASME.jpg
https://art.gametdb.com/ds/cover/JA/ASME.jpg
https://art.gametdb.com/ds/coverHQ/EN/IPKD.png
```

### Common Variations

Different emulators and tools may use slightly different endpoint patterns or file extensions, but the core structure remains consistent.

## Available Cover Types and Sizes

GameTDB provides multiple cover variants to support different use cases and display requirements:

| Cover Type | Size/Quality | Use Case | Format |
|-----------|-------------|----------|--------|
| `coverS` | Small | Menu thumbnails, compact displays | PNG |
| `coverM` | Medium | Standard display size | JPG/PNG |
| `cover` | Full/Standard | General purpose cover display | JPG/PNG |
| `coverHQ` | High Quality | High-resolution, full-size artwork | PNG |

### Size Standards

For Nintendo DS menu systems like TWiLight Menu++, the standard cover dimensions are **128x115 pixels**. Covers smaller than this are typically the "S" variant, while the HQ versions provide high-resolution artwork suitable for larger displays or upscaling.

### Clickable Sizes

When browsing GameTDB's web interface, users can click on any artwork to display the largest available size variant, allowing preview of the highest-quality version available for a given game and region.

## Region Codes and Fallback Strategies

### Primary Region Codes

GameTDB uses standardized two-letter region codes for organizing covers:

| Code | Region |
|------|--------|
| US | United States |
| EN | English (Europe/General) |
| JA | Japan |
| JP | Japan (alternative) |
| EU | Europe (general) |
| DE | Germany |
| FR | France |
| IT | Italy |
| ES | Spain |
| NL | Netherlands |
| PT | Portugal |
| KO | Korea |
| ZH | Chinese |

### Region Detection in Gamecodes

The fourth character of a DS game's 4-character ID indicates its origin region:

- **J** = Japan
- **E** = North America (USA/Canada)
- **P** = Europe/Australia
- **K** = Korea
- **D** = Germany
- **F** = France
- **I** = Italy
- **S** = Spain

### Fallback Strategy

When fetching covers programmatically, implement a fallback strategy:

1. Try the primary region code for the detected region
2. Fall back to `EN` for European games when region-specific codes are unavailable
3. Fall back to `US` for North American games when unavailable
4. Try alternative region codes if primary variants don't exist
5. Some German games may have artwork under `EN` even when not available under `DE`

The default cover region for European games is now `EN` when only a single EUR version of cover artwork exists in the database.

## Gamecode System for Lookups

### Game ID Structure

Each DS ROM contains a 4-character game ID in its header, not added manually. This ID is also printed on cartridges and case boxes as the "NTR-" prefixed code (e.g., NTR-ASME-USA becomes ASME in the database).

### Gamecode Format

- **Characters 1-2**: Game identifier (e.g., AS for Super Mario 64 DS)
- **Character 3**: Additional variant/version code
- **Character 4**: Region code (J, E, P, K, D, F, I, S, etc.)

Example: `ASME` breaks down as:
- `AS` = Game ID
- `M` = Variant
- `E` = Region (North America)

### Extracting Gamecodes

Gamecodes are extracted from:
1. The ROM header (bytes 0x00-0x03 for NTR-mode DS games)
2. Database XML files (such as ADVANsCEne_NDScrc.xml)
3. GameTDB's web database when searching by title

### Database Files

GameTDB maintains XML database files containing comprehensive game metadata including gamecodes, regions, and available artwork. These are available through the Downloads page for offline use.

## Programmatic Access Patterns

### Direct HTTP Requests

The simplest programmatic approach is direct URL construction and HTTP requests:

```python
import requests

def get_cover(gamecode, region="EN", cover_type="coverS"):
    url = f"https://art.gametdb.com/ds/{cover_type}/{region}/{gamecode}.png"
    response = requests.get(url)
    return response.content if response.status_code == 200 else None
```

### Fallback Logic

Implement fallback chains to handle missing covers:

```python
def get_cover_with_fallback(gamecode, preferred_region="EN"):
    fallback_regions = [preferred_region, "EN", "US", "JA"]
    cover_types = ["coverS", "coverM", "cover"]

    for region in fallback_regions:
        for cover_type in cover_types:
            url = f"https://art.gametdb.com/ds/{cover_type}/{region}/{gamecode}.png"
            if url_exists(url):
                return url
    return None
```

### Existing Tools and Libraries

Several open-source projects implement GameTDB access:

- **TwilightBoxart** (GitHub: KirovAir/TwilightBoxart) - C# tool for fetching covers
- **TWBoxPy** (GitHub: pablouser1/TWBoxPy) - Python box art downloader for TwilightMenu++
- **box-art-downloader** (GitHub: FM1337/box-art-downloader) - Focused on DSIMenu++ compatibility
- **PicoCover** (GitHub: Scaletta/PicoCover) - Cross-platform cover art generator using GameTDB

### Integration Points

Tools typically integrate GameTDB by:

1. Scanning game ROM headers to extract gamecodes
2. Constructing URLs based on detected region and gamecode
3. Downloading covers with fallback logic
4. Storing covers locally in standard formats

## Rate Limiting and Usage Policies

### Terms of Service for Software Integration

According to GameTDB's usage guidelines, the database and artwork may be used in software with the following provisions:

- Adding a link to GameTDB in your application/tool allows users to contribute new games and artwork
- Contact GameTDB to inform them of your project so they can provide support
- No signup required to download game databases and cover packs

### Website Usage Restrictions

GameTDB artwork and databases should not be used on websites without explicit permission from the GameTDB administrators.

### Rate Limiting

No explicit rate limiting documentation is published, however best practices for respectful access include:

- Batch downloads during off-peak hours when possible
- Implement reasonable delays between requests if downloading many covers
- Cache downloaded artwork locally rather than fetching repeatedly
- Respect HTTP caching headers returned by the art server

### Fair Use Considerations

The cover artwork in GameTDB comes from multiple sources, including official Nintendo material and community contributions. Users should be aware of copyright considerations when using artwork outside of GameTDB's intended purpose.

## Offline Access

For systems without internet connectivity or unreliable connections, GameTDB provides complete offline packs:

- **Full Cover Packs**: Download complete regional collections from the DS Downloads page
- **Database XML**: Download game metadata XML files for offline game lookup
- **No Signup Required**: Both packs are freely available without account creation

## Platform Support

GameTDB covers multiple Nintendo platforms:

- **DS** (Nintendo DS)
- **3DS** (Nintendo 3DS)
- **Wii**
- **Wii U**
- **PS3**

This document focuses on the DS platform, but similar URL patterns and region codes apply to other supported platforms.

## Summary

GameTDB provides a robust, community-driven infrastructure for accessing DS game cover artwork. The consistent URL pattern system makes programmatic access straightforward, while multiple region codes and cover types accommodate diverse use cases. The system's flexibility, combined with freely available offline packs and active community maintenance, makes it the de facto standard for DS box art in custom menu systems and game launchers.
