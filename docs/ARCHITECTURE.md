# Tercode — System Architecture

**Project:** Tercode  
**Meaning:** Terminal Code  
**Status:** Draft v1.0  
**Primary language:** Go  
**Architecture style:** Modular, layered, event-driven, provider-agnostic, local-first

---

# 1. Purpose

This document defines the engineering architecture of Tercode.

The PRD defines what Tercode should be. This document defines how its major systems fit together and the boundaries they must preserve.

The architecture must support:

- agentic coding
- multiple AI providers
- model discovery/testing
- powerful terminal UI
- deep user customization
- multiple agent contexts
- Git integration
- local-first operation
- LAN collaboration
- future multi-agent workflows
- extensibility
- cross-platform operation

---

# 2. Architectural North Star

The central architectural principle is:

> **Tercode Core provides capabilities. The user defines the experience.**

The core must remain independent from the user's visual configuration.

Conceptually:

```text
                         TERCODE CORE
                              │
       ┌──────────────────────┼──────────────────────┐
       │                      │                      │
   AGENT ENGINE          WORKSPACE ENGINE       PROVIDER ENGINE
       │                      │                      │
       └──────────────────────┼──────────────────────┘
                              │
                         EVENT SYSTEM
                              │
                              ▼
                       ┌──────────────┐
                       │   UI ENGINE  │
                       └──────┬───────┘
                              │
             ┌────────────────┼────────────────┐
             ▼                ▼                ▼
        Components         Layout           Behavior
             │                │                │
             └────────────────┼────────────────┘
                              │
                              ▼
                     USER CONFIGURATION
                              │
                              ▼
                      USER'S TERCODE
```

The same binary must be able to produce radically different user environments through configuration and extensions.

---

# 3. Architectural Goals

Tercode should prioritize:

1. security
2. correctness
3. maintainability
4. reliability
5. user experience
6. performance
7. extensibility

Technical requirements include:

- fast startup
- low memory overhead
- efficient filesystem operations
- efficient networking
- streaming
- cancellation
- bounded concurrency
- cross-platform support
- testability

---

# 4. Architectural Style

Tercode combines:

### Layered architecture

Separates presentation, application, domain, and infrastructure concerns.

### Ports and adapters

External systems are accessed through interfaces.

### Event-driven architecture

Agent, workspace, and collaboration events flow through structured event interfaces.

### Local-first architecture

The local client remains useful without the collaboration subsystem.

### Component/composition architecture

The TUI is built from composable primitives rather than a fixed hardcoded layout.

### Extension architecture

Providers, tools, agents, commands, themes, and future plugins extend the system without coupling them to core internals.

---

# 5. High-Level Architecture

```text
┌───────────────────────────────────────────────────────────────────┐
│                         TERCODE PROCESS                            │
│                                                                   │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────────┐  │
│  │     TUI     │  │     CLI      │  │ Configuration / State   │  │
│  └──────┬──────┘  └──────┬───────┘  └────────────┬────────────┘  │
│         │                 │                       │               │
│         └─────────────────┼───────────────────────┘               │
│                           ▼                                       │
│                 ┌─────────────────────┐                           │
│                 │  Application Layer  │                           │
│                 └──────────┬──────────┘                           │
│                            ▼                                      │
│                 ┌─────────────────────┐                           │
│                 │    Agent Runtime    │                           │
│                 └──────────┬──────────┘                           │
│                            │                                      │
│          ┌─────────────────┼──────────────────┐                   │
│          ▼                 ▼                  ▼                   │
│     Context Engine     Tool Runtime      Provider Runtime        │
│          │                 │                  │                   │
│          ▼                 ▼                  ▼                   │
│      Search/AST        FS/Shell/Git     OpenRouter/NIM/etc.       │
│                                                                   │
│                 ┌─────────────────────┐                           │
│                 │   Workspace Layer   │                           │
│                 └──────────┬──────────┘                           │
│                            │                                      │
│                      ┌─────┴─────┐                                │
│                      ▼           ▼                                │
│                   SQLite        Git                               │
└───────────────────────────────────────────────────────────────────┘
                             │
                       optional network
                             ▼
                  ┌────────────────────────┐
                  │ Collaboration Service  │
                  └────────────┬───────────┘
                               │
                              LAN
                               │
                    ┌──────────┼──────────┐
                    ▼          ▼          ▼
                 Client A   Client B   Client C
```

---

# 6. Process Model

Tercode has two primary executable modes.

## 6.1 `tercode`

Main interactive application.

Responsibilities:

- TUI
- CLI
- agent runtime
- provider communication
- local tools
- local workspace
- configuration
- collaboration client
- optional collaboration host

## 6.2 `tercoded`

Optional collaboration daemon.

Responsibilities:

- workspace hosting
- authentication
- synchronization
- event distribution
- presence

`tercoded` must not be required for ordinary local development.

---

# 7. Runtime Modes

### Local

```text
tercode
  ↓
Agent
  ↓
Local Workspace
```

### Host

```text
tercode / tercoded
  ↓
Workspace
  ↓
Collaboration Server
  ↓
LAN
```

### Client

```text
tercode
  ↓
Collaboration Client
  ↓
LAN
  ↓
Remote Workspace
```

---

# 8. Module Boundaries

Recommended modules:

```text
internal/
├── app/
├── agent/
├── context/
├── provider/
├── tools/
├── workspace/
├── collaboration/
├── tui/
├── cli/
├── config/
├── security/
├── persistence/
├── events/
├── runtime/
└── platform/
```

Each module must have a clearly defined responsibility.

---

# 9. Application Layer

The application layer is the orchestration boundary.

Responsibilities:

- initialize services
- construct dependencies
- manage lifecycle
- route commands
- manage sessions
- start/stop agents
- coordinate workspace operations
- coordinate collaboration

It must not contain provider-specific or renderer-specific logic.

---

# 10. Dependency Injection

Dependencies should be explicit.

Conceptually:

```go
type App struct {
    Agent      AgentRuntime
    Workspace  WorkspaceService
    Providers  ProviderRegistry
    Tools      ToolRegistry
    Config     ConfigService
    Events     EventBus
}
```

Avoid global mutable state.

Dependencies should be constructed near application startup.

---

# 11. Agent Runtime

The Agent Runtime is the core execution engine.

Responsibilities:

- agent sessions
- tasks
- prompt processing
- context requests
- model requests
- tool calls
- tool results
- streaming
- retries
- fallback
- approvals
- cancellation
- execution limits
- structured events

The Agent Runtime must not depend on:

- TUI rendering
- theme configuration
- Bubble Tea components
- provider-specific UI behavior

---

# 12. Agent Execution Model

```text
User Request
      ↓
Context Preparation
      ↓
Model Request
      ↓
Model Response
      │
      ├───────────────┐
      ▼               ▼
Final Answer       Tool Call
                      ↓
                Permission Check
                      │
                 ┌────┴────┐
                 ▼         ▼
               Allow      Ask
                 │         │
                 ▼         ▼
            Execute Tool  Approval
                 │         │
                 └────┬────┘
                      ▼
                 Tool Result
                      ↓
                 Context Update
                      ↓
                 Model Request
```

Execution continues until completion, failure, cancellation, or configured limits.

---

# 13. Agent Sessions

A session contains:

```text
Session
├── ID
├── Workspace
├── User
├── Agent
├── Provider
├── Model
├── Conversation
├── Context
├── Permissions
└── Events
```

Sessions should be resumable where practical.

---

# 14. Agent Tasks

A task represents a concrete unit of work.

```text
Task
├── ID
├── Session
├── Request
├── Plan
├── Status
├── Tool Calls
├── Files Modified
└── Result
```

States:

```text
PENDING
PLANNING
RUNNING
WAITING_APPROVAL
COMPLETED
FAILED
CANCELLED
```

---

# 15. Event Architecture

Structured events decouple producers from consumers.

```text
Agent
Provider
Tool
Workspace
Collaboration
      │
      ▼
   Event Bus
      │
 ┌────┼───────────┐
 ▼    ▼           ▼
TUI  Logging   Collaboration
```

Example:

```go
type Event struct {
    ID        string
    Type      EventType
    Timestamp time.Time
    SessionID string
    TaskID    string
    Payload   any
}
```

Events may include:

```text
SESSION_STARTED
SESSION_ENDED
TASK_CREATED
TASK_STARTED
TASK_COMPLETED
TASK_FAILED
MODEL_REQUEST
MODEL_RESPONSE
MODEL_ERROR
TOOL_REQUEST
TOOL_STARTED
TOOL_OUTPUT
TOOL_COMPLETED
TOOL_FAILED
APPROVAL_REQUIRED
APPROVAL_GRANTED
APPROVAL_DENIED
FILE_CREATED
FILE_MODIFIED
FILE_DELETED
GIT_CHANGED
USER_JOINED
USER_LEFT
```

---

# 16. TUI Architecture

The TUI is a **UI engine**, not a hardcoded screen.

It consumes application state/events and sends user actions back through application services.

```text
Agent / Workspace / Provider
             │
             ▼
          Event Bus
             │
             ▼
          TUI State
             │
             ▼
        UI Composition
             │
       ┌─────┼─────┐
       ▼     ▼     ▼
 Components Layout Behavior
       │     │     │
       └─────┼─────┘
             ▼
          Renderer
             │
             ▼
          Terminal
```

Input:

```text
Terminal
   ↓
Input Handler
   ↓
Command / Action
   ↓
Application Layer
   ↓
Domain Services
```

---

# 17. Tercode UI Engine

The UI Engine is one of Tercode's primary differentiators.

It should provide reusable primitives rather than a fixed interface.

Core primitives may include:

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
Scroll
Agent
Tool
Session
Workspace
Git
Model
Provider
User
```

The component system must support composition.

---

# 18. Layout Engine

Layouts should be declarative.

Possible primitives:

```text
row
column
split
stack
tabs
overlay
panel
scroll
```

Example:

```yaml
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

The layout engine should not know about OpenRouter, Git, or model-specific behavior.

---

# 19. Component Data Binding

Components should be able to consume stable application state/events.

For example:

```text
agent
→ active agent state

git
→ Git service state

model
→ active model/provider state

workspace
→ workspace state

session
→ active session state
```

The UI should consume stable selectors/view models rather than reaching into arbitrary internal structures.

---

# 20. UI Customization Model

Customization has four levels.

### Level 1 — Theme

```text
colors
symbols
fonts
borders
spacing
styles
```

### Level 2 — Layout

```text
panels
splits
tabs
position
size
visibility
```

### Level 3 — Components

Users choose which components exist.

### Level 4 — Behavior

```text
keybindings
commands
events
hooks
workflows
agents
scripts
```

This makes customization functional, not merely visual.

---

# 21. Theme System

Themes are data-driven.

Example:

```yaml
theme:
  name: custom

colors:
  primary: "#7C3AED"
  success: "#22C55E"
  warning: "#F59E0B"
  error: "#EF4444"

symbols:
  user: ">"
  assistant: "◆"
  tool: "→"
  success: "✓"
  error: "✗"
```

Themes must not require recompilation.

---

# 22. Configuration Architecture

Configuration precedence:

```text
Built-in Defaults
       ↓
Global Config
       ↓
Project Config
       ↓
Environment Variables
       ↓
CLI Arguments
       ↓
Interactive Session Overrides
```

Recommended paths:

```text
~/.config/tercode/
├── config.yaml
├── ui.yaml
├── theme.yaml
├── keybindings.yaml
├── commands/
├── agents/
├── layouts/
└── themes/
```

Project-specific configuration lives under:

```text
.config/
├── config.yaml
├── ui.yaml
├── theme.yaml
├── keybindings.yaml
├── themes/
├── layouts/
├── commands/
└── agents/
```

---

# 23. Declarative UI vs Scripting

YAML/TOML should be used for stable declarative configuration.

A future embedded scripting layer may be used for advanced behavior:

- dynamic commands
- event hooks
- workflows
- UI behavior

Lua is a candidate, but scripting should be introduced only after the declarative UI API is stable.

The scripting API must expose stable public primitives rather than internal Go implementation details.

---

# 24. UI API Stability

Customization must be treated as a public API.

Breaking internal changes should not unnecessarily break user configurations.

The public customization API should expose:

```text
components
layouts
state selectors
actions
events
commands
keybindings
hooks
```

Internal implementation remains private.

---

# 25. Command System

Commands are independent from the TUI.

```text
Command Registry
├── Built-in Commands
├── User Commands
└── Project Commands
```

Each command defines:

```text
Name
Description
Arguments
Completion
Permissions
Handler
```

The same command service can be used by:

```text
TUI
CLI
Scripts
Future API clients
```

---

# 26. Provider Architecture

Providers are adapters behind a stable interface.

```text
                    Provider Registry
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
      OpenRouter        NVIDIA NIM     OpenAI-Compatible
          │                │                │
          ▼                ▼                ▼
       Adapter          Adapter          Adapter
```

The Agent Runtime must never contain provider-specific branches.

Bad:

```text
if provider == "openrouter" { ... }
```

Good:

```text
Agent → Provider Interface → Adapter
```

---

# 27. Provider Interface

Conceptually:

```go
type Provider interface {
    ID() string
    Name() string

    Authenticate(ctx context.Context) error

    ListModels(ctx context.Context) ([]Model, error)

    Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)

    Stream(ctx context.Context, req ChatRequest) (<-chan Event, error)

    TestModel(ctx context.Context, model string) TestResult

    Capabilities(model string) Capabilities
}
```

The exact behavioral contract belongs in `API_SPEC.md`.

---

# 28. Model Registry

The Model Registry normalizes provider model information.

```text
Model
├── Provider
├── ID
├── Name
├── Context Window
├── Modalities
├── Capabilities
├── Pricing
├── Availability
└── Usage
```

Provider APIs remain authoritative for live availability.

Cached metadata is only a convenience.

---

# 29. Model Selection

Model selection should be its own service.

It may consider:

```text
user preference
task requirements
availability
quota
context size
tool support
latency
reliability
cost
```

The Agent Runtime asks the selector for a model rather than implementing selection itself.

---

# 30. Fallback

Fallback is a dedicated service.

```text
Request
  ↓
Model Selector
  ↓
Primary Model
  ├── success → result
  └── failure
        ↓
   Fallback Manager
        ↓
   Secondary Model
```

Failures should be classified before fallback.

Not every failure should trigger fallback, especially permanent authentication errors.

---

# 31. Tool Architecture

Tools use a common interface:

```go
type Tool interface {
    Name() string
    Description() string
    Schema() Schema
    Execute(ctx context.Context, input any) Result
}
```

Execution:

```text
Agent
 ↓
Tool Registry
 ↓
Permission Manager
 ↓
Tool
 ↓
Result
```

---

# 32. Filesystem Architecture

```text
Agent
 ↓
Filesystem Tool
 ↓
Workspace Filesystem
 ↓
Operating System
```

The filesystem service must enforce workspace boundaries.

---

# 33. Shell Architecture

```text
Agent
 ↓
Shell Tool
 ↓
Permission Manager
 ↓
Process Manager
 ↓
Operating System
```

The Process Manager handles:

- process creation
- stdout/stderr
- exit code
- timeout
- cancellation
- cleanup

---

# 34. Git Architecture

```text
Agent
 ↓
Git Tool
 ↓
Git Service
 ↓
Git CLI / library
 ↓
Repository
```

Git remains the source of truth for version history.

Tercode collaboration does not replace Git.

---

# 35. Context Engine

The Context Engine finds relevant information efficiently.

```text
                 Context Engine
                       │
       ┌───────────────┼───────────────┐
       ▼               ▼               ▼
   Repository        Search            Git
     Scanner           │             Context
       │               │
       ▼               ▼
   Tree-sitter      ripgrep
       │
       ▼
   Symbol Index
```

Pipeline:

```text
User Request
 ↓
Task Classification
 ↓
Repository State
 ↓
Search
 ↓
Relevant Files
 ↓
Symbols / AST
 ↓
Git Context
 ↓
Ranking
 ↓
Model Context
```

The context engine must avoid unnecessarily loading entire repositories.

---

# 36. Context Budget

The Context Manager tracks:

```text
model context window
system instructions
user messages
tool results
project files
conversation history
reserved output
```

When necessary:

```text
Conversation
 ↓
Compaction
 ↓
Summary
 ↓
Relevant context
```

---

# 37. Workspace Architecture

The Workspace is the project boundary.

```text
Workspace
├── Identity
├── Filesystem
├── Git
├── Sessions
├── Users
├── Agents
├── Tasks
├── Events
└── Collaboration
```

The Workspace Service abstracts whether the workspace is local or collaboratively hosted.

---

# 38. Persistence

SQLite should store local application metadata.

Potential data:

```text
workspaces
sessions
tasks
agents
events
providers
models
usage
settings
```

Project files are not duplicated unnecessarily into SQLite.

Credentials must not be stored directly in plaintext SQLite records.

---

# 39. Secrets

Secrets use an abstraction:

```text
Provider
 ↓
Secret Manager
 ↓
Environment / OS Keychain / Future Secret Backend
```

The secret manager must prevent credentials from reaching logs or agent context.

---

# 40. Collaboration Architecture

The collaboration subsystem contains:

```text
server
client
protocol
sync
presence
conflict
authentication
```

Initial model:

```text
                  Workspace Server
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
      Client A        Client B        Client C
```

---

# 41. LAN Protocol

Initial transport:

```text
TCP
 ↓
TLS
 ↓
HTTP / WebSocket
 ↓
Tercode Protocol
```

HTTP is suitable for request/response operations.

WebSocket is suitable for real-time events.

QUIC may be evaluated later.

---

# 42. Collaboration Events

Possible events:

```text
HELLO
AUTH
JOIN_WORKSPACE
USER_JOINED
USER_LEFT
FILE_CREATED
FILE_MODIFIED
FILE_DELETED
FILE_LOCK
FILE_UNLOCK
AGENT_STARTED
AGENT_EVENT
AGENT_FINISHED
SYNC_REQUEST
SYNC_RESPONSE
```

The complete wire contract belongs in `PROTOCOL.md`.

---

# 43. Synchronization

File synchronization should use:

```text
path
file identity
hash
version
modified by
modified time
content
```

Large files should be streamed.

Conflicts must be detected using version/hash information.

---

# 44. Conflict Management

When multiple users modify the same file:

```text
Base Version
    │
 ┌──┴──┐
 ▼     ▼
 A     B
```

Tercode reports a conflict and may offer:

```text
Keep Local
Keep Remote
Merge
View Diff
Cancel
```

Git should be used for complex merges when appropriate.

CRDTs are not required for the MVP.

---

# 45. Collaboration Authentication

LAN access must not automatically imply trust.

The system should use:

```text
Workspace identity
+
Authentication token/access code
+
TLS
```

The host may require explicit user approval.

---

# 46. Multi-Agent Architecture

Future multi-agent operation:

```text
Main Task
    ↓
Planner
    ├── Backend Agent
    ├── Frontend Agent
    └── Test Agent
```

Agents should have:

- isolated tasks
- scoped context
- scoped permissions
- clear ownership of changes

The workspace must detect conflicting modifications.

---

# 47. Security Boundaries

Trust boundaries include:

```text
Internet
   │
AI Providers
   │
────────────
Tercode
   │
 ┌─┼─────────────┐
 ▼ ▼             ▼
Files Shell     LAN Users
```

Potentially untrusted inputs:

- AI output
- project files
- README files
- Git content
- MCP responses
- network clients
- external API responses
- shell output

All dangerous actions require appropriate validation and permission checks.

---

# 48. Prompt Injection

Project content must be treated as untrusted data.

The runtime should distinguish:

```text
system instructions
user instructions
tool output
project content
external content
```

The architecture should support safeguards against malicious repository instructions and injected tool commands.

---

# 49. Concurrency

Go concurrency should be used for:

- model streaming
- tool execution
- agent tasks
- file watching
- LAN connections
- event processing
- background indexing

Long-running operations must support cancellation.

Unbounded goroutine creation is forbidden.

---

# 50. Cancellation

Cancellation should propagate:

```text
User
 ↓
TUI
 ↓
Application
 ↓
Agent
 ↓
Provider / Tool
 ↓
Network / OS process
```

`context.Context` should be the primary cancellation mechanism.

---

# 51. Performance Architecture

Performance-sensitive paths:

- startup
- filesystem
- search
- repository indexing
- model streaming
- TUI rendering
- LAN synchronization

Use:

- native filesystem APIs
- streaming
- lazy loading
- incremental indexing
- caching
- bounded concurrency

Heavy operations should not block TUI responsiveness.

---

# 52. Repository Indexing

Indexing should be incremental.

```text
Initial Index
 ↓
File Watcher
 ↓
Changed Files
 ↓
Incremental Update
```

Full repository rescans should not occur for every user prompt.

---

# 53. File Watching

File watching should detect:

```text
CREATE
MODIFY
DELETE
RENAME
```

Events should be debounced and fed into workspace/event systems.

---

# 54. Provider Networking

Networking should be centralized where practical.

```text
Provider Adapter
 ↓
HTTP Client
 ↓
Connection Pool
 ↓
TLS
 ↓
Internet
```

Support:

- connection reuse
- streaming
- timeout
- cancellation
- retries
- backoff

---

# 55. Error Model

Errors should be structured.

```go
type Error struct {
    Code      ErrorCode
    Component string
    Message   string
    Cause     error
    Retryable bool
}
```

Categories:

```text
AUTH_ERROR
RATE_LIMIT
QUOTA_EXCEEDED
MODEL_UNAVAILABLE
NETWORK_ERROR
TIMEOUT
TOOL_ERROR
PERMISSION_DENIED
WORKSPACE_ERROR
SYNC_CONFLICT
CONFIG_ERROR
```

The UI can render errors based on their category rather than parsing strings.

---

# 56. Extensibility

Core functionality is compiled into Tercode.

Configurable extensions include:

```text
themes
layouts
keybindings
commands
agents
prompts
workflows
provider definitions
MCP servers
```

Future plugins may be introduced only after stable extension boundaries exist.

---

# 57. Public Extension API

The stable public customization API should expose:

```text
components
layouts
state selectors
actions
events
commands
keybindings
hooks
```

Internal Go structs must not automatically become public APIs.

This protects user configurations from internal refactoring.

---

# 58. Community Configuration Ecosystem

Tercode should support shareable configuration repositories.

A complete "rice" may contain:

```text
ui.yaml
theme.yaml
layout.yaml
keybindings.yaml
commands/
agents/
workflows/
```

The goal is that users can share complete Tercode environments without forking the core project.

---

# 59. Cross-Platform Architecture

Target:

```text
Linux
Windows
macOS
```

Platform-specific functionality should be isolated behind interfaces.

Examples:

```text
filesystem
process execution
keychain
terminal capabilities
file watching
network interfaces
```

---

# 60. Suggested Repository Structure

```text
tercode/
├── cmd/
│   ├── tercode/
│   │   └── main.go
│   └── tercoded/
│       └── main.go
│
├── internal/
│   ├── app/
│   ├── agent/
│   │   ├── runtime/
│   │   ├── session/
│   │   ├── task/
│   │   └── executor/
│   ├── context/
│   │   ├── manager/
│   │   ├── search/
│   │   ├── index/
│   │   └── ranking/
│   ├── provider/
│   │   ├── registry/
│   │   ├── selector/
│   │   ├── fallback/
│   │   └── adapters/
│   │       ├── openrouter/
│   │       ├── nvidia/
│   │       └── compatible/
│   ├── tools/
│   │   ├── registry/
│   │   ├── filesystem/
│   │   ├── shell/
│   │   ├── git/
│   │   ├── search/
│   │   └── mcp/
│   ├── workspace/
│   ├── collaboration/
│   │   ├── server/
│   │   ├── client/
│   │   ├── protocol/
│   │   ├── sync/
│   │   ├── presence/
│   │   ├── conflict/
│   │   └── auth/
│   ├── tui/
│   │   ├── app/
│   │   ├── renderer/
│   │   ├── components/
│   │   ├── layout/
│   │   ├── theme/
│   │   ├── commands/
│   │   ├── input/
│   │   └── state/
│   ├── cli/
│   ├── config/
│   ├── security/
│   ├── persistence/
│   ├── events/
│   └── platform/
│       ├── linux/
│       ├── windows/
│       └── darwin/
│
├── pkg/
│   └── protocol/
├── docs/
├── themes/
├── examples/
├── scripts/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── SECURITY.md
├── CONTRIBUTING.md
├── ARCHITECTURE.md
└── CHANGELOG.md
```

This is a target structure, not a requirement to create every package immediately.

---

# 61. Dependency Rules

### Rule 1 — TUI isolation

The TUI must not directly depend on providers.

### Rule 2 — Provider isolation

Providers must not depend on the TUI.

### Rule 3 — Agent isolation

The Agent Runtime must not depend on rendering.

### Rule 4 — Tool isolation

Tools must not directly manipulate TUI state.

### Rule 5 — Workspace isolation

Workspace logic must not depend on the TUI.

### Rule 6 — Collaboration isolation

Collaboration must communicate through workspace/application interfaces and events.

### Rule 7 — Provider-specific logic

Provider-specific behavior belongs inside provider adapters.

### Rule 8 — Platform isolation

Platform-specific code must remain isolated.

### Rule 9 — Configuration isolation

User configuration must interact through stable public APIs, not internal implementation details.

### Rule 10 — Core independence

The Tercode core must remain functional without optional customization scripts, MCP servers, or LAN collaboration.

---

# 62. Dependency Direction

Preferred direction:

```text
Presentation
(TUI / CLI)
      ↓
Application
      ↓
Domain / Core Services
      ↓
Interfaces / Ports
      ↓
Infrastructure Adapters
```

Infrastructure implementations must not force dependencies upward into the core.

---

# 63. State Ownership

| State | Owner |
|---|---|
| Agent state | Agent Runtime |
| Session state | Session Manager |
| Task state | Task Manager |
| Model metadata | Model Registry |
| Provider state | Provider Runtime |
| Workspace state | Workspace Service |
| Git state | Git Service |
| Configuration | Config Service |
| Credentials | Secret Manager |
| Collaboration state | Collaboration Service |
| TUI presentation state | TUI |
| Persistent metadata | Persistence |

Components should not arbitrarily mutate another subsystem's state.

---

# 64. Source of Truth

```text
Project files
→ Filesystem

Git history
→ Git

Credentials
→ Secret Manager

Workspace metadata
→ Workspace/Persistence

Agent state
→ Agent Runtime

Provider/model availability
→ Provider APIs

Visual state
→ TUI
```

Caches must not silently override authoritative state.

---

# 65. Failure Isolation

### Provider failure

```text
Provider A fails
 ↓
Provider classification
 ↓
Fallback if appropriate
```

TUI remains available.

### LAN failure

```text
LAN disconnect
 ↓
Collaboration subsystem fails
 ↓
Local workspace continues
```

### Tool failure

```text
Tool fails
 ↓
Structured tool error
 ↓
Agent receives result
 ↓
Agent decides next action
```

---

# 66. Graceful Degradation

Tercode should continue to provide useful functionality when optional subsystems fail.

Examples:

```text
No LAN
→ local development works

No model discovery
→ manually configured model works

No Tree-sitter
→ search/text context works

No MCP
→ native tools work

Provider unavailable
→ configured fallback may work

No customization files
→ stock UI works
```

---

# 67. Testing Architecture

Each layer must be independently testable.

```text
Unit
 ├── Agent
 ├── Provider
 ├── Tools
 ├── Context
 ├── Config
 └── TUI
        ↓
Integration
 ├── Provider + Agent
 ├── Agent + Tools
 ├── Workspace + Git
 └── Collaboration
        ↓
End-to-End
        ↓
Security
```

A mock provider must exist for deterministic CI tests.

---

# 68. Mock Provider

The mock provider should support deterministic scenarios:

```text
successful response
tool call
streaming
timeout
rate limit
authentication failure
malformed response
provider outage
```

Tests must not require real API keys.

---

# 69. Startup Flow

```text
main()
  ↓
Load configuration
  ↓
Initialize logger
  ↓
Initialize secrets
  ↓
Initialize persistence
  ↓
Discover workspace
  ↓
Initialize providers
  ↓
Initialize tools
  ↓
Initialize agent runtime
  ↓
Initialize TUI
  ↓
Start event loop
```

Heavy indexing should be deferred/asynchronous.

---

# 70. Agent Request Flow

```text
User
 ↓
TUI Input
 ↓
Application
 ↓
Agent Session
 ↓
Context Manager
 ↓
Model Selector
 ↓
Provider
 ↓
Model
 ↓
Response
 ├── final answer
 └── tool call
       ↓
   permission
       ↓
   tool runtime
       ↓
   tool result
       ↓
   agent
```

---

# 71. LAN Request Flow

```text
Client
 ↓
Collaboration Client
 ↓
TLS/WebSocket
 ↓
Workspace Server
 ↓
Authentication
 ↓
Workspace Service
 ↓
Event Bus
 ↓
Other Clients
```

---

# 72. Agent + LAN Flow

```text
Agent
 ↓
Event Bus
 ├── Local TUI
 └── Collaboration Server
       ├── Client A
       ├── Client B
       └── Client C
```

Remote users can therefore observe relevant agent activity without directly coupling their UI to the agent runtime.

---

# 73. Architecture Anti-Patterns

Avoid:

### God object

One giant `TercodeApp` containing every subsystem.

### Provider leakage

Agent code branching on provider names.

### TUI leakage

Agent code directly calling rendering functions.

### Global mutable state

Hidden shared state that makes testing difficult.

### Unbounded concurrency

Unlimited goroutines or network operations.

### Raw credential handling

Passing secrets through unnecessary layers.

### Hardcoded UI

Embedding the default layout throughout component code.

### Configuration leakage

Exposing unstable internal implementation details as public configuration APIs.

---

# 74. Implementation Order

Recommended sequence:

```text
1. Project foundation
2. Configuration
3. Provider abstraction
4. OpenAI-compatible provider
5. OpenRouter
6. Model registry/testing
7. Agent runtime
8. Tool runtime
9. Filesystem
10. Shell
11. Git
12. Context engine
13. Event system
14. TUI engine
15. Customization system
16. Session/context system
17. Security/permissions
18. LAN collaboration
19. Multi-agent
20. MCP/extensions
```

The UI engine should be designed early, but deep customization should be implemented after the underlying application events/state APIs are stable enough to bind to.

---

# 75. Architecture Evolution

### Phase 1

```text
Local Tercode
Provider + Agent + Tools
```

### Phase 2

```text
TUI + Event System + Context
```

### Phase 3

```text
UI Composition + Customization
```

### Phase 4

```text
LAN Workspace
```

### Phase 5

```text
Multi-Agent
```

### Phase 6

```text
MCP + Extensions + Community Ecosystem
```

The architecture must allow these stages without rewriting the core.

---

# 76. Final Architecture

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

The architectural identity of Tercode is therefore:

> **A modular agentic development core surrounded by a composable, user-defined terminal environment.**

The agent engine provides intelligence.

The provider layer provides model freedom.

The tool system provides capabilities.

The workspace provides project ownership.

The collaboration layer provides teamwork.

The UI engine provides the building blocks.

The configuration system lets the developer assemble those building blocks into their own environment.

Therefore:

> **Tercode should be to agentic software development what Neovim and Hyprland are to their respective domains: powerful by default, deeply configurable, composable, extensible, and ultimately shaped by its users.**
