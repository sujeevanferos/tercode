# Tercode — API Specification

**Status:** Draft v1.0  
**Scope:** Internal service contracts and future public extension boundaries

## 1. API Principles

Interfaces must:

- be small
- represent stable concepts
- avoid leaking implementation details
- support cancellation
- return structured errors
- remain testable

## 2. Provider API

Conceptual interface:

```go
type Provider interface {
    ID() string
    Name() string
    Authenticate(context.Context) error
    ListModels(context.Context) ([]Model, error)
    Chat(context.Context, ChatRequest) (ChatResponse, error)
    Stream(context.Context, ChatRequest) (<-chan StreamEvent, error)
    TestModel(context.Context, string) TestResult
    Capabilities(string) Capabilities
}
```

## 3. Tool API

```go
type Tool interface {
    Name() string
    Description() string
    Schema() Schema
    Execute(context.Context, ToolInput) ToolResult
}
```

## 4. Agent API

```go
type AgentRuntime interface {
    CreateSession(context.Context, SessionRequest) (Session, error)
    Run(context.Context, TaskRequest) (TaskResult, error)
    Cancel(context.Context, TaskID) error
}
```

## 5. Workspace API

The workspace service should expose operations for:

```text
Open
Close
Info
ReadFile
WriteFile
List
GitStatus
GitDiff
```

All filesystem operations must respect workspace boundaries.

## 6. Event API

```go
type EventBus interface {
    Publish(Event) error
    Subscribe(EventFilter, EventHandler) Subscription
}
```

Subscription lifecycle must be explicit.

## 7. Permission API

```text
Check
RequestApproval
Grant
Deny
Revoke
```

Permissions must be scoped to an appropriate resource/action.

## 8. Model API

Models should expose normalized metadata:

```text
ID
Provider
Name
ContextWindow
Capabilities
Modalities
Pricing
```

## 9. Configuration API

Configuration must expose semantic values rather than raw parser structures.

## 10. Collaboration API

The collaboration service should separate:

```text
connection
authentication
workspace synchronization
presence
events
conflict resolution
```

## 11. API Stability

Internal APIs may change during development.

Once an interface is declared public or used by user extensions, it becomes a compatibility surface and must be versioned/documented.
