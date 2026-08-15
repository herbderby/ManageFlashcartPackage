---
name: reference-mcpb-packaging
description: Constraints on `.mcpb` extension manifests, Cowork plugin limitations as of March 2026, and the relevant filesystem locations on Claude Desktop.
metadata:
  type: reference
---

# MCPB Packaging, Cowork, and Claude Desktop

## MCPB manifest (`dxt_version: "0.1"`)

- **`prompts` array is NOT supported** in `dxt_version: "0.1"`
  manifests. Claude Desktop rejects with
  `Invalid manifest: prompts: Required` — which actually means
  each prompt entry is missing required fields the schema
  expects, not that prompts are required. Workaround: **omit**
  `prompts` from the manifest entirely and serve them
  dynamically via `prompts/list` and `prompts/get` from the
  server itself.
- **`tools` array works fine** in the manifest. Static tool
  declarations are accepted.
- **Re-installing a `.mcpb`** overwrites the existing extension at
  `~/Library/Application Support/Claude/Claude Extensions/local.mcpb.herb-derby.flashcart-tools/`.
  No uninstall step needed; the new install replaces the old.

## Cowork plugin limitations (as of March 2026)

- Plugin `.mcp.json` format uses **flat keys** — no `mcpServers`
  wrapper object. Different from Claude Desktop's
  `claude_desktop_config.json` which does use the wrapper.
- **Bundled binary plugins are NOT launched by Cowork** (known bug
  as of March 2026). Only **system PATH commands** (`npx`, `uvx`,
  `php`) get spawned correctly.
- `claude --plugin-dir <path>` works fine for local testing in
  Claude Code — confirms a plugin works before debugging why
  Cowork won't run it.
- `claude_desktop_config.json` with absolute paths is a fallback
  pathway for users who can't get Cowork to launch the binary.

## Claude Desktop filesystem locations

| Resource | Path |
| --- | --- |
| Desktop config | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| MCP log | `~/Library/Logs/Claude/mcp.log` |
| Plugin session cache | `~/Library/Application Support/Claude/local-agent-mode-sessions/` |
| Installed extensions | `~/Library/Application Support/Claude/Claude Extensions/` |

MCPB testing in this project has been done in **Claude Chat**
(the desktop app), not Cowork, because of the binary-launch bug
above.
