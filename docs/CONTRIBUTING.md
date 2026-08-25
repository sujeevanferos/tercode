# Contributing to Tercode

Thank you for contributing to Tercode.

Tercode is an open-source project focused on agentic development, provider freedom, collaboration, and deep terminal customization.

## 1. Before Contributing

Read:

```text
AGENTS.md
PRD.md
ARCHITECTURE.md
```

for the project's engineering and product principles.

## 2. Contribution Types

Contributions may include:

- bug fixes
- tests
- provider adapters
- tools
- UI components
- themes
- layouts
- documentation
- performance improvements
- security improvements

## 3. Code Changes

Keep changes:

- focused
- reviewable
- tested
- documented where user-facing

Avoid unrelated refactoring in feature/bug-fix commits.

## 4. Architecture

Do not bypass architectural boundaries to achieve a short-term result.

If a change requires a new dependency direction, explain it in the pull request.

## 5. UI Contributions

New UI components should be:

- reusable
- composable
- configurable where appropriate
- independent from provider-specific logic

Do not hardcode the default layout into reusable components.

## 6. Providers

Provider adapters must not leak provider-specific behavior into the agent runtime.

## 7. Security

Never submit:

- API keys
- passwords
- private tokens
- personal credentials
- private project data

Report security vulnerabilities privately.

## 8. Tests

Add tests for meaningful behavior.

Run the relevant test suite before submitting changes.

## 9. Commits

Prefer focused commits:

```text
feat: add model discovery
fix: handle provider timeout
test: add tool permission cases
docs: document layout configuration
```

## 10. Pull Requests

A pull request should explain:

- what changed
- why it changed
- architectural impact
- security impact
- testing performed
- documentation changes

## 11. Design Changes

For significant architectural changes, discuss the design before implementing a large change.

Do not silently replace an established subsystem.

## 12. Philosophy

Tercode should remain:

> Powerful by default. Limitless by configuration.

Contributions should strengthen both halves of that principle.
