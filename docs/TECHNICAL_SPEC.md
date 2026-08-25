# Tercode — Technical Specification

**Status:** Draft v1.0  
**Scope:** Core implementation contract

## 1. Purpose

This document translates the PRD and architecture into implementation-level requirements.

It defines the major runtime contracts without prescribing unnecessary internal implementation details.

## 2. Runtime

Tercode is implemented in Go and distributed primarily as native binaries.

Targets:

- Linux amd64/arm64
- macOS amd64/arm64
- Windows amd64/arm64 where terminal capabilities permit

Primary executables:

```text
tercode
tercoded
```

`tercoded` is optional and exists for collaboration hosting.

## 3. Core Services

The application should construct these services:

```text
ConfigService
SecretManager
WorkspaceService
ProviderRegistry
ModelRegistry
ModelSelector
AgentRuntime
ToolRegistry
PermissionManager
ContextManager
EventBus
Persistence
CollaborationService
UIRuntime
```

Services communicate through explicit interfaces.

## 4. Agent Contract

An agent task receives:

```text
user request
session
workspace
agent configuration
model selection
permission policy
context
```

It produces:

```text
events
tool calls
tool results
final response
structured task result
```

A task must support cancellation.

## 5. Provider Contract

Providers must support, where available:

```text
Authenticate
ListModels
Chat
Stream
TestModel
Capabilities
```

Provider-specific request/response translation belongs entirely inside the adapter.

## 6. Tool Contract

Every tool exposes:

```text
Name
Description
InputSchema
Permissions
Execute
```

Tool execution returns structured results containing:

```text
success
output
error
metadata
```

## 7. Event Contract

Events are immutable facts.

Minimum event metadata:

```text
id
type
timestamp
session_id
task_id
payload
```

Consumers must not mutate published events.

## 8. Context Contract

The context manager must:

- accept task requirements
- identify relevant project information
- enforce model context limits
- prioritize relevant information
- support compaction
- expose deterministic context decisions for debugging

## 9. Configuration

Configuration precedence:

```text
defaults
→ global
→ project
→ environment
→ CLI
→ session
```

Invalid configuration must produce actionable diagnostics.

## 10. Persistence

SQLite is the initial local metadata store.

Persist:

- sessions
- tasks
- workspace metadata
- model/provider metadata where useful
- configuration metadata
- events where required for recovery/history

Do not persist project files unnecessarily.

## 11. UI Runtime

The UI runtime must separate:

```text
domain/application state
presentation state
layout
theme
input
rendering
```

Components consume stable state selectors/view models.

## 12. Collaboration

Initial collaboration uses host/client semantics.

Required capabilities:

- authentication
- encrypted transport
- workspace discovery/connection
- event streaming
- file synchronization
- presence
- conflict detection

The protocol must be versioned.

## 13. Error Handling

Errors must be typed/categorized.

No business logic should depend on parsing human-readable error strings.

## 14. Cancellation

All network, tool, indexing, and agent operations that may run for an extended period must accept `context.Context`.

## 15. Performance Requirements

The application must avoid blocking the TUI on:

- model requests
- repository indexing
- file synchronization
- shell processes
- model discovery

Streaming operations should render incrementally.

## 16. Compatibility

Public configuration, provider interfaces, and collaboration protocol versions must be treated as compatibility surfaces.

Internal packages remain private unless deliberately promoted to a public API.

## 17. Implementation Rule

Implement the smallest architecture that satisfies the current requirement while preserving the documented boundaries. Do not implement future infrastructure speculatively.
