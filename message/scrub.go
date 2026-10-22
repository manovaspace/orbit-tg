package message

import (
	"context"
	"time"
)

// Scrub attempts to delete a trigger message with bounded retries and exponential backoff.
func Scrub(ctx context.Context, chatID int64, messageID int, deleteFn DeleteFunc) error {
	if messageID <= 0 || deleteFn == nil {
		return nil
	}

	const maxAttempts = 4
	backoff := 100 * time.Millisecond

	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}

		err := deleteFn(ctx, chatID, messageID)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return lastErr
}
