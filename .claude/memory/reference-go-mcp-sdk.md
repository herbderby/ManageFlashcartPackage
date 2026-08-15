---
name: reference-go-mcp-sdk
description: API shape for the `github.com/modelcontextprotocol/go-sdk` v1.4.0 — server construction, tool and prompt handlers, the JSON-schema struct tag gotcha, and in-process testing transports.
metadata:
  type: reference
---

# Go MCP SDK (`github.com/modelcontextprotocol/go-sdk` v1.4.0)

## Server construction

- `mcp.NewServer(impl, opts)` returns a server.
- `mcp.AddTool[In, Out](server, &mcp.Tool{...}, handler)` registers
  a typed tool. `In` and `Out` are user-defined structs.
- `server.AddPrompt(&mcp.Prompt{...}, handler)` registers a prompt.
- `server.Run(ctx, transport)` runs the server on a transport.

## Transport

- `&mcp.StdioTransport{}` — newline-delimited JSON over stdin/stdout.
  **Not** Content-Length-framed like the LSP-style transports.

## Handler signatures

- **Tool handler:**
  `func(ctx, *CallToolRequest, In) (*CallToolResult, Out, error)`
- **Prompt handler:**
  `func(ctx, *GetPromptRequest) (*GetPromptResult, error)`
- **No-input tools** use `struct{}` as the `In` type.

## `jsonschema` struct-tag gotcha

The `jsonschema` tag value **is** the description directly:

- **WRONG:** `jsonschema:"description=some text"` — panics with
  `tag must not begin with WORD=`.
- **RIGHT:** `jsonschema:"some text"` — the entire tag value is
  taken as the description.

## Notification draining in tests

Tools and prompts emit `notifications/tools/list_changed` and
`notifications/prompts/list_changed` when registered. Tests that
talk to the server must drain these notifications before reading
the actual responses they're testing — otherwise the next read
returns the notification frame instead of the expected reply.

## In-process testing

The SDK exposes a connected-pair transport for tests:

- `mcp.NewInMemoryTransports()` returns two transports connected
  to each other.
- **Server side:** `server.Connect(ctx, t1, nil)` →
  `*ServerSession`.
- **Client side:** `mcp.NewClient(impl, nil)` then
  `client.Connect(ctx, t2, nil)` → `*ClientSession`.
- **Client APIs to drive the server:**
  - `session.ListPrompts(ctx, nil)`
  - `session.GetPrompt(ctx, &mcp.GetPromptParams{...})`
  - `session.ListTools(ctx, nil)`
- Pattern in this repo: server setup is extracted to `newServer()`
  in `main.go` and shared between `main()` and the test file
  (`mcp_test.go`). Run with `go test ./...`.
