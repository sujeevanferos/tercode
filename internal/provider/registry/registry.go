package registry

import (
	"fmt"
	"sync"

	"github.com/sujeevanferos/tercode/internal/provider"
)

// Registry holds all registered providers and allows lookup by ID.
// It is the single source of truth for which providers are available at runtime.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]provider.Provider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]provider.Provider)}
}

// Register adds a provider. Returns an error if a provider with the same ID
// is already registered.
func (r *Registry) Register(p provider.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.providers[p.ID()]; exists {
		return fmt.Errorf("provider registry: %q already registered", p.ID())
	}
	r.providers[p.ID()] = p
	return nil
}

// Get returns the provider with the given ID, or an error if not found.
func (r *Registry) Get(id string) (provider.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	if !ok {
		return nil, fmt.Errorf("provider registry: %q not found", id)
	}
	return p, nil
}

// All returns a snapshot of all registered providers.
func (r *Registry) All() []provider.Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]provider.Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result
}

// IDs returns all registered provider IDs.
func (r *Registry) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.providers))
	for id := range r.providers {
		ids = append(ids, id)
	}
	return ids
}
