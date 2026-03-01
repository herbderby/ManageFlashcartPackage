---
name: hello-sdcard
description: |
  Use the list_volumes MCP tool to show mounted drives, including
  SD cards and removable media. Reports volume name, mount point,
  filesystem type, total size, and free space.
user-invocable: false
---

# HelloSDCard

This plugin provides a single MCP tool for listing mounted volumes.

## Available Tool

### list_volumes

Lists all mounted volumes under /Volumes/ on macOS. Returns JSON
with these fields for each volume:

- **name** -- volume label
- **mountPoint** -- full mount path (e.g., /Volumes/SDCARD)
- **fsType** -- filesystem type (apfs, hfs, msdos for FAT32, exfat)
- **totalBytes** / **totalHuman** -- total capacity
- **freeBytes** / **freeHuman** -- available space

The tool takes no input parameters.

## Usage

When the user asks about connected drives, mounted volumes, or SD
cards, call `list_volumes` and present the results in a readable
table or list. SD cards typically appear as FAT32 (msdos) or exFAT
volumes.
