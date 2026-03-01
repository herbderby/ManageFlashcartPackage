---
title: "Go MCP SDK Options for Building MCP Servers"
sources:
  - url: "https://modelcontextprotocol.io"
    description: "Official MCP specification and documentation"
  - url: "https://github.com/modelcontextprotocol"
    description: "MCP GitHub organization"
  - url: "https://github.com/modelcontextprotocol/go-sdk"
    description: "Official Go SDK for MCP (maintained with Google)"
  - url: "https://github.com/mark3labs/mcp-go"
    description: "Popular community Go MCP implementation"
  - url: "https://github.com/metoro-io/mcp-golang"
    description: "Metoro's minimalist Go MCP implementation"
  - url: "https://modelcontextprotocol.io/specification/2025-06-18/basic/transports"
    description: "MCP Transport specifications"
date: "2026-03-01"
---

## Executive Summary

The Model Context Protocol (MCP) is an open protocol enabling seamless integration between LLM applications and external data sources and tools. Go has become a strong choice for implementing MCP servers due to its compiled binaries, minimal dependencies, and excellent concurrency support. This document surveys available Go SDK options, protocol architecture, and implementation considerations.

---

## Table of Contents

1. [MCP Protocol Overview](#mcp-protocol-overview)
2. [Available Go SDK Options](#available-go-sdk-options)
3. [SDK Comparison](#sdk-comparison)
4. [Protocol Architecture Deep Dive](#protocol-architecture-deep-dive)
5. [Implementation Approaches](#implementation-approaches)
6. [Dependencies and Binary Size](#dependencies-and-binary-size)
7. [Recommendations](#recommendations)

---

## MCP Protocol Overview

### What is MCP?

The Model Context Protocol is a standardized way for AI applications to connect with external data sources and tools. It follows a client-server architecture where:

- **Servers** expose capabilities (tools, resources, prompts) through the MCP protocol
- **Clients** (LLM applications) discover and invoke these capabilities
- Communication happens via JSON-RPC 2.0 messages

### Core Concepts

#### 1. **JSON-RPC 2.0 Foundation**

MCP is fundamentally JSON-RPC 2.0 sent over various transports. A typical message structure includes:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 1
}
```

Message types:
- **Request**: Contains `jsonrpc`, `id`, `method`, and optional `params`
- **Response**: Contains `jsonrpc`, `id`, and either `result` or `error`
- **Notification**: Contains `jsonrpc`, `method`, and optional `params` (no `id`)

#### 2. **Transport Layer**

MCP supports multiple transport mechanisms:

| Transport | Use Case | Advantages |
|-----------|----------|-----------|
| **STDIO** | Local integration | Simple, no network, OS-level sandboxing |
| **HTTP + SSE** | Remote servers | Stateless, cloud-deployable |
| **Streamable HTTP** | Streaming responses | Long-lived connections |

The **STDIO transport** is most common for local tools:
- Client launches server as subprocess
- Client writes JSON-RPC messages to server's stdin
- Server writes responses to stdout
- Messages are newline-delimited

### MCP Handshake Protocol

All MCP servers follow this initialization sequence:

```
1. Client sends "initialize" request with client info and capabilities
2. Server responds with capabilities (tools, resources, prompts)
3. Client sends "initialized" notification
4. Server is now ready to handle tool calls and other requests
```

Example initialize request:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "clientInfo": {
      "name": "claude-desktop",
      "version": "1.0.0"
    },
    "capabilities": {}
  }
}
```

### Core Server Operations

#### **tools/list** - Advertise Available Tools

Servers must implement `tools/list` to advertise their capabilities:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

Response includes tool definitions with JSON Schema descriptions:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "get_weather",
        "description": "Get weather for a location",
        "inputSchema": {
          "type": "object",
          "properties": {
            "location": { "type": "string" },
            "unit": { "enum": ["celsius", "fahrenheit"] }
          },
          "required": ["location"]
        }
      }
    ]
  }
}
```

#### **tools/call** - Execute Tools

When the LLM selects a tool to use:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_weather",
    "arguments": {
      "location": "San Francisco",
      "unit": "celsius"
    }
  }
}
```

Server responds with execution result:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Weather in San Francisco: 18°C, partly cloudy"
      }
    ]
  }
}
```

---

## Available Go SDK Options

### 1. Official Go SDK (modelcontextprotocol/go-sdk)

**Status**: Official, maintained by MCP team in collaboration with Google

**Repository**: https://github.com/modelcontextprotocol/go-sdk

**Key Characteristics**:
- Fully typed API with Go structs
- Automatic JSON Schema generation from struct tags
- Comprehensive feature support
- Stable v1.0.0 release (compatibility guaranteed)
- Multiple built-in transports

**When to Choose**: When you want official support and a feature-complete SDK

**Example**:
```go
package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp/server"
)

func main() {
	s := server.NewServer()

	// Add a tool with automatic schema generation
	s.AddTool("get_weather", "Get weather for a location",
		func(params struct {
			Location string `json:"location"`
		}) (string, error) {
			return "Sunny, 72°F", nil
		})

	// Run over stdio
	s.Run()
}
```

### 2. mark3labs/mcp-go

**Status**: Community-maintained, highly popular, actively developed

**Repository**: https://github.com/mark3labs/mcp-go

**Key Characteristics**:
- Fast and minimal overhead
- Clean, intuitive API
- Complete MCP specification support
- Supports stdio, SSE, HTTP, and in-process transports
- Well-documented with examples
- Good community adoption

**When to Choose**: When you prefer community-driven development and want flexible transport options

**Example**:
```go
package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewServer()

	s.AddTool("get_weather", mcp.Tool{
		Description: "Get weather for a location",
		InputSchema: mcp.NewObjectSchema(map[string]interface{}{
			"location": map[string]interface{}{
				"type": "string",
			},
		}, []string{"location"}),
		Handler: func(params map[string]interface{}) interface{} {
			location := params["location"].(string)
			return "Sunny, 72°F in " + location
		},
	})

	s.Serve()
}
```

### 3. metoro-io/mcp-golang

**Status**: Community-maintained, minimalist approach

**Repository**: https://github.com/metoro-io/mcp-golang

**Key Characteristics**:
- Minimal boilerplate required
- Type-safe with native Go struct definitions
- Automatic schema generation from Go types
- Built-in stdio and HTTP transports
- Low dependencies philosophy
- Quickstart documentation at mcpgolang.com

**When to Choose**: When you want minimal code and rapid development

**Example**:
```go
package main

import (
	"github.com/metoro-io/mcp-golang"
)

type WeatherRequest struct {
	Location string `json:"location" description:"The location for weather"`
}

type WeatherResponse struct {
	Result string `json:"result"`
}

func getWeather(req WeatherRequest) (*WeatherResponse, error) {
	return &WeatherResponse{
		Result: "Sunny, 72°F in " + req.Location,
	}, nil
}

func main() {
	server := mcp_golang.NewServer(
		mcp_golang.WithTool("get_weather", getWeather),
	)
	server.Serve()
}
```

---

## SDK Comparison

| Aspect | Official SDK | mcp-go | mcp-golang |
|--------|--------------|--------|-----------|
| **Maturity** | v1.0.0 stable | Production-ready | Production-ready |
| **Support** | Official/Google | Community | Community |
| **API Style** | Functional builders | Structured | Type-driven |
| **Schema Gen** | Automatic from structs | Manual definition | Automatic from structs |
| **Transport Options** | Multiple | 4+ transports | 2+ transports |
| **Learning Curve** | Moderate | Gentle | Very gentle |
| **Code Boilerplate** | Low-Medium | Low-Medium | Very low |
| **Documentation** | Official docs | Good | Good |
| **Community Size** | Growing | Large | Growing |
| **Maintenance** | Long-term committed | Active | Active |

---

## Protocol Architecture Deep Dive

### Message Lifecycle

A complete request-response cycle for a tool call:

```
Client                                Server
  |                                     |
  |-- tools/list request (id=1) ------->|
  |                                     |
  |<------ tools/list response ---------|
  |                                     |
  |-- tools/call request (id=2) ------->|
  |   {tool: "foo", args: {...}}        |
  |                                     |
  |   [Server processes...]             |
  |                                     |
  |<------ tools/call response ---------|
  |   {content: [...]}                  |
```

### Capability Negotiation

During initialization, both client and server advertise capabilities:

```go
// Server capabilities example
{
  "capabilities": {
    "tools": {},      // Supports tools
    "resources": {},  // Supports resources
    "prompts": {},    // Supports prompt templates
    "logging": {}     // Supports logging
  }
}
```

### Error Handling

MCP uses JSON-RPC error responses:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "error": {
    "code": -32600,
    "message": "Invalid Request",
    "data": { "details": "Tool not found" }
  }
}
```

Standard error codes:
- `-32700`: Parse error
- `-32600`: Invalid Request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error
- `-32000 to -32099`: Server error (custom)

---

## Implementation Approaches

### Approach 1: Using an Official/Popular SDK (Recommended)

**Pros**:
- Handles JSON-RPC serialization automatically
- Schema validation and generation
- Type safety
- Fewer bugs to worry about
- Good error handling
- Community support

**Cons**:
- Additional dependency
- Slight startup overhead
- Less control over wire format

**Best for**: Production systems, rapid development, most new projects

**Estimated time to working server**: 10-30 minutes

### Approach 2: Hand-Rolling the Protocol

**Overview**: Implementing MCP by hand forces understanding of the entire wire protocol, message lifecycle, and schema negotiation down to the byte.

**Core components you must implement**:

1. **JSON-RPC Message Reader**
```go
type Message struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      int             `json:"id,omitempty"`
    Method  string          `json:"method,omitempty"`
    Params  json.RawMessage `json:"params,omitempty"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *ErrorObj       `json:"error,omitempty"`
}
```

2. **Request Dispatcher**
```go
switch msg.Method {
case "initialize":
    handleInitialize(msg)
case "tools/list":
    handleToolsList(msg)
case "tools/call":
    handleToolsCall(msg)
}
```

3. **Schema Generator**
```go
// Convert Go types to JSON Schema for tools/list response
func generateSchema(input interface{}) Schema {
    // ... reflect on struct fields, build schema
}
```

4. **STDIO Transport Handler**
```go
// Read newline-delimited JSON from stdin
// Parse into Message struct
// Execute handler
// Write JSON response to stdout
```

**Pros**:
- Full control over implementation
- Minimal dependencies
- Educational
- Can optimize for specific needs

**Cons**:
- Significant implementation work (~500-1000 lines)
- Must maintain MCP compatibility as protocol evolves
- Easy to introduce subtle bugs
- Reinventing well-tested code
- Not recommended for production

**Best for**: Educational purposes, understanding the protocol, very specialized requirements

**Estimated time to working server**: 2-4 hours (or more)

---

## Dependencies and Binary Size

### Go's Compilation Advantages

Go offers significant advantages over other MCP server languages:

| Aspect | Go | Python | Node.js |
|--------|----|---------|---------|
| **Binary** | Single static binary | Requires Python + venv | Requires Node runtime + node_modules |
| **Startup Time** | Milliseconds | 1-3 seconds | 500ms-2s |
| **Dependencies** | Statically linked | Runtime dependencies | Massive node_modules |
| **Deployment** | Copy executable | Install venv, requirements | Install node_modules (hundreds MB) |
| **Concurrency** | Native goroutines | GIL-limited threads | Callback/async |
| **Binary Size** | 5-20 MB typical | N/A (interpreter) | N/A (interpreter) |

### SDK Dependency Comparison

**Official SDK (modelcontextprotocol/go-sdk)**:
- Core dependencies: minimal standard library usage
- `encoding/json` for JSON handling
- `io` for transport abstractions
- Total module size: ~500KB (source)
- Compiled binary: +2-3 MB to base Go program

**mcp-go**:
- Additional dependencies for HTTP/SSE support
- Optional transports reduce required imports
- Slightly larger than official SDK
- Compiled binary: +3-4 MB to base Go program

**mcp-golang**:
- Minimal dependencies philosophy
- Focused on essentials
- Slightly smaller than mcp-go
- Compiled binary: +2-3 MB to base Go program

### Minimal Hello World Server

```go
package main

import "fmt"

// Bare minimum MCP server reading from stdin, writing to stdout
func main() {
    // Would need ~200-300 lines to implement:
    // - JSON parsing/generation
    // - RPC method dispatch
    // - Tool schema definition
    // - Basic error handling
}
```

**Bare minimum implementation size**: 300-500 lines of Go code

**With SDK**: ~50-100 lines of Go code

---

## Maturity and Maintenance Status

### Official SDK (modelcontextprotocol/go-sdk)

- **Status**: v1.0.0 stable (as of 2025)
- **Compatibility Guarantee**: Breaking changes will not be introduced
- **Maintenance**: Officially supported by MCP team + Google partnership
- **Update Frequency**: Regular updates aligned with MCP specification
- **Issue Response**: Professional maintenance standards
- **Long-term Viability**: High (official project)

### mcp-go (mark3labs)

- **Status**: Production-ready, actively maintained
- **Maintenance**: Active community maintainer
- **Update Frequency**: Regular (responsive to protocol updates)
- **Issue Response**: Community-driven but responsive
- **Long-term Viability**: Good (popular community project)
- **Maturity Level**: Well-tested in production systems

### mcp-golang (metoro-io)

- **Status**: Production-ready
- **Maintenance**: Actively maintained by Metoro.io
- **Update Frequency**: Regular
- **Issue Response**: Responsive maintainers
- **Long-term Viability**: Good (backed by a company)
- **Maturity Level**: Used in production at Metoro

---

## Minimal MCP Server Structure

### Using Official SDK

```go
package main

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp/server"
    "github.com/modelcontextprotocol/go-sdk/mcp/transport/stdio"
)

func main() {
    // Create server
    s := server.NewServer(&server.Options{
        Name: "hello-world",
        Version: "1.0.0",
    })

    // Add a simple tool
    s.AddTool("echo", "Echo back the input", func(ctx context.Context, input struct {
        Message string `json:"message"`
    }) (string, error) {
        return input.Message, nil
    })

    // Create stdio transport and run
    stdio := stdio.NewStdioTransport()
    s.Serve(stdio)
}
```

### File Structure

```
my-mcp-server/
├── go.mod
├── go.sum
├── main.go              (~40 lines)
└── tools/
    ├── echo.go          (optional: extracted tool definitions)
    └── weather.go       (optional: more tools)
```

### Build and Run

```bash
go build -o my-server main.go

# Server now ready to connect to MCP clients
./my-server
```

The resulting binary is:
- Statically compiled (no external dependencies needed)
- Fast starting (~milliseconds)
- Small (~8-15 MB depending on complexity)
- Cross-platform (rebuild for any OS)

---

## Protocol Implementation Checklist

If hand-rolling the protocol, ensure:

- [ ] **Initialization handshake** implemented correctly
- [ ] **JSON-RPC 2.0 compliance** (proper error codes, message format)
- [ ] **tools/list** method returns correct schema format
- [ ] **tools/call** method executes and returns results properly
- [ ] **Error handling** with proper JSON-RPC error responses
- [ ] **STDIO transport** with newline-delimited JSON messages
- [ ] **Message ID tracking** for request-response correlation
- [ ] **Type validation** of tool arguments against schema
- [ ] **Concurrent request handling** (goroutines for each request)
- [ ] **Graceful shutdown** on client disconnect

---

## Recommendations

### For Most Projects: Use Official SDK

**Recommendation**: `github.com/modelcontextprotocol/go-sdk`

**Rationale**:
- Stability guaranteed through v1.0.0+ compatibility promise
- Official support from MCP team and Google
- Well-documented with examples
- Fully featured
- Long-term investment protection

### For Rapid Development: Use mcp-golang

**Recommendation**: `github.com/metoro-io/mcp-golang`

**Rationale**:
- Extremely minimal boilerplate
- Type-safe through Go structs
- Automatic schema generation
- Great for prototypes and MVPs
- Excellent quickstart documentation

### For Flexibility: Use mcp-go

**Recommendation**: `github.com/mark3labs/mcp-go`

**Rationale**:
- Most transport options out of the box
- Mature community project
- Good balance of features and simplicity
- Excellent for advanced use cases

### Do NOT Hand-Roll Unless:

- You're doing this for educational purposes
- You have specific performance requirements that SDKs can't meet
- You're operating in an extremely constrained environment
- You're certain you can maintain compatibility with protocol updates

---

## Conclusion

Go is an excellent choice for implementing MCP servers due to:

1. **Compilation**: Single, stateless binary with zero runtime dependencies
2. **Performance**: Millisecond startup, excellent concurrency handling
3. **SDK Maturity**: Multiple production-ready options available
4. **Operational Simplicity**: Easy deployment, no environment setup required

For new projects, **use the official SDK** unless you have specific needs that another implementation better serves. The three major options (official SDK, mcp-go, mcp-golang) are all production-ready and well-maintained.

The MCP protocol itself is straightforward (JSON-RPC 2.0), making Go implementations simple and reliable. Most projects can have a working server in under an hour with an SDK.

---

## References

- MCP Specification: https://modelcontextprotocol.io
- Official Go SDK: https://github.com/modelcontextprotocol/go-sdk
- mcp-go: https://github.com/mark3labs/mcp-go
- mcp-golang: https://github.com/metoro-io/mcp-golang
- MCP Transport Docs: https://modelcontextprotocol.io/specification/2025-06-18/basic/transports
- JSON-RPC 2.0 Spec: https://www.jsonrpc.org/specification
