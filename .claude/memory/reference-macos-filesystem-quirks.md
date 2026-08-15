---
name: reference-macos-filesystem-quirks
description: Darwin Statfs field types, the macOS firmlinks-on-FAT32 trap with os.MkdirAll, and the AppleDouple companion-file rule for FAT32 cards.
metadata:
  type: reference
---

# macOS Filesystem Quirks

## Darwin volume detection via `syscall.Statfs_t`

- **`Fstypename` is `[16]int8`** — Go's `syscall` package types
  it as `int8` rather than `byte`. Converting to a Go string
  requires an int8-to-byte pass.
- **`Bsize` is `uint32`.**
- **`Blocks` and `Bavail` are `uint64`.**
- Total bytes: `uint64(Bsize) * Blocks`
- Free bytes: `uint64(Bsize) * Bavail`

The `uint64` cast on `Bsize` is required even though the product
fits in 32 bits for small volumes — Go does not promote the
multiplication automatically.

## Firmlinks + FAT32 trap with `os.MkdirAll`

- **`/Volumes` on Big Sur and later is a firmlink** to a
  synthetic data volume, not a real directory.
- **`os.MkdirAll` fails with `ENOTDIR`** because its stat-based
  path walk gets confused by firmlink resolution when crossing
  the `/Volumes/*` boundary onto a FAT32 SD card.
- **Fix: `resolvedMkdirAll`** in `pathutil.go` resolves the
  existing portion of the path via `filepath.EvalSymlinks` first,
  then reconstructs the full target path on the resolved side.
- **`os.Create` works fine** — it goes straight to the `creat()`
  syscall without an intermediate stat-based path walk, so the
  firmlink confusion doesn't apply.

## AppleDouple on FAT32

Every macOS write to a FAT32 volume creates a `._*` companion
file alongside the real file. These hold the extended attributes
and resource forks that FAT32 has no native concept for.

- **Field test:** 66 AppleDouple files generated from a single
  kernel extraction onto a flashcart.
- **Rule:** `clean_dot_files` must be called after **every** write
  batch onto the card. Forgetting one batch leaves orphan
  AppleDouples that confuse the cart's firmware.
- **Also skip these macOS-system directories** when listing or
  copying:
  - `.Trashes`
  - `.Spotlight-V100`
  - `.fseventsd`
