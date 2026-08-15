---
name: reference-nointro-database
description: How the embedded No-Intro hash database is shaped and generated, the Myrient DAT prefix list it consumes, and the parallel flashcart-model registry used by tool prompts.
metadata:
  type: reference
---

# Embedded No-Intro Database

## File and shape

- **Filename:** `nointro.json.gz` (~828 KB).
- **Schema:** `SHA1 → [ext, name]` — flat map from a ROM's SHA1
  hex string to a two-element tuple of its filename extension
  and its canonical name.
- **Coverage:** ~21,000 ROMs across supported non-NDS systems.

## Loading

- Embedded into the binary via `//go:embed`.
- Lazy-loaded on first access using `sync.Once` in `nointro.go`.
- A single source-of-truth `consoleToLibretro` map (also in
  `nointro.go`) translates No-Intro console names into libretro
  thumbnail-set names.

## Generator

- Source: `tools/gen_nointro.go`. Run with
  `go run ./tools` (or via `make gen-nointro`).
- Inputs:
  - `~/Desktop/Games/NoIntro.db` — gzipped JSON of all No-Intro
    DATs.
  - Myrient DATs at `/Volumes/Media/myrient/DATs/`.
- Generator prefers the **main DAT** over Aftermarket / Private
  variants when both exist for a system.

## Myrient DAT prefixes the generator consumes

The exact prefix strings, **with trailing open-paren** (this is
the disambiguator that distinguishes "Atari 2600" from "Atari
2600 Homebrew"):

- `Atari - Atari 2600 (`
- `Atari - Atari 5200 (`
- `Atari - Atari 7800 (BIN) (` — the A78 DAT itself is empty;
  use the BIN variant.
- `NEC - PC Engine - TurboGrafx-16 (`
- `Coleco - ColecoVision (`
- `Mattel - Intellivision (`
- `SNK - NeoGeo Pocket (` — **no space** before the paren on
  this one. Different from the others.
- `Bandai - WonderSwan (`
- `Nintendo - Pokemon Mini (`

## Flashcart model registry (`models.go`)

Parallel registry, lives alongside the No-Intro DB but is **not**
generated:

- `FlashcartModel` struct, `flashcartModels` map, `lookupModel()`
  function.
- **Supported models** (full workflow, Wood R4 kernel):
  - `ace3ds_plus`
  - `r4ils`
- **Recognized models** (identifiable but no automated workflow):
  - `gateway_blue`, `r4_original`, `r4sdhc`, `ex4ds`, `dstt`,
    `r4i_sdhc_demon`, `acekard_2i`, `r4i_gold`, `supercard_dsone`
- **`substituteModel()` in `skill.go`** replaces these placeholders
  in prompt templates:
  `{{model_name}}`, `{{kernel_url}}`, `{{kernel_archive}}`,
  `{{autoboot_dir}}`, `{{loader_dir}}`, `{{forwarder_file}}`.
- **`promptError()` returns the error as prompt content**, not as a
  protocol-level error. This keeps the chat UX usable when a
  lookup fails — the user sees a sentence explaining the problem
  rather than an opaque MCP error.
