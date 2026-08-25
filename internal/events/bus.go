package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Handler is a function called when a matching event is published.
// Handlers are called synchronously on the publishing goroutine;
// they must not block for extended periods.
type Handler func(Event)

// Filter determines which events a subscription receives.
// A nil Filter matches all events.
type Filter func(Event) bool

// Subscription represents an active event subscription.
// Call Cancel to stop receiving events and free resources.
type Subscription struct {
	id     string
	cancel func()
}

// Cancel removes the subscription from the bus.
func (s *Subscription) Cancel() { s.cancel() }

// Bus is a synchronous, in-process event bus.
//
// Events are delivered to matching subscribers on the goroutine that calls
// Publish. Subscribers should offload heavy work to separate goroutines.
//
// The Bus is safe for concurrent use.
type Bus struct {
	mu   sync.RWMutex
	subs map[string]subscription
}

type subscription struct {
	id      string
	filter  Filter
	handler Handler
}

// NewBus creates an initialized event bus.
func NewBus() *Bus {
	return &Bus{subs: make(map[string]subscription)}
}

// New creates a new Event with a unique ID and current timestamp, setting the
// provided type and payload. WorkspaceID, SessionID, and TaskID must be set
// by the caller when applicable.
func New(t Type, payload any) Event {
	return Event{
		ID:        uuid.New().String(),
		Type:      t,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

// Publish dispatches an event to all matching subscribers.
// It returns the event for convenience (so callers can log it).
func (b *Bus) Publish(ev Event) Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.subs {
		if s.filter == nil || s.filter(ev) {
			s.handler(ev)
		}
	}
	return ev
}

// Subscribe registers a handler to receive events matching filter.
// If filter is nil, every event is delivered. Returns a Subscription that
// must be cancelled when no longer needed to prevent leaks.
func (b *Bus) Subscribe(filter Filter, handler Handler) *Subscription {
	id := uuid.New().String()
	sub := subscription{id: id, filter: filter, handler: handler}

	b.mu.Lock()
	b.subs[id] = sub
	b.mu.Unlock()

	return &Subscription{
		id: id,
		cancel: func() {
			b.mu.Lock()
			delete(b.subs, id)
			b.mu.Unlock()
		},
	}
}

// SubscribeType returns a Subscription that only receives events of the given type.
func (b *Bus) SubscribeType(t Type, handler Handler) *Subscription {
	return b.Subscribe(func(ev Event) bool { return ev.Type == t }, handler)
}

// Wait blocks until an event matching filter is published, the context is
// cancelled, or the optional timeout elapses. It returns the event and nil
// on success, or a zero Event and an error if the context was cancelled.
func (b *Bus) Wait(ctx context.Context, filter Filter) (Event, error) {
	ch := make(chan Event, 1)
	sub := b.Subscribe(filter, func(ev Event) {
		select {
		case ch <- ev:
		default:
		}
	})
	defer sub.Cancel()

	select {
	case ev := <-ch:
		return ev, nil
	case <-ctx.Done():
		return Event{}, fmt.Errorf("events: wait cancelled: %w", ctx.Err())
	}
}
