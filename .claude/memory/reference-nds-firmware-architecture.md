---
name: reference-nds-firmware-architecture
description: How the lower-layer Wood R4 kernel and the upper-layer TWiLight Menu++ environment coexist on a flashcart, the chainloading mechanism, and the two distinct `globalsettings.ini` files that must agree.
metadata:
  type: reference
---

# Dual-Firmware Architecture

## Layer separation

The flashcart's firmware stack is two layers:

- **Base kernel — Wood R4 1.62.** Provides the boot files,
  DLDI driver (for SD-card access), and the NDS game loader.
  This is the layer Nintendo DS hardware boots into directly.
- **Upper layer — TWiLight Menu++.** Provides the menu UI,
  emulators (for non-DS systems like GB, NES, etc.), and the
  box-art renderer.

## Chainloading

- **Autoboot files** named `_DSMENU.dat` and `_DS_MENU.dat` sit at
  the root of the card. Wood R4 finds them on boot and **chainloads
  TWiLight Menu++** instead of staying in its own menu.
- This is why a fresh Wood R4 install needs TWiLight artifacts
  copied alongside — the autoboot files reference paths into the
  TWiLight kernel.

## TWiLight's stripped-down kernel loader

- TWiLight maintains its own stripped-down kernel loader directory
  named `_wfwd/`. This is **separate** from the main `__rpg/`
  directory that holds the active Wood R4 kernel.
- Both directories load independently depending on which layer is
  in control.

## The two `globalsettings.ini` files

This is the easy-to-miss part:

- **`__rpg/globalsettings.ini`** — Wood R4's settings.
- **`_wfwd/globalsettings.ini`** — TWiLight's settings.
- **Both** need `showHiddenFiles = 0` for hidden-file behavior to
  be consistent. Setting it in one but not the other produces a
  confusing UX where some screens hide dot files and others don't.

Any setting that should be consistent across the cart must be
written to both files.
