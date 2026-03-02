package main

import "fmt"

// FlashcartModel describes a Nintendo DS flashcart model with the
// information needed to select the correct kernel, autoboot files,
// and flashcart loader from the TWiLight Menu++ archive.
type FlashcartModel struct {
	ID            string   // internal key, e.g. "ace3ds_plus"
	DisplayName   string   // e.g. "Ace3DS+"
	Aliases       []string // other names users might say
	VisualHints   []string // URLs/text printed on the cart sticker
	KernelType    string   // "wood_r4" or "ysmenu"
	KernelURL     string   // download URL for base kernel archive
	KernelArchive string   // filename for /tmp/ download
	AutobootDir   string   // TWiLight Autoboot/ subdirectory name
	LoaderDir     string   // TWiLight "Flashcart Loader/" subdirectory name
	ForwarderFile string   // e.g. "Wfwd.dat" (empty if none)
	Supported     bool     // full workflow support implemented?
	Notes         string   // model-specific notes for the user
}

// flashcartModels maps model IDs to their definitions. Supported
// models have complete workflow prompts; recognized models can be
// identified but lack automated setup.
var flashcartModels = map[string]*FlashcartModel{
	// --- R4/R4SDHC family (Wood R4 kernel) -- supported ---

	"ace3ds_plus": {
		ID:            "ace3ds_plus",
		DisplayName:   "Ace3DS+",
		Aliases:       []string{"Ace3DS Plus", "ace3ds"},
		VisualHints:   []string{"ace3ds.com", "ace3ds+"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		KernelArchive: "Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		AutobootDir:   "Ace3DS+",
		LoaderDir:     "Ace3DS+",
		ForwarderFile: "Wfwd.dat",
		Supported:     true,
		Notes:         "Most common R4 clone. Uses Wood R4 1.62 kernel.",
	},
	"r4ils": {
		ID:            "r4ils",
		DisplayName:   "R4iLS",
		Aliases:       []string{"R4i LS", "R4iLS"},
		VisualHints:   []string{"r4ils.com"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		KernelArchive: "Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		AutobootDir:   "R4iLS",
		LoaderDir:     "R4iLS",
		ForwarderFile: "Wfwd.dat",
		Supported:     true,
		Notes:         "Same kernel as Ace3DS+, different autoboot/loader directories.",
	},
	"gateway_blue": {
		ID:            "gateway_blue",
		DisplayName:   "Gateway Blue",
		Aliases:       []string{"Gateway Blue Card"},
		VisualHints:   []string{"gateway-3ds.com"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/Gateway_Blue/Gateway_Blue_Wood_R4_1.62.zip",
		KernelArchive: "Gateway_Blue_Wood_R4_1.62.zip",
		AutobootDir:   "Gateway Blue",
		LoaderDir:     "Gateway Blue",
		ForwarderFile: "Wfwd.dat",
		Supported:     false,
		Notes:         "Uses Wood R4 kernel. Not yet field-tested.",
	},
	"r4_original": {
		ID:            "r4_original",
		DisplayName:   "Original R4",
		Aliases:       []string{"R4 DS", "R4DS"},
		VisualHints:   []string{"r4ds.com", "r4ds.cn"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/R4DS/R4DS_Wood_R4_1.62.zip",
		KernelArchive: "R4DS_Wood_R4_1.62.zip",
		AutobootDir:   "Original R4",
		LoaderDir:     "Original R4",
		ForwarderFile: "Wfwd.dat",
		Supported:     false,
		Notes:         "FAT16 only, max 2 GB SD card. Not yet field-tested.",
	},
	"r4sdhc": {
		ID:            "r4sdhc",
		DisplayName:   "R4 SDHC",
		Aliases:       []string{"R4 SDHC Dual-Core", "R4SDHC"},
		VisualHints:   []string{"r4sdhc.com"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/R4_SDHC/R4_SDHC_Wood_R4_1.62.zip",
		KernelArchive: "R4_SDHC_Wood_R4_1.62.zip",
		AutobootDir:   "R4 SDHC",
		LoaderDir:     "R4 SDHC",
		ForwarderFile: "Wfwd.dat",
		Supported:     false,
		Notes:         "Not yet field-tested.",
	},
	"ex4ds": {
		ID:            "ex4ds",
		DisplayName:   "EX4DS",
		Aliases:       []string{"EX4DS SDHC"},
		VisualHints:   []string{"ex4ds.com"},
		KernelType:    "wood_r4",
		KernelURL:     "https://archive.flashcarts.net/Ace3DS+_R4iLS/Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		KernelArchive: "Ace3DS+_R4iLS_Wood_R4_1.62.zip",
		AutobootDir:   "EX4DS",
		LoaderDir:     "EX4DS",
		ForwarderFile: "Wfwd.dat",
		Supported:     false,
		Notes:         "Ace3DS+ clone. Not yet field-tested.",
	},

	// --- DSTT family (YSMenu kernel) -- recognized, not supported ---

	"dstt": {
		ID:          "dstt",
		DisplayName: "DSTT",
		Aliases:     []string{"DSTTi", "DSTT/DSTTi"},
		VisualHints: []string{"dstt.net", "ndstt.com", "ndstt.net"},
		KernelType:  "ysmenu",
		Supported:   false,
		Notes:       "Uses YSMenu kernel instead of Wood R4. Not yet supported.",
	},
	"r4i_sdhc_demon": {
		ID:          "r4i_sdhc_demon",
		DisplayName: "R4i-SDHC DEMON",
		Aliases:     []string{"R4i SDHC 2014+", "DSTTi DEMON-HW"},
		VisualHints: []string{"r4isdhc.com", "r4isdhc.com 2014"},
		KernelType:  "ysmenu",
		Supported:   false,
		Notes:       "DSTTi DEMON-HW clone. Uses YSMenu. Not yet supported.",
	},

	// --- Other families -- recognized, not supported ---

	"acekard_2i": {
		ID:          "acekard_2i",
		DisplayName: "Acekard 2i",
		Aliases:     []string{"AK2i"},
		VisualHints: []string{"acekard.com"},
		KernelType:  "akaio",
		Supported:   false,
		Notes:       "Uses AKAIO kernel. Not yet supported.",
	},
	"r4i_gold": {
		ID:          "r4i_gold",
		DisplayName: "R4i Gold 3DS Plus",
		Aliases:     []string{"R4i Gold", "R4i Gold 3DS+"},
		VisualHints: []string{"r4ids.cn"},
		KernelType:  "wood_r4",
		Supported:   false,
		Notes:       "Uses Wood R4 kernel variant. Not yet field-tested.",
	},
	"supercard_dsone": {
		ID:          "supercard_dsone",
		DisplayName: "SuperCard DSONE",
		Aliases:     []string{"DSONE", "SCDS1"},
		VisualHints: []string{"supercard.sc"},
		KernelType:  "evolution_os",
		Supported:   false,
		Notes:       "Uses Evolution OS kernel. Not yet supported.",
	},
}

// lookupModel returns the FlashcartModel for the given ID, or an
// error if the model is not in the registry.
func lookupModel(id string) (*FlashcartModel, error) {
	m, ok := flashcartModels[id]
	if !ok {
		return nil, fmt.Errorf("unknown flashcart model %q", id)
	}
	return m, nil
}

// modelListText builds a formatted text listing of all known
// flashcart models for use in the identification prompt. Models
// are grouped by support status.
func modelListText() string {
	var supported, recognized []string
	for _, m := range sortedModels() {
		line := fmt.Sprintf("- **%s** (id: `%s`): sticker text %s. %s",
			m.DisplayName, m.ID, formatHints(m.VisualHints), m.Notes)
		if m.Supported {
			supported = append(supported, line)
		} else {
			recognized = append(recognized, line)
		}
	}

	text := "### Fully Supported (automated setup available)\n\n"
	for _, s := range supported {
		text += s + "\n"
	}
	text += "\n### Recognized (identification only, no automated setup yet)\n\n"
	for _, s := range recognized {
		text += s + "\n"
	}
	return text
}

// sortedModels returns all models sorted by ID for stable output.
func sortedModels() []*FlashcartModel {
	keys := make([]string, 0, len(flashcartModels))
	for k := range flashcartModels {
		keys = append(keys, k)
	}
	// Simple insertion sort; the list is small.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	models := make([]*FlashcartModel, len(keys))
	for i, k := range keys {
		models[i] = flashcartModels[k]
	}
	return models
}

// formatHints joins visual hints into a quoted, comma-separated string.
func formatHints(hints []string) string {
	if len(hints) == 0 {
		return "(none)"
	}
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = `"` + h + `"`
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ", " + parts[i]
	}
	return result
}
