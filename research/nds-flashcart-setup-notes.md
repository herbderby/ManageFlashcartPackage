# NDS Flashcart SD Card Setup Notes

## Goal

Prepare a microSD card for use with an Ace3DS+/R4iLS clone flashcart in a Nintendo DS Lite.

## Hardware

- **Flashcart:** Ace3DS+/R4iLS clone (R4 and R4SDHC Family per Sanrax)
- **SD Card:** 128 GB, FAT32 formatted, mounted at `/Volumes/NDS`
- **Console:** Nintendo DS Lite

## Kernel Choice

**Ace WoodR4 1.62** -- the standard kernel recommended by the Sanrax guide for this cart family. It supports cheats (needed for anti-piracy bypass in games like Pokemon Black/White 2), soft-reset back to the menu, and custom themes.

Three kernel options were available per the Sanrax guide:

1. **Ace WoodR4 1.62** (chosen) -- full-featured, cheat + soft-reset support
2. **Pico-Launcher** -- modern material UI, fast loader, but no cheats or soft-reset
3. **AceOS 2.13** -- WoodR4 + bundled emulators (GBARunner2/3, GameYob, NesDS, Moonshell2)

## Sanrax Guide Reference

- **Main guide page:** <https://sanrax.github.io/flashcart-guides/>
- **Ace3DS+/R4iLS guide:** <https://sanrax.github.io/flashcart-guides/cart-guides/ace3ds_r4ils/>
- **GitHub repo:** <https://github.com/Sanrax/flashcart-guides>
- **Flashcard archive repo:** <https://github.com/flashcarts/flashcard-archive>

## Setup Steps (from Sanrax guide)

1. Format the SD card (FAT32). Already done -- card is mounted at `/Volumes/NDS`.
2. Download the kernel zip from:
   `https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip`
3. Extract the zip and copy **the contents** (not the zip itself) into the root of the SD card.
   The root should contain `__rpg/` directory and `_DS_MENU.DAT` and `_DSMENU.DAT` files.
4. (Optional) Download a cheat database from:
   <https://gbatemp.net/threads/deadskullzjrs-nds-i-cheat-databases.488711>
   Copy `usrcheat.dat` to `__rpg/cheats/` on the SD card (create the `cheats` folder if needed).
5. Create a `Games` folder in the SD card root and place `.nds` ROM files inside.
6. Insert the SD card into the cart, plug the cart into the DS Lite, and boot.

### Expected SD card layout

```
/Volumes/NDS/
├── __rpg/
│   └── cheats/
│       └── usrcheat.dat    (optional)
├── _DS_MENU.DAT
├── _DSMENU.DAT
└── Games/
    ├── game1.nds
    └── game2.nds
```

### Note on the kernel binary patch

The `_DS_MENU.DAT` and `_DSMENU.DAT` in the archive have been patched by davidmorom to fix a NAND save bug affecting WarioWare D.I.Y. and Jam with the Band. The fix is a single ARM instruction change: `mov r6,#0x9` changed to `mov r6,#0x0`, preventing an erroneous 9-bit left shift on save file sector addresses. See: <https://gbatemp.net/threads/fixed-investigating-wood-1-62-nand-save-problems.659509/>

### Anti-piracy note

Ace WoodR4 1.62 is missing some AP patches for newer games (notably Pokemon Black/White 2). To work around this, press Y on the game in the WoodR4 menu, open cheats, and enable "Bypass Anti-Piracy" before launching.

## Post-Setup Cleanup

After writing files to the FAT32 volume from macOS, run the Flashcart Tools `clean_dot_files` function on `/Volumes/NDS` to remove the `._*` AppleDouble resource fork files that macOS scatters on FAT32 volumes. The DS doesn't need them and they clutter the file browser.

## Problem: `download_file` Tool Failure

The Flashcart Tools MCP server's `download_file` tool could not download the kernel zip. Every attempt failed with:

```
download "https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip":
  Get "https://us.dl.blobfrii.com/flashcard-archive//Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip":
  unexpected EOF
```

### What's happening

- `archive.flashcarts.net` redirects to `us.dl.blobfrii.com` (a CDN).
- The CDN begins the TLS/HTTP connection but drops it before delivering the payload (`unexpected EOF`).
- The same URL works fine in a standard web browser.

### Likely cause

The `download_file` tool's HTTP client sends minimal headers (probably no `User-Agent`, no `Accept-Encoding`, no typical browser headers). The blobfrii CDN likely rejects or deprioritizes bare clients -- possibly a Cloudflare-style bot mitigation, or the CDN requires specific headers to serve the file. A browser succeeds because it sends a full header set and handles any challenge pages transparently.

### Alternate mirrors tried

| URL | Result |
|-----|--------|
| `https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip` | `unexpected EOF` (redirects to blobfrii) |
| `https://flashcard-archive.ds-homebrew.com/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip` | `unexpected EOF` (same CDN) |
| `https://sourceforge.net/projects/flashcard-archive/files/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip/download` | Not confirmed (attempted but unclear result) |

### Workaround

Download the zip manually in a browser or with `curl`/`wget` (which send proper headers), then use the Flashcart Tools to extract and copy files to the SD card. For example:

```fish
curl -L -o /tmp/WoodR4_1.62.zip \
  "https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip"
```

### Suggested fix for `download_file`

The tool should send a standard `User-Agent` header (e.g., `Mozilla/5.0 ...` or at minimum something like `Flashcart-Tools/1.0`) and an `Accept-Encoding: gzip, deflate` header. Following redirects (which it appears to do, since it reaches blobfrii) is necessary but not sufficient if the CDN is filtering on headers.

## Remaining Steps

Once the zip is downloaded locally:

1. Extract the zip to the SD card root using `extract_archive`
2. Create the `Games` directory using `create_directory`
3. Optionally download and install a cheat database
4. Run `clean_dot_files` on `/Volumes/NDS` to remove macOS resource fork litter
5. Eject and test in the DS Lite
