---
name: reference-box-art-lookup
description: Sources and lookup ordering for cartridge box art — GameTDB for NDS titles, SHA1-based libretro lookups for emulated non-NDS systems, with fallbacks and resize requirements.
metadata:
  type: reference
---

# Box-Art Lookup Workflow

## NDS box art

- **Source:** GameTDB at
  `https://art.gametdb.com/ds/coverS/US/{TID}.png`
- **Image size:** 128 × 115 pixels — already the target size, no
  resize needed.
- **Key:** `{TID}` is the Title ID read from the ROM header.

## Non-NDS box art (emulated systems)

The workflow uses **content hashing first**, falling back to ROM
filename:

1. `compute_sha1(rom_file)` — hash the ROM bytes.
2. `lookup_nointro(sha1)` — query the embedded No-Intro database
   (see [reference-nointro-database.md](reference-nointro-database.md))
   for the canonical name + system.
3. Use the canonical name to fetch the libretro thumbnail URL.
4. **If `lookup_nointro` returns `found=false`,** fall back to
   the raw ROM filename — strip the extension and query libretro
   with that string.

## Non-NDS image conventions

- **Filename on disk:** `{full_rom_filename}.png` — derived from
  the ROM's filename, **not** the ROM header. Important: the
  filename matters because the cart's menu matches box art to
  ROM by filename, not by content.
- **Image size:** libretro images are **oversized** (typically
  several hundred pixels wide). They must be **resized to
  128 × 115** before being written to the card, or the menu's
  preview pane misrenders them.
