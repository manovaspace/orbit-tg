package telego

import (
	"context"
	"log/slog"
	"time"

	"github.com/mymmrac/telego"

	tg "github.com/manovaspace/orbit-tg"
	infra "github.com/manovaspace/orbit-tg/internal/infrastructure/telego"
)

// TelegoAdapter implements tg.Messenger wrapping the telego Telegram Bot client.
type TelegoAdapter struct {
	inner *infra.Adapter
}

var _ tg.Messenger = (*TelegoAdapter)(nil)

// New initializes a TelegoAdapter with a bot token.
func New(token string, log *slog.Logger) (*TelegoAdapter, error) {
	api, err := telego.NewBot(token, telego.WithRequestConstructor(inlineQueryConstructor{}))
	if err != nil {
		return nil, err
	}
	return NewFromBot(api, log), nil
}

// NewFromBot wraps an existing *telego.Bot into a TelegoAdapter.
func NewFromBot(api *telego.Bot, log *slog.Logger) *TelegoAdapter {
	return &TelegoAdapter{inner: infra.NewFromBot(api, log)}
}

// API exposes the underlying *telego.Bot if advanced driver features are needed.
func (a *TelegoAdapter) API() *telego.Bot {
	return a.inner.API()
}

// SendText treats text as raw plain text, escaping it to prevent HTML injection errors.
func (a *TelegoAdapter) SendText(ctx context.Context, chatID int64, text string, opts ...tg.SendOption) (*tg.Message, error) {
	return a.inner.SendText(ctx, chatID, text, opts...)
}

// SendHTML treats html as pre-formatted Telegram HTML, auto-balancing open tags.
func (a *TelegoAdapter) SendHTML(ctx context.Context, chatID int64, html string, opts ...tg.SendOption) (*tg.Message, error) {
	return a.inner.SendHTML(ctx, chatID, html, opts...)
}

// EditText edits an existing message with escaped plain text.
func (a *TelegoAdapter) EditText(ctx context.Context, chatID int64, messageID int, text string, opts ...tg.EditOption) (*tg.Message, error) {
	return a.inner.EditText(ctx, chatID, messageID, text, opts...)
}

// EditHTML edits an existing message with sanitized Telegram HTML.
func (a *TelegoAdapter) EditHTML(ctx context.Context, chatID int64, messageID int, html string, opts ...tg.EditOption) (*tg.Message, error) {
	return a.inner.EditHTML(ctx, chatID, messageID, html, opts...)
}

// EditKeyboard updates the inline keyboard markup of an existing message.
func (a *TelegoAdapter) EditKeyboard(ctx context.Context, chatID int64, messageID int, markup tg.Markup) error {
	return a.inner.EditKeyboard(ctx, chatID, messageID, markup)
}

// DeleteSingle deletes a single message.
func (a *TelegoAdapter) DeleteSingle(ctx context.Context, chatID int64, messageID int) error {
	return a.inner.DeleteSingle(ctx, chatID, messageID)
}

// Delete deletes one or more messages.
func (a *TelegoAdapter) Delete(ctx context.Context, chatID int64, messageIDs ...int) error {
	return a.inner.Delete(ctx, chatID, messageIDs...)
}

// AnswerCallback sends a response to an inline button callback query.
func (a *TelegoAdapter) AnswerCallback(ctx context.Context, queryID string, opts ...tg.CallbackOption) error {
	return a.inner.AnswerCallback(ctx, queryID, opts...)
}

// Pin pins a message in a chat or supergroup topic.
func (a *TelegoAdapter) Pin(ctx context.Context, chatID int64, messageID int, silent bool) error {
	return a.inner.Pin(ctx, chatID, messageID, silent)
}

// Unpin unpins a message in a chat or supergroup topic.
func (a *TelegoAdapter) Unpin(ctx context.Context, chatID int64, messageID int) error {
	return a.inner.Unpin(ctx, chatID, messageID)
}

// ReplyEphemeral sends an HTML message and automatically queues its deletion after ttl.
func (a *TelegoAdapter) ReplyEphemeral(ctx context.Context, chatID int64, html string, ttl time.Duration, opts ...tg.SendOption) (*tg.Message, error) {
	return a.inner.ReplyEphemeral(ctx, chatID, html, ttl, opts...)
}

// Scrub safely deletes a trigger message using bounded exponential retry.
func (a *TelegoAdapter) Scrub(ctx context.Context, chatID int64, messageID int) error {
	return a.inner.Scrub(ctx, chatID, messageID)
}

// Close gracefully closes the adapter and waits for pending ephemeral deletions.
func (a *TelegoAdapter) Close() error {
	return a.inner.Close()
}
