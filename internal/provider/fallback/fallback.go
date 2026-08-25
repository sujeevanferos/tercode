package fallback

import (
	"context"
	"fmt"

	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/provider/registry"
)

// FallbackRule defines an alternative provider/model if the primary fails.
type FallbackRule struct {
	PrimaryProvider string
	FallbackProvider string
	FallbackModel    string
}

// Manager handles routing failed model requests to configured alternatives.
type Manager struct {
	reg   *registry.Registry
	rules []FallbackRule
}

// New creates a FallbackManager.
func New(reg *registry.Registry, rules []FallbackRule) *Manager {
	return &Manager{reg: reg, rules: rules}
}

// GetFallback returns an alternate provider and model for a failed primary request.
func (m *Manager) GetFallback(ctx context.Context, failedProvider, failedModel string) (provider.Provider, string, bool) {
	for _, r := range m.rules {
		if r.PrimaryProvider == failedProvider {
			p, err := m.reg.Get(r.FallbackProvider)
			if err == nil {
				return p, r.FallbackModel, true
			}
		}
	}
	// Default: if openrouter failed, try generic openai if registered
	if failedProvider == "openrouter" {
		if p, err := m.reg.Get("openai"); err == nil {
			return p, "gpt-4o", true
		}
	}
	return nil, "", false
}

// ExecuteWithFallback attempts a chat call with the primary provider and falls back if retryable.
func (m *Manager) ExecuteWithFallback(ctx context.Context, primary provider.Provider, req provider.ChatRequest) (provider.ChatResponse, error) {
	resp, err := primary.Chat(ctx, req)
	if err == nil {
		return resp, nil
	}

	fbProvider, fbModel, ok := m.GetFallback(ctx, primary.ID(), req.Model)
	if !ok {
		return provider.ChatResponse{}, fmt.Errorf("primary failed (%v) and no fallback available", err)
	}

	req.Model = fbModel
	fbResp, fbErr := fbProvider.Chat(ctx, req)
	if fbErr != nil {
		return provider.ChatResponse{}, fmt.Errorf("primary error: %v; fallback error: %w", err, fbErr)
	}
	return fbResp, nil
}
