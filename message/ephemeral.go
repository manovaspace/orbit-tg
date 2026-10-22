package message

import (
	"context"
	"sync"
	"time"
)

// DeleteFunc deletes a specific message in a chat.
type DeleteFunc func(ctx context.Context, chatID int64, messageID int) error

// EphemeralManager coordinates timed message deletions tied to the application lifecycle.
type EphemeralManager struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	deleteFn DeleteFunc
}

// NewEphemeralManager creates a lifecycle-bound ephemeral message manager.
func NewEphemeralManager(parentCtx context.Context, deleteFn DeleteFunc) *EphemeralManager {
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, cancel := context.WithCancel(parentCtx)
	return &EphemeralManager{
		ctx:      ctx,
		cancel:   cancel,
		deleteFn: deleteFn,
	}
}

// Schedule queues a message to be deleted after the given TTL.
func (m *EphemeralManager) Schedule(chatID int64, messageID int, ttl time.Duration) {
	if messageID <= 0 || m.deleteFn == nil {
		return
	}
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		timer := time.NewTimer(ttl)
		defer timer.Stop()

		select {
		case <-m.ctx.Done():
			return
		case <-timer.C:
			delCtx, delCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer delCancel()
			_ = m.deleteFn(delCtx, chatID, messageID)
		}
	}()
}

// Close cancels all pending deletion timers and waits for workers to terminate cleanly.
func (m *EphemeralManager) Close() {
	m.cancel()
	m.wg.Wait()
}
