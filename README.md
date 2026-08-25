# Tercode (Terminal Code)

**Tercode** is an open-source, provider-agnostic, local-first agentic development environment for the terminal.

> **Powerful by default. Limitless by configuration.**  
> *The Neovim/Hyprland philosophy applied to agentic software engineering.*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev/)
[![Status](https://img.shields.io/badge/Status-Draft_v1.0-orange.svg)](docs/ROADMAP.md)

---

## Overview

Unlike traditional IDEs where manual code editing is the primary workflow, Tercode is **agent-centric**. You interact directly with an autonomous AI development agent that understands your project structure, inspects code, runs tests, performs search, executes shell operations, edits files, and manages Git workflows.

Tercode is built around a fundamental architectural boundary:
- **Tercode Core** provides capabilities: the agent runtime, provider adapters, context engine, tools, workspace management, and local-network collaboration.
- **Tercode UI Engine** provides composable terminal primitives (panels, splits, tabs, markdown, diffs, chat, prompt).
- **The Developer** defines the experience: configuring layouts, themes, keybindings, workflows, and custom agents without modifying the core binary.

---

## Core Capabilities

### Agent-Centric Workflow
Specialized agent modes tailored to the software engineering lifecycle:
- **Chat**: Context-aware technical dialogue and code exploration.
- **Plan**: Structured implementation planning and architectural design without modifying files.
- **Code**: Autonomous multi-step implementation and file modification loops.
- **Review**: Automated inspection of diffs, code quality, and security invariants.
- **Debug**: Systematic test-driven investigation and failure resolution.

### Provider Independence
Bring your own models, API keys, and endpoints. First-class adapters for:
- **OpenRouter**
- **NVIDIA NIM**
- **OpenAI-Compatible APIs** (Local LLMs, Ollama, vLLM, self-hosted gateways)
- Automatic model discovery, capability probing (tool calling, vision, streaming, reasoning), latency benchmarking, and resilient fallback mechanisms.

### Composable Terminal UI
A declarative terminal interface built from reusable primitives rather than hardcoded screens:
- **Theme**: Data-driven color palettes, typography, borders, and status symbols.
- **Layout**: Declarative YAML-based splits, flexible sizing, tabs, and status bars.
- **Components**: Selectable data widgets (Agent status, Workspace tree, Git diff, Model inspector).
- **Behavior**: Custom keybindings, slash commands (`/config`, `/diff`, `/model`, `/agents`), and workflow hooks.

### Security and Permission Pipeline
Strict trust boundaries for executing AI actions in production environments:
- Granular permission policies (`deny`, `ask`, `allow`) for file modifications and shell processes.
- Workspace sandboxing to prevent path traversal outside the project directory.
- Defense against prompt injection by strictly isolating repository content from system instructions.
- Zero credential leakage: API keys are managed through dedicated secret stores and never committed, logged, or included in LLM context.

### Context Engine
Repository intelligence without exhaustive token consumption:
- Lexical search integration (`ripgrep`).
- Abstract Syntax Tree (AST) and symbol indexing (`Tree-sitter`).
- Git state tracking (uncommitted diffs, commit history).
- Dynamic token budgeting and context compaction.

### Local-First LAN Collaboration
Real-time peer collaboration over local networks (`tercode host` / `tercode connect`):
- Encrypted TLS and WebSocket transport.
- Presence tracking and agent activity streaming.
- Versioned file synchronization and concurrent edit conflict detection.
- Completely self-contained—no third-party cloud services required.

---

## System Architecture

```text
                               TERCODE
                                  │
                  ┌───────────────┼────────────────┐
                  │               │                │
                 TUI             CLI          Automation
                  │               │                │
                  └───────────────┼────────────────┘
                                  ▼
                          Application Layer
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
              ▼                   ▼                   ▼
        Agent Runtime         Workspace          Collaboration
              │                   │                   │
       ┌──────┼──────┐            │            ┌──────┴──────┐
       ▼      ▼      ▼            ▼            ▼             ▼
    Context Tools Provider       Git         Server        Client
       │      │      │
       │   ┌──┼──┐   │
       │   │  │  │   │
       │   ▼  ▼  ▼   ▼
       │  FS Shell Git Providers
       │                 │
       │          ┌──────┼────────┐
       │          ▼      ▼        ▼
       │       OpenRouter NIM Compatible
       │
       └──────────────┐
                      ▼
               Context / Index
```

---

## Customization Model

Tercode customization follows four distinct tiers:

| Tier | Scope | Description |
| :--- | :--- | :--- |
| **Level 1 — Theme** | Visuals | Colors, borders, symbols, typography, and spacing. |
| **Level 2 — Layout** | Structure | Declarative splits, panel dimensions, tabs, and visibility toggles. |
| **Level 3 — Components** | Information | Selection and placement of widgets (Git status, active agents, logs). |
| **Level 4 — Behavior** | Workflows | Keybindings, slash commands, event hooks, and custom agent definitions. |

---

## Documentation

The project specification and technical design are documented in [`docs/`](./docs):

- [Product Requirements Document (PRD)](./docs/PRD.md)
- [System Architecture](./docs/ARCHITECTURE.md)
- [Agent Development Instructions (AGENTS.md)](./docs/AGENTS.md)
- [Technical Specification](./docs/TECHNICAL_SPEC.md)
- [API Specification](./docs/API_SPEC.md)
- [Data Model](./docs/DATA_MODEL.md)
- [Collaboration Protocol](./docs/PROTOCOL.md)
- [UI Design & Customization Specification](./docs/UI_DESIGN.md)
- [Security Model](./docs/SECURITY.md)
- [Development Guide](./docs/DEVELOPMENT.md)
- [Product Roadmap](./docs/ROADMAP.md)
- [Contributing Guidelines](./docs/CONTRIBUTING.md)

---

## Getting Started

### Prerequisites
- **Go**: 1.22 or higher
- **Git**: 2.30 or higher

### Building from Source
```bash
# Clone the repository
git clone git@github.com:sujeevanferos/tercode.git
cd tercode

# Build the binary
go build -o bin/tercode ./cmd/tercode

# Launch the interactive terminal environment
./bin/tercode
```

---

## Development Roadmap

- **Phase 0 — Foundation**: Repository structure, configuration engine, structured logger, CLI skeleton.
- **Phase 1 — Provider Layer**: Unified provider abstraction, OpenRouter / NIM / OpenAI adapters, model discovery.
- **Phase 2 — Agent Core**: Session state machines, task runners, autonomous tool execution loop, permissions.
- **Phase 3 — Native Tools**: Sandboxed filesystem, process execution (shell), Git service.
- **Phase 4 — Context Engine**: Repository search, symbol indexing, token budgeting, context compaction.
- **Phase 5 & 6 — UI Engine & Customization**: Composable TUI, declarative YAML layouts, themes, slash commands.
- **Phase 7 — LAN Collaboration**: Host/Client workspace sync over encrypted TLS/WebSockets.
- **Phase 8 to 10 — Multi-Agent & Ecosystem**: Multi-agent planners, MCP support, community configuration ecosystem.

---

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](./docs/CONTRIBUTING.md) and [AGENTS.md](./docs/AGENTS.md) before submitting pull requests or proposing architectural changes.

---

## License

Tercode is open-source software licensed under the [MIT License](LICENSE).
