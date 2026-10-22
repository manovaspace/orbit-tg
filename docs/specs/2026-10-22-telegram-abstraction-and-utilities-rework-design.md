# Telegram Abstraction Layer & Utilities Rework Design

**Date:** 2026-10-22  
**Status:** Approved  
**Scope:** `orbit/orbit-tg` (Go toolkit), `orbit/orbit-staff`, `orbit/orbit-cli`, `clients/manova/tgbot` (TypeScript)

---

## 1. Executive Summary & Problem Statement

Telegram integration across the Manova workspace currently suffers from architectural drift, duplication, missing abstractions, and lifecycle leaks:

1. **Zero Abstractions in `orbit/orbit-tg`**: Method signatures and structs are coupled directly to concrete third-party types (`github.com/mymmrac/telego`). There is no abstract `Messenger`, `Markup`, `Update`, or `Router`, making unit testing impossible without live network access or brittle mocking.
2. **Brittle Formatting & Concurrency Leaks**: Methods like `Scrub` and `ReplyEphemeral` launch unmanaged goroutines (`go m.scrubWithRetry(...)`) with `context.Background()` or `time.Sleep`, bypassing caller lifecycles and leaking goroutines on shutdown. Furthermore, tying ephemeral timers to request contexts would prematurely abort deletions.
3. **Severe Feature Gaps in Keyboards**: The keyboard package is 32 lines with two hardcoded inline functions (`DismissOK`, `ConfirmRow`). There is zero support for `ReplyKeyboardMarkup`, `ReplyKeyboardRemove`, `ForceReply`, URL buttons, web app buttons, or fluent builder APIs.
4. **Consumer Fragmentation**:
   - `orbit/orbit-staff` defined an ad-hoc function signature (`type SendText func(...) error`) to insulate its domain.
   - `orbit/orbit-cli` bypassed `orbit-tg` completely, writing a manual `net/http` JSON POST client in `pkg/owner/telegram.go`.
   - `clients/manova/tgbot` implemented comprehensive HTML sanitization and chunking in TypeScript (`shared/telegram-html.ts`), but these utilities had no Go equivalent.
5. **Stale Documentation**: Existing Cursor rules (`telegram-bot-architecture.mdc`) reference legacy frameworks (`aiogram 3`) that were replaced months ago.

This design establishes a clean, unified Telegram abstraction layer and utility suite across Go and TypeScript, standardizing on **Telegram HTML** for all message bodies.

---

## 2. Text Formatting Decision: 100% Telegram HTML Standard

Telegram supports two text styling formats: **MarkdownV2** and **HTML**.

### Why We Standardize 100% on Telegram HTML
- **Escape Safety**: Telegram MarkdownV2 requires escaping 18 punctuation characters (`_`, `*`, `[`, `]`, `(`, `)`, `~`, `` ` ``, `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`). Dynamic strings (e.g. `v1.2.0`, email addresses, IP addresses, tokens) frequently break MarkdownV2 parsing, throwing fatal `400 Bad Request: can't parse entities` errors. In contrast, Telegram HTML only requires standard XML entity escaping (`&amp;`, `&lt;`, `&gt;`).
- **Feature Parity**: Telegram HTML supports all modern styling primitives: `<b>`, `<i>`, `<u>`, `<s>`, `<code>`, `<pre><code class="language-*">`, `<blockquote>`, `<blockquote expandable>`, `<tg-spoiler>`, `<tg-emoji>`, and `<a href="...">`.
- **Elimination of Telego AST**: Telego's `InputRichMessage` block builder is proprietary, verbose, and non-portable. Standardizing on Telegram HTML aligns Go with TypeScript, Python, and the native Telegram Bot API.

### Plain Text vs HTML Semantics
- **`SendText(ctx, chatID, text, ...)`**: Treats `text` as raw plain text. It automatically executes `message.HTMLEscape(text)` before dispatching with `ParseModeHTML`. This guarantees 100% crash immunity for arbitrary user or system strings (e.g. stack traces, error messages, unescaped user inputs).
- **`SendHTML(ctx, chatID, html, ...)`**: Treats `html` as pre-formatted Telegram HTML. Runs `message.SanitizeHTML(html)` to auto-balance open tags and strip unsupported tags.

---

## 3. Go Architecture: `orbit/orbit-tg`

### Package Structure
```
orbit/orbit-tg/
├── tg.go                  # Core contracts: Messenger, Message, Chat, User, Markup, Update, CallbackQuery
├── options.go             # Functional send/edit/callback options (WithThreadID, WithMarkup, etc.)
├── keyboard/
│   ├── types.go           # Abstract models: InlineKeyboard, ReplyKeyboard, Button, RemoveKeyboard, ForceReply
│   ├── builder.go         # Fluent builder APIs: NewInline(), NewReply(), Grid()
│   ├── common.go          # Standard dialogs: DismissOK, ConfirmRow, PaginationRow
│   └── keyboard_test.go   # Comprehensive builder tests
├── message/
│   ├── html.go            # Telegram HTML sanitizer, entity escaper, tag balancer, chunker
│   ├── ephemeral.go       # Lifecycle-managed auto-delete message scheduler
│   ├── scrub.go           # Safe trigger deletion with bounded exponential retries
│   └── message_test.go    # Sanitizer, chunking, and lifecycle tests
├── rich/
│   ├── rich.go            # HTML string template builders: Notice, BulletList, NumberedList, CodeBlock
│   └── rich_test.go       # Formatting and escaping tests
├── mention/
│   ├── mention.go         # UTF-16 code unit offset helpers (UTF16Len, LineUTF16Offset)
│   └── mention_test.go    # Offset tests
├── chat/
│   ├── chat.go            # Thread, chat type, and display name helpers on tg types
│   └── chat_test.go       # Chat helper tests
├── bot/
│   ├── bot.go             # Bot engine runner: long-polling & webhook support
│   ├── router.go          # Route matching: commands (/cmd), regex, callback prefixes
│   └── bot_test.go        # Lifecycle and router tests
├── adapter/
│   └── telego/            # Telego driver implementing tg.Messenger & tg.Bot
│       ├── client.go      # Telego messenger adapter
│       ├── mapper.go      # Converts abstract tg types <-> telego Bot API types
│       └── client_test.go # Integration adapter tests
└── mock/
    ├── messenger.go       # In-memory mock messenger for caller unit tests
    └── messenger_test.go  # Mock verification tests
```

---

### Core Interfaces & Models (`tg.go`)

```go
package tg

import (
	"context"
	"time"
)

// ParseMode defines the text parsing format (defaults to HTML).
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

// Chat represents an abstract Telegram chat or supergroup topic.
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"` // "private", "group", "supergroup", "channel"
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
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
```

---

### Functional Options (`options.go`)

```go
package tg

type SendConfig struct {
	ThreadID           int
	Markup             Markup
	Silent             bool
	ProtectContent     bool
	ReplyToMessageID   int
	DisableLinkPreview bool
}

type SendOption func(*SendConfig)

func WithThreadID(id int) SendOption {
	return func(c *SendConfig) { c.ThreadID = id }
}

func WithMarkup(m Markup) SendOption {
	return func(c *SendConfig) { c.Markup = m }
}

func WithSilent(silent bool) SendOption {
	return func(c *SendConfig) { c.Silent = silent }
}

func WithProtectContent(protect bool) SendOption {
	return func(c *SendConfig) { c.ProtectContent = protect }
}

func WithReplyTo(id int) SendOption {
	return func(c *SendConfig) { c.ReplyToMessageID = id }
}

func WithDisableLinkPreview(disable bool) SendOption {
	return func(c *SendConfig) { c.DisableLinkPreview = disable }
}

type EditConfig struct {
	Markup             Markup
	DisableLinkPreview bool
}

type EditOption func(*EditConfig)

func WithEditMarkup(m Markup) EditOption {
	return func(c *EditConfig) { c.Markup = m }
}

func WithEditDisableLinkPreview(disable bool) EditOption {
	return func(c *EditConfig) { c.DisableLinkPreview = disable }
}

type CallbackConfig struct {
	Text      string
	ShowAlert bool
	URL       string
	CacheTime int
}

type CallbackOption func(*CallbackConfig)

func WithAlert(alert bool) CallbackOption {
	return func(c *CallbackConfig) { c.ShowAlert = alert }
}

func WithCallbackText(text string) CallbackOption {
	return func(c *CallbackConfig) { c.Text = text }
}

func WithCallbackURL(url string) CallbackOption {
	return func(c *CallbackConfig) { c.URL = url }
}

func WithCacheTime(seconds int) CallbackOption {
	return func(c *CallbackConfig) { c.CacheTime = seconds }
}
```

---

### Keyboards & Fluent Builders (`keyboard/`)

```go
package keyboard

type InlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
	WebAppURL    string `json:"web_app_url,omitempty"`
	CopyText     string `json:"copy_text,omitempty"`
}

type InlineKeyboard struct {
	Rows [][]InlineButton `json:"rows"`
}

func (k *InlineKeyboard) MarkupType() string { return "inline" }

type ReplyButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

type ReplyKeyboard struct {
	Rows         [][]ReplyButton `json:"rows"`
	Resize       bool            `json:"resize"`
	OneTime      bool            `json:"one_time"`
	Placeholder  string          `json:"placeholder,omitempty"`
	Selective    bool            `json:"selective"`
	IsPersistent bool            `json:"is_persistent"`
}

func (k *ReplyKeyboard) MarkupType() string { return "reply" }

type RemoveKeyboard struct {
	Selective bool `json:"selective"`
}

func (k *RemoveKeyboard) MarkupType() string { return "remove" }

type ForceReply struct {
	Placeholder string `json:"placeholder,omitempty"`
	Selective   bool   `json:"selective"`
}

func (k *ForceReply) MarkupType() string { return "force_reply" }
```

#### Fluent Builder Usage
```go
// Fluent inline keyboard
kb := keyboard.NewInline().
	Row(
		keyboard.Callback("Confirm ✅", "order:confirm"),
		keyboard.Callback("Cancel ❌", "order:cancel"),
	).
	Row(
		keyboard.URL("View Invoice 📄", "https://billing.internal.manova.space/inv/123"),
	).
	Build()

// Dynamic grid builder
items := []string{"Apples", "Bananas", "Cherries", "Dates"}
grid := keyboard.NewInline().
	Grid(items, 2, func(item string) keyboard.InlineButton {
		return keyboard.Callback(item, "select:"+item)
	}).
	Build()
```

---

### Message Utilities & Lifecycle (`message/`, `rich/`, `mention/`, `chat/`)

#### 1. HTML Sanitizer, Balancer & Chunker (`message/html.go`)
- **`HTMLEscape(text string) string`**: Replaces `&`, `<`, `>` with standard entities.
- **`SanitizeHTML(html string) string`**: Strips unauthorized tags, preserves Telegram tags (`b`, `i`, `u`, `s`, `code`, `pre`, `blockquote`, `tg-spoiler`, `a`, `tg-emoji`), and automatically balances unclosed tags to prevent API rejection.
- **`ChunkHTML(html string, maxLength int) []string`**: Splits large messages (default 4000 characters) along semantic boundaries (`\n\n` -> `\n` -> `' '` -> hard cut) while balancing open/closed tags across chunk splits.

#### 2. Ephemeral Messages (`message/ephemeral.go`)
- Timed deletions are managed by a background cleanup supervisor bound to the application/messenger lifecycle, not short-lived turn contexts.
- Pending timers are tracked and safely drained or cancelled upon `Messenger.Close()`.

#### 3. Bounded-Retry Scrub (`message/scrub.go`)
- Deletes user trigger messages (commands, secrets, sensitive inputs) with configurable exponential backoff and bounded attempts (e.g. 4 attempts), managed via the lifecycle context.

#### 4. Clean HTML String Templates (`rich/rich.go`)
- Replaces proprietary `telego.InputRichMessage` with pure, portable HTML string generation:
  - `Notice(title string, sections ...Section) string`
  - `BulletList(items []string) string`
  - `NumberedList(items []string) string`
  - `CodeBlock(lang, code string) string`
  - `Blockquote(text string, expandable bool) string`

#### 5. Code Unit Offsets (`mention/mention.go`)
- Retains `UTF16Len` and `LineUTF16Offset` for accurate Telegram message entity calculations.
- Purges domain-specific checklist parsing logic (`MatchLine`, `parseLeadingIndex`).

#### 6. Chat & User Helpers (`chat/chat.go`)
- Adapts `IsGroup(c tg.Chat)`, `MessageThread(msg tg.Message)`, `CallbackThread(q tg.CallbackQuery)`, and `DisplayName(u tg.User)` to operate purely on abstract `tg` models.

---

### In-Memory Mock Driver (`mock/messenger.go`)

```go
package mock

import (
	"context"
	"sync"
	"github.com/manovaspace/orbit-tg"
)

type SentRecord struct {
	ChatID int64
	Text   string
	IsHTML bool
	Config tg.SendConfig
}

type EditRecord struct {
	ChatID    int64
	MessageID int
	Text      string
	IsHTML    bool
	Config    tg.EditConfig
}

type MockMessenger struct {
	mu           sync.Mutex
	Sent         []SentRecord
	Edited       []EditRecord
	Deleted      []int
	NextID       int
	FailWith     error
	FailCallback error
}

func NewMessenger() *MockMessenger {
	return &MockMessenger{NextID: 1000}
}

// Implements tg.Messenger fully for zero-network consumer testing.
```

---

## 4. TypeScript Architecture: `clients/manova/tgbot`

To achieve symmetry and cross-language consistency:

1. **`src/shared/telegram/types.ts`**:
   - `TelegramMessenger`: Standard interface for sending, editing, deleting, answering callbacks.
   - `InlineKeyboardMarkup`, `ReplyKeyboardMarkup`, `ReplyKeyboardRemove`, `ForceReply`: Pure data models.
2. **`src/shared/telegram/builder.ts`**:
   - `InlineKeyboardBuilder`: Fluent builder (`.row()`, `.text()`, `.url()`, `.grid()`).
   - `ReplyKeyboardBuilder`: Fluent builder (`.row()`, `.text()`, `.resize()`, `.placeholder()`).
3. **`src/shared/telegram/html.ts`**:
   - Canonical implementation of `escapeHtml`, `sanitizeTelegramHtml`, `chunkTelegramHtml`, `formatTelegramRichText`, and `formatTelegramStreamText`.
4. **`src/shared/telegram/dialogs.ts`**:
   - Standard reusable prompt keyboards (`createConfirmKeyboard`, `createDismissKeyboard`, `createPaginationKeyboard`).
5. **`src/shared/telegram-html.ts`**:
   - Thin re-export file pointing to `src/shared/telegram/html.ts` to preserve 100% backward compatibility across 20+ existing consumer files.

---

## 5. Consumer Refactoring

### 1. `orbit/orbit-staff`
- **File**: `orbit/orbit-staff/internal/infrastructure/telegram/notify.go`
  - Replace `type SendText func(...)` with direct injection of `tg.Messenger`.
  - Use `tg.SendHTML` with `tg.WithThreadID`.
- **File**: `orbit/orbit-staff/cmd/staff/main.go`
  - Wire `telego.New(...)` adapter into `telegram.New(telegram.Config{ Messenger: adapter })`.
- **File**: `orbit/orbit-staff/internal/infrastructure/telegram/notify_test.go`
  - Replace ad-hoc mock functions with `mock.NewMessenger()`.

### 2. `orbit/orbit-cli`
- **File**: `orbit/orbit-cli/pkg/owner/telegram.go`
  - Modernize manual JSON payload builder to use structured HTML string escaping and formatting.
  - Maintain standalone CI compliance without requiring sibling workspace checkouts during `test-standalone-build.sh`.

### 3. Handbook & Cursor Rules
- **File**: `handbook/cursor/rules/telegram-bot-architecture.mdc`
  - Remove references to `aiogram 3` and Telethon Python scripts.
  - Formally document the zero-spam in-place state machine, 100% Telegram HTML standard, and `tg.Messenger` interface pattern.

---

## 6. Verification & Testing Strategy

1. **Unit Testing `orbit/orbit-tg`**:
   - `keyboard`: Test fluent builders, row chaining, grid distribution, JSON serialization.
   - `message`: Test HTML escaping, illegal tag sanitization, tag balancing, message chunking at 4000 chars, ephemeral and scrub lifecycle tracking.
   - `rich`: Test notice generation, bullet/numbered lists, and code block formatting.
   - `mock`: Test recording and assertions of `MockMessenger`.
   - `adapter/telego`: Test adapter mappings from abstract types to Telego types.
2. **Regression Testing Consumers**:
   - `orbit/orbit-tg`: `go test -v ./...`
   - `orbit/orbit-staff`: `go test -v ./...`
   - `orbit/orbit-cli`: `go test -v ./...` and `bash scripts/test-standalone-build.sh`
3. **TypeScript Testing**:
   - `clients/manova/tgbot`: `bun test`
