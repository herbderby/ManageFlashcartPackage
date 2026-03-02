# MCPB Desktop Extension Format

Research compiled 2026-03-01. Sources listed at end.

## What Is a .mcpb File?

A `.mcpb` file is a **ZIP archive** with a renamed extension. You can
rename it to `.zip` and extract it with any unzip tool. The format is
analogous to Chrome's `.crx` or VS Code's `.vsix`.

The format was originally called `.dxt` ("Desktop Extensions") and
developed internally at Anthropic. In late 2025, Anthropic open-sourced
the specification under the Model Context Protocol organization on
GitHub, renaming to `.mcpb` (MCP Bundle). Existing `.dxt` files
continue to work.

**Spec version:** `manifest_version: "0.3"` (updated 2025-12-02).
However, both extensions installed on this Mac use `"dxt_version": "0.1"`,
suggesting Claude Desktop accepts the older format as well.

**CLI toolchain:** `@anthropic-ai/mcpb` v2.1.2 (npm). This is a
build-time-only tool. Since `.mcpb` is just a ZIP, we can build
without it using `zip`.

## Minimum Bundle Structure

```
my-extension.mcpb (ZIP)
 +-- manifest.json          # REQUIRED
 +-- server/
 |    +-- my-binary         # compiled executable
 +-- icon.png               # optional, 512x512 recommended, PNG
```

## manifest.json Schema

### Required Fields

| Field              | Type   | Description                          |
|--------------------|--------|--------------------------------------|
| `manifest_version` | string | `"0.3"` (or `dxt_version: "0.1"`)   |
| `name`             | string | Machine-readable identifier          |
| `version`          | string | Semver (e.g. `"0.1.0"`)             |
| `description`      | string | Brief explanation                    |
| `author`           | object | `{ "name": "...", "url": "..." }`    |
| `server`           | object | Server configuration (see below)     |

### Server Configuration

```json
"server": {
  "type": "binary",
  "entry_point": "server/my-tool",
  "mcp_config": {
    "command": "${__dirname}/server/my-tool",
    "args": [],
    "env": {}
  }
}
```

**Server types:** `"node"`, `"python"`, `"binary"`, `"uv"` (experimental).

For `"type": "binary"`:
- `entry_point` is relative to the bundle root
- Claude Desktop auto-appends `.exe` on Windows
- All shared libraries must be included (static linking preferred)
- Go with `CGO_ENABLED=0` is ideal

**Variable substitution in `mcp_config`:**
- `${__dirname}` -- full path to the unpacked extension directory
- `${user_config.fieldname}` -- user-supplied configuration values
- `${HOME}` -- home directory

**Platform-specific overrides:**
```json
"mcp_config": {
  "command": "${__dirname}/server/my-tool",
  "args": [],
  "platforms": {
    "win32": {
      "command": "${__dirname}/server/my-tool.exe"
    }
  }
}
```

### Optional Fields

| Field              | Type    | Description                                |
|--------------------|---------|--------------------------------------------|
| `display_name`     | string  | User-friendly name shown in UI             |
| `long_description` | string  | Detailed markdown description              |
| `icon`             | string  | Icon asset path                            |
| `repository`       | object  | `{ "type": "git", "url": "..." }`         |
| `homepage`         | string  | URL                                        |
| `tools`            | array   | `[{ "name": "...", "description": "..." }]` |
| `tools_generated`  | boolean | If true, tools discovered at runtime       |
| `prompts`          | array   | Static prompt declarations                 |
| `prompts_generated`| boolean | If true, prompts discovered at runtime     |
| `keywords`         | array   | Search terms                               |
| `license`          | string  | License identifier                         |
| `compatibility`    | object  | Platform/version requirements              |
| `user_config`      | object  | Schema for user-provided config values     |

### Compatibility Section

```json
"compatibility": {
  "claude_desktop": ">=0.10.0",
  "platforms": ["darwin"],
}
```

Platform values: `"darwin"` (macOS), `"win32"` (Windows), `"linux"` (Linux).
No architecture field exists in the spec.

### User Configuration

```json
"user_config": {
  "api_key": {
    "type": "string",
    "title": "API Key",
    "required": true,
    "sensitive": true
  }
}
```

Field types: `string`, `number`, `boolean`, `directory`, `file`.
`sensitive: true` stores values in the OS keychain.

## How Binary Extensions Launch

1. User installs `.mcpb` (double-click, drag-and-drop, or menu)
2. Claude Desktop validates manifest and prompts for `user_config`
3. ZIP unpacked to internal extensions directory
4. `${__dirname}` set to the unpacked path
5. Desktop reads `server.mcp_config`, substitutes variables
6. **Binary spawned as subprocess**, communicates over **stdio** (JSON-RPC)
7. No restart required after install

Extensions unpack to:
`~/Library/Application Support/Claude/Claude Extensions/<extension-id>/`

## Platforms and Architecture

- macOS (`darwin`) and Windows (`win32`) are supported by Claude Desktop
- Linux is valid in the spec but Claude Desktop has no Linux release
- **No architecture field.** The spec is arch-agnostic
- For Apple Silicon Macs, ship an `arm64` binary (or universal via `lipo`)
- Go does not produce universal binaries natively; build both and combine

Known issue: ARM64 binary replaced with x86_64 during install on Apple
Silicon (GitHub issue anthropics/claude-code#13617).

## Resources, Prompts, and Tools

**Tools:** Declared statically in manifest AND/OR discovered dynamically
via `tools/list`. Use `tools_generated: true` for dynamic discovery.

**Prompts:** Same pattern. Static in manifest or dynamic via
`prompts/list` and `prompts/get`. Support argument substitution
(`${arguments.key}`).

**Resources:** No manifest declaration. Only available dynamically
through the running MCP server (`resources/list`, `resources/read`).

## Building Without npm

Since `.mcpb` is a ZIP, the packaging workflow is:

```bash
# Build the Go binary
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
  go build -ldflags="-s -w" -o ext/server/flashcart-tools

# Create manifest.json in ext/

# Pack as .mcpb (just a zip)
cd ext && zip -r ../flashcart-tools.mcpb manifest.json server/
```

The `mcpb` CLI adds convenience (`mcpb init`, `mcpb validate`,
`mcpb sign`) but is not required for basic packaging.

Auto-excluded by the CLI: `.git`, `node_modules/.cache`, `.DS_Store`,
`*.log`, `.env.local`. When building manually, just don't include them.

## Example Binary manifest.json

```json
{
  "manifest_version": "0.3",
  "name": "flashcart-tools",
  "display_name": "Flashcart Tools",
  "version": "0.1.0",
  "description": "Manage Nintendo DS flashcart SD cards",
  "author": {
    "name": "Herb Derby",
    "url": "https://github.com/herbderby"
  },
  "repository": {
    "type": "git",
    "url": "https://github.com/herbderby/ManageFlashcartPackage"
  },
  "server": {
    "type": "binary",
    "entry_point": "server/flashcart-tools",
    "mcp_config": {
      "command": "${__dirname}/server/flashcart-tools",
      "args": [],
      "env": {}
    }
  },
  "tools": [
    {
      "name": "list_volumes",
      "description": "List mounted volumes with filesystem details"
    }
  ],
  "compatibility": {
    "claude_desktop": ">=0.10.0",
    "platforms": ["darwin"]
  },
  "license": "MIT"
}
```

## Local Observations

Two extensions currently installed on this Mac:

1. **Apple Notes** (`ant.dir.ant.anthropic.notes`) -- by Anthropic,
   `dxt_version: "0.1"`, type `node`
2. **osascript** (`ant.dir.gh.k6l3.osascript`) -- by Kenneth Lien,
   `dxt_version: "0.1"`, type `node`

Both use `dxt_version: "0.1"`, not `manifest_version: "0.3"`.
Neither is a binary extension. The extension ID format is
`ant.dir.<source>.<author>.<name>`.

Extensions directory: `~/Library/Application Support/Claude/Claude Extensions/`
Installation metadata: `~/Library/Application Support/Claude/extensions-installations.json`

## Open Questions

1. **`dxt_version` vs `manifest_version`:** Which should we use?
   The installed extensions use `dxt_version: "0.1"`. The GitHub
   spec says `manifest_version: "0.3"`. Need to test which Claude
   Desktop currently accepts. Safest: try `dxt_version: "0.1"` first
   since it's proven to work on this Mac.

2. **Binary extension validation:** No binary extensions are installed
   locally, so the `"type": "binary"` path is untested on this Mac.
   The spec says it works; needs real validation.

3. **Architecture bug:** The reported ARM64-to-x86_64 replacement
   (issue #13617) could affect us. Monitor.

## Sources

- [Anthropic Engineering Blog: Desktop Extensions](https://www.anthropic.com/engineering/desktop-extensions)
- [MCP Blog: Adopting the MCP Bundle format](https://blog.modelcontextprotocol.io/posts/2025-11-20-adopting-mcpb/)
- [GitHub: modelcontextprotocol/mcpb](https://github.com/modelcontextprotocol/mcpb)
  - [MANIFEST.md specification](https://github.com/modelcontextprotocol/mcpb/blob/main/MANIFEST.md)
  - [CLI.md](https://github.com/modelcontextprotocol/mcpb/blob/main/CLI.md)
- [MCPBundles: What is an .mcpb file?](https://www.mcpbundles.com/docs/concepts/mcpb-files)
- [Claude Support: Building Desktop Extensions](https://support.claude.com/en/articles/12922929-building-desktop-extensions-with-mcpb)
- [GitHub issue: ARM64 binary replaced](https://github.com/anthropics/claude-code/issues/13617)
- Local inspection: `~/Library/Application Support/Claude/Claude Extensions/`
- Local inspection: `~/Library/Application Support/Claude/extensions-installations.json`
