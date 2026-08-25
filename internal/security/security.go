// Package security implements the PermissionManager and SecretManager for Tercode.
//
// Security architecture principles:
//   - All tool execution passes through the PermissionManager before running.
//   - Credentials are managed exclusively through the SecretManager.
//   - No secret value ever reaches a log, the TUI output, or the agent context.
package security

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Policy controls how a permission is handled.
type Policy string

const (
	// PolicyDeny prevents the action unconditionally.
	PolicyDeny Policy = "deny"
	// PolicyAllow permits the action unconditionally.
	PolicyAllow Policy = "allow"
	// PolicyAsk prompts the user for approval before the action executes.
	PolicyAsk Policy = "ask"
)

// Risk categorises the potential impact of a tool action.
type Risk string

const (
	RiskLow    Risk = "low"
	RiskMedium Risk = "medium"
	RiskHigh   Risk = "high"
)

// PermissionRequest describes a proposed tool action requiring a policy decision.
type PermissionRequest struct {
	ToolName    string
	Action      string
	Description string
	Risk        Risk
	// SessionID ties the request to the agent session that issued it.
	SessionID string
}

// PermissionResult is the outcome of a permission evaluation.
type PermissionResult struct {
	Granted bool
	// Reason is a human-readable explanation for display in the TUI.
	Reason string
}

// ApprovalFunc is called when a tool requires interactive approval (PolicyAsk).
// It receives the request and must return true to grant or false to deny.
// Implementations must be cancellable via the context.
type ApprovalFunc func(ctx context.Context, req PermissionRequest) (bool, error)

// PermissionManager evaluates tool permission policies and coordinates
// interactive approval when required.
type PermissionManager struct {
	mu       sync.RWMutex
	policies map[string]Policy // toolName → Policy
	approval ApprovalFunc
}

// NewPermissionManager creates a PermissionManager with the given default
// policies and approval callback.
//
// defaultPolicies maps tool names to their policy. Tools not present in the
// map default to PolicyAsk.
func NewPermissionManager(defaultPolicies map[string]Policy, approval ApprovalFunc) *PermissionManager {
	if defaultPolicies == nil {
		defaultPolicies = make(map[string]Policy)
	}
	return &PermissionManager{policies: defaultPolicies, approval: approval}
}

// Check evaluates the policy for the given request and returns a PermissionResult.
// For PolicyAsk, it invokes the approval function.
func (pm *PermissionManager) Check(ctx context.Context, req PermissionRequest) (PermissionResult, error) {
	pm.mu.RLock()
	policy, ok := pm.policies[req.ToolName]
	pm.mu.RUnlock()

	if !ok {
		policy = PolicyAsk // Safe default: ask for unknown tools.
	}

	switch policy {
	case PolicyDeny:
		return PermissionResult{Granted: false, Reason: "denied by policy"}, nil

	case PolicyAllow:
		return PermissionResult{Granted: true, Reason: "allowed by policy"}, nil

	case PolicyAsk:
		if pm.approval == nil {
			// No approval function configured — deny by default.
			return PermissionResult{Granted: false, Reason: "no approval handler configured"}, nil
		}
		granted, err := pm.approval(ctx, req)
		if err != nil {
			return PermissionResult{}, fmt.Errorf("permission: approval failed: %w", err)
		}
		reason := "user approved"
		if !granted {
			reason = "user denied"
		}
		return PermissionResult{Granted: granted, Reason: reason}, nil
	}

	return PermissionResult{Granted: false, Reason: "unknown policy"}, nil
}

// SetPolicy overrides the policy for a specific tool name at runtime.
func (pm *PermissionManager) SetPolicy(toolName string, p Policy) {
	pm.mu.Lock()
	pm.policies[toolName] = p
	pm.mu.Unlock()
}

// ---

// SecretManager retrieves credentials from environment variables and, in future,
// the OS credential store. It ensures secrets are never passed through layers
// unnecessarily.
type SecretManager struct {
	mu      sync.RWMutex
	runtime map[string]string // ephemeral in-memory store for session overrides
}

// NewSecretManager creates an initialized SecretManager.
func NewSecretManager() *SecretManager {
	return &SecretManager{runtime: make(map[string]string)}
}

// Get retrieves a secret by key. Resolution order:
//  1. In-memory runtime store (session overrides, never persisted)
//  2. Environment variable matching key (uppercased)
//
// Returns an empty string if the secret is not found; callers should treat
// an empty result as "not configured" rather than an error.
func (sm *SecretManager) Get(key string) string {
	sm.mu.RLock()
	if v, ok := sm.runtime[key]; ok {
		sm.mu.RUnlock()
		return v
	}
	sm.mu.RUnlock()

	// Try environment variable: key → TERCODE_<KEY_UPPER>
	envKey := "TERCODE_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	// Also try the key directly as an environment variable name for
	// conventional keys like OPENROUTER_API_KEY.
	return os.Getenv(key)
}

// SetRuntime stores a secret in the ephemeral in-memory store.
// This value lasts only for the current process lifetime and is never persisted.
func (sm *SecretManager) SetRuntime(key, value string) {
	sm.mu.Lock()
	sm.runtime[key] = value
	sm.mu.Unlock()
}

// Redact returns a safe representation of a secret for logging:
// it shows only the first 4 characters followed by "****".
func Redact(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:4] + "****"
}
