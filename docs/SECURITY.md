# Tercode — Security Model

**Status:** Draft v1.0

## 1. Security Objective

Tercode executes AI-generated actions against real developer environments.

Security therefore takes priority over convenience.

## 2. Trust Model

Potentially untrusted inputs include:

- model output
- repository files
- README files
- Git content
- MCP responses
- shell output
- external API responses
- LAN clients

Do not treat project instructions as equivalent to trusted system instructions.

## 3. Credential Protection

API keys and tokens must:

- never be committed
- never be logged
- never be included in model context
- never be shown in full in the TUI
- be stored through a Secret Manager where possible

## 4. Filesystem Protection

Agent access must be workspace-scoped by default.

Sensitive paths should not be accessible without explicit authorization.

## 5. Shell Security

Shell execution is privileged behavior.

The permission system must be able to:

```text
deny
ask
allow
```

Destructive commands should require authorization according to policy.

## 6. Git Security

Destructive Git operations must require explicit authorization.

Tercode must not silently force-reset, delete branches, or rewrite history.

## 7. Prompt Injection

Repository content can contain instructions designed to manipulate an agent.

The runtime must preserve the distinction between:

```text
system instructions
user instructions
tool output
project content
external content
```

Project content must not automatically gain higher instruction priority.

## 8. LAN Security

Collaboration must use:

- authenticated connections
- encrypted transport
- explicit workspace authorization

Do not expose workspaces anonymously by default.

## 9. Logging

Never log:

```text
API keys
passwords
session tokens
private credentials
```

Sensitive values must be redacted.

## 10. Provider Security

Provider adapters must not persist credentials unless explicitly handled by the Secret Manager.

Provider responses must be treated as external data.

## 11. Dependency Security

Dependencies should be:

- necessary
- maintained
- reviewed
- pinned through Go module tooling

Run vulnerability checks in CI where practical.

## 12. Security Reporting

Security vulnerabilities should be reported privately through the project's documented security contact rather than public issue trackers.

## 13. Security Principle

When convenience and security conflict, preserve security and make the safe behavior understandable to the user.
