# Tercode — Development Guide

## 1. Development Environment

Required:

- Go
- Git

Optional tooling depends on the subsystem being developed.

The basic Tercode build should not require Python, Node.js, Java, or Docker.

## 2. Build

Use the repository's standard Go commands.

Typical development commands:

```bash
go build ./...
go test ./...
go vet ./...
```

Additional lint/security tooling should be documented by the repository as it is adopted.

## 3. Run

During development:

```bash
go run ./cmd/tercode
```

If the collaboration daemon exists:

```bash
go run ./cmd/tercoded
```

## 4. Test Before Commit

At minimum:

```bash
go test ./...
go vet ./...
```

Run targeted integration/e2e tests when modifying external boundaries.

## 5. Development Order

For a new subsystem:

```text
Requirement
→ design
→ interface
→ implementation
→ tests
→ documentation
```

## 6. Local Provider Testing

Use mock providers for normal development.

Real API provider tests should be explicit and must never require credentials for the default test suite.

## 7. UI Development

Develop the UI using deterministic state and event inputs.

Avoid coupling UI tests to a live provider.

## 8. Collaboration Development

Use local test hosts/clients.

Test:

- authentication
- connection failure
- synchronization
- conflict detection
- reconnect behavior

## 9. Performance Testing

Benchmark only after correctness.

Potential benchmark targets:

```text
startup
file search
repository indexing
context selection
TUI rendering
model streaming
file synchronization
```

## 10. Debugging

Prefer structured logs and reproducible tests.

Never ask users to provide raw logs containing credentials.

## 11. Generated Files

Do not commit:

```text
local databases
temporary logs
API credentials
build binaries
editor-specific temporary files
```

unless explicitly required by the repository.

## 12. Documentation

Update documentation when changing:

- user-facing commands
- configuration
- provider behavior
- UI customization
- protocol
- security behavior
- public APIs

## 13. Branch Workflow

- `production`: Protected production release branch.
- `main`: Integration branch for stable, tested features.
- `dev`: Active development branch for daily work and feature engineering.
- `user`: Dedicated sandbox for user experiments and custom configurations.
- `feat/*`, `fix/*`: Short-lived topic branches branched off and merged back into `dev`.

