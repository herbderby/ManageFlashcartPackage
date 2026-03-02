# Flashcart Tools -- User Manual

Flashcart Tools is a Claude Desktop extension that helps you set up
and maintain Nintendo DS flashcart SD cards. Plug in your card, open
Claude, and ask -- Claude handles the firmware, ROMs, box art, and
cleanup.

## Getting Started

1. Install the extension by double-clicking flashcart-tools.mcpb.
2. Format your micro SD card as FAT32 (use Disk Utility or the
   Nintendo DS Flashcart Tool).
3. Insert the card into your computer.
4. Open Claude Desktop and say: "I want to set up my Ace3DS+
   flashcart SD card."

Claude will detect the card and walk you through the rest.

## Available Workflows

Run these in order for a fresh card, or individually as needed.
Ask Claude by name (e.g., "Run flashcart_init") or just describe
what you want and Claude will pick the right one.

| Prompt                     | What It Does                             |
|----------------------------|------------------------------------------|
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
