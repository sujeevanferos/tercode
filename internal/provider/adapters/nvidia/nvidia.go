// Package nvidia implements the NVIDIA NIM (Inference Microservices) adapter.
package nvidia

import (
	"context"
	"fmt"
	"strings"

	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/provider/adapters/openai"
)

const defaultBaseURL = "https://integrate.api.nvidia.com/v1"

// Adapter implements provider.Provider for NVIDIA NIM endpoints.
// It wraps the OpenAI-compatible adapter with NVIDIA specific defaults.
type Adapter struct {
	*openai.Adapter
}

// New creates an NVIDIA NIM adapter.
func New(apiKey string, baseURL string) *Adapter {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	underlying := openai.New("nvidia", "NVIDIA NIM", baseURL, apiKey)
	return &Adapter{Adapter: underlying}
}

// Capabilities returns capabilities tailored for NVIDIA hosted models.
func (a *Adapter) Capabilities(modelID string) provider.Capabilities {
	return provider.Capabilities{
		Streaming:   true,
		ToolCalling: strings.Contains(strings.ToLower(modelID), "instruct") || strings.Contains(strings.ToLower(modelID), "nemotron"),
		ContextWindow: 131072,
	}
}

func (a *Adapter) Authenticate(ctx context.Context) error {
	if err := a.Adapter.Authenticate(ctx); err != nil {
		return fmt.Errorf("nvidia: authentication failed: %w", err)
	}
	return nil
}
