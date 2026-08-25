# AGENTS.md — Tercode Development Instructions

> This file defines the engineering rules, development principles, and operating instructions for AI agents contributing to the Tercode codebase.

---

## 1. Mission

You are an engineering agent working on **Tercode (Terminal Code)**.

Tercode is an open-source, provider-agnostic, highly customizable agentic development environment for the terminal.

The product vision is:

> **Powerful by default. Limitless by configuration.**

Tercode should combine:

- a capable agentic coding runtime
- multiple AI providers
- safe tool execution
- project-aware context
- Git integration
- LAN-first collaboration
- a composable terminal UI
- deep user customization
- future extensibility

The goal is not to build a clone of Claude Code, Codex CLI, or another existing product.

The goal is to build a distinct platform with its own architecture, interaction model, and customization philosophy.

---

# 2. Source of Truth

Before implementing a feature, understand these documents:

```text
PRD
    ↓
defines WHAT Tercode must do

ARCHITECTURE
    ↓
defines HOW the system is structured

AGENTS.md
    ↓
defines HOW an engineering agent must work
```

When requirements conflict:

1. Resolve contradictions explicitly.
2. Prefer the most recent approved requirement.
3. Preserve architectural boundaries.
4. Do not silently invent product behavior.
5. If a decision materially changes architecture or public APIs, document it.

Relevant documents:

```text
README.md
PRD.md
ARCHITECTURE.md
API_SPEC.md
PROTOCOL.md
UI_DESIGN.md
SECURITY.md
```

Not all of these files need to exist initially. Create them when their corresponding subsystem becomes mature enough to require a formal contract.

---

# 3. Core Engineering Philosophy

Tercode follows these principles:

### 3.1 Core vs Experience

The core provides capabilities.

The user defines the experience.

Do not hardcode the user interface into the agent runtime.

### 3.2 Provider Independence

The agent must never depend on a specific AI provider.

### 3.3 Local First

Tercode must remain useful without LAN collaboration or optional services.

### 3.4 Explicit Boundaries

Subsystems communicate through stable interfaces rather than implementation details.

### 3.5 Strong Defaults

A new user must get a useful experience without configuration.

### 3.6 Composability

Users should be able to compose Tercode primitives into their own environment.

### 3.7 Security Before Convenience

AI-generated commands, project content, network users, and external responses are untrusted.

### 3.8 Simplicity Before Abstraction

Do not introduce an abstraction merely because it is theoretically elegant.

Introduce abstractions when they protect a real architectural boundary or enable a real extension point.

---

# 4. Technology Rules

## Primary Language

Use **Go** for the core application.

Go is selected for:

- native compilation
- low operational overhead
- efficient networking
- filesystem performance
- concurrency
- cross-platform distribution
- simple deployment

Do not introduce another runtime language into the core without an explicit architectural decision.

---

# 5. General Coding Standards

Write production-quality Go.

Prefer:

- small cohesive packages
- explicit dependencies
- interfaces at architectural boundaries
- context-aware operations
- structured errors
- deterministic behavior
- readable code
- tests close to the code they verify

Avoid:

- unnecessary global state
- giant structs
- giant functions
- hidden side effects
- reflection-heavy designs
- unnecessary generics
- premature optimization
- excessive abstraction
- provider-specific logic leaking into core packages

---

# 6. Package Design

The intended architecture is approximately:

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
└── platform/
```

Do not create packages simply to satisfy this tree.

A package should exist because it owns a coherent responsibility.

---

# 7. Dependency Rules

These rules are mandatory.

### TUI → Provider

Forbidden.

The TUI must communicate through application/domain interfaces.

### Provider → TUI

Forbidden.

Provider adapters must have no knowledge of rendering.

### Agent → TUI

Forbidden.

The agent emits events/state; the TUI renders them.

### Tool → TUI

Forbidden.

Tools return structured results.

### Workspace → TUI

Forbidden.

Workspace state is exposed through services/events.

### Provider-specific logic → Agent

Forbidden.

Provider-specific behavior belongs inside adapters.

### Platform-specific code → Core

Avoid.

Platform-specific behavior belongs behind platform interfaces.

### User configuration → Internal implementation

Forbidden.

Configuration should use stable public APIs.

---

# 8. Agent Runtime Rules

The Agent Runtime is one of the most important parts of Tercode.

It must remain UI-independent.

The runtime is responsible for:

- sessions
- tasks
- context
- model requests
- tool execution
- streaming
- approvals
- retries
- cancellation
- fallback
- structured events

The runtime must not contain rendering logic.

Bad:

```go
agent.ShowSpinner()
agent.RenderToolCall()
agent.PrintMarkdown()
```

Good:

```go
events.Publish(ToolStarted{...})
events.Publish(ToolCompleted{...})
```

The TUI decides how those events are presented.

---

# 9. Agent Execution

The normal agent lifecycle is:

```text
User Request
    ↓
Context Preparation
    ↓
Model Selection
    ↓
Model Request
    ↓
Model Response
    ↓
Tool Call?
 ┌──┴────┐
No      Yes
 │       │
 ▼       ▼
Done   Permission
         ↓
      Execute Tool
         ↓
      Tool Result
         ↓
      Context Update
         ↓
      Model Request
```

Do not bypass the normal execution pipeline for convenience.

---

# 10. Provider Rules

Providers must implement a stable provider contract.

Initial provider targets:

- OpenRouter
- NVIDIA NIM
- generic OpenAI-compatible APIs

Future providers must be addable without modifying the agent runtime.

The agent should never contain code such as:

```go
if provider == "openrouter" {
    ...
}
```

Instead:

```text
Agent
  ↓
Provider Interface
  ↓
Provider Adapter
```

Provider adapters are responsible for:

- API authentication
- request formatting
- response parsing
- streaming
- model discovery
- capability detection
- provider-specific errors

---

# 11. Model Handling

Do not assume every model supports the same capabilities.

Model capabilities may include:

```text
streaming
tool calling
vision
reasoning
structured output
large context
```

The runtime must inspect capabilities before using optional features.

Never assume tool calling is universally available.

---

# 12. API Keys and Secrets

Never:

- hardcode API keys
- commit credentials
- log credentials
- include credentials in agent context
- print full credentials in errors
- store credentials in plaintext project configuration

Use a Secret Manager abstraction.

Potential sources:

```text
environment variables
OS credential store
future secret backends
```

When debugging authentication, redact credentials.

---

# 13. Tool System

Tools must use a common interface.

Every tool should define:

```text
name
description
input schema
execution
result
permissions
```

Initial native tools include:

```text
filesystem
search
shell
git
test/build
```

Future tools may include MCP and external integrations.

---

# 14. Tool Safety

Treat every tool invocation as potentially dangerous.

The permission pipeline is:

```text
Agent
 ↓
Tool Registry
 ↓
Permission Manager
 ↓
Tool
```

Permission modes should support:

```text
deny
ask
allow
```

Do not silently bypass permission checks.

Do not automatically weaken security because a tool is inconvenient to authorize.

---

# 15. Shell Execution

Shell commands are especially sensitive.

The shell runtime must support:

- timeout
- cancellation
- stdout
- stderr
- exit status
- process cleanup

Never assume AI-generated shell commands are safe.

Never execute arbitrary commands merely because the model requested them without passing through the configured permission policy.

---

# 16. Filesystem Safety

All project file access should go through the Workspace Filesystem abstraction.

Respect workspace boundaries.

Do not allow an agent to access unrelated sensitive files by default.

Examples of sensitive locations include:

```text
SSH keys
credential stores
browser profiles
system configuration
private tokens
unrelated repositories
```

Do not add broad filesystem permissions merely to make a test pass.

---

# 17. Git Rules

Git is a first-class subsystem.

Use a Git service rather than embedding Git logic throughout the agent.

The agent may:

- inspect status
- inspect diffs
- inspect history
- create branches
- stage/commit when authorized
- revert changes when authorized

Never silently rewrite user history.

Destructive Git operations require explicit permission.

---

# 18. Context Engineering

Do not send an entire repository to a model unless explicitly necessary.

Prefer:

```text
search
 ↓
relevant files
 ↓
symbols
 ↓
Git context
 ↓
focused model context
```

Use incremental indexing where practical.

Context management should be model-aware.

Track:

```text
system instructions
conversation
tool results
files
context window
reserved output
```

---

# 19. TUI Philosophy

The TUI is not merely a collection of screens.

It is a **composable UI engine**.

The default Tercode interface should be clean and useful.

Advanced users should be able to completely transform it.

The design philosophy is inspired by:

- Hyprland
- Neovim
- Kitty
- Alacritty

Do not interpret this as permission to copy their interfaces.

Tercode must develop its own identity.

---

# 20. UI Composition

Prefer reusable primitives:

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

Components should be composable.

Do not hardcode the default Tercode layout into individual components.

---

# 21. Customization Levels

Tercode customization has four levels.

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

Users choose what information exists on screen.

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

Customization is therefore a product feature, not cosmetic configuration.

---

# 22. UI Configuration

Prefer declarative configuration for stable UI definitions.

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

Do not expose internal implementation fields as configuration options.

If users need a configuration capability, expose a stable semantic property.

Bad:

```yaml
internal_renderer_buffer_strategy_v2: ...
```

Good:

```yaml
agent:
  show_tool_output: true
```

---

# 23. Scripting

An embedded scripting language may eventually provide:

- custom commands
- hooks
- workflows
- dynamic behavior
- advanced UI interactions

Lua is a candidate.

Do not add scripting to the MVP merely because it is possible.

First establish a stable declarative UI and event/action API.

The scripting API must expose public Tercode concepts rather than raw internal Go structures.

---

# 24. Public Customization API

Treat customization as a public API.

Stable concepts include:

```text
components
layouts
events
actions
state selectors
commands
keybindings
hooks
```

Internal refactoring should not unnecessarily break user configurations.

When changing the customization API:

1. identify affected configs
2. document the change
3. provide migration guidance where practical
4. update examples
5. update tests

---

# 25. Configuration Precedence

Use:

```text
Built-in Defaults
 ↓
Global Config
 ↓
Project Config
 ↓
Environment
 ↓
CLI Arguments
 ↓
Session Overrides
```

Do not create arbitrary precedence rules.

Document every new configuration source.

---

# 26. Event System

Use structured events for cross-subsystem communication.

Events should be immutable after publication.

Examples:

```text
TASK_STARTED
MODEL_REQUEST
MODEL_RESPONSE
TOOL_STARTED
TOOL_OUTPUT
TOOL_COMPLETED
FILE_MODIFIED
GIT_CHANGED
USER_JOINED
USER_LEFT
```

Do not use the event bus as an uncontrolled global state store.

Events represent facts that happened, not arbitrary mutable state.

---

# 27. Collaboration Rules

LAN collaboration must be:

- authenticated
- encrypted
- explicit
- local-first
- failure-tolerant

Do not automatically expose a workspace to the entire LAN.

The host must know when another user joins.

The initial architecture is host/client, not P2P CRDT.

Do not introduce CRDTs unless a concrete requirement justifies them.

---

# 28. Collaboration Protocol

Keep transport and domain protocol separate.

Transport:

```text
TCP/TLS
HTTP
WebSocket
```

Protocol:

```text
authentication
workspace events
file synchronization
presence
agent events
conflict information
```

Do not spread wire-format logic throughout business logic.

---

# 29. Synchronization

Use explicit file versions/hashes.

A file state should conceptually contain:

```text
path
hash
version
modified by
modified at
```

Detect concurrent modifications.

Never silently overwrite another user's changes.

---

# 30. Offline/Network Failure

LAN failure must not destroy the local development experience.

When possible:

```text
LAN disconnected
 ↓
Collaboration unavailable
 ↓
Local Tercode continues
```

Do not make the agent runtime depend on a healthy collaboration connection.

---

# 31. Concurrency

Use Go concurrency deliberately.

Appropriate uses:

- model streaming
- tool execution
- file watching
- event processing
- network connections
- background indexing

Every long-running operation should support cancellation.

Avoid:

```go
go func() {
    // no cancellation
}()
```

for persistent/background work.

Track lifecycle ownership of goroutines.

---

# 32. Error Handling

Use structured errors.

Do not make callers parse human-readable error strings to determine behavior.

Errors should distinguish:

```text
authentication
rate limiting
quota
network
timeout
model availability
tool failure
permission
workspace
configuration
sync conflict
```

Preserve underlying causes where useful.

Do not expose secrets in errors.

---

# 33. Logging

Use structured logging.

Logs should contain useful metadata such as:

```text
component
event
session ID
task ID
provider
model
duration
error code
```

Never log:

```text
API keys
passwords
tokens
private credentials
unnecessary source code
```

Provide useful debug logging without creating a privacy hazard.

---

# 34. Performance

Do not optimize blindly.

First identify the actual bottleneck.

Performance-sensitive areas:

```text
startup
filesystem
search
repository indexing
model streaming
TUI rendering
LAN synchronization
```

Prefer:

- streaming
- lazy loading
- incremental indexing
- bounded concurrency
- caching
- efficient standard-library functionality

Avoid premature caching that creates stale-state complexity.

---

# 35. Resource Management

Every resource must have clear ownership.

This includes:

- files
- network connections
- goroutines
- timers
- processes
- subscriptions
- database connections

Use cleanup mechanisms such as:

```go
defer
context cancellation
Close()
```

where appropriate.

---

# 36. Testing Requirements

Every meaningful feature should include tests.

Use:

### Unit tests

For individual components and services.

### Integration tests

For subsystem boundaries.

### End-to-end tests

For critical user workflows.

### Security tests

For permissions, secrets, workspace boundaries, and collaboration authentication.

Do not rely exclusively on manual testing.

---

# 37. Provider Testing

Never require real API keys in normal CI.

Create a mock provider capable of simulating:

```text
success
streaming
tool calls
timeout
rate limit
authentication failure
malformed response
provider outage
```

Real provider tests may exist separately and must be opt-in.

---

# 38. TUI Testing

Separate UI logic from terminal rendering so that most behavior can be tested without a real terminal.

Test:

- state transitions
- commands
- keybindings
- layout parsing
- configuration
- component behavior
- event handling

Do not make visual snapshot tests the only form of TUI testing.

---

# 39. Configuration Testing

Test:

- defaults
- precedence
- invalid configuration
- migration
- missing values
- type errors
- project/global overrides

A malformed user configuration should produce a clear error rather than a panic.

---

# 40. Backward Compatibility

For public APIs and user configuration:

- avoid unnecessary breaking changes
- version stable protocols
- document migrations
- preserve compatibility where practical

Internal APIs may evolve more freely if they are not exposed to users.

---

# 41. Documentation Requirements

When adding a public feature, update the appropriate documentation.

At minimum, consider:

```text
README.md
PRD.md
ARCHITECTURE.md
configuration docs
API_SPEC.md
PROTOCOL.md
SECURITY.md
examples
```

Do not leave major public behavior undocumented.

---

# 42. CLI UX Requirements

CLI output should be:

- concise
- readable
- scriptable where appropriate
- consistent

Interactive TUI output can be visually expressive.

Non-interactive commands should avoid unnecessary decorative output unless requested.

Support machine-readable output for suitable commands in the future.

---

# 43. UX Consistency

The same concept should have the same name everywhere.

For example, if the concept is called:

```text
workspace
```

do not call it:

```text
project
repository
environment
```

interchangeably unless those terms represent genuinely different concepts.

Maintain a product terminology glossary as the project grows.

---

# 44. Do Not Copy Existing Products

Tercode is inspired by existing developer tools, but must not clone their:

- branding
- visual identity
- exact layouts
- proprietary workflows
- source code
- textual content

Learn from their engineering philosophies, not their implementation.

---

# 45. Design Review Before Implementation

Before implementing a significant subsystem, answer:

1. What problem does it solve?
2. Which architectural layer owns it?
3. Which interfaces does it require?
4. Which components depend on it?
5. Does it introduce a new public API?
6. Does it affect configuration?
7. Does it affect security?
8. Does it affect persistence?
9. Does it affect the collaboration protocol?
10. How will it be tested?

For substantial changes, update architecture documentation before or alongside implementation.

---

# 46. Avoid Overengineering

Do not build:

- plugin systems before extension points are understood
- distributed consensus before collaboration requires it
- CRDTs before conflict requirements justify them
- vector databases before retrieval requirements justify them
- complex dependency injection frameworks
- microservices
- unnecessary background daemons

Tercode should remain a single efficient application until there is a concrete reason to split functionality.

---

# 47. Preferred Development Sequence

When implementing a new feature:

```text
Requirement
    ↓
Design
    ↓
Interface
    ↓
Implementation
    ↓
Tests
    ↓
Documentation
    ↓
Review
```

Do not jump directly from a vague requirement to a large implementation.

---

# 48. Change Discipline

Before changing an existing subsystem:

1. Read its package documentation.
2. Read relevant tests.
3. Identify its public interfaces.
4. Identify its consumers.
5. Understand persistence/protocol implications.
6. Make the smallest coherent change.
7. Run relevant tests.
8. Run the full test suite when appropriate.

Do not rewrite stable code simply because a different implementation looks cleaner.

---

# 49. Git Workflow

Keep commits focused.

Prefer:

```text
feat: add model discovery
feat: add provider test command
fix: handle streaming disconnect
refactor: isolate provider transport
test: add mock provider failures
docs: document theme configuration
```

Avoid giant commits containing unrelated changes.

Do not commit:

```text
API keys
secrets
local machine configuration
generated temporary files
personal workspace data
```

---

# 50. Agent Behavior During Development

As an AI engineering agent:

- inspect before editing
- understand surrounding code
- prefer existing abstractions
- keep changes scoped
- explain architectural decisions
- test changes
- do not fabricate test results
- do not claim a feature works without verification
- do not silently change requirements
- do not remove functionality without justification

If a requirement is ambiguous and the decision affects architecture or user-visible behavior, stop and request clarification rather than inventing a permanent design.

---

# 51. When You Find a Better Design

You are allowed to challenge the existing design.

However:

Do not silently replace it.

Instead:

1. Explain the current design.
2. Explain the proposed alternative.
3. Explain trade-offs.
4. Identify affected components.
5. Recommend one option.
6. Update the architecture only after the decision is accepted.

This is especially important for:

- provider APIs
- collaboration protocol
- UI customization API
- persistence
- security
- public configuration

---

# 52. Definition of Done

A feature is not complete merely because the code compiles.

A feature is complete when appropriate:

```text
implementation
+
tests
+
error handling
+
security considerations
+
documentation
+
configuration support
+
compatibility considerations
```

have been addressed.

For user-facing features, verify the actual user workflow.

---

# 53. MVP Discipline

The first goal is a strong, stable core.

Do not attempt to implement the entire long-term Tercode vision at once.

Prioritize:

```text
Provider
 ↓
Agent
 ↓
Tools
 ↓
Context
 ↓
Event System
 ↓
TUI
 ↓
Customization
 ↓
LAN
 ↓
Multi-Agent
 ↓
Extensions
```

Each layer should be useful before the next major layer is added.

---

# 54. North-Star Architecture

The final architecture should preserve this relationship:

```text
                       TERCODE
                          │
              ┌───────────┼───────────┐
              ▼           ▼           ▼
           Agent       Workspace   Providers
              │           │           │
              └───────────┼───────────┘
                          ▼
                     Event System
                          │
                          ▼
                     UI Engine
                          │
          ┌───────────────┼────────────────┐
          ▼               ▼                ▼
      Components        Layout          Behavior
          │               │                │
          └───────────────┼────────────────┘
                          ▼
                  User Configuration
                          │
                          ▼
                  Developer's Tercode
```

The agent should be powerful.

The provider system should be open.

The workspace should be safe.

The collaboration system should be reliable.

The UI should be composable.

The configuration should be expressive.

The entire system should remain understandable.

---

# 55. Final Instruction to the Engineering Agent

Build Tercode as a **long-lived open-source platform**, not as a disposable prototype.

Every architectural decision should preserve the ability for Tercode to grow from:

```text
a terminal AI coding tool
```

into:

```text
a customizable agentic development environment
```

without requiring a fundamental rewrite.

Remember the project's central philosophy:

> **Powerful by default. Limitless by configuration.**

And the most important implementation rule:

> **Do not build the experience for the user. Build the primitives that allow the user to build their experience.**
