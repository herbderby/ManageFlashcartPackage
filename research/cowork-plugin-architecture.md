---
title: "Claude Cowork Plugin Architecture and Development"
sources:
  - url: "https://code.claude.com/docs/en/plugins"
    description: "Create plugins - Claude Code Docs"
  - url: "https://code.claude.com/docs/en/plugins-reference"
    description: "Plugins reference documentation"
  - url: "https://code.claude.com/docs/en/discover-plugins"
    description: "Discover and install prebuilt plugins through marketplaces"
  - url: "https://code.claude.com/docs/en/slash-commands"
    description: "Slash commands documentation"
  - url: "https://code.claude.com/docs/en/skills"
    description: "Extend Claude with skills"
  - url: "https://code.claude.com/docs/en/mcp"
    description: "Connect Claude Code to tools via MCP"
  - url: "https://support.claude.com/en/articles/13345190-get-started-with-cowork"
    description: "Getting started with Cowork"
  - url: "https://support.claude.com/en/articles/13837440-use-plugins-in-cowork"
    description: "Use plugins in Cowork"
  - url: "https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop"
    description: "Getting Started with Local MCP Servers on Claude Desktop"
  - url: "https://github.com/anthropics/knowledge-work-plugins"
    description: "Open source repository of plugins for knowledge workers"
  - url: "https://github.com/anthropics/skills"
    description: "Public repository for Agent Skills"
date: "2026-03-01"
---

# Claude Cowork Plugin Architecture and Development

## Overview

Claude Cowork plugins are extensible packages that bundle together skills, connectors, slash commands, and sub-agents into single installable units. Rather than configuring each component individually, plugins provide ready-to-go setups optimized for specific roles, teams, and companies. Each plugin bundles together domain knowledge, tool integrations, automated commands, and specialized agents to create comprehensive solutions for knowledge work.

All plugin components are file-based, making them easy to build, edit, share, and version control. The architecture is built on open standards including the Model Context Protocol (MCP) and leverages a well-defined directory structure.

## What Are Cowork Plugins?

Cowork plugins are structured packages that encode methodology, workflows, tool connections, and domain knowledge into files that Claude can understand and use. A single plugin can bundle:

- **Skills**: Reusable procedures and domain knowledge that Claude applies contextually
- **MCP Server Connectors**: Integrations to external services, APIs, and databases
- **Slash Commands**: Special commands that trigger specific workflows (e.g., `/commit`, `/deploy`)
- **Sub-agents**: Specialized agents configured for specific tasks
- **Hooks**: Shell scripts that execute automatically on specific events
- **Reference Materials**: Documentation and standards that guide Claude's decision-making

When installed, a plugin teaches Claude how to perform specific jobs for your role, applying your methodology, tool connections, and standards consistently.

## Plugin Directory Structure

Every Claude Code plugin follows a standardized directory structure:

```
plugin-name/
├── .claude-plugin/
│   └── plugin.json          # Plugin manifest (metadata & configuration)
├── commands/                # Slash commands (optional)
│   └── my-command.md
├── agents/                  # Specialized agents (optional)
│   └── code-reviewer.md
├── skills/                  # Agent Skills (optional)
│   └── my-skill/
│       ├── SKILL.md
│       ├── scripts/
│       └── references/
├── hooks/                   # Event handlers (optional)
│   └── pre-commit.sh
├── .mcp.json               # MCP server configuration (optional)
├── README.md               # Plugin documentation
└── assets/                 # Supporting files
```

### Key Directory Details

- **`.claude-plugin/`**: Required directory containing only `plugin.json`. This is the single location where manifest metadata resides.
- **`commands/`**: Directory containing markdown files defining slash commands. Each `.md` file in this directory becomes a `/command` available in Claude Code.
- **`skills/`**: Directory containing skill subdirectories. Each skill folder contains a `SKILL.md` file with frontmatter and optional `scripts/` and `references/` subdirectories.
- **`agents/`**: Directory for specialized agents that can be invoked for specific tasks.
- **`hooks/`**: Directory for shell scripts that execute on specific lifecycle events.
- **`.mcp.json`**: Optional configuration file at plugin root defining MCP server connections.

Components can reference custom paths via `plugin.json` if they're located outside default directories.

## The plugin.json Manifest

The `plugin.json` manifest file is located at `.claude-plugin/plugin.json` and defines plugin metadata and component locations. While technically optional (Claude auto-discovers components in default locations), including a manifest provides metadata, enables custom component paths, and ensures proper distribution.

### Manifest Format

```json
{
  "name": "sales-pro",
  "version": "1.0.0",
  "description": "Comprehensive sales workflow plugin with CRM integration, deal tracking, and proposal generation",
  "author": {
    "name": "Sales Team",
    "email": "sales@company.com",
    "url": "https://company.com"
  },
  "homepage": "https://github.com/company/sales-pro-plugin",
  "repository": {
    "type": "git",
    "url": "https://github.com/company/sales-pro-plugin"
  },
  "license": "MIT",
  "keywords": ["sales", "crm", "deals", "proposals"],
  "commands": "custom-commands/",
  "agents": ["agents/", "specialized-agents/"],
  "skills": "skills/",
  "hooks": "hooks/",
  "mcpServers": ["salesforce-server", "hubspot-server"]
}
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | String | Plugin identifier (kebab-case), becomes the plugin's directory name |
| `version` | String | Semantic version (e.g., "1.0.0") |
| `description` | String | Brief description of plugin functionality |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `author` | Object | Author information with `name`, `email`, `url` |
| `homepage` | String | URL to plugin documentation |
| `repository` | Object | Git repository information with `type` and `url` |
| `license` | String | License identifier (e.g., "MIT", "Apache-2.0") |
| `keywords` | Array | Tags for discoverability in marketplaces |
| `commands` | String/Array | Path(s) to custom slash commands |
| `agents` | String/Array | Path(s) to specialized agents |
| `skills` | String | Path to skills directory |
| `hooks` | String | Path to hooks directory |
| `mcpServers` | Array | References to MCP server configurations |

## MCP Server Configuration

MCP (Model Context Protocol) servers provide Claude with connections to external tools, APIs, databases, and services. Plugin MCP configuration is typically defined in a `.mcp.json` file at the plugin root or can be embedded in `plugin.json`.

### .mcp.json File Format

```json
{
  "mcpServers": {
    "salesforce": {
      "command": "node",
      "args": ["${CLAUDE_PLUGIN_ROOT}/servers/salesforce-server.js"],
      "env": {
        "SALESFORCE_KEY": "${SALESFORCE_API_KEY}",
        "PLUGIN_HOME": "${CLAUDE_PLUGIN_ROOT}"
      }
    },
    "database": {
      "command": "${CLAUDE_PLUGIN_ROOT}/bin/db-server",
      "args": ["--config", "${CLAUDE_PLUGIN_ROOT}/config.json"],
      "cwd": "${CLAUDE_PLUGIN_ROOT}",
      "env": {
        "DB_PATH": "${CLAUDE_PLUGIN_ROOT}/data"
      }
    },
    "npm-package": {
      "command": "npx",
      "args": ["@company/mcp-server", "--plugin-mode"]
    }
  }
}
```

### Configuration Elements

| Element | Description | Example |
|---------|-------------|---------|
| **command** | Executable to run (can be `npx`, `node`, `python`, or direct path) | `"node"` or `"${CLAUDE_PLUGIN_ROOT}/server.js"` |
| **args** | Arguments passed to the command | `["--config", "config.json"]` |
| **env** | Environment variables (can use `${CLAUDE_PLUGIN_ROOT}`) | `{"DB_PATH": "${CLAUDE_PLUGIN_ROOT}/data"}` |
| **cwd** | Working directory for command execution | `"${CLAUDE_PLUGIN_ROOT}"` |

### The ${CLAUDE_PLUGIN_ROOT} Variable

The `${CLAUDE_PLUGIN_ROOT}` variable is automatically substituted at runtime with the absolute path to the plugin's root directory. This enables:

- **Relative Path References**: Use relative paths instead of absolute ones
- **Portability**: Plugins work regardless of installation location
- **Configuration File Access**: Reference `${CLAUDE_PLUGIN_ROOT}/config.json`
- **Binary Execution**: Reference `${CLAUDE_PLUGIN_ROOT}/bin/server`
- **Data Directory Access**: Reference `${CLAUDE_PLUGIN_ROOT}/data`

Example expansions:
- `${CLAUDE_PLUGIN_ROOT}/servers/db-server` → `/home/user/.claude/plugins/my-plugin/servers/db-server`
- `${CLAUDE_PLUGIN_ROOT}/config.json` → `/home/user/.claude/plugins/my-plugin/config.json`

### MCP Server Lifecycle

MCP servers defined in plugins run on the **host machine** (where Claude Code/Cowork is installed), not in a remote environment. This allows servers to:

- Access local files and directories
- Control system processes
- Manage local databases
- Execute scripts with full system permissions

When Claude needs to use a tool provided by an MCP server, the server is initialized if not already running, the tool is executed, and results are returned to Claude for further processing.

### Multiple MCP Servers

A single plugin can configure multiple MCP servers, each providing different tool integrations:

```json
{
  "mcpServers": {
    "github": {
      "command": "node",
      "args": ["./servers/github-server.js"]
    },
    "database": {
      "command": "${CLAUDE_PLUGIN_ROOT}/bin/db",
      "args": ["--port", "9000"]
    },
    "file-system": {
      "command": "python",
      "args": ["${CLAUDE_PLUGIN_ROOT}/servers/fs_server.py"]
    }
  }
}
```

Each server operates independently and provides tools accessible to Claude throughout the plugin session.

## Skills: Domain Knowledge and Procedures

Skills are a core plugin component that bundles procedural knowledge and domain expertise. Each skill consists of a `SKILL.md` file with YAML frontmatter and markdown content, optionally accompanied by supporting scripts and reference materials.

### SKILL.md Structure

```markdown
---
name: code-review
description: |
  Conduct thorough code reviews following our code standards, security best
  practices, and performance guidelines. Checks for bugs, style violations,
  architectural issues, and suggests improvements.
disable-model-invocation: false
user-invocable: true
---

# Code Review Skill

## Overview
Performs comprehensive code reviews using our company standards...

## Code Review Checklist
1. Security vulnerabilities
2. Performance issues
3. Code style compliance
4. Documentation completeness
5. Test coverage

## Standard Patterns
- Use async/await instead of callbacks
- Implement error handling with try-catch
- Add JSDoc comments to public functions

## Anti-patterns to Watch For
- Missing null checks
- Inefficient database queries
...
```

### Frontmatter Metadata

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `name` | String | Skill identifier; becomes the slash command name (e.g., `/code-review`) | Required |
| `description` | String | Brief description Claude uses to determine when to invoke the skill | Required |
| `disable-model-invocation` | Boolean | If `true`, skill can only be invoked by user (`/name` command), not automatically by Claude | `false` |
| `user-invocable` | Boolean | If `false`, skill is only for Claude's use, not available as a slash command | `true` |

### Skill Organization

The recommended directory structure for a skill with supporting materials:

```
my-skill/
├── SKILL.md              # Main skill with frontmatter and instructions
├── scripts/
│   ├── analyze.py
│   ├── validate.sh
│   └── generate-report.js
└── references/
    ├── company-standards.md
    ├── api-reference.md
    └── best-practices.md
```

### How Skills Work

1. **Discovery**: When Claude starts, it reads available skill descriptions
2. **Context Matching**: Claude evaluates whether a skill's description matches the current task
3. **Loading**: If relevant, Claude loads the full SKILL.md content into context
4. **Progressive Disclosure**: If instructions reference other files, Claude reads them on-demand using bash commands
5. **Application**: Claude applies the skill instructions to the current work

### Size and Performance

Skills should be **kept under 500 lines** in SKILL.md itself. Detailed reference material should be in separate files within the `references/` directory. This enables progressive disclosure: Claude loads information in stages as needed, rather than loading everything upfront.

### Reference Content

Reference content in skills provides domain knowledge that Claude applies inline to current work:

- **Code Style Guides**: Programming language conventions, naming patterns, file organization
- **Company Standards**: Development practices, security requirements, deployment procedures
- **API Documentation**: Endpoint specifications, authentication methods, response formats
- **Architectural Patterns**: Design patterns your company standardizes on, technology decisions
- **Process Workflows**: How to handle edge cases, escalation procedures, approval workflows

When Claude loads a skill with references, it can read and apply these materials to ensure consistent, standards-compliant work.

## Slash Commands

Slash commands provide special operations that control Claude's behavior. Commands are defined as markdown files in the `commands/` directory with YAML frontmatter.

### Command Structure

```markdown
---
description: Review code changes and provide feedback
argument-hint: "[file-path] [review-type]"
allowed-tools: Bash, Read, Write, Grep
disable-model-invocation: false
---

# Code Review Command

You are conducting a code review...

## Review Criteria
1. Security: Check for vulnerabilities
2. Performance: Identify bottlenecks
3. Style: Verify consistency
...
```

### Frontmatter Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | String | Short description shown in `/help` |
| `argument-hint` | String | Optional hint showing expected arguments |
| `allowed-tools` | String | Comma-separated list of tools the command can use |
| `disable-model-invocation` | Boolean | If `true`, command only available via explicit `/name` invocation |

### Command Naming

- Files are automatically recognized as commands when placed in `commands/` directory
- Filename becomes the command name (e.g., `code-review.md` → `/code-review`)
- Commands follow the pattern `/namespace:command-name` for organization
- Commands are available immediately in the command palette after startup

### Argument Handling

Commands support dynamic arguments using placeholders:

```markdown
---
description: Create a new feature branch and set up scaffolding
argument-hint: "[feature-name] [description]"
---

# Create Feature Branch

Using arguments: $1 for feature-name, $2 for description

Create branch: feature/$1
Description: $2
```

## Installing and Managing Plugins

### Installation Methods

Plugins can be installed through multiple channels:

1. **Cowork Interface**: Open Cowork, click "Customize" → "Browse plugins", then click "Install"
2. **GitHub Repositories**: Install from any public GitHub repository
3. **Git URLs**: Install from GitLab, Bitbucket, or self-hosted git servers
4. **Local Paths**: Install from local file system directories during development
5. **Remote URLs**: Install from hosted `marketplace.json` files

### Plugin Marketplaces

Official and community marketplaces provide curated plugin collections:

- **Official Anthropic Marketplace**: Anthropic-managed directory of high-quality plugins
- **Build with Claude**: Community plugin marketplace for code and automation
- **Community Repositories**: GitHub and other repositories with plugin collections
- **Custom Marketplaces**: Organizations can host private `marketplace.json` files

### Creating Plugins in Cowork

To build a plugin directly in Cowork:

1. Switch to the "Cowork" tab in Claude Desktop
2. Click "Customize" in the left sidebar
3. Describe the plugin you want to build (even one sentence is sufficient)
4. Claude will ask questions about your workflow, tools, standards, and edge cases
5. Claude creates the plugin structure with appropriate files

### Plugin Create Assistant

Anthropic provides "Plugin Create", a specialized plugin that guides you through plugin development, helping you define skills, configure MCP servers, create commands, and structure your plugin correctly.

## Advanced Plugin Concepts

### Sub-agents

Plugins can include specialized sub-agents configured for specific tasks. These agents operate with custom instructions, tool access, and domain knowledge relevant to their specialization.

Sub-agents are useful for:
- Delegating complex multi-step tasks
- Maintaining specialized knowledge and procedures
- Handling specific workflows with custom logic
- Scaling work across multiple specialized processes

### Hooks

Hooks are shell scripts in the `hooks/` directory that execute automatically on lifecycle events:

```bash
#!/bin/bash
# hooks/pre-install.sh - Runs before plugin installation
# Pre-flight checks, dependency validation

#!/bin/bash
# hooks/post-install.sh - Runs after installation
# Setup, initialization, configuration

#!/bin/bash
# hooks/pre-update.sh - Runs before plugin updates
# Backup, validation

#!/bin/bash
# hooks/post-update.sh - Runs after updates
# Migration, re-initialization
```

Hooks enable automated setup, cleanup, and migration tasks when plugins are installed or updated.

### Plugin Dependencies

Plugins can declare dependencies on other plugins through `plugin.json`. This ensures that required base plugins are installed before the dependent plugin becomes available.

## Plugin Size and Distribution Constraints

### General Size Guidance

- **SKILL.md Files**: Keep under 500 lines; move detailed reference material to separate files
- **Total Plugin Size**: Most plugins range from 1MB to 100MB depending on included binaries and data
- **Configuration Files**: Keep JSON manifests and configuration files minimal

### Bundling Native Binaries

Plugins can bundle native binaries and executables by:

1. **Including in plugin directory**: Place binary files in a `bin/` or `servers/` subdirectory
2. **Referencing with ${CLAUDE_PLUGIN_ROOT}**: Use the plugin root variable in configuration
3. **Setting executable permissions**: Ensure binaries have execute permissions on installation
4. **Platform-specific binaries**: Include separate binaries for different platforms (macOS, Linux, Windows)

Example binary structure:

```
plugin-name/
├── .claude-plugin/plugin.json
├── bin/
│   ├── darwin/
│   │   └── my-server          # macOS binary
│   ├── linux/
│   │   └── my-server          # Linux binary
│   └── windows/
│       └── my-server.exe      # Windows binary
└── .mcp.json
```

The `.mcp.json` can reference the appropriate binary:

```json
{
  "mcpServers": {
    "my-server": {
      "command": "${CLAUDE_PLUGIN_ROOT}/bin/my-server"
    }
  }
}
```

### Plugin Portability

By using `${CLAUDE_PLUGIN_ROOT}` and relative paths:

- Plugins work on any installation path
- No hard-coded absolute paths needed
- Plugins can be easily shared and distributed
- Version control friendly (no path-specific configuration)

## MCP Server vs. Skills vs. Commands

Understanding the distinctions between these components helps in designing effective plugins:

### Skills
- **Purpose**: Encode procedural knowledge and domain expertise
- **Activation**: Context-aware; Claude loads when task matches description
- **Size**: Moderate (under 500 lines recommended)
- **Token Cost**: ~30-50 tokens per skill loaded
- **Use Case**: Best practices, standards, workflows, reference materials

### MCP Servers
- **Purpose**: Connect to external tools, APIs, and services
- **Activation**: On-demand when tools are needed
- **Size**: Can include binary executables and large datasets
- **Token Cost**: Only consumed when tools are actually used (Tool Search reduces by ~85%)
- **Use Case**: Database access, API integrations, file systems, external services

### Slash Commands
- **Purpose**: Trigger specific workflows manually
- **Activation**: Explicit user invocation via `/command`
- **Size**: Flexible, can be scripts or procedures
- **Token Cost**: Only consumed when executed
- **Use Case**: Complex multi-step operations, dev workflows, deployments

### Integration Pattern

An effective plugin typically combines:

1. **Skills** for domain knowledge and decision-making
2. **MCP Servers** for tool access
3. **Slash Commands** for triggered workflows
4. **Hooks** for lifecycle management

This layered approach provides comprehensive functionality while optimizing token usage and performance.

## Best Practices

### Plugin Development Guidelines

1. **Modular Organization**: Keep concerns separated in skills, commands, and scripts
2. **Clear Documentation**: Provide comprehensive README and inline skill descriptions
3. **Size Management**: Keep SKILL.md under 500 lines; use references for detailed content
4. **Error Handling**: Include fallback procedures and validation in commands and scripts
5. **Testing**: Test plugins in different Cowork contexts and workstreams
6. **Versioning**: Use semantic versioning and maintain a changelog

### Configuration Best Practices

1. **Use ${CLAUDE_PLUGIN_ROOT}**: Always use this variable for path references
2. **Environment Separation**: Use environment variables for credentials and configuration
3. **Clear Naming**: Use descriptive names for MCP servers, skills, and commands
4. **Minimal Manifest**: Only include necessary fields in plugin.json
5. **Documentation**: Comment complex configuration in .mcp.json files

### Distribution Best Practices

1. **Repository Setup**: Include proper .gitignore, LICENSE, and documentation
2. **Marketplace Metadata**: Ensure keywords and descriptions are accurate for discovery
3. **Version Management**: Tag releases in git for version tracking
4. **Changelog**: Maintain CHANGELOG.md for version history and updates
5. **License Selection**: Choose appropriate open-source license (MIT, Apache-2.0, etc.)

## Summary

The Claude Cowork plugin architecture provides a flexible, file-based system for extending Claude's capabilities. By combining well-defined components—manifest files, MCP servers, skills, commands, and hooks—developers can create comprehensive, reusable packages that teach Claude how to work effectively within specific domains and organizations.

The key architectural principles are:

- **File-based structure** for version control and transparency
- **Component modularity** allowing independent development and testing
- **Configuration variables** like `${CLAUDE_PLUGIN_ROOT}` enabling portability
- **Progressive disclosure** where Claude loads information on-demand
- **Open standards** based on Model Context Protocol for tool integration
- **Marketplace distribution** enabling easy sharing and installation

Whether building for individual teams, organizations, or public distribution, the plugin system provides the foundation for extending Claude's capabilities across knowledge work domains.
