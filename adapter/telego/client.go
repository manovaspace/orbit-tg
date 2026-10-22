package telego

import (
	"context"
	"log/slog"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/message"
)

// TelegoAdapter implements tg.Messenger wrapping the telego Telegram Bot client.
type TelegoAdapter struct {
	api       *telego.Bot
	log       *slog.Logger
	ephemeral *message.EphemeralManager
}

var _ tg.Messenger = (*TelegoAdapter)(nil)

// New initializes a TelegoAdapter with a bot token.
func New(token string, log *slog.Logger) (*TelegoAdapter, error) {
	api, err := telego.NewBot(token)
	if err != nil {
		return nil, err
	}
	return NewFromBot(api, log), nil
}

// NewFromBot wraps an existing *telego.Bot into a TelegoAdapter.
func NewFromBot(api *telego.Bot, log *slog.Logger) *TelegoAdapter {
	if log == nil {
		log = slog.Default()
	}
	adapter := &TelegoAdapter{
		api: api,
		log: log,
	}
	adapter.ephemeral = message.NewEphemeralManager(context.Background(), adapter.DeleteSingle)
	return adapter
}

// API exposes the underlying *telego.Bot if advanced driver features are needed.
func (a *TelegoAdapter) API() *telego.Bot {
	return a.api
}

// SendText treats text as raw plain text, escaping it to prevent HTML injection errors.
func (a *TelegoAdapter) SendText(ctx context.Context, chatID int64, text string, opts ...tg.SendOption) (*tg.Message, error) {
	return a.sendFormatted(ctx, chatID, message.HTMLEscape(text), opts...)
}

// SendHTML treats html as pre-formatted Telegram HTML, auto-balancing open tags.
func (a *TelegoAdapter) SendHTML(ctx context.Context, chatID int64, html string, opts ...tg.SendOption) (*tg.Message, error) {
	return a.sendFormatted(ctx, chatID, message.SanitizeHTML(html), opts...)
}

func (a *TelegoAdapter) sendFormatted(ctx context.Context, chatID int64, content string, opts ...tg.SendOption) (*tg.Message, error) {
	cfg := tg.ApplySendOptions(opts...)

	params := &telego.SendMessageParams{
		ChatID:    tu.ID(chatID),
		Text:      content,
		ParseMode: string(tg.ParseModeHTML),
	}

	if cfg.ThreadID > 0 {
		params.MessageThreadID = cfg.ThreadID
	}
	if cfg.Markup != nil {
		params.ReplyMarkup = MapMarkup(cfg.Markup)
	}
	if cfg.Silent {
		params.DisableNotification = true
	}
	if cfg.ProtectContent {
		params.ProtectContent = true
	}
	if cfg.ReplyToMessageID > 0 {
		params.ReplyParameters = &telego.ReplyParameters{
			MessageID: cfg.ReplyToMessageID,
		}
	}
	if cfg.DisableLinkPreview {
		params.LinkPreviewOptions = &telego.LinkPreviewOptions{
			IsDisabled: true,
		}
	}

	sent, err := a.api.SendMessage(ctx, params)
	if err != nil {
		return nil, err
	}
	return MapMessage(sent), nil
}

// EditText edits an existing message with escaped plain text.
func (a *TelegoAdapter) EditText(ctx context.Context, chatID int64, messageID int, text string, opts ...tg.EditOption) (*tg.Message, error) {
	return a.editFormatted(ctx, chatID, messageID, message.HTMLEscape(text), opts...)
}

// EditHTML edits an existing message with sanitized Telegram HTML.
func (a *TelegoAdapter) EditHTML(ctx context.Context, chatID int64, messageID int, html string, opts ...tg.EditOption) (*tg.Message, error) {
	return a.editFormatted(ctx, chatID, messageID, message.SanitizeHTML(html), opts...)
}

func (a *TelegoAdapter) editFormatted(ctx context.Context, chatID int64, messageID int, content string, opts ...tg.EditOption) (*tg.Message, error) {
	cfg := tg.ApplyEditOptions(opts...)

	params := &telego.EditMessageTextParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Text:      content,
		ParseMode: string(tg.ParseModeHTML),
	}

	if cfg.Markup != nil {
		if kb, ok := MapMarkup(cfg.Markup).(*telego.InlineKeyboardMarkup); ok {
			params.ReplyMarkup = kb
		}
	}
	if cfg.DisableLinkPreview {
		params.LinkPreviewOptions = &telego.LinkPreviewOptions{
			IsDisabled: true,
		}
	}

	edited, err := a.api.EditMessageText(ctx, params)
	if err != nil {
		return nil, err
	}
	return MapMessage(edited), nil
}

// EditKeyboard updates the inline keyboard markup of an existing message.
func (a *TelegoAdapter) EditKeyboard(ctx context.Context, chatID int64, messageID int, markup tg.Markup) error {
	params := &telego.EditMessageReplyMarkupParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	}
	if markup != nil {
		if kb, ok := MapMarkup(markup).(*telego.InlineKeyboardMarkup); ok {
			params.ReplyMarkup = kb
		}
	}
	_, err := a.api.EditMessageReplyMarkup(ctx, params)
	return err
}

// DeleteSingle deletes a single message.
func (a *TelegoAdapter) DeleteSingle(ctx context.Context, chatID int64, messageID int) error {
	if messageID <= 0 {
		return nil
	}
	return a.api.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	})
}

// Delete deletes one or more messages.
func (a *TelegoAdapter) Delete(ctx context.Context, chatID int64, messageIDs ...int) error {
	ids := make([]int, 0, len(messageIDs))
	seen := make(map[int]struct{}, len(messageIDs))
	for _, id := range messageIDs {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil
	}
	if len(ids) == 1 {
		return a.DeleteSingle(ctx, chatID, ids[0])
	}

	err := a.api.DeleteMessages(ctx, &telego.DeleteMessagesParams{
		ChatID:     tu.ID(chatID),
		MessageIDs: ids,
	})
	if err != nil {
		// Fallback to one-by-one deletion
		var lastErr error
		for _, id := range ids {
			if dErr := a.DeleteSingle(ctx, chatID, id); dErr != nil {
				lastErr = dErr
			}
		}
		return lastErr
	}
	return nil
}

// AnswerCallback sends a response to an inline button callback query.
func (a *TelegoAdapter) AnswerCallback(ctx context.Context, queryID string, opts ...tg.CallbackOption) error {
	cfg := tg.ApplyCallbackOptions(opts...)
	params := &telego.AnswerCallbackQueryParams{
		CallbackQueryID: queryID,
		Text:            cfg.Text,
		ShowAlert:       cfg.ShowAlert,
		URL:             cfg.URL,
		CacheTime:       cfg.CacheTime,
	}
	return a.api.AnswerCallbackQuery(ctx, params)
}

// Pin pins a message in a chat or supergroup topic.
func (a *TelegoAdapter) Pin(ctx context.Context, chatID int64, messageID int, silent bool) error {
	return a.api.PinChatMessage(ctx, &telego.PinChatMessageParams{
		ChatID:              tu.ID(chatID),
		MessageID:           messageID,
		DisableNotification: silent,
	})
}

// Unpin unpins a message in a chat or supergroup topic.
func (a *TelegoAdapter) Unpin(ctx context.Context, chatID int64, messageID int) error {
	return a.api.UnpinChatMessage(ctx, &telego.UnpinChatMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	})
}

// ReplyEphemeral sends an HTML message and automatically queues its deletion after ttl.
func (a *TelegoAdapter) ReplyEphemeral(ctx context.Context, chatID int64, html string, ttl time.Duration, opts ...tg.SendOption) (*tg.Message, error) {
	sent, err := a.SendHTML(ctx, chatID, html, opts...)
	if err != nil {
		return nil, err
	}
	a.ephemeral.Schedule(chatID, sent.ID, ttl)
	return sent, nil
}

// Scrub safely deletes a trigger message using bounded exponential retry.
func (a *TelegoAdapter) Scrub(ctx context.Context, chatID int64, messageID int) error {
	return message.Scrub(ctx, chatID, messageID, a.DeleteSingle)
}

// Close gracefully closes the adapter and waits for pending ephemeral deletions.
func (a *TelegoAdapter) Close() error {
	if a.ephemeral != nil {
		a.ephemeral.Close()
	}
	return nil
}
