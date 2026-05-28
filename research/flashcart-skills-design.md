# Flashcart SD Card Setup -- Skill Decomposition

This document breaks the end-to-end process of preparing a Nintendo DS
flashcart SD card into discrete, encapsulated skills. Each skill has a
single responsibility, clearly defined inputs and outputs, and can be
invoked independently or chained in sequence.

The decomposition is based on a real session where every step was done
manually. The problems encountered are noted inline.

---

## Dependency Graph

```
flashcart_detect
       |
flashcart_identify
       |
flashcart_init ──────────────────┐
       |                         |
flashcart_twilight_install       |
       |                         |
flashcart_emulators              |
       |                         |
flashcart_add_game (repeatable)  |
       |                         |
flashcart_boxart                 |
       |                         |
flashcart_cleanup ───────────────┘
```

`flashcart_detect` and `flashcart_identify` must run first.
Everything downstream of `flashcart_init` depends on knowing the
cart family. `flashcart_cleanup` can run at any point after init.

---

## Skill 1: flashcart_detect

**Purpose:** Find the target SD card and validate it is ready for use.

**Inputs:** None.

**Procedure:**

1. Call `list_volumes`.
2. Filter out system volumes (APFS, the boot disk, Time Machine).
3. Identify candidate removable volumes (FAT32/msdos on non-system
   mount points, or exFAT volumes that look like SD cards).
4. If exactly one candidate: select it.
   If zero: tell the user no card was found.
   If more than one: ask the user which volume to use.
5. Validate the filesystem is FAT32 (`fsType == "msdos"`).
   If exFAT or another type: warn the user and explain how to
   reformat. Do not proceed until the card is FAT32.
6. Record the mount point (e.g. `/Volumes/NDS2`) as the working
   target for all subsequent skills.

**Outputs:** A validated mount point string and volume metadata
(name, size, free space).

**Problem this solves:** In the real session, the card was exFAT and
I did not catch it until I manually inspected the `list_volumes`
output. The user had to reformat and ask me to try again.

---

## Skill 2: flashcart_identify

**Purpose:** Determine the hardware family of the flashcart from
photographs.

**Inputs:** One or two photos of the cart (front and back).

**Procedure:**

1. Call `flashcart_identify` to load the visual identification guide.
2. Follow the guide's procedure strictly: PCB color first, then shell
   indents, then label text. Never identify from the front label alone.
3. Map the identification to one of the known hardware families:
   - `ace3ds_plus` (Ace3DS+ / R4iLS clone family)
   - `demon` (DSTTi DEMON clone family)
   - `r4_original` (Original R4 DS)
   - `dstt` (Original DSTT / DSTTi)
   - `dspico` (DSpico)
   - `other` (recognized but not automated)
4. Record the cart family identifier for use by downstream skills.

**Outputs:** A cart family identifier string and a confidence level.

**Problem this solves:** In the real session, I skipped identification
entirely and assumed Wood R4 for original R4. The cart was actually an
Ace3DS+ clone with a misleading DEMON-style label. Installing the
wrong kernel means the cart simply will not boot.

---

## Skill 3: flashcart_init

**Purpose:** Install the correct base kernel for the identified cart.

**Inputs:** Mount point (from `flashcart_detect`), cart family
identifier (from `flashcart_identify`).

**Procedure:**

1. Dispatch on cart family:

   | Cart Family    | Kernel                    | Download URL |
   |----------------|---------------------------|-------------|
   | `ace3ds_plus`  | Ace Wood R4 1.62          | `https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip` |
   | `r4_original`  | Wood R4 1.62 (original)   | `https://archive.flashcarts.net/R4_original_M3_Simply/R4DS_Wood_R4_1.62.zip` |
   | `demon`        | YSMenu 7.06 + DEMON R4.dat | (two separate downloads, see identification guide Section 8.2) |
   | `dstt`         | YSMenu 7.06 (DSTTi-Clone) | (from RetroGameFan multi-cart update) |
   | `dspico`       | Pico-Launcher             | (from LNH team GitHub) |

2. Download the kernel archive to `/tmp/`.
3. Extract it.
4. Copy the kernel files to the card root. The exact file set varies
   by kernel (e.g. `__rpg/` + `_DS_MENU.DAT` for Ace Wood R4;
   `TTMenu/` + `TTMenu.dat` + `R4.dat` for DEMON YSMenu).
5. Create a `Games/` directory on the card root.
6. Run `clean_dot_files` on the card.

**Outputs:** A bootable card with the base kernel installed.

**Problem this solves:** The existing `flashcart_init` workflow only
knows about one kernel. This skill encodes the dispatch table so the
correct kernel is installed for the actual hardware.

---

## Skill 4: flashcart_twilight_install

**Purpose:** Install TWiLight Menu++ over the base kernel, with the
correct autoboot and flashcart loader for the cart family.

**Inputs:** Mount point, cart family identifier.

**Procedure:**

1. Download `TWiLightMenu-Flashcard.7z` from the latest GitHub
   release:
   `https://github.com/DS-Homebrew/TWiLightMenu/releases/latest/download/TWiLightMenu-Flashcard.7z`
2. Extract it to `/tmp/twilight/`.
3. Copy `_nds/`, `BOOT.NDS`, `roms/`, and `snemul.cfg` to the card.
4. Map the cart family to the correct autoboot folder:

   | Cart Family   | Autoboot Folder                    |
   |---------------|------------------------------------|
   | `ace3ds_plus` | `Autoboot/Ace3DS+/`                |
   | `r4_original` | `Autoboot/Original R4 & M3 Simply/`|
   | `demon`       | `Autoboot/R4i-SDHC/`              |
   | `dstt`        | `Autoboot/DSTT/`                   |
   | `dspico`      | `Autoboot/DSpico/`                 |

5. Copy the contents of the autoboot folder to the card root.
   These files replace the base kernel boot files so the cart
   boots directly into TWiLight Menu++.
6. Map the cart family to the correct flashcart loader folder:

   | Cart Family   | Loader Folder                |
   |---------------|------------------------------|
   | `ace3ds_plus` | `Flashcart Loader/Ace3DS+/`  |
   | `r4_original` | (use YSMenu tab procedure)   |
   | `demon`       | (use YSMenu tab procedure)   |
   | `dstt`        | (use YSMenu tab procedure)   |
   | `dspico`      | (not needed)                 |

7. Copy the contents of the loader folder to the card root.
8. Run `clean_dot_files` on the card.

**Outputs:** TWiLight Menu++ installed with autoboot and kernel
loader configured for the specific cart.

**Problem this solves:** In the real session, I had to manually
determine which of the many autoboot and loader folders matched the
Ace3DS+. The README.txt inside the archive helps, but parsing it and
selecting the right folder should be automated by the dispatch table.

---

## Skill 5: flashcart_emulators

**Purpose:** Install the TWiLight Menu++ Virtual Console add-on
(retro console emulators).

**Inputs:** Mount point.

**Procedure:**

1. Download `AddOn-VirtualConsole.7z` from the latest GitHub release:
   `https://github.com/DS-Homebrew/TWiLightMenu/releases/latest/download/AddOn-VirtualConsole.7z`
2. Extract it to `/tmp/vc/`.
3. Copy `_nds/TWiLightMenu/emulators/` to the card (merge into
   existing `_nds/` tree).
4. Copy `_nds/TWiLightMenu/addons/` to the card.
5. Copy `roms/` subdirectories to the card (these are empty
   placeholder folders like `roms/gb/`, `roms/nes/`, etc.).
6. Run `clean_dot_files` on the card.

**Outputs:** 22 emulator .nds binaries installed, ROM directory
structure created.

**Notes:** This skill is simple and cart-family-independent. It only
requires that TWiLight Menu++ is already installed (i.e.
`flashcart_twilight_install` has run).

---

## Skill 6: flashcart_add_game

**Purpose:** Add a single ROM to the card with a clean filename.

**Inputs:** Mount point, path to the ROM file on the host.

**Procedure:**

1. Determine the file extension to identify the system.
2. Copy the ROM to `Games/` on the card.
3. Clean the filename:
   - Strip region tags: `(USA)`, `(Europe)`, `(World)`, etc.
   - Strip revision tags: `(Rev 1)`, `(Rev 2)`, etc.
   - Strip language lists: `(En,Fr,De,Es,It)`, etc.
   - Strip `(GB Compatible)` and similar qualifier tags.
   - Clean up `The` prefix clutter:
     `Lord of the Rings, The - The Fellowship of the Ring`
     becomes `Lord of the Rings - Fellowship of the Ring`.
   - Preserve the file extension.
4. Run `clean_dot_files` on the `Games/` directory.

**Outputs:** A cleanly named ROM on the card.

**Design note:** Filename cleaning should be a pure function that
takes a No-Intro canonical name and returns a human-friendly display
name. The mapping from canonical to clean name should be deterministic
and testable.

---

## Skill 7: flashcart_boxart

**Purpose:** Download cover art for all ROMs on the card.

**Inputs:** Mount point.

**Procedure:**

1. Ensure `_nds/TWiLightMenu/boxart/` exists on the card.
2. Scan `Games/` (and optionally `roms/` subdirectories) for ROM
   files.
3. For each ROM, determine the art source by file extension:

   | Extension | Source   | Naming Convention     |
   |-----------|----------|-----------------------|
   | `.nds`    | GameTDB  | `{TID}.png`           |
   | all others| libretro | `{romfilename}.png`   |

4. **NDS games (GameTDB path):**
   a. Read 4 bytes at ROM offset 0x0C to get the title ID.
   b. Download from `https://art.gametdb.com/ds/coverS/US/{TID}.png`.
      The `coverS` path returns 128x115 images -- already the right
      size. Do NOT use `/cover/` (returns 404 for most titles).
   c. Save as `{TID}.png` in the boxart directory.

5. **Retro games (libretro path):**
   a. Compute the SHA1 hash of the ROM.
   b. Call `lookup_nointro` to get the canonical No-Intro name and
      system string.
   c. URL-encode the No-Intro name (spaces to `%20`, commas to
      `%2C`, parentheses to `%28`/`%29`).
   d. Download from:
      `https://thumbnails.libretro.com/{system}/Named_Boxarts/{encoded_name}.png`
   e. The libretro images are 512x512. Resize to 128x115 using
      `resize_image`.
   f. Save as `{romfilename}.png` (matching the cleaned filename
      on the card, with the extension included, e.g.
      `Tetris.gb.png`).

6. Run `clean_dot_files` on the boxart directory.
7. Report which games got art and which did not (homebrew and
   obscure titles may not have covers).

**Outputs:** PNG cover art for each ROM, correctly named and sized.

**Key details discovered in the real session:**
- GameTDB URL must use `coverS` not `cover`.
- libretro thumbnail URLs use the full No-Intro canonical name (with
  all the region/revision tags), even though we cleaned the filename.
  The SHA1 lookup is how we recover the canonical name.
- TWiLight recognizes two naming conventions: title ID for NDS, and
  `{filename}.png` for everything else. Both must be correct or the
  art simply does not appear.
- Images over 44 KB will not display in cached box art mode. The
  128x115 resize from GameTDB is already under this limit. The
  libretro resize usually lands around 28-34 KB.

---

## Skill 8: flashcart_cleanup

**Purpose:** Remove macOS junk files and verify the card layout.

**Inputs:** Mount point.

**Procedure:**

1. Run `clean_dot_files` on the entire card.
2. Verify expected directory structure exists:
   - `__rpg/` or `_wfwd/` (base kernel, depends on cart family)
   - `_nds/TWiLightMenu/` (if TWiLight is installed)
   - `Games/`
   - Boot files (`_DS_MENU.DAT` / `_DSMENU.DAT` or equivalent)
3. Report any unexpected files at root level.
4. Report total card usage and free space.

**Outputs:** A clean card with verified structure.

**Problem this solves:** Every file copy on macOS to a FAT32 volume
creates AppleDouble `._*` resource fork files. These show up as
ghost entries in the flashcart menu. In the real session I had to
call `clean_dot_files` four separate times. Ideally this runs
automatically after every card-writing operation, but as a fallback
it should be easy to invoke as a standalone pass.

---

## Cross-Cutting Concerns

### State Passing

Every skill after `flashcart_detect` needs the mount point. Every
skill after `flashcart_identify` needs the cart family. These two
values should be established early and threaded through the entire
chain. A simple approach: store them in a dict or a pair of
variables that each skill receives as input.

### AppleDouble Cleanup

The `clean_dot_files` call appears in every skill that writes to the
card. Two options:

1. **Post-hook:** Have each skill call `clean_dot_files` on the card
   after it finishes writing. This is what I did manually.
2. **Wrap copy_file:** Create a `copy_to_card` helper that calls
   `copy_file` followed by `clean_dot_files` on the destination.
   This prevents junk from accumulating between skills.

Option 2 is cleaner but more invasive. Option 1 is simpler to
implement.

### Error Recovery

If a download fails (CDN drops, 404, etc.), the skill should:
1. Retry once.
2. If still failing, try a mirror URL if one is known.
3. If no mirrors, report the failure clearly and do not leave
   partial files on the card.

### Idempotency

Each skill should be safe to run multiple times. If the kernel is
already installed, `flashcart_init` should overwrite cleanly. If
box art already exists for a game, `flashcart_boxart` should skip
it (or offer to re-download).

---

## Full Setup Sequence

For a user starting from scratch with an empty card:

```
flashcart_detect        -- find and validate the SD card
flashcart_identify      -- identify the cart from photos
flashcart_init          -- install the base kernel
flashcart_twilight_install -- install TWiLight Menu++
flashcart_emulators     -- install Virtual Console emulators
flashcart_add_game      -- (repeat for each ROM)
flashcart_boxart        -- download cover art for all games
flashcart_cleanup       -- final cleanup pass
```

A user who already has a working card and just wants to add a game
would run:

```
flashcart_detect
flashcart_add_game
flashcart_boxart
flashcart_cleanup
```
