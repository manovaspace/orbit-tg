package tg

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ParseMode defines the text parsing format (standardized to HTML).
type ParseMode string

const (
	ParseModeHTML ParseMode = "HTML"
)

// User represents an abstract Telegram user.
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	IsBot     bool   `json:"is_bot"`
}

// DisplayName returns username if present, otherwise concatenated name or user ID string.
func (u User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	if name != "" {
		return name
	}
	return fmt.Sprintf("%d", u.ID)
}

// Chat represents an abstract Telegram chat, channel, or supergroup topic.
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"` // "private", "group", "supergroup", "channel"
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
}

// IsGroup returns true if the chat is a group or supergroup.
func (c Chat) IsGroup() bool {
	return c.Type == "group" || c.Type == "supergroup"
}

// Message represents an abstract sent or received Telegram message.
type Message struct {
	ID       int       `json:"id"`
	Chat     Chat      `json:"chat"`
	ThreadID int       `json:"thread_id,omitempty"`
	From     *User     `json:"from,omitempty"`
	Date     time.Time `json:"date"`
	Text     string    `json:"text"`
}

// CallbackQuery represents an incoming callback query from an inline keyboard.
type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data,omitempty"`
}

// Update represents an incoming Telegram update.
type Update struct {
	ID            int            `json:"id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// HandlerFunc handles an abstract Telegram update.
type HandlerFunc func(ctx context.Context, u *Update) error

// Markup represents any attached keyboard or interaction element.
type Markup interface {
	MarkupType() string // "inline", "reply", "remove", "force_reply"
}

// Messenger defines the primary contract for interacting with Telegram.
type Messenger interface {
	SendText(ctx context.Context, chatID int64, text string, opts ...SendOption) (*Message, error)
	SendHTML(ctx context.Context, chatID int64, html string, opts ...SendOption) (*Message, error)
	EditText(ctx context.Context, chatID int64, messageID int, text string, opts ...EditOption) (*Message, error)
	EditHTML(ctx context.Context, chatID int64, messageID int, html string, opts ...EditOption) (*Message, error)
	EditKeyboard(ctx context.Context, chatID int64, messageID int, markup Markup) error
	Delete(ctx context.Context, chatID int64, messageIDs ...int) error
	AnswerCallback(ctx context.Context, queryID string, opts ...CallbackOption) error
	Pin(ctx context.Context, chatID int64, messageID int, silent bool) error
	Unpin(ctx context.Context, chatID int64, messageID int) error
	Close() error
}
