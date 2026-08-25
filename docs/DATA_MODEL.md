# Tercode — Data Model

**Status:** Draft v1.0

## 1. Purpose

This document defines the conceptual entities used across Tercode.

It is not a final database schema.

## 2. Workspace

```text
Workspace
├── id
├── name
├── root_path
├── repository
├── created_at
└── metadata
```

## 3. User

```text
User
├── id
├── display_name
└── permissions
```

A LAN user identity is separate from an AI provider identity.

## 4. Provider

```text
Provider
├── id
├── name
├── type
└── configuration
```

Credentials are not part of the ordinary provider record.

## 5. Model

```text
Model
├── id
├── provider_id
├── name
├── context_window
├── capabilities
├── modalities
└── metadata
```

## 6. Agent

```text
Agent
├── id
├── name
├── configuration
├── model_preference
├── tools
└── permissions
```

## 7. Session

```text
Session
├── id
├── workspace_id
├── agent_id
├── provider_id
├── model_id
├── created_at
└── state
```

## 8. Task

```text
Task
├── id
├── session_id
├── request
├── plan
├── status
├── created_at
└── completed_at
```

## 9. Event

```text
Event
├── id
├── type
├── timestamp
├── workspace_id
├── session_id
├── task_id
└── payload
```

## 10. File State

For collaboration:

```text
FileState
├── path
├── hash
├── version
├── modified_by
└── modified_at
```

## 11. Persistence Principle

The data model stores application metadata.

Project source files remain in the filesystem/Git repository.

Do not duplicate source files into the database without a clear requirement.
