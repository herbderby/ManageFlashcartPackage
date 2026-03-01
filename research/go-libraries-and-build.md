---
title: "Go Libraries and Build Techniques for Plugin Development"
sources:
  - url: "https://github.com/bodgit/sevenzip"
    description: "Pure Go 7z archive library"
  - url: "https://pkg.go.dev/github.com/bodgit/sevenzip"
    description: "bodgit/sevenzip package documentation"
  - url: "https://opensource.com/article/21/1/go-cross-compiling"
    description: "Cross-compiling with Golang"
  - url: "https://freshman.tech/snippets/go/cross-compile-go-programs/"
    description: "How to cross-compile for Windows, macOS, and Linux"
  - url: "https://pkg.go.dev/golang.org/x/image"
    description: "Go extended image processing library"
  - url: "https://dev.to/aryaprakasa/the-trade-offs-of-optimizing-and-compressing-go-binaries-492d"
    description: "Trade-offs of optimizing and compressing Go binaries"
  - url: "https://words.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/"
    description: "Shrinking Go binaries techniques"
  - url: "https://goreleaser.com/customization/upx/"
    description: "UPX compression with GoReleaser"
  - url: "https://www.arp242.net/static-go.html"
    description: "Statically compiling Go programs"
  - url: "https://eli.thegreenplace.net/2024/building-static-binaries-with-go-on-linux/"
    description: "Building static binaries with Go on Linux"
date: "2026-03-01"
---

# Go Libraries and Build Techniques for Plugin Development

## Section 1: Go 7z Extraction (bodgit/sevenzip)

### Overview

The `bodgit/sevenzip` library is a pure Go reader for 7-zip archives, inspired by Go's built-in `archive/zip` package. It provides a complete implementation for handling 7-zip format files without requiring any external dependencies, C libraries, or external binaries.

### Pure Go Implementation

This is a critical advantage for cross-platform distribution and native binary compilation. Because it is pure Go (no CGO), it:

- Compiles to fully static binaries by default
- Cross-compiles trivially across platforms without C toolchain dependencies
- Maintains a single code path across all architectures (arm64, amd64, etc.)
- Does not require glibc or other system C libraries

### Compression Methods Supported

The library supports a comprehensive set of compression algorithms used in 7-zip archives:

| Compression Method | Type | Use Case |
|---|---|---|
| Copy | Uncompressed | Raw storage |
| LZMA | Dictionary-based | Default 7-zip compression |
| LZMA2 | Modern variant | Improved LZMA for modern systems |
| Deflate | Stream compression | ZIP-compatible format |
| Bzip2 | Block-based | Good compression ratios |
| LZ4 | Fast compression | Low CPU overhead |
| Brotli | Modern compression | High compression ratios |
| Zstandard | Meta-compression | Balanced speed/compression |
| BCJ / BCJ2 | Binary transforms | Executable preprocessing |
| Delta | Delta filtering | Reduces redundancy |
| ARM / PPC / SPARC | Architecture filters | Platform-specific optimization |

### Key Features

- **Archive Format Support**: Handles both uncompressed and compressed headers, password-protected archives, multi-volume archives (split into multiple files), and self-extracting archives (SFX)
- **Data Validation**: Validates CRC checksums as it parses file content, ensuring data integrity
- **File System Interface**: Implements Go's `fs.FS` interface, allowing 7-zip archives to be treated like filesystem hierarchies
- **Simple API**: Provides a familiar interface similar to `archive/zip`, making it easy for developers accustomed to Go's standard library

### API Usage Patterns

The library follows Go conventions for archive handling:

```go
// Open a 7-zip archive
reader, err := sevenzip.OpenReader("archive.7z")
if err != nil {
    // handle error
}
defer reader.Close()

// Iterate through files in archive
for _, file := range reader.File {
    // Open individual file for reading
    rc, err := file.Open()
    if err != nil {
        // handle error
    }
    defer rc.Close()
    
    // Read file contents
    // ...
}
```

### Known Limitations

- **Read-Only**: The library is designed for extraction only; it cannot create or modify 7-zip archives
- **No Compression Writing**: Cannot be used to create 7-zip files programmatically
- **Archive Format Only**: Does not implement the 7z executable format (only archive format)

### Alternative Libraries

While `bodgit/sevenzip` is the primary pure Go implementation, other approaches exist:

| Library | Type | Notes |
|---|---|---|
| `golift/xtractr` | Go wrapper | Higher-level abstraction over multiple archive formats |
| `itchio/sevenzip-go` | CGO bindings | Binds to native 7-zip library (requires external dependency) |
| System `p7zip` command | Shelling out | Invoke external 7z executable (requires tool installation) |

For cross-platform native binary distribution, `bodgit/sevenzip` is the best choice due to its pure Go implementation.

---

## Section 2: Go Cross-Compilation and Binary Distribution

### How Go Cross-Compilation Works

Go's design philosophy includes first-class support for cross-compilation. Unlike C/C++ toolchains that require platform-specific compilers and linkers, Go compiles directly to machine code for any supported platform from any build machine.

The Go compiler generates the target architecture's machine code without relying on external assemblers or platform-specific tools. This means you can build a Windows binary on macOS, a Linux ARM64 binary on Windows, or any other combination.

### GOOS and GOARCH Environment Variables

Cross-compilation is controlled through two environment variables:

| Variable | Purpose | Example Values |
|---|---|---|
| `GOOS` | Target operating system | `linux`, `windows`, `darwin` (macOS), `freebsd` |
| `GOARCH` | Target processor architecture | `amd64` (64-bit Intel), `arm64` (ARM 64-bit), `386` (32-bit Intel) |

The valid combinations are not all permutations. You can list all supported combinations with:

```bash
go tool dist list
```

### Common Build Targets for Native Distribution

For macOS and Windows distribution, these are the primary targets:

| Target | Command |
|---|---|
| **macOS ARM64** (Apple Silicon M1/M2/M3) | `GOOS=darwin GOARCH=arm64 go build -o app-macos-arm64` |
| **macOS AMD64** (Intel Macs) | `GOOS=darwin GOARCH=amd64 go build -o app-macos-amd64` |
| **Windows AMD64** | `GOOS=windows GOARCH=amd64 go build -o app-windows-amd64.exe` |
| **Linux AMD64** | `GOOS=linux GOARCH=amd64 go build -o app-linux-amd64` |
| **Linux ARM64** | `GOOS=linux GOARCH=arm64 go build -o app-linux-arm64` |

### Static Linking (Default Behavior)

Go produces static binaries by default, meaning all dependencies are compiled directly into the executable. This is one of Go's significant advantages for distribution:

- **No Runtime Dependencies**: The binary runs on any system with the correct OS/architecture, regardless of installed libraries
- **Portable**: No need for users to install Go runtime, C libraries, or other dependencies
- **Simple Deployment**: Copy single executable to target system

**Exception - CGO**: If your code uses `cgo` (calls C code) or imports packages that use cgo (like `os/user`, `net`), Go will create a dynamically linked binary. For pure Go code, static linking is automatic.

### Binary Size Reduction Techniques

Go binaries are typically larger than minimal C programs due to the included runtime (garbage collector, goroutines, reflection, etc.). However, several techniques reduce size significantly:

#### Technique 1: Stripping Debug Information

Use `go build -ldflags="-s -w"`:

- `-s` flag: Omit symbol table (function names, variable names used by debuggers)
- `-w` flag: Omit DWARF debugging information (line numbers, type information)

**Important**: Stripping does NOT remove information needed for panic stack traces. Runtime panics remain readable and useful.

**Size Impact**: ~28% reduction from a typical binary

**Example**:
```bash
go build -ldflags="-s -w" -o myapp main.go
```

#### Technique 2: UPX Compression

UPX is an executable packer that compresses the binary and decompresses it at startup.

**Size Impact**: Can reduce binary to ~15% of original size (combined with -s -w)

**Compression Modes**:
```bash
upx --best app-linux                    # Standard best compression
upx --ultra-brute app-linux            # Maximum compression (slower)
```

**Trade-offs**:

| Advantage | Disadvantage |
|---|---|
| Extremely small file size | Startup decompression overhead |
| Reduced download time | Slower application startup (15-160ms depending on size) |
| Better for storage-constrained environments | Not ideal for CLI tools requiring fast startup |
| | Each run decompresses the binary |

**When to Use**: Server applications where startup happens once, or for distribution over slow networks. Not recommended for command-line utilities.

#### Technique 3: Build Tag Conditionals

For binaries with optional features, use build tags to exclude unused code:

```bash
go build -tags=minimal -o app-small main.go
```

This allows smaller binaries when certain optional features aren't needed.

### Typical Binary Sizes

Reference sizes for common scenarios:

| Scenario | Size | Notes |
|---|---|---|
| Empty `main()` function | 863 KB | Go runtime overhead baseline |
| Simple "Hello World" program | ~1.3 MB | Minimal functionality |
| Web server with SQLite3 | 11 MB | Full-featured service |
| After `-ldflags="-s -w"` | 7.7 MB | ~30% reduction |
| After UPX compression | 2.2 MB | ~80% reduction |

**Note**: For most modern deployment scenarios, 11 MB is acceptable. The convenience of static linking and built-in runtime features (no garbage collection configuration needed, concurrency primitives, etc.) typically outweighs the size overhead.

### Multi-Architecture Distribution Pattern

A common pattern for distributing Go applications for multiple architectures:

```bash
#!/bin/bash
# build.sh - Build for multiple platforms

TARGETS=(
    "darwin/arm64"    # macOS ARM64
    "darwin/amd64"    # macOS Intel
    "windows/amd64"   # Windows
    "linux/amd64"     # Linux Intel
    "linux/arm64"     # Linux ARM
)

for target in "${TARGETS[@]}"; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    output="dist/myapp-${GOOS}-${GOARCH}"
    [ "$GOOS" = "windows" ] && output="${output}.exe"
    
    echo "Building $output..."
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w" \
        -o "$output" \
        main.go
done
```

This produces portable executables like:
- `dist/myapp-darwin-arm64`
- `dist/myapp-darwin-amd64`
- `dist/myapp-windows-amd64.exe`
- `dist/myapp-linux-amd64`
- `dist/myapp-linux-arm64`

Users can download the appropriate binary for their platform and run it directly.

### Static Linking with CGO and musl

If your code uses CGO (calls C libraries), you can still create static binaries using musl libc instead of glibc:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-linux-musl-gcc \
  go build -ldflags="-linkmode external -extldflags '-static'" \
  -o app-linux-amd64 main.go
```

This requires musl development tools but produces fully static binaries that work on any Linux distribution.

---

## Section 3: Go Image Processing

### golang.org/x/image Library Overview

The `golang.org/x/image` package is the official Go extended library for image processing. It provides supplementary functionality beyond Go's standard `image` package, including scaling, transformation, and format-specific operations.

### Supported Image Formats

Go's standard library provides format support:

| Format | Import | Features |
|---|---|---|
| PNG | `image/png` | Encode/decode, various bit depths |
| JPEG | `image/jpeg` | Encode/decode, quality control |
| GIF | `image/gif` | Animate, encode/decode |

### Image Resizing Approaches

#### Standard Library Approach

Go's standard `image` package allows basic pixel manipulation:

```go
import (
    "image"
    "image/png"
    "os"
)

// Open and decode PNG
file, _ := os.Open("input.png")
defer file.Close()
img, _ := png.Decode(file)

// Get dimensions and create new resized image
bounds := img.Bounds()
width, height := bounds.Max.X, bounds.Max.Y

// Create new RGBA image with target dimensions
newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

// Manual pixel copying or interpolation required
```

#### golang.org/x/image/draw

The `golang.org/x/image/draw` package provides interpolation algorithms:

```go
import (
    "golang.org/x/image/draw"
)

// Supported scalers (from highest to lowest quality)
draw.CatmullRom.Scale(dst, dstRect, src, srcBounds, draw.Over, nil)
draw.BiLinear.Scale(dst, dstRect, src, srcBounds, draw.Over, nil)
draw.ApproxBiLinear.Scale(dst, dstRect, src, srcBounds, draw.Over, nil)
draw.NearestNeighbor.Scale(dst, dstRect, src, srcBounds, draw.Over, nil)
```

**Scaling Algorithms**:

| Algorithm | Quality | Speed | Use Case |
|---|---|---|---|
| NearestNeighbor | Lowest (pixelated) | Fastest | Pixel art, upscaling |
| ApproxBiLinear | Good | Fast | General purpose |
| BiLinear | Better | Moderate | Smooth scaling |
| CatmullRom | Highest | Slowest | High-quality downscaling |

#### Alternative Libraries

Community packages offer more streamlined APIs:

| Package | Features |
|---|---|
| `github.com/disintegration/imaging` | High-level image operations, multiple filters |
| `github.com/nfnt/resize` | Pure Go resizing with advanced filters |

### PNG Encoding and File Size

PNG files support different encoding options that affect file size:

```go
import (
    "image/png"
)

// Simple encoding (default compression)
png.Encode(outputFile, img)

// For better compression, use manual png.Encoder
encoder := &png.Encoder{
    CompressionLevel: png.DefaultCompression,
}
encoder.Encode(outputFile, img)
```

**Compression Levels**:

| Level | Speed | File Size | Use |
|---|---|---|---|
| `DefaultCompression` | Fast | Standard | General purpose |
| `BestSpeed` | Very fast | Larger | When size doesn't matter |
| `BestCompression` | Slower | Smaller | For distribution |

**File Size Factors**:

- **Color Depth**: Grayscale (8-bit) vs full RGB (24-bit) vs RGBA (32-bit)
- **Palette Usage**: Indexed color (256 colors) reduces size for simple images
- **Filter Method**: PNG uses predictive filters that compress redundant data
- **Content Complexity**: Random pixels compress poorly; structured patterns compress well

### Typical Workflow for Image Processing

A complete example combining library usage:

```go
package main

import (
    "image"
    "image/png"
    "os"
    "golang.org/x/image/draw"
)

func main() {
    // Open source PNG
    src, _ := os.Open("source.png")
    defer src.Close()
    img, _ := png.Decode(src)
    
    // Create destination image (scaled to half size)
    bounds := img.Bounds()
    newW := bounds.Max.X / 2
    newH := bounds.Max.Y / 2
    dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
    
    // Scale with bilinear interpolation
    draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
    
    // Encode to output PNG
    out, _ := os.Create("output.png")
    defer out.Close()
    png.Encode(out, dst)
}
```

---

## Summary and Recommendations

### For 7z Archive Handling

Use `bodgit/sevenzip` when:
- You need to extract 7-zip files in a Go application
- Cross-platform distribution is important
- You want to avoid external dependencies
- Pure Go implementation is preferred

### For Binary Distribution

Best practices for multi-platform Go binaries:

1. **Always build with `GOOS` and `GOARCH`** for each target platform
2. **Use `-ldflags="-s -w"`** to reduce size by ~30%
3. **Consider UPX compression** for network distribution, but avoid for CLI tools
4. **Distribute one binary per platform** (darwin-arm64, darwin-amd64, windows-amd64, etc.)
5. **Use build scripts** to automate multi-platform compilation
6. **Test on target platforms** before releasing binaries

### For Image Processing

- Use **standard library `image/png`** for basic PNG reading/writing
- Use **`golang.org/x/image/draw`** for resizing with quality control
- Choose **BiLinear or CatmullRom** for best visual results
- Consider **community packages** if you need advanced image manipulation
- **Compress PNG files** when distributing to minimize download sizes
