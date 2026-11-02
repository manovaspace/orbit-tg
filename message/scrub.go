package message

import (
	"context"

	"github.com/manovaspace/orbit-tg/internal/application/scrub"
)

// Scrub attempts to delete a trigger message with bounded retries and exponential backoff.
func Scrub(ctx context.Context, chatID int64, messageID int, deleteFn DeleteFunc) error {
	if deleteFn == nil {
		return scrub.Scrub(ctx, chatID, messageID, nil)
	}
	return scrub.Scrub(ctx, chatID, messageID, scrub.DeleteFunc(deleteFn))
}
