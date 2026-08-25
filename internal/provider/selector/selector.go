package selector

import (
	"context"
	"fmt"

	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/provider/registry"
)

// Criteria specifies requirements for picking a model.
type Criteria struct {
	RequiresTools     bool
	RequiresStreaming bool
	MinContextWindow  int
}

// Selector resolves the best provider and model given user config or criteria.
type Selector struct {
	reg *registry.Registry
}

// New creates a ModelSelector.
func New(reg *registry.Registry) *Selector {
	return &Selector{reg: reg}
}

// Select resolves the requested provider and model.
func (s *Selector) Select(ctx context.Context, preferredProvider, preferredModel string, crit Criteria) (provider.Provider, string, error) {
	p, err := s.reg.Get(preferredProvider)
	if err != nil {
		return nil, "", fmt.Errorf("selector: provider %q: %w", preferredProvider, err)
	}

	modelID := preferredModel
	if modelID == "" {
		// Fallback to provider's first discovered model or a reasonable default
		models, err := p.ListModels(ctx)
		if err == nil && len(models) > 0 {
			modelID = models[0].ID
		} else {
			switch p.ID() {
			case "openrouter":
				modelID = "anthropic/claude-3.5-sonnet"
			case "nvidia":
				modelID = "nvidia/llama-3.1-nemotron-70b-instruct"
			default:
				modelID = "gpt-4o"
			}
		}
	}

	caps := p.Capabilities(modelID)
	if crit.RequiresTools && !caps.ToolCalling {
		return nil, "", fmt.Errorf("model %q does not support tool calling", modelID)
	}

	return p, modelID, nil
}
