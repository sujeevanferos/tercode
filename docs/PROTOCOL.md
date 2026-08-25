# Tercode — Collaboration Protocol

**Status:** Draft v1.0  
**Purpose:** Define the initial LAN collaboration protocol.

## 1. Scope

The protocol connects Tercode clients to a Tercode workspace host.

It is not intended to replace Git.

## 2. Transport

Initial transport:

```text
TCP
→ TLS
→ HTTP/WebSocket
```

WebSocket is used for real-time events.

HTTP may be used for request/response operations.

## 3. Protocol Version

Every connection must negotiate a protocol version.

Example:

```text
tercode-protocol: 1
```

Incompatible versions must be rejected clearly.

## 4. Connection Flow

```text
Client
 ↓
Connect
 ↓
Protocol negotiation
 ↓
Authentication
 ↓
Workspace authorization
 ↓
Initial synchronization
 ↓
Event stream
```

## 5. Authentication

The host must authenticate clients before workspace access.

Authentication material must never be transmitted over plaintext transport.

## 6. Message Envelope

Conceptually:

```json
{
  "version": 1,
  "type": "file.modified",
  "id": "event-id",
  "timestamp": "...",
  "workspace": "...",
  "payload": {}
}
```

## 7. Core Message Types

```text
hello
auth.request
auth.response
workspace.join
workspace.state
file.create
file.modify
file.delete
file.lock
file.unlock
presence.join
presence.leave
agent.event
sync.request
sync.response
conflict.detected
error
```

## 8. File State

A synchronized file should track:

```text
path
hash
version
modified_by
modified_at
```

## 9. Conflict Detection

A client must not silently overwrite a file whose base version differs from the host's current version.

The host should return a conflict event.

## 10. Presence

Presence should expose only information needed for collaboration, such as:

```text
user identifier
connected state
active workspace
active session/task where appropriate
```

Avoid unnecessary personal data.

## 11. Agent Events

Agent activity may be broadcast as structured events.

Remote clients should not receive secrets or private provider credentials.

## 12. Error Codes

Initial categories:

```text
AUTH_FAILED
UNAUTHORIZED
WORKSPACE_NOT_FOUND
PROTOCOL_MISMATCH
SYNC_CONFLICT
INVALID_MESSAGE
RATE_LIMITED
INTERNAL_ERROR
```

## 13. Future Evolution

Do not add CRDT or peer-to-peer semantics until a concrete requirement demands them.

The protocol should remain simple and observable in the initial releases.
