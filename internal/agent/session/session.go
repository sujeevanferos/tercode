// Package session implements agent session state and history tracking.
package session

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sujeevanferos/tercode/internal/provider"
)

// Mode is the operating mode of the agent (Chat, Plan, Code, Review, Debug).
type Mode string

const (
	ModeChat   Mode = "chat"
	ModePlan   Mode = "plan"
	ModeCode   Mode = "code"
	ModeReview Mode = "review"
	ModeDebug  Mode = "debug"
)

// Session represents an active conversational and task session.
type Session struct {
	ID          string             `json:"id"`
	WorkspaceID string             `json:"workspace_id"`
	Mode        Mode               `json:"mode"`
	ProviderID  string             `json:"provider_id"`
	ModelID     string             `json:"model_id"`
	History     []provider.Message `json:"history"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// Manager manages in-memory active sessions.
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewManager creates a session manager.
func NewManager() *Manager {
	return &Manager{sessions: make(map[string]*Session)}
}

// Create creates a new Session.
func (m *Manager) Create(workspaceID, providerID, modelID string, mode Mode) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mode == "" {
		mode = ModeCode
	}

	s := &Session{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Mode:        mode,
		ProviderID:  providerID,
		ModelID:     modelID,
		History:     make([]provider.Message, 0),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	m.sessions[s.ID] = s
	return s
}

// Get returns a session by ID.
func (m *Manager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

// Append adds messages to session history.
func (m *Manager) Append(id string, msgs ...provider.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[id]; ok {
		s.History = append(s.History, msgs...)
		s.UpdatedAt = time.Now().UTC()
	}
}
