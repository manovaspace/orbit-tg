package message

import (
	"context"
	"log/slog"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Messenger struct {
	api *telego.Bot
	log *slog.Logger
}

func New(api *telego.Bot, log *slog.Logger) *Messenger {
	if log == nil {
		log = slog.Default()
	}
	return &Messenger{api: api, log: log}
}

func (m *Messenger) SendText(ctx context.Context, chatID int64, threadID int, text string, markup telego.ReplyMarkup) (*telego.Message, error) {
	params := tu.Message(tu.ID(chatID), text)
	if threadID > 0 {
		params = params.WithMessageThreadID(threadID)
	}
	if markup != nil {
		params = params.WithReplyMarkup(markup)
	}
	return m.api.SendMessage(ctx, params)
}

func (m *Messenger) SendRich(ctx context.Context, chatID int64, threadID int, rich telego.InputRichMessage, markup telego.ReplyMarkup) (*telego.Message, error) {
	params := (&telego.SendRichMessageParams{}).
		WithChatID(tu.ID(chatID)).
		WithRichMessage(rich)
	if threadID > 0 {
		params = params.WithMessageThreadID(threadID)
	}
	if markup != nil {
		params = params.WithReplyMarkup(markup)
	}
	return m.api.SendRichMessage(ctx, params)
}

func (m *Messenger) EditRich(ctx context.Context, chatID int64, messageID int, rich telego.InputRichMessage, markup telego.ReplyMarkup) error {
	params := (&telego.EditMessageTextParams{}).
		WithChatID(tu.ID(chatID)).
		WithMessageID(messageID).
		WithRichMessage(&rich)
	if kb, ok := markup.(*telego.InlineKeyboardMarkup); ok && kb != nil {
		params = params.WithReplyMarkup(kb)
	}
	_, err := m.api.EditMessageText(ctx, params)
	return err
}

func (m *Messenger) AnswerCallback(ctx context.Context, queryID string, text string, alert bool) error {
	if text == "" && !alert {
		return m.api.AnswerCallbackQuery(ctx, tu.CallbackQuery(queryID))
	}
	return m.api.AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: queryID,
		Text:            text,
		ShowAlert:       alert,
	})
}

func (m *Messenger) TryDelete(ctx context.Context, chatID int64, messageID int) {
	if messageID <= 0 {
		return
	}
	if err := m.api.DeleteMessage(ctx, &telego.DeleteMessageParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	}); err != nil {
		m.log.Warn("delete message failed (bot needs admin + Delete messages)",
			"chat_id", chatID, "message_id", messageID, "error", err)
	}
}

func (m *Messenger) TryDeleteMany(ctx context.Context, chatID int64, messageIDs ...int) {
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
		return
	}
	if len(ids) == 1 {
		m.TryDelete(ctx, chatID, ids[0])
		return
	}
	if err := m.api.DeleteMessages(ctx, &telego.DeleteMessagesParams{
		ChatID:     tu.ID(chatID),
		MessageIDs: ids,
	}); err != nil {
		m.log.Warn("delete messages failed; retrying one-by-one", "chat_id", chatID, "error", err)
		for _, id := range ids {
			m.TryDelete(ctx, chatID, id)
		}
	}
}

// Scrub deletes a user's trigger message when the bot has delete rights.
// ponytail: use background ctx — telegohandler may cancel handler ctx before delete finishes.
func (m *Messenger) Scrub(_ context.Context, msg telego.Message) {
	if msg.MessageID <= 0 {
		return
	}
	go m.scrubWithRetry(msg.Chat.ID, msg.MessageID)
}

func (m *Messenger) scrubWithRetry(chatID int64, messageID int) {
	const attempts = 4
	for i := range attempts {
		if i > 0 {
			time.Sleep(150 * time.Millisecond)
		}
		err := m.api.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    tu.ID(chatID),
			MessageID: messageID,
		})
		if err == nil {
			return
		}
		if i == attempts-1 {
			m.log.Error("delete user message failed (grant bot admin: Delete messages)",
				"chat_id", chatID, "message_id", messageID, "error", err)
		}
	}
}

// ReplyEphemeral sends text and deletes it after ttl (errors, hints).
func (m *Messenger) ReplyEphemeral(ctx context.Context, chatID int64, threadID int, text string, ttl time.Duration) error {
	sent, err := m.SendText(ctx, chatID, threadID, text, nil)
	if err != nil {
		return err
	}
	go func(messageID int) {
		time.Sleep(ttl)
		m.TryDelete(context.Background(), chatID, messageID)
	}(sent.MessageID)
	return nil
}
