# Tercode — Product Requirements Document

**Project:** Tercode  
**Meaning:** Terminal Code  
**Status:** Draft v1.0  
**Product type:** Open-source, local-first, agentic terminal development environment

---

## 1. Product Overview

Tercode is an open-source, provider-agnostic, highly customizable agentic coding environment for the terminal.

Its primary interface is an interactive TUI launched with:

```bash
tercode
```

The user enters natural-language requests through the primary prompt:

```text
> _
```

Tercode's AI agent can understand a project, inspect files, search code, modify files, execute commands, run tests, inspect Git state, and use external tools.

Tercode is **agent-centric rather than editor-centric**. It is not intended to become a traditional IDE where manual code editing is the primary workflow.

The product is inspired by the philosophy of highly customizable developer tools such as Hyprland, Kitty, Alacritty, and Neovim:

> **Powerful by default. Limitless by configuration.**

Tercode should provide strong functionality out of the box while allowing developers to build their own development environment from Tercode's primitives.

---

# 2. Vision

Tercode aims to become:

> **The Neovim/Hyprland philosophy applied to agentic software development.**

The core provides the capabilities. The user defines the experience.

A stock Tercode installation should be clean, functional, lightweight, and usable immediately. Advanced users should be able to customize its appearance, layout, components, keybindings, commands, agents, workflows, and behavior without modifying the Tercode source code.

Two developers using the same Tercode binary should be able to have radically different interfaces and workflows.

---

# 3. Product Philosophy

### 3.1 Provider freedom

Users should be able to bring their own AI providers, API keys, and models.

### 3.2 Agent first

The user primarily communicates with an AI development agent rather than manually editing code.

### 3.3 Local first

Project state, configuration, credentials, and collaboration should remain under user control wherever possible.

### 3.4 Customization first

Customization is a core product capability, not a cosmetic theme feature.

### 3.5 Composability

Tercode should expose reusable UI, agent, tool, command, and workflow primitives that users can compose.

### 3.6 Collaboration

Developers should be able to work together over a local network without requiring a cloud collaboration service.

### 3.7 Strong defaults

Customization must never come at the expense of usability. A fresh Tercode installation must work well without configuration.

---

# 4. Target Users

- Individual developers
- AI-assisted developers
- Terminal/Linux power users
- Open-source developers
- Developers using multiple AI providers
- Local/self-hosted AI users
- Small development teams
- Universities and laboratories
- Developers who enjoy customizing development environments

---

# 5. Core User Experience

Running:

```bash
tercode
```

launches the interactive Tercode environment.

Conceptually:

```text
┌──────────────────────────────────────────────────────────────┐
│ Tercode                                      Model / Provider │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│                    Welcome to Tercode                        │
│                                                              │
│             Your AI-powered development workspace            │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│ > _                                                          │
└──────────────────────────────────────────────────────────────┘
```

The TUI is the primary product interface.

CLI subcommands remain available for scripting and automation.

---

# 6. Agent-Centric Development

Tercode is not a traditional IDE.

The intended workflow is:

```text
User request
    ↓
Agent understands task
    ↓
Plan
    ↓
Inspect project
    ↓
Use tools
    ↓
Modify files
    ↓
Run commands/tests
    ↓
Inspect results
    ↓
Iterate
    ↓
Report result
```

The TUI may display source files, diffs, project structure, Git changes, tool output, and agent activity, but manual code editing is not the primary workflow.

---

# 7. Interactive Prompt

The default prompt is:

```text
> _
```

It must support:

- multiline input
- history
- cancellation
- autocomplete
- slash commands
- context references
- file references

Natural-language requests are sent to the agent.

Example:

```text
> Analyze this repository and explain its architecture.
```

---

# 8. Slash Command System

Application controls should use slash commands.

Example:

```text
> /
```

opens or filters a command palette.

Initial commands may include:

```text
/add-dir
/agents
/clear
/compact
/config
/context
/diff
/help
/history
/model
/provider
/review
/session
/status
/theme
/tools
```

Natural-language prompts and application commands must remain distinct.

---

# 9. CLI Layer

The interactive TUI is the primary experience, but Tercode also provides CLI commands.

Examples:

```bash
tercode --help
tercode models
tercode provider list
tercode provider test openrouter
tercode host
tercode connect <host>
```

CLI commands and TUI commands should use the same application services.

---

# 10. Agent Modes

Initial modes:

### Chat

Project-aware conversation.

### Plan

Create an implementation plan without modifying files.

### Code

Implement requested changes.

### Review

Review current changes.

### Debug

Investigate and resolve failures.

All modes use the same agent runtime.

---

# 11. Agent Runtime Requirements

The runtime must support:

- tool calls
- tool results
- streaming
- cancellation
- retries
- timeouts
- context management
- execution limits
- user approval
- structured events
- multiple sessions
- multiple agent contexts
- model fallback

The runtime must remain independent from the TUI.

---

# 12. Tool System

Initial tools:

```text
read_file
write_file
edit_file
delete_file
list_directory
search
shell
git
run_tests
run_build
```

Future tools may include:

- MCP
- browser
- databases
- Docker
- cloud services
- custom tools

Every tool must have:

- name
- description
- input schema
- output
- permissions
- execution state
- error behavior

---

# 13. Permissions and Safety

Tool permissions should support:

```text
deny
ask
allow
```

Potentially destructive shell commands and filesystem operations should require authorization according to policy.

Tercode must:

- protect credentials
- restrict workspace access
- prevent unauthorized LAN access
- sanitize logs
- treat project content as untrusted input
- provide prompt-injection safeguards
- authenticate collaboration clients

---

# 14. Provider System

Tercode must not be locked to one AI provider.

The initial provider architecture should support:

1. OpenRouter
2. Generic OpenAI-compatible APIs
3. NVIDIA NIM

Future providers may include:

- OpenAI
- Anthropic
- Google
- Ollama
- vLLM
- local inference servers
- other compatible services

A provider abstraction must support, where available:

```text
authentication
model discovery
chat
streaming
tool calling
capabilities
model testing
```

---

# 15. API Key Management

Credentials should not normally be stored as plaintext project configuration.

Supported sources:

- environment variables
- OS keychain/credential store
- interactive configuration
- external secret managers in future

Credentials must never be intentionally exposed in:

- logs
- Git
- agent context
- crash reports
- TUI output

---

# 16. Model Discovery and Testing

Users must be able to discover models:

```bash
tercode models
```

or:

```text
> /model
```

Model metadata may include:

- provider
- model ID
- context size
- input/output modalities
- tool support
- vision support
- reasoning support
- streaming support
- pricing where available

Users must be able to test a provider/model.

Testing should verify:

1. credentials
2. network connectivity
3. model availability
4. inference
5. streaming
6. tool calling where supported
7. latency

---

# 17. Context Engine

Tercode must efficiently understand repositories.

Initial context sources:

- project files
- directory structure
- Git status
- Git diff
- README/configuration
- search results
- command output
- user-provided context

Potential technologies:

- ripgrep
- Tree-sitter
- language servers
- AST/symbol indexing

The MVP does not require a vector database.

The system should avoid sending entire repositories to the model unnecessarily.

---

# 18. Multiple Agent Contexts

Tercode should support multiple simultaneous contexts.

Example:

```text
Tercode
├── Context 1 → Backend Agent
├── Context 2 → Frontend Agent
└── Context 3 → Test Agent
```

Contexts may have independent:

- conversations
- tasks
- models
- providers
- permissions
- agents

The TUI should provide ways to switch between them.

---

# 19. Customization Philosophy

Customization is a first-class architectural and product requirement.

Tercode should not merely provide a list of themes.

Instead, it should expose primitives that users can compose.

Conceptually:

```text
Tercode Core
    ↓
UI / Agent / Workspace primitives
    ↓
User configuration
    ↓
User-defined environment
```

The goal is similar to the customization philosophy of Hyprland and Neovim.

---

# 20. Customizable UI

Tercode should allow customization of:

### Appearance

- colors
- symbols
- borders
- typography
- spacing
- animations where supported
- information density

### Layout

- header
- footer
- sidebar
- panels
- tabs
- splits
- component positions
- component sizing

### Components

Users should be able to choose which components are displayed.

Possible components:

```text
Text
Box
Panel
List
Table
Tree
Input
Prompt
Chat
Markdown
Diff
Code
Progress
Status
Notification
Tabs
Split
Agent
Tool
Session
Workspace
Git
Model
Provider
User
```

### Behavior

- keybindings
- slash commands
- hooks
- events
- workflows
- custom agents

---

# 21. Tercode UI Philosophy

The default Tercode UI should be intentionally clean and functional rather than excessively designed.

Advanced users should be able to transform it into a completely different interface.

For example:

```text
Minimal
> _
```

or:

```text
┌──────────────┬────────────────────────────┐
│ Workspace    │ Agent                      │
│              │                            │
│ src/         │ ◆ Implementing API         │
│ tests/       │                            │
│ README       │ → reading routes.go        │
│              │ ✓ tests passed             │
├──────────────┴────────────────────────────┤
│ >                                          │
└───────────────────────────────────────────┘
```

Both are valid Tercode environments.

---

# 22. Configuration and UI DSL

Configuration should be data-driven.

YAML is appropriate for configuration, themes, layouts, and keybindings.

Example:

```yaml
ui:
  layout:
    type: split

    children:
      - type: workspace
        width: 30%

      - type: vertical
        children:
          - type: agent
            flex: 1
          - type: prompt
            height: 3
```

Configuration should eventually support a stable scripting layer for advanced behavior.

Lua or another embedded language may be introduced after the declarative UI API is mature.

The scripting layer must not be required for normal use.

---

# 23. User-Defined Commands

Users should be able to define commands such as:

```text
/review
/deploy
/test-all
/security
```

A custom command may invoke:

- prompts
- agents
- tools
- workflows
- external integrations

---

# 24. Custom Agents

Users should eventually be able to define specialized agents.

Example:

```yaml
name: reviewer

model: openrouter/model-name

system_prompt: |
  You are a strict senior software engineer.
  Review code for correctness and security.

tools:
  - read_file
  - search
  - git

permissions:
  write: deny
```

---

# 25. Tercode Community / "Rice" Ecosystem

Tercode should eventually support a community ecosystem similar in spirit to Linux/Neovim configuration communities.

Users may share:

- themes
- layouts
- components
- keybindings
- commands
- agents
- workflows
- complete configurations ("Tercode rices")

A shared configuration should be able to transform the Tercode experience without changing the core binary.

---

# 26. LAN Collaboration

Multiple users on different PCs should be able to work on the same project over a local network.

Example:

```text
PC 1 ──┐
PC 2 ──┼── Tercode Workspace
PC 3 ──┘
```

Users should be able to:

- join a workspace
- see connected users
- synchronize changes
- observe agent activity
- collaborate on tasks
- detect conflicts

---

# 27. Collaboration Architecture

Initial collaboration should use a host/client architecture.

Host:

```bash
tercode host
```

Client:

```bash
tercode connect <host>
```

The host owns the authoritative collaborative workspace state.

The MVP should use authenticated TLS networking with HTTP/WebSocket communication.

CRDT/P2P collaboration is not required for the initial version.

---

# 28. Synchronization

Synchronization should be event-based.

Events may include:

```text
FILE_CREATED
FILE_MODIFIED
FILE_DELETED
USER_JOINED
USER_LEFT
AGENT_STARTED
AGENT_EVENT
AGENT_FINISHED
SYNC_REQUEST
SYNC_RESPONSE
```

The system must detect conflicting modifications.

Initial conflict handling may use:

- file locking
- version tracking
- conflict detection
- diff/merge
- Git

---

# 29. Multi-Agent Collaboration

Future versions should support multiple agents working on different tasks.

Example:

```text
Main Task
    ↓
Planner
    ├── Backend Agent
    ├── Frontend Agent
    └── Test Agent
```

Agents must have isolated tasks and appropriate permissions.

Conflicting modifications must be detected.

---

# 30. Git

Git remains Tercode's version-control foundation.

Tercode should be able to:

- inspect status
- inspect diffs
- inspect history
- create branches
- commit with authorization
- revert changes

Tercode collaboration does not replace Git or GitHub/GitLab.

---

# 31. MCP

MCP should be supported after the native tool system is stable.

MCP tools must use the same permission and security framework as native tools.

---

# 32. Performance

Tercode should prioritize:

- fast startup
- low memory usage
- efficient filesystem access
- streaming
- bounded concurrency
- incremental repository indexing
- minimal unnecessary background services

Go is selected because it provides:

- native binaries
- efficient networking
- strong concurrency primitives
- fast filesystem APIs
- low operational overhead
- easy cross-platform distribution

---

# 33. Cross-Platform Support

Target platforms:

- Linux
- Windows
- macOS

The basic client should not require Python, Node.js, Java, or Docker.

---

# 34. Non-Goals

Tercode will not initially attempt to:

- become a traditional graphical IDE
- replace Git
- replace GitHub/GitLab
- train foundation models
- build an inference engine
- implement CRDT-based P2P collaboration
- require a vector database
- execute destructive operations without authorization

---

# 35. MVP

The MVP must provide:

### Core

- Go application
- configuration
- logging
- workspace detection
- fast startup

### Providers

- OpenRouter
- generic OpenAI-compatible API
- NVIDIA NIM
- API key management
- model discovery
- model testing
- streaming
- tool calling where supported

### Agent

- chat
- plan
- code
- review
- debug
- tool loop
- context management
- approval system
- cancellation

### Tools

- filesystem
- search
- shell
- Git
- test/build

### TUI

- `tercode` launches interactive TUI
- `>` prompt
- streaming output
- slash commands
- command completion
- agent/tool visualization
- diff viewer
- model/provider indicator
- session management
- initial theme/layout customization
- keybindings

### Collaboration

- host
- client
- authentication
- workspace synchronization
- presence
- conflict detection

---

# 36. Development Roadmap

## v0.1 — Foundation

- Go project
- CLI
- configuration
- provider abstraction
- OpenAI-compatible provider
- OpenRouter
- model discovery
- secret management

## v0.2 — Agent Runtime

- interactive prompt
- streaming
- filesystem
- shell
- Git
- plan/code/review/debug
- permissions
- event system

## v0.3 — Context

- repository search
- ripgrep
- Tree-sitter
- symbol indexing
- Git context
- context compaction

## v0.4 — Tercode UI Engine

- interactive TUI
- component system
- layout engine
- slash commands
- autocomplete
- themes
- keybindings
- configurable layouts
- multiple contexts

## v0.5 — LAN Workspace

- host
- client
- authentication
- synchronization
- presence
- conflict handling

## v0.6 — Multi-Agent

- task management
- agent spawning
- delegation
- coordination

## v0.7 — Extensibility

- MCP
- custom agents
- custom commands
- scripting
- plugins
- community configuration ecosystem

## v1.0 — Stable Platform

- stable provider API
- stable UI customization API
- stable workspace protocol
- cross-platform releases
- security review
- production-ready collaboration
- comprehensive documentation

---

# 37. Success Criteria

Tercode succeeds when:

1. A developer can install it and immediately run `tercode`.
2. The interactive agent experience is useful without configuration.
3. Multiple providers can be used without changing the agent core.
4. Models can be discovered and tested.
5. Agents can safely modify and test real projects.
6. The TUI is substantially customizable without modifying source code.
7. Developers can create their own Tercode layouts and workflows.
8. Developers can share complete Tercode configurations.
9. Multiple users can collaborate over a LAN.
10. The architecture remains modular as the project grows.

---

# 38. Product Identity

**Name:** Tercode  
**Meaning:** Terminal Code

**Core philosophy:**

> Powerful by default. Limitless by configuration.

**Positioning:**

> Tercode is an open-source, highly customizable agentic development environment for the terminal.

**Long-term vision:**

> Build the Neovim/Hyprland philosophy for AI-assisted software development: a powerful core that gives developers the freedom to create their own development environment.

---

# 39. Final Product Definition

Tercode is not merely another AI coding CLI.

It is a platform composed of:

```text
Agent Engine
+
Provider System
+
Tool System
+
Workspace
+
Git
+
Collaboration
+
UI Engine
+
Customization System
```

The core provides the capabilities.

The UI engine provides composable primitives.

The user defines the experience.

Therefore:

> **Tercode should be to agentic development what Neovim and Hyprland are to their respective domains: highly capable out of the box, deeply configurable, composable, extensible, and ultimately shaped by its users.**
