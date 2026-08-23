# ManageFlashcartPackage — Memory Index

## SDK and runtime

- [Go MCP SDK](reference-go-mcp-sdk.md) — `mcp.NewServer`, typed `AddTool`, prompt handlers, `jsonschema` tag gotcha, in-process testing transports.

## Filesystem and platform

- [macOS filesystem quirks](reference-macos-filesystem-quirks.md) — `Statfs_t` field types, firmlinks vs `os.MkdirAll`, AppleDouple `._*` files on FAT32, system dirs to skip.

## Flashcart domain

- [NDS dual-firmware architecture](reference-nds-firmware-architecture.md) — Wood R4 1.62 base, TWiLight Menu++ upper layer, autoboot chainloading, the two `globalsettings.ini` files.
- [Box-art lookup](reference-box-art-lookup.md) — GameTDB for NDS, SHA1 + libretro for non-NDS, filename-based on-disk naming, 128×115 resize.
- [Embedded No-Intro DB + model registry](reference-nointro-database.md) — `nointro.json.gz` schema and generator, Myrient DAT prefixes, `flashcartModels` map, prompt-template substitution.

## Session handoff

Step progress, commit history, and tool-count progression live in
the project's `NEXT.md` and `git log`, not here. `MEMORY.md` is an
index of durable references.
