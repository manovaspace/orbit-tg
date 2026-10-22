# Telegram Abstraction Layer & Utilities Rework Design

**Date:** 2026-10-22  
**Status:** Approved  
**Scope:** `orbit/orbit-tg` (Go toolkit), `orbit/orbit-staff`, `orbit/orbit-cli`, `clients/manova/tgbot` (TypeScript)

---

## 1. Executive Summary & Problem Statement

Telegram integration across the Manova workspace currently suffers from significant architectural drift, duplication, and missing abstractions:

1. **Zero Abstractions in `orbit/orbit-tg`**: Every struct, method parameter, and return type is coupled directly to third-party concrete types (`github.com/mymmrac/telego`). There is no `Messenger` or `Keyboard` interface, making unit testing impossible without live network access or complicated mocking.
2. **Brittle Formatting & Concurrency Leaks**: Methods like `Scrub` and `ReplyEphemeral` launch unmanaged goroutines (`go m.scrubWithRetry(...)`) with `context.Background()` and `time.Sleep`, bypassing caller lifecycles and leaking goroutines on shutdown.
3. **Severe Feature Gaps in Keyboards**: The keyboard package is 32 lines with two hardcoded inline functions (`DismissOK`, `ConfirmRow`). There is zero support for `ReplyKeyboardMarkup`, `ReplyKeyboardRemove`, `ForceReply`, URL buttons, web app buttons, or fluent builder APIs.
4. **Consumer Fragmentation**:
   - `orbit/orbit-staff` had to define an ad-hoc function signature (`type SendText func(...) error`) to decouple its domain.
   - `orbit/orbit-cli` bypassed `orbit-tg` completely, writing a manual `net/http` JSON POST client in `pkg/owner/telegram.go`.
   - `clients/manova/tgbot` implemented comprehensive HTML sanitization and chunking in TypeScript (`shared/telegram-html.ts`), but these battle-tested utilities have no equivalent in Go.
5. **Stale Documentation**: Existing Cursor rules (`telegram-bot-architecture.mdc`) reference legacy frameworks (`aiogram 3`) that were replaced months ago.

This design establishes a clean, unified Telegram abstraction layer and utility suite across both Go and TypeScript, standardizing on **Telegram HTML** for all message bodies.

---

## 2. Text Formatting Decision: 100% Telegram HTML Standard

Telegram supports two text styling formats: **MarkdownV2** and **HTML**.

### Why We Standardize 100% on Telegram HTML
- **Escape Safety**: Telegram MarkdownV2 requires escaping 18 punctuation characters (`_`, `*`, `[`, `]`, `(`, `)`, `~`, `` ` ``, `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`). Dynamic strings (e.g. `v1.2.0`, email addresses, IP addresses, tokens) frequently break MarkdownV2 parsing, throwing fatal `400 Bad Request: can't parse entities` errors. In contrast, Telegram HTML only requires standard XML entity escaping (`&amp;`, `&lt;`, `&gt;`).
- **Feature Parity**: Telegram HTML supports all modern styling primitives: `<b>`, `<i>`, `<u>`, `<s>`, `<code>`, `<pre><code class="language-*">`, `<blockquote>`, `<blockquote expandable>`, `<tg-spoiler>`, `<tg-emoji>`, and `<a href="...">`.
- **Elimination of Telego AST**: Telego's `InputRichMessage` block builder is proprietary, verbose, and non-portable. Standardizing on Telegram HTML aligns Go with TypeScript, Python, and the native Telegram Bot API.

**Contract Rule**: All internal bot messages, operational alerts, notifications, and templates MUST use Telegram HTML. Markdown input from agents or LLMs is parsed and converted to Telegram HTML prior to transmission.

---

## 3. Go Architecture: `orbit/orbit-tg`

### Package Structure
```
orbit/orbit-tg/
├── tg.go                  # Core contracts: Messenger, Message, Chat, User, Markup
├── options.go             # Functional send/edit options (WithThreadID, WithKeyboard, etc.)
├── keyboard/
│   ├── types.go           # Abstract models: InlineKeyboard, ReplyKeyboard, Button
│   ├── builder.go         # Fluent builder APIs: NewInline(), NewReply(), Grid()
│   ├── common.go          # Standard dialogs: Confirm, Dismiss, Pagination
│   └── keyboard_test.go   # Comprehensive builder tests
├── message/
│   ├── html.go            # Telegram HTML sanitizer, entity escaper, tag balancer, chunker
│   ├── ephemeral.go       # Context-bound auto-delete message lifecycle
│   ├── scrub.go           # Safe trigger deletion with bounded exponential retries
│   └── message_test.go    # Sanitizer, chunking, and lifecycle tests
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

// Markup represents any attached keyboard or interaction element.
type Markup interface {
    isMarkup()
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
}
```

---

### Keyboards & Fluent Builders (`keyboard/`)

```go
package keyboard

// InlineButton represents an individual interactive inline button.
type InlineButton struct {
    Text         string `json:"text"`
    CallbackData string `json:"callback_data,omitempty"`
    URL          string `json:"url,omitempty"`
    WebAppURL    string `json:"web_app_url,omitempty"`
    CopyText     string `json:"copy_text,omitempty"`
}

// InlineKeyboard represents an inline keyboard markup.
type InlineKeyboard struct {
    Rows [][]InlineButton `json:"rows"`
}

func (k *InlineKeyboard) isMarkup() {}

// ReplyButton represents a persistent chat keyboard button.
type ReplyButton struct {
    Text            string `json:"text"`
    RequestContact  bool   `json:"request_contact,omitempty"`
    RequestLocation bool   `json:"request_location,omitempty"`
}

// ReplyKeyboard represents a regular chat keyboard.
type ReplyKeyboard struct {
    Rows            [][]ReplyButton `json:"rows"`
    Resize          bool            `json:"resize"`
    OneTime         bool            `json:"one_time"`
    Placeholder     string          `json:"placeholder,omitempty"`
    Selective       bool            `json:"selective"`
    IsPersistent    bool            `json:"is_persistent"`
}

func (k *ReplyKeyboard) isMarkup() {}

// RemoveKeyboard instructs clients to dismiss the custom reply keyboard.
type RemoveKeyboard struct {
    Selective bool `json:"selective"`
}

func (k *RemoveKeyboard) isMarkup() {}

// ForceReply prompts the user to reply directly to the message.
type ForceReply struct {
    Placeholder string `json:"placeholder,omitempty"`
    Selective   bool   `json:"selective"`
}

func (k *ForceReply) isMarkup() {}
```

#### Fluent Builder Usage
```go
// Chained inline keyboard construction
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

### Telegram Utilities (`message/`)

#### 1. HTML Sanitizer, Balancer & Chunker (`html.go`)
- **`HTMLEscape(text string) string`**: Escapes bare `&`, `<`, `>`.
- **`SanitizeHTML(html string) string`**: Strips illegal tags, preserves allowed Telegram tags (`b`, `i`, `u`, `s`, `code`, `pre`, `blockquote`, `tg-spoiler`, `a`, `tg-emoji`), and balances unclosed open tags automatically.
- **`ChunkHTML(html string, maxLength int) []string`**: Splits long messages (default 4000 characters) along semantic boundaries (`\n\n` -> `\n` -> `' '` -> hard split) while preserving and balancing tags across chunk boundaries.

#### 2. Ephemeral Messages (`ephemeral.go`)
- **Context-Bound Lifecycle**: When sending an ephemeral hint or status update, the message deletion timer is tied to the parent context. If the server shuts down or the context is cancelled, pending timers are cleared cleanly without leaking background goroutines.

#### 3. Bounded-Retry Scrub (`scrub.go`)
- Deletes user trigger messages (commands, secrets, unwanted triggers) with configurable bounded retries and exponential backoff, respecting context cancellation.

---

### Mock Driver for Consumer Unit Tests (`mock/`)

```go
package mock

import (
    "context"
    "sync"
    "github.com/manovaspace/orbit-tg"
)

type MockMessenger struct {
    mu       sync.Mutex
    Sent     []SentRecord
    Edited   []EditRecord
    Deleted  []int
    FailWith error
}

func NewMessenger() *MockMessenger {
    return &MockMessenger{}
}

// Implements tg.Messenger completely for zero-network testing in orbit-staff and orbit-cli.
```

---

## 4. TypeScript Architecture: `clients/manova/tgbot`

To achieve symmetry and cross-language consistency:

1. **`src/shared/telegram/types.ts`**:
   - `TelegramMessenger`: Standard interface for sending, editing, deleting, answering callbacks.
   - `InlineKeyboard`: Pure data model representation.
   - `ReplyKeyboard`: Pure data model representation.
2. **`src/shared/telegram/builder.ts`**:
   - `InlineKeyboardBuilder`: Fluent builder (`.row()`, `.button()`, `.url()`, `.grid()`).
   - `ReplyKeyboardBuilder`: Fluent builder (`.row()`, `.button()`, `.resize()`, `.placeholder()`).
3. **`src/shared/telegram/html.ts`**:
   - Standard re-export of `sanitizeTelegramHtml`, `chunkTelegramHtml`, `formatTelegramRichText`, and `formatTelegramStreamText` from `shared/telegram-html.ts`.
4. **`src/shared/telegram/dialogs.ts`**:
   - Common dialog components (`createConfirmKeyboard`, `createDismissKeyboard`, `createPaginationKeyboard`).

---

## 5. Consumer Refactoring

### 1. `orbit/orbit-staff`
- **File**: `orbit/orbit-staff/internal/infrastructure/telegram/notify.go`
  - Replace the custom `type SendText func(...)` with direct injection of `tg.Messenger`.
  - Implement notification dispatch using `tg.SendHTML`.
- **File**: `orbit/orbit-staff/cmd/staff/main.go`
  - Instantiate `tg.Messenger` via `telegoadapter.New(token, log)`.
- **File**: `orbit/orbit-staff/internal/infrastructure/telegram/notify_test.go`
  - Replace custom dummy functions with `mock.NewMessenger()`.

### 2. `orbit/orbit-cli`
- **File**: `orbit/orbit-cli/pkg/owner/telegram.go`
  - Remove manual `net/http` client making JSON requests to `https://api.telegram.org/...`.
  - Wire `orbit-tg`'s `tg.Messenger`.
  - Provide fallback configuration resolution via `ResolveTelegramConfig()`.

### 3. Handbook & Cursor Rules
- Update `handbook/cursor/rules/telegram-bot-architecture.mdc`:
  - Remove references to `aiogram 3` and Telethon Python scripts.
  - Formally document the zero-spam in-place state machine, Telegram HTML standards, and `tg.Messenger` interfaces.

---

## 6. Verification & Testing Strategy

1. **Unit Testing `orbit/orbit-tg`**:
   - `keyboard`: Test fluent builders, row chaining, grid distribution, JSON serialization.
   - `message`: Test HTML escaping, illegal tag sanitization, tag balancing, message chunking at 4000 chars.
   - `mock`: Test recording and assertions of `MockMessenger`.
   - `adapter/telego`: Test adapter mappings from abstract types to Telego types.
2. **Regression Testing Consumers**:
   - Run `go test -v ./...` in `orbit/orbit-tg`.
   - Run `go test -v ./...` in `orbit/orbit-staff`.
   - Run `go test -v ./...` in `orbit/orbit-cli`.
3. **TypeScript Testing**:
   - Run `bun test` in `clients/manova/tgbot`. Verify new builder and dialog tests pass.
