---
title: "TWiLight Menu++ - Flashcart Menu System"
sources:
  - url: "https://wiki.ds-homebrew.com/twilightmenu/installing-flashcard"
    description: "Official installation guide for flashcards"
  - url: "https://wiki.ds-homebrew.com/twilightmenu/how-to-get-box-art"
    description: "Official box art guide"
  - url: "https://github.com/DS-Homebrew/TWiLightMenu/releases"
    description: "GitHub releases page"
  - url: "https://github.com/DS-Homebrew/TWiLightMenu"
    description: "Main GitHub repository"
  - url: "https://www.gamebrew.org/wiki/TWiLight_Menu++"
    description: "GameBrew encyclopedia entry"
  - url: "https://emulation.gametechwiki.com/index.php/TWiLight_Menu++"
    description: "Emulation General Wiki"
date: "2026-03-01"
---

## Overview

TWiLight Menu++ is an open-source DSi Menu upgrade and replacement designed for Nintendo DS flashcards, the Nintendo DSi, the Nintendo 3DS, and Nintendo 2DS systems. It serves as a comprehensive menu system and game launcher that provides a modern interface for running games and applications on these devices.

## What is TWiLight Menu++?

TWiLight Menu++ is a free, community-developed menu system that enhances the functionality of Nintendo handheld consoles by providing:

- A customizable graphical user interface for launching games and applications
- Support for multiple game formats and console emulation
- Advanced features like automatic patching, overclocking, and display enhancements
- Multi-language support with extensive customization options

The software is actively maintained by the DS-Homebrew community and distributed through GitHub under an open-source license.

## Supported Systems and Emulators

TWiLight Menu++ can launch games from the following systems:

### Game Consoles
- Nintendo DS and DSi (primary support)
- Nintendo Game Boy and Game Boy Color
- Nintendo Game Boy Advance
- Nintendo Entertainment System (NES)
- Super Nintendo Entertainment System (SNES)
- Sega Master System and Game Gear
- Sega Genesis/Mega Drive
- Atari 2600, 5200, 7800, and XEGS
- MSX
- Intellivision
- Neo Geo Pocket
- Sord M5
- PC Engine/TurboGrafx-16
- WonderSwan
- ColecoVision
- SG-1000 and SC-3000

### Additional Features
- DSTWO plugins (for DSTWO flashcard users)
- Video playback capabilities
- DS(i)Ware ROM support
- On-the-fly AP (anti-piracy) patch application

## Installation on Flashcarts

### System Requirements

To install TWiLight Menu++ on a DS flashcart, you will need:

- A compatible DS flashcart with microSD card slot
- A microSD card (capacity depends on your game library)
- A computer with 7-Zip or compatible archive extraction software
- The latest TWiLight Menu++ release package

### Installation Steps

1. **Download the Latest Release**
   - Visit the GitHub releases page and download the latest TWiLight Menu++ package
   - Files are distributed in `.7z` archive format (see Release Format section below)

2. **Extract the Archive**
   - Use 7-Zip to extract the downloaded `.7z` file
   - On Windows, right-click the file, select "Show more options", hover over 7-zip, then click "Open archive"

3. **Copy Core Files to SD Card**
   - Copy the `_nds` folder from the extracted files to the microSD card root
   - Copy `boot_fc.nds` (flashcard version) from the Flashcard users directory to the microSD card root
   - Rename `boot_fc.nds` to `boot.nds` on the microSD card

4. **Flashcard-Specific Setup**
   - Extract flashcard-specific files from `Flashcard users/Autoboot/(your-flashcard-model)/` to the microSD root
   - Consult your specific flashcard documentation if your model is not listed

5. **DS Phat/DS Lite Configuration**
   - Access the DS system menu settings
   - Enable auto-start functionality so the flashcard launches automatically on boot

6. **Game Placement**
   - Create a `roms` folder on the microSD card root
   - Organize your game ROM files within the appropriate subdirectories for each system

### SD Card Directory Structure

The SD card root should contain the following after installation:

```
microSD root/
├── _nds/                              (Main application folder)
│   ├── TWiLightMenu/
│   │   ├── boxart/                   (Box art images)
│   │   ├── dsimenu/themes/           (Theme files)
│   │   ├── icons/                    (Custom game icons/banners)
│   │   ├── extras/                   (Cheats, additional data)
│   │   └── [other configuration files]
│   └── [emulator and application files]
├── boot.nds                           (Main bootable file, renamed from boot_fc.nds)
└── roms/                              (Games and ROM files)
```

### Game Loader Options

TWiLight Menu++ supports multiple game loaders for different compatibility and performance profiles:

#### nds-bootstrap (Default)
- Used by default for launching DS games
- On DS flashcards, uses B4DS mode
- Provides moderate compatibility with slight performance overhead

#### Pico Loader
- Alternative loader option for DS games
- Offers faster boot times compared to nds-bootstrap
- Improved compatibility with most games
- Some games may be incompatible

#### YSMenu
- Alternative menu system approach
- Installation: Copy `TTMenu` folder and `YSMenu.nds` to microSD root
- Important: Do not copy `TTMenu.dat` as this breaks autobooting
- Provides different interface and compatibility profile

## Release Distribution Format

### Archive Format

TWiLight Menu++ releases are distributed as `.7z` (7-Zip) compressed archives. This format offers:

- Better compression than traditional ZIP files
- Smaller download sizes
- Widespread support across platforms

### Release Structure

Each release archive contains:

- `_nds/` directory with all application and configuration files
- Platform-specific subdirectories including:
  - `DSi users/` - For DSi installation
  - `Flashcard users/` - For flashcart installation (includes `boot_fc.nds` and model-specific autoboot files)
  - `3DS users/` - For 3DS installation
- Documentation and release notes

### Release Variants

Releases are available for different target platforms:
- DSi
- Flashcard
- 3DS
- Additional utility packages

## Box Art Integration

Box art adds visual appeal to game libraries by displaying game cover artwork in the TWiLight Menu++ interface.

### Box Art Specifications

Box art files must meet specific requirements for proper display:

| Format | Extension | Recommended Size | Maximum Size | Cache Limit |
|--------|-----------|------------------|--------------|-------------|
| DS Games | .png | 128x115 | 208x143 | 44 KiB* |
| FDS | .png | 115x115 | - | 44 KiB* |
| NES/MS/GG/Genesis | .png | 84x115 | - | 44 KiB* |
| SNES | .png | 158x115 | - | 44 KiB* |

*Size limit applies only if "Box Art viewing" is set to "Cached" in TWiLight Menu++ settings.

### File Naming Conventions

Box art files can be named using either of two methods:

1. **By Game TID (Title ID)**
   - Example: `ASME.png` for a game with TID ASME

2. **By ROM Filename**
   - Example: `SM64DS.nds.png` (includes .nds extension in filename)

### File Storage Location

All box art files must be placed in:

```
sd:/_nds/TWiLightMenu/boxart/
```

### Obtaining Box Art

#### GameTDB Method
- Download box art from GameTDB database
- Select "S Covers (png)" category
- Download desired game covers as PNG images
- Place in `boxart/` folder

#### Automated Downloader Tool
- Use **TwilightBoxart** tool for automatic download and organization
- Supports: NDS, Game Boy (Advance), SNES, NES, Game Gear, Genesis
- Written in C# for Windows systems
- Simplifies bulk organization of cover art

#### Manual Collection
- Manually source PNG images from various online databases
- Ensure images meet size and format specifications
- Rename files according to TID or ROM filename convention

## Configuration and Settings

### Main Configuration Files

Key configuration files are stored in the `_nds/TWiLightMenu/` directory:

- **Theme configuration** - Located in `dsimenu/themes/[theme-name]/`
- **Cheat data** - Located in `extras/usrcheat.dat`
- **Custom icons** - Place in `icons/` folder
- **Box art** - Stored in `boxart/` folder

### Customization Options

TWiLight Menu++ provides extensive customization:

- **Theme selection** - Switch between pre-installed or custom themes
- **Language** - Support for multiple interface languages
- **Game list appearance** - Icon display, sorting options, filtering
- **Performance tuning** - Overclocking support on DSi/3DS (up to 133MHz)
- **Audio enhancement** - Sound frequency adjustment from 32kHz to 48kHz (DSi/3DS)
- **Widescreen mode** - 16:10 aspect ratio support for compatible games on 3DS

### Performance Features

#### Overclocking (DSi/3DS only)
- Increase CPU speed to 133MHz
- Improves performance in demanding games
- Does not work on original DS hardware

#### Audio Enhancement (DSi/3DS only)
- Upgrade sound output from 32kHz to 48kHz
- Results in higher quality audio for supported games

#### Widescreen Support (3DS only)
- Play DS games in 16:10 widescreen format
- Requires specific setup and game compatibility
- Not all games support widescreen mode

## Advanced Features

### Automatic AP-Patching

TWiLight Menu++ includes anti-piracy patch functionality that:

- Automatically detects and applies necessary patches to ROMs
- Performs patching in RAM (does not modify original ROM files)
- Ensures compatibility with games that check for legitimate cartridges

### Multiple Kernel Options

Users can configure TWiLight Menu++ as:
- **Primary kernel** - Default system menu replacement
- **Secondary kernel** - Dual-boot configuration alongside original firmware

### Add-on System

Additional features can be installed including:
- **Multimedia player** - View photos and play videos
- **Virtual Console** - Play games from retro systems
- **Additional emulators** - Extend system support

## Community and Support

### Development
- Active open-source project maintained by DS-Homebrew community
- Regular updates and feature additions
- GitHub repository for bug reports and contributions

### Documentation
- Comprehensive wiki with installation guides
- Troubleshooting FAQ
- Language support documentation

### Related Tools
- **TWiLight Menu Updater** - Automated update tool for easy upgrading
- **TwilightBoxart** - Automated box art downloader
- **nds-bootstrap** - Core game loader engine

## Version History

The project maintains a release history available on GitHub. Notable release milestones include major feature additions, compatibility improvements, and platform expansions. Check the official GitHub releases page for current version information and changelog details.

## Installation Alternatives

For users on other platforms:
- **Nintendo DSi** - Uses `BOOT.NDS` instead of `boot_fc.nds`
- **Nintendo 3DS/2DS** - Uses `.cia` format installation
- **Original DS** - Limited functionality, recommended for DS lite or phat models only

## Troubleshooting

Common issues and solutions:

- **Games not appearing**: Ensure ROM files are extracted from archives and placed in proper directories
- **Autoboot not working**: Verify DS system settings have auto-start enabled
- **Box art not displaying**: Check image format (must be PNG), size limits, and file location
- **Performance issues**: Consider using alternative game loaders (Pico Loader, YSMenu) for compatibility testing

## Conclusion

TWiLight Menu++ represents the most comprehensive and actively maintained menu system for DS flashcards, offering extensive customization, multi-system emulation support, and community-driven development. Its flexible architecture accommodates both casual users seeking a modern interface and advanced users requiring specific compatibility configurations.
