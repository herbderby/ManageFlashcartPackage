# DS Flashcart Visual Identification Guide

**Purpose:** This document is a machine-readable reference for identifying Nintendo DS
flashcarts from photographs. It is designed to be consumed by an LLM (e.g., Claude Code)
that assists users in identifying their carts and selecting the correct kernel. The emphasis
is on **visual and physical markers** that disambiguate carts whose labels are deliberately
misleading.

**Date compiled:** 2026-03-02

**Primary sources:**

- https://www.flashcarts.net/ds-quick-start-guide
- https://www.flashcarts.net/ysmenu-compat-ext
- https://sanrax.github.io/flashcart-guides/cart-guides/r4i-sdhc/
- https://sanrax.github.io/flashcart-guides/cart-guides/ace3ds_r4ils/
- https://wiki.gbatemp.net/wiki/Ultimate_Flashcart_Download_Index
- https://www.reddit.com/r/flashcarts/comments/rdl55f/things_i_look_at_when_i_identify_flashcarts/

---

## 1. Why Labels Cannot Be Trusted

The DS flashcart market is dominated by clones of clones. Manufacturers routinely
re-label one hardware design with stickers that claim to be a completely different product.
The most common case today is **Ace3DS+/R4iLS clones sold in shells labeled as DSTTi
DEMON carts** (e.g., "R4 SDHC Gold Pro," "R4 SDHC Dual-Core," "R4 SDHC RTS Lite").

**Rule of thumb:** Never trust the front label alone. Always inspect the back of the cart
for PCB color, shell geometry, and other physical markers described below.

---

## 2. The Major Hardware Families (Currently Sold)

There are three hardware families you will encounter in carts sold today. Each requires
a different kernel. Installing the wrong kernel can range from "it just won't boot" to
a **permanent, unrecoverable brick**.

### 2.1 Ace3DS+ / R4iLS Clone Family

**Internal hardware:** Original R4 DS derivative with SDHC support.
**Kernel:** Ace Wood R4 v1.62 (or AceOS, or Pico-Launcher for Ace3DS+).
**Setup guide:** https://sanrax.github.io/flashcart-guides/cart-guides/ace3ds_r4ils/
**Kernel download:** https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip

**Known label variants (all use the same kernel):**

- Ace3DS+ (original, often unlabeled or with an "Ace3DS+" sticker)
- Ace3DS X (has a physical ntrboot switch)
- r4isdhc.com.cn (2020+)
- r4isdhc.hk (2020+ with star-shaped year outline, silver label, "Smart Update" text)
- "208-in-1" or "XXX-in-1" multicarts
- Unlabeled carts with a red PCB
- **Relabeled DEMON shells** (labeled as r4isdhc.com Gold Pro, Dual-Core, RTS Lite, etc.)

**Boot error on wrong kernel:** `Can't open _dsmenu.dat` or `Couldn't find _DS_MENU.DAT`

### 2.2 DSTTi DEMON Clone Family (AKA "Timebomb Carts")

**Internal hardware:** DSTTi clone with custom DEMON firmware. Expects an encrypted
`R4.dat` boot file.
**Kernel (primary):** YSMenu 7.06 (DSTTi-Clone YSMenu folder + DEMON bootstrap R4.dat).
**Kernel (alternate, r4isdhc.com only):** R4iMenu v4.3b (timebomb removed).
**Kernel (alternate, r4i-sdhc.com only):** R4iMenu v1.87b (timebomb removed).
**Setup guide:** https://sanrax.github.io/flashcart-guides/cart-guides/r4i-sdhc/

**Known label variants (all use the same kernel family):**

- r4isdhc.com 2014+ (Gold Pro, RTS Lite, Dual-Core, Snoopy, Upgrade -- no hardware
  difference between these, only label/color)
- r4i-sdhc.com (various models)
- r4i-gold.eu
- r4i-gold.cc (R4i Gold 3DS variant)
- Various other DEMON-firmware carts listed at https://www.flashcarts.net/ysmenu-compat-ext#demon-dstti-clones

**Boot error on wrong kernel:** Shows `MENU?` on screen with an empty or incorrectly
set up SD card. Will not show `_dsmenu.dat` errors (that indicates Ace3DS+ hardware).

### 2.3 DSpico

**Internal hardware:** Open-source RP2040-based design by LNH Team. First open-source
DS(i) flashcart with full DSi mode support.
**Kernel:** Pico-Launcher + Pico-Loader.
**Setup guide:** https://sanrax.github.io/flashcart-guides/cart-guides/dspico/

**Visual identification:** Distinctive shell (3D-printed or injection-moulded depending on
seller). Has a visible USB port (micro-USB or USB-C) and a development port. Cannot be
confused with the R4-style carts above.

---

## 3. Visual Identification Decision Tree

When presented with a photo of a DS flashcart, work through these checks in order.
The first two checks (PCB color and shell indent pattern) are the most reliable
discriminators between Ace3DS+ clones and genuine DEMON carts.

### 3.1 PCB Color (Most Reliable Single Indicator)

Examine the **back** of the cart. The PCB (printed circuit board) is often visible
through the shell plastic, especially through translucent or thin white shells, and
is always visible through the contact window at the bottom edge.

| PCB Color | Hardware Family | Confidence |
|-----------|----------------|------------|
| **Red** | Ace3DS+ / R4iLS clone | **Very high** |
| **Green or dark green** | Genuine DEMON (DSTTi clone) **or** original DSTT | High |
| **Blue** | Various (some DSTTi clones, some R4iTT) | Medium -- needs further checks |
| **Black** | Various older carts | Medium -- needs further checks |

**Key rule:** If the PCB is red, it is almost certainly an Ace3DS+/R4iLS clone regardless
of what the front label says. This is the single most reliable visual indicator.

### 3.2 Shell Indent Pattern (Second Most Reliable)

Examine the **sides** of the cart shell. DS flashcarts have small indentations along the
left and right edges (these help with grip and are part of the mold design).

| Indent Pattern | Hardware Family | Notes |
|----------------|----------------|-------|
| **Shorter, deeper indents** | Ace3DS+ / R4iLS clone | Indents are more pronounced, more closely spaced |
| **Longer, shallower indents** | Genuine DEMON cart | Indents are more subtle, spread further apart |

**Reference image:** https://www.flashcarts.net/assets/images/ds_carts/demonr4ils.png
(Side-by-side comparison of DEMON vs R4iLS shell indents from flashcarts.net)

### 3.3 Shell Quality and Fit

| Observation | Likely Hardware |
|-------------|----------------|
| Higher quality shell, tight fit, smooth finish | Genuine DEMON cart |
| Lower quality shell, loose fit, rough edges, may not seat well in console | Ace3DS+ / R4iLS clone |

This is a softer indicator and harder to assess from photos alone, but noticeable
differences in shell finish and seam quality can help.

### 3.4 Front Label Analysis

Labels are the **least reliable** identifier due to deliberate mislabeling, but they
still provide useful context when combined with the physical indicators above.

#### 3.4.1 Labels That Suggest Ace3DS+ / R4iLS Clone

- "r4isdhc.com.cn" (note the `.cn` -- different from `r4isdhc.com`)
- "r4isdhc.hk" with year 2020+ inside a **star-shaped** outline
- "XXX-in-1" (e.g., "208 in 1")
- No label at all (bare white or colored shell)
- Silver label with "Smart Update" text (on r4isdhc.hk variants)

#### 3.4.2 Labels That Suggest Genuine DEMON Cart (But Verify with PCB/Shell)

- "r4isdhc.com" (no `.cn`, no `.hk`) with 2014+ branding
- "r4i-sdhc.com"
- "r4i-gold.eu"
- Branding as "Gold Pro," "RTS Lite," "Dual-Core," "Snoopy," or "Upgrade"
- "Dual-Core / SMART UPDATE / 3DS DSi" on a brushed-metal label

**Critical caveat:** All of the DEMON-style labels in section 3.4.2 are now commonly
found on relabeled Ace3DS+ clones. You **must** verify with PCB color and shell indents.

#### 3.4.3 Labels That Suggest r4isdhc.hk -- Distinguishing Clone From DEMON

Carts labeled r4isdhc.hk can be either Ace3DS+ clones or DEMON carts depending on
the era and specific variant:

- **Ace3DS+ clone indicators:** Star-shaped year outline, silver label, "Smart Update"
  text, NO "Gold Pro" / "Real Time Save" / "RTS LITE" text, red PCB.
- **DEMON cart indicators:** Ribbon-style year outline, gold or brushed metal label,
  "Gold Pro" / "RTS" text present, green PCB.

### 3.5 URL on Label

If a URL is printed on the cart label, it is the single most important piece of text
on the label (more important than the product name). Map the URL to a hardware family:

| URL on Label | Hardware Family | Kernel |
|-------------|----------------|--------|
| r4isdhc.com.cn | Ace3DS+ / R4iLS clone | Ace Wood R4 1.62 |
| r4isdhc.hk (2020+, star year) | Ace3DS+ / R4iLS clone | Ace Wood R4 1.62 |
| r4xmp.com | Ace3DS+ / R4iLS clone | Ace Wood R4 1.62 |
| ace3ds.com | Ace3DS+ / R4iLS clone | Ace Wood R4 1.62 |
| r4isdhc.com (2014+) | DEMON (DSTTi clone) | R4iMenu 4.3b or YSMenu |
| r4i-sdhc.com | DEMON (DSTTi clone) | R4iMenu 1.87b or YSMenu |
| r4i-gold.eu | DEMON (DSTTi clone) | YSMenu only |
| ndstt.com / ndstt.net | Original DSTT | YSMenu (DSTTi-Clone folder) |
| r4ds.com | Original R4 | Wood R4 1.62 |
| (no URL) | **Ambiguous** -- must use physical indicators | Depends on hardware |

### 3.6 Software-Based Identification (Boot Behavior)

If physical inspection is inconclusive, the boot behavior with a blank or incorrectly
configured SD card can disambiguate:

| Boot Behavior | Hardware Family |
|---------------|----------------|
| `Can't open _dsmenu.dat` | Ace3DS+ / R4iLS clone |
| `Couldn't find _DS_MENU.DAT` | Ace3DS+ / R4iLS clone |
| `MENU?` on screen | DEMON (DSTTi clone) or original DSTT |
| White screen (no output) | Could be various -- may indicate V2.10T/V2.20T DEMON variant missing `_DS_MENU.DAT` |
| Ace3DS-style menu interface | Ace3DS+ / R4iLS clone (regardless of label) |

---

## 4. Complete Identification Procedure

Given one or more photos of a flashcart, apply these steps in order:

```
1. CHECK FOR DSPICO
   - Does it have a USB port (micro-USB or USB-C)?
   - Does it have a distinctive non-R4-style shell?
   → If yes: DSpico. Kernel = Pico-Launcher. DONE.

2. CHECK PCB COLOR
   - Look at the back of the cart, especially through the contact window
     or through translucent shell plastic.
   → Red PCB: Ace3DS+/R4iLS clone. Kernel = Ace Wood R4 1.62. DONE.
   → Green/dark PCB: Proceed to step 3.
   → Cannot determine: Proceed to step 3.

3. CHECK SHELL INDENT PATTERN
   - Compare the side indentations to reference images.
   → Shorter, deeper indents: Ace3DS+/R4iLS clone. Kernel = Ace Wood R4 1.62. DONE.
   → Longer, shallower indents: Likely genuine DEMON. Proceed to step 4.
   → Cannot determine: Proceed to step 4.

4. READ LABEL URL
   - Is there a URL printed on the label?
   → r4isdhc.com.cn, r4isdhc.hk (2020+/star), r4xmp.com, ace3ds.com:
     Ace3DS+/R4iLS clone. Kernel = Ace Wood R4 1.62. DONE.
   → r4isdhc.com (2014+, not .cn, not .hk):
     DEMON cart. Kernel = R4iMenu 4.3b or YSMenu. DONE.
   → r4i-sdhc.com:
     DEMON cart. Kernel = R4iMenu 1.87b or YSMenu. DONE.
   → r4i-gold.eu or other DEMON-listed URL:
     DEMON cart. Kernel = YSMenu only. DONE.
   → No URL: Proceed to step 5.

5. ANALYZE LABEL TEXT AND STYLE
   - "XXX-in-1" → Ace3DS+/R4iLS clone.
   - No label at all → Likely Ace3DS+/R4iLS clone (check PCB).
   - r4isdhc.hk with ribbon-style year + "Gold Pro"/"RTS" → Likely DEMON.
   - r4isdhc.hk with star-shaped year + "Smart Update" → Ace3DS+/R4iLS clone.
   → If still ambiguous: Proceed to step 6.

6. SOFTWARE TEST
   - Insert a blank SD card and power on the cart.
   → "_dsmenu.dat" error: Ace3DS+/R4iLS clone.
   → "MENU?" screen: DEMON or DSTT.
   → White screen: Check for V2.10T/V2.20T variant or damaged cart.
```

---

## 5. Kernel Quick Reference

### 5.1 Ace3DS+ / R4iLS Clone Family

| Kernel | Download | Notes |
|--------|----------|-------|
| Ace Wood R4 v1.62 | https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip | Recommended. Known AP issue with Pokemon B/W 2 (use cheat bypass). |
| AceOS 2.13 | https://github.com/flashcarts/AOS/#setup | Wood R4 + bundled emulators (GBARunner2/3, GameYob, NesDS, etc.) |
| Pico-Launcher + Pico-Loader | https://picoarchive.cdn.blobfrii.com/pico_package_ACE3DS.zip?picoloader=v1.5.0&picolauncher=v1.1.0&fcnetrev=0 | Modern UI. No cheats or soft-reset. |

**SD card layout for Ace Wood R4 1.62:**
```
SD root/
├── __rpg/          (kernel system files)
│   └── cheats/
│       └── usrcheat.dat   (optional cheat database)
├── _DS_MENU.DAT
├── Games/          (place .nds ROMs here)
└── (other kernel files)
```

### 5.2 DSTTi DEMON Clone Family

| Kernel | Download | Compatibility |
|--------|----------|---------------|
| R4iMenu v4.3b (r4isdhc.com) | https://archive.flashcarts.net/r4isdhc.com/r4isdhc.com_2014-and-above_DEMON_4.3.zip | r4isdhc.com 2014+ only. Has RTS. No timebomb. |
| R4iMenu v1.87b (r4i-sdhc.com) | https://archive.flashcarts.net/r4i-sdhc.com/r4i-sdhc.com_DEMON_1.87b.zip | r4i-sdhc.com only. Has RTS. No timebomb. |
| YSMenu 7.06 | https://gbatemp.net/download/retrogamefan-multi-cart-update.35737/download | All DEMON carts. Use DSTTi-Clone YSMenu folder + DEMON bootstrap R4.dat. |
| DEMON bootstrap R4.dat | https://archive.flashcarts.net/YSMenu/DEMON_common/R4.dat | Required for YSMenu on DEMON carts. |
| Pico-Launcher + Pico-Loader | https://picoarchive.cdn.blobfrii.com/pico_package_DSTT.zip?picoloader=v1.5.0&picolauncher=v1.1.0&fcnetrev=0 | All DEMON/DSTT carts. No cheats or soft-reset. |

**SD card layout for YSMenu on DEMON carts:**
```
SD root/
├── Games/          (place .nds ROMs here)
├── TTMenu/         (YSMenu system files)
├── TTMenu.dat
└── R4.dat          (DEMON bootstrap binary -- NOT the same as a normal R4.dat)
```

**SD card layout for R4iMenu (r4isdhc.com v4.3b):**
```
SD root/
├── Games/          (place .nds ROMs here)
├── _rpg/           (kernel system files)
├── _dsmenu.dat
└── R4.dat
```

---

## 6. Dangerous Misidentifications (Brick Risks)

The following kernel/cart mismatches can cause **permanent, unrecoverable bricking**:

| Cart | Dangerous Kernel | Result |
|------|-----------------|--------|
| R4 DS Pro (r4dspro.com) | YSMenu | **PERMANENT BRICK** |
| R4iTT without screw on back | DSTTi-Clone YSMenu | **PERMANENT BRICK** |
| R4i Gold 3DS RTS with year/"PRO" (r4i-gold.com) | DSTTi-Clone YSMenu | **PERMANENT BRICK** |
| R4i3D with year/"NEW" (r4i3d.com) | DSTTi-Clone YSMenu | **PERMANENT BRICK** |
| R4i SDHC Silver RTS Lite with year/"NEW" (r4isdhc.com) | DSTTi-Clone YSMenu renaming | **PERMANENT BRICK** |
| Ace3DS+ clone | YSMenu | Will not boot (no brick, but won't work) |
| DEMON cart | Ace Wood R4 | Will not boot (no brick, but won't work) |

**General safety rule:** If you are unsure what hardware you have, try the Ace Wood R4
kernel first. If the cart is actually a DEMON cart, it simply won't boot (no damage).
The reverse is also safe: loading DEMON kernels on an Ace3DS+ clone won't brick it
either, it just won't work. The dangerous cases are the specific carts listed above
where YSMenu writes to the flash chip.

---

## 7. Case Study: R4 SDHC Dual-Core PLUS (No URL)

This section documents a real-world identification performed on 2026-03-02.

**Front label text:** "R4 SDHC / 3DS DSi / Dual-Core / SMART UPDATE" with an orange
"PLUS" banner in the upper-right corner. Brushed-metal label. No URL printed anywhere.

**Initial (incorrect) assessment based on label alone:** DSTTi DEMON clone from the
r4isdhc.com 2014+ family. The "Dual-Core" branding and "SMART UPDATE" text are
standard DEMON-family identifiers.

**Physical inspection (back of cart):**
1. **PCB color:** Red, clearly visible through the translucent white shell. This is the
   definitive indicator of an Ace3DS+/R4iLS clone.
2. **Shell indents:** Shorter, deeper indents on the sides -- consistent with Ace3DS+/R4iLS
   clone shell molds, not the longer/shallower DEMON shell.

**Correct identification:** Ace3DS+/R4iLS clone, relabeled with a DEMON-style sticker.

**Correct kernel:** Ace Wood R4 v1.62.

**Lesson:** The absence of a URL on the label, combined with "Dual-Core / SMART UPDATE"
text, initially pointed toward a DEMON cart. But the red PCB and shell indent pattern
overrode the label-based identification. **Physical markers always take precedence over
label text.**

---

## 8. Reference Images

The following URLs contain reference images useful for visual comparison:

- **DEMON vs R4iLS shell comparison:** https://www.flashcarts.net/assets/images/ds_carts/demonr4ils.png
- **Genuine DEMON cart (gold label):** https://www.flashcarts.net/assets/images/ds_carts/r4isdhc_com_front.png
- **Genuine DEMON cart (back, showing indent pattern):** https://www.flashcarts.net/assets/images/ds_carts/r4isdhc_com_back.png
- **Unlabeled DEMON clone (back):** https://www.flashcarts.net/assets/images/ds_carts/timebomb_back.png
- **r4isdhc.com.cn R4iLS clone:** https://www.flashcarts.net/assets/images/ds_carts/r4isdhc_com_cn.png
- **r4isdhc.hk R4iLS clone:** https://www.flashcarts.net/assets/images/ds_carts/r4isdhc_hk.png
- **208-in-1 Ace3DS+ clone:** https://www.flashcarts.net/assets/images/ds_carts/208in1.png
- **Unlabeled Ace3DS+ (red PCB):** https://www.flashcarts.net/assets/images/ds_carts/ace3ds-nolabel.png
- **Ace3DS X (with ntrboot switch):** https://www.flashcarts.net/assets/images/ds_carts/ace3dsx.png
- **DSpico:** https://www.flashcarts.net/assets/images/ds_carts/dspico.png
- **Original R4 DS clone (grey, green PCB "ROHS Card 7a"):** Described in flashcarts.net quick start guide.
- **R4iLS original:** https://sanrax.github.io/flashcart-guides/images/r4ils.png

---

## 9. Older / Less Common Hardware Families

These are less likely to be encountered in new purchases but may appear in used carts
or existing collections.

### 9.1 Original R4 DS (r4ds.com) and M3 Simply

**Kernel:** Wood R4 v1.62 (official, not Ace fork).
**Download:** https://archive.flashcarts.net/R4_original_M3_Simply/R4DS_Wood_R4_1.62.zip
**Limitation:** Supports only 2 GB or smaller microSD cards. Does not work on stock DSi/3DS.
**Visual ID:** Grey unlabeled shell, green PCB with "ROHS Card 7a" text.

### 9.2 Original DSTT / DSTTi (ndstt.com)

**Kernel:** YSMenu 7.06 (DSTTi-Clone YSMenu folder, NO DEMON R4.dat needed).
**Visual ID:** Various shell colors. If an empty SD produces a `MENU?` screen, likely DSTT-based.
**Note:** Stock TTMenu kernel had malware that bricked clone carts. Always use YSMenu or
RetroGameFan's repacked TTMenu.

### 9.3 Original R4SDHC (r4sdhc.com)

**Kernel:** R4-Clone YSMenu + `_DS_MENU.DAT` bootstrap.
**Note:** Some OG R4SDHC carts are actually DEMON clones internally. If it shows `MENU?`
on an empty SD, use DEMON YSMenu instead.
**Limitation:** Unstable with SD cards larger than 4 GB.

### 9.4 Acekard 2/2i and Clones

**Kernel:** AKAIO or BL2CK.
**Warning:** Do NOT use official AKAIO v1.9.0 on clones (may brick).

### 9.5 R4i Gold 3DS Plus (r4ids.cn)

**Kernel:** Wood R4 v1.64 (last batch faulty, use BL2CK or TWiLight Menu++ instead).
**Note:** Production halted 2020. Last batch cannot play NDS ROMs on stock kernel.

### 9.6 EZ-Flash Parallel

**Visual ID:** Distinctive stylish shell design, wider than typical DS carts.
**Known issues:** Slow SD speeds, sleep mode broken, fitment issues.
**Not recommended** due to GPL violations and hardware issues.

### 9.7 Stargate 3DS

**Visual ID:** 3DS cart form factor (thicker than DS carts), has both DS and 3DS modes.
**Note:** Only usable on 3DS (not DS/DS Lite/DSi due to physical shape).

---

## 10. SD Card Formatting Notes

All flashcart families work best with FAT32-formatted SD cards. On macOS, use
Disk Utility or the `diskutil` command to format:

```fish
# Find the disk identifier (e.g., disk4)
diskutil list

# Format as FAT32 with MBR partition scheme
sudo diskutil eraseDisk FAT32 NDS MBRFormat /dev/diskN
```

After writing files to a FAT32 volume on macOS, clean up AppleDouble resource fork
files (`._*` files) which can confuse some flashcart kernels:

```fish
# Using dot_clean (built into macOS)
dot_clean /Volumes/NDS

# Or using the Flashcart Tools MCP clean_dot_files tool
```

**SD card size notes:**
- Original R4 DS: 2 GB maximum
- Original R4SDHC: 4 GB recommended maximum (unstable above)
- All other currently-sold carts: 32 GB FAT32 recommended. Larger cards work but
  must be formatted as FAT32 (not exFAT).

---

## 11. Confidence Levels for Image-Based Identification

When identifying a flashcart from photographs, assign confidence based on which
indicators are visible:

| Indicators Available | Confidence Level | Recommendation |
|---------------------|-----------------|----------------|
| PCB color + shell indents + label URL | **Very high** | Proceed with kernel setup |
| PCB color + shell indents (no URL) | **High** | Proceed with kernel setup |
| PCB color only | **High** (if red) / **Medium** (if green) | Proceed if red; verify with other indicators if green |
| Shell indents only | **Medium** | Ask for back-of-cart photo if possible |
| Label text only | **Low** | Ask for back-of-cart photo; warn about mislabeling |
| No photo available | **Very low** | Request photos; suggest blank-SD boot test |

---

## 12. Summary Decision Matrix

| PCB Color | Shell Indents | Label URL | → Hardware | → Kernel |
|-----------|--------------|-----------|-----------|---------|
| Red | Short/deep | Any | Ace3DS+/R4iLS | Ace Wood R4 1.62 |
| Red | Any/unknown | Any | Ace3DS+/R4iLS | Ace Wood R4 1.62 |
| Green | Long/shallow | r4isdhc.com | DEMON | R4iMenu 4.3b or YSMenu |
| Green | Long/shallow | r4i-sdhc.com | DEMON | R4iMenu 1.87b or YSMenu |
| Green | Long/shallow | r4i-gold.eu | DEMON | YSMenu only |
| Green | Long/shallow | (other DEMON URL) | DEMON | YSMenu + DEMON R4.dat |
| Green | Long/shallow | (no URL) | Likely DEMON | YSMenu + DEMON R4.dat |
| Green | Short/deep | Any | **Contradiction** | Ask for more photos; likely relabeled |
| Unknown | Unknown | r4isdhc.com.cn | Ace3DS+/R4iLS | Ace Wood R4 1.62 |
| Unknown | Unknown | r4isdhc.hk (2020+/star) | Ace3DS+/R4iLS | Ace Wood R4 1.62 |
| Unknown | Unknown | (no URL, DEMON-style label) | **Ambiguous** | Request back photo; suggest blank-SD test |
| N/A | N/A | Has USB port | DSpico | Pico-Launcher |
