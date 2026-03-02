// Package main generates nointro.json.gz by merging two data sources:
//
//  1. ~/Desktop/Games/NoIntro.db (gzipped JSON) -- GB, GBC, GBA, NES,
//     SNES, SMS, GG, Genesis entries.
//  2. Myrient No-Intro DATs (XML) -- Atari 2600/5200/7800, PC Engine,
//     ColecoVision, Intellivision, Neo Geo Pocket, WonderSwan, Pokemon Mini.
//
// The output maps SHA1 hex strings to [extension, No-Intro name] pairs.
// Run with: go run ./tools
package main

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// consoleTypeToExt maps NoIntro.db ConsoleType values to the file
// extensions used by TWiLight Menu++.
var consoleTypeToExt = map[string]string{
	"GameBoy":                          "gb",
	"GameBoyColor":                     "gbc",
	"GameBoyAdvance":                   "gba",
	"NintendoEntertainmentSystem":      "nes",
	"SuperNintendoEntertainmentSystem": "sfc",
	"SegaMasterSystem":                 "sms",
	"SegaGameGear":                     "gg",
	"SegaGenesis":                      "gen",
}

// myrientSystem describes a system whose No-Intro DAT should be
// downloaded from the Myrient archive.
type myrientSystem struct {
	prefix string // filename prefix to match in directory listing
	ext    string // file extension for ROMs of this system
}

// myrientSystems lists the systems not covered by NoIntro.db that
// TWiLight Menu++ emulates.
var myrientSystems = []myrientSystem{
	{"Atari - Atari 2600 (", "a26"},
	{"Atari - Atari 5200 (", "a52"},
	{"Atari - Atari 7800 (BIN) (", "a78"},
	{"NEC - PC Engine - TurboGrafx-16 (", "pce"},
	{"Coleco - ColecoVision (", "col"},
	{"Mattel - Intellivision (", "int"},
	{"SNK - NeoGeo Pocket (", "ngp"},
	{"Bandai - WonderSwan (", "ws"},
	{"Nintendo - Pokemon Mini (", "mini"},
}

const myrientBaseURL = "https://myrient.erista.me/dats/No-Intro/"

// noIntroDB is the path to the local NoIntro.db gzipped JSON file.
var noIntroDB = filepath.Join(os.Getenv("HOME"), "Desktop", "Games", "NoIntro.db")

// datFile represents the top-level element of a No-Intro DAT XML file.
type datFile struct {
	XMLName xml.Name  `xml:"datafile"`
	Games   []datGame `xml:"game"`
}

// datGame represents a single game entry in a No-Intro DAT.
type datGame struct {
	Name string   `xml:"name,attr"`
	ROMs []datROM `xml:"rom"`
}

// datROM represents a ROM file entry inside a game.
type datROM struct {
	Name string `xml:"name,attr"`
	SHA1 string `xml:"sha1,attr"`
}

func main() {
	db := make(map[string][2]string) // sha1 -> [ext, name]

	// Phase A: Parse local NoIntro.db.
	fmt.Println("Phase A: Parsing NoIntro.db...")
	if err := parseNoIntroDB(db); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing NoIntro.db: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  %d entries after NoIntro.db\n", len(db))

	// Phase B: Download and parse Myrient DATs.
	fmt.Println("Phase B: Downloading Myrient DATs...")
	if err := parseMyrientDATs(db); err != nil {
		fmt.Fprintf(os.Stderr, "error with Myrient DATs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  %d entries after Myrient\n", len(db))

	// Phase C: Write compressed output.
	fmt.Println("Phase C: Writing nointro.json.gz...")
	if err := writeDB(db); err != nil {
		fmt.Fprintf(os.Stderr, "error writing DB: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Done: %d entries written to nointro.json.gz\n", len(db))
}

// noIntroDBEntry matches the JSON structure of entries in NoIntro.db.
type noIntroDBEntry struct {
	ConsoleType string `json:"ConsoleType"`
	Name        string `json:"Name"`
	SHA1        string `json:"Sha1"`
}

// parseNoIntroDB reads ~/Desktop/Games/NoIntro.db (gzipped JSON array)
// and adds entries for the 8 supported console types to db.
func parseNoIntroDB(db map[string][2]string) error {
	f, err := os.Open(noIntroDB)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer r.Close()

	var entries []noIntroDBEntry
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		return fmt.Errorf("json: %w", err)
	}

	for _, e := range entries {
		ext, ok := consoleTypeToExt[e.ConsoleType]
		if !ok {
			continue
		}
		if e.SHA1 == "" || e.Name == "" {
			continue
		}
		db[strings.ToLower(e.SHA1)] = [2]string{ext, e.Name}
	}
	return nil
}

// parseMyrientDATs fetches the Myrient directory listing, finds the
// current DAT file for each system, downloads it, and parses the XML
// to add entries to db.
func parseMyrientDATs(db map[string][2]string) error {
	resp, err := http.Get(myrientBaseURL)
	if err != nil {
		return fmt.Errorf("fetch directory listing: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read listing: %w", err)
	}
	listing := string(body)

	for _, sys := range myrientSystems {
		filename := findDAT(listing, sys.prefix)
		if filename == "" {
			fmt.Printf("  warning: no DAT found for %q\n", sys.prefix)
			continue
		}

		fmt.Printf("  downloading %s...\n", filename)
		datURL := myrientBaseURL + url.PathEscape(filename)
		n, err := fetchAndParseDAT(datURL, sys.ext, db)
		if err != nil {
			fmt.Printf("  warning: %s: %v\n", sys.prefix, err)
			continue
		}
		fmt.Printf("  added %d entries for %s\n", n, sys.prefix)
	}
	return nil
}

// hrefRe matches href attributes in HTML anchor tags.
var hrefRe = regexp.MustCompile(`href="([^"]*)"`)

// findDAT searches the HTML directory listing for a .dat file whose
// name starts with prefix. It prefers the main DAT (without
// "Aftermarket" or "Private" qualifiers) over variant DATs. Returns
// the unescaped filename, or empty string if not found.
func findDAT(listing, prefix string) string {
	matches := hrefRe.FindAllStringSubmatch(listing, -1)
	var main, fallback string
	for _, m := range matches {
		href := m[1]
		decoded, err := url.PathUnescape(href)
		if err != nil {
			decoded = href
		}
		if !strings.HasPrefix(decoded, prefix) || !strings.HasSuffix(decoded, ".dat") {
			continue
		}
		// Prefer DATs without "Aftermarket" or "Private" qualifiers.
		rest := decoded[len(prefix):]
		if strings.Contains(rest, "Aftermarket") || strings.Contains(rest, "Private") {
			fallback = decoded
		} else {
			main = decoded
		}
	}
	if main != "" {
		return main
	}
	return fallback
}

// fetchAndParseDAT downloads a No-Intro DAT XML file and adds its
// ROM entries to db. Returns the number of entries added.
func fetchAndParseDAT(datURL, ext string, db map[string][2]string) (int, error) {
	resp, err := http.Get(datURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP %d for %s", resp.StatusCode, datURL)
	}

	var dat datFile
	if err := xml.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return 0, fmt.Errorf("parse XML: %w", err)
	}

	count := 0
	for _, game := range dat.Games {
		for _, rom := range game.ROMs {
			if rom.SHA1 == "" {
				continue
			}
			// Strip the file extension from the ROM name.
			name := rom.Name
			if idx := strings.LastIndex(name, "."); idx >= 0 {
				name = name[:idx]
			}
			sha1 := strings.ToLower(rom.SHA1)
			if _, exists := db[sha1]; !exists {
				db[sha1] = [2]string{ext, name}
				count++
			}
		}
	}
	return count, nil
}

// writeDB gzip-compresses the merged database as JSON and writes
// it to nointro.json.gz in the current directory.
func writeDB(db map[string][2]string) error {
	f, err := os.Create("nointro.json.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	if err := json.NewEncoder(w).Encode(db); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}
