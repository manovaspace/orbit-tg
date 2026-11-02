package message

import (
	"context"
	"time"

	"github.com/manovaspace/orbit-tg/internal/application/ephemeral"
)

// DeleteFunc deletes a specific message in a chat.
type DeleteFunc func(ctx context.Context, chatID int64, messageID int) error

// EphemeralManager coordinates timed message deletions tied to the application lifecycle.
type EphemeralManager struct {
	inner *ephemeral.EphemeralManager
}

// NewEphemeralManager creates a lifecycle-bound ephemeral message manager.
func NewEphemeralManager(parentCtx context.Context, deleteFn DeleteFunc) *EphemeralManager {
	var inner ephemeral.DeleteFunc
	if deleteFn != nil {
		inner = ephemeral.DeleteFunc(deleteFn)
	}
	return &EphemeralManager{inner: ephemeral.NewEphemeralManager(parentCtx, inner)}
}

// Schedule queues a message to be deleted after the given TTL.
func (m *EphemeralManager) Schedule(chatID int64, messageID int, ttl time.Duration) {
	if m == nil || m.inner == nil {
		return
	}
	m.inner.Schedule(chatID, messageID, ttl)
}

// Close cancels all pending deletion timers and waits for workers to terminate cleanly.
func (m *EphemeralManager) Close() {
	if m == nil || m.inner == nil {
		return
	}
	m.inner.Close()
}
