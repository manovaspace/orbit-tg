# Telegram Abstraction Layer & Utilities Rework Implementation Plan

**Spec Reference:** `docs/specs/2026-10-22-telegram-abstraction-and-utilities-rework-design.md`  
**Date:** 2026-10-22  
**Target Branches:**
- `orbit/orbit-tg`: `refactor/telegram-abstraction-and-utils`
- `orbit/orbit-staff`: `refactor/staff-telegram-abstraction` (from `feat/keycloak-staff-identity`)
- `orbit/orbit-cli`: `refactor/cli-telegram-abstraction` (from `main`)
- `clients/manova/tgbot`: `refactor/telegram-shared-utils` (from `feat/xraycraft-inbound-naming`)

---

## Phase 1: `orbit/orbit-tg` Core Contracts & Utilities

### Task 1: Core Contracts & Functional Options (`tg.go`, `options.go`)
- **Files to create/modify**:
  - `orbit/orbit-tg/tg.go`
  - `orbit/orbit-tg/options.go`
- **Responsibilities**:
  - Define `User`, `Chat`, `Message`, `CallbackQuery`, `Update`, `Markup` interface (`MarkupType() string`), and `HandlerFunc`.
  - Define `Messenger` interface with `SendText`, `SendHTML`, `EditText`, `EditHTML`, `EditKeyboard`, `Delete`, `AnswerCallback`, `Pin`, `Unpin`, `Close`.
  - Define `SendConfig`, `EditConfig`, `CallbackConfig`, and option functions (`WithThreadID`, `WithMarkup`, `WithSilent`, `WithProtectContent`, `WithReplyTo`, `WithDisableLinkPreview`, `WithAlert`, `WithCallbackText`, `WithCallbackURL`, `WithCacheTime`).
- **Verification**: `go vet ./...` in `orbit/orbit-tg`.

### Task 2: Keyboard Models, Fluent Builders & Dialogs (`keyboard/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/keyboard/types.go`
  - `orbit/orbit-tg/keyboard/builder.go`
  - `orbit/orbit-tg/keyboard/common.go`
  - `orbit/orbit-tg/keyboard/keyboard.go` (compatibility / entrypoint)
  - `orbit/orbit-tg/keyboard/keyboard_test.go`
- **Responsibilities**:
  - Implement `InlineKeyboard`, `InlineButton`, `ReplyKeyboard`, `ReplyButton`, `RemoveKeyboard`, `ForceReply`. All implement `tg.Markup`.
  - Implement `NewInline()`, `Row()`, `Button()`, `Callback()`, `URL()`, `WebApp()`, `CopyText()`, `Grid()`, and `Build()`.
  - Implement `NewReply()`, `Row()`, `Button()`, `RequestContact()`, `RequestLocation()`, and `Build()`.
  - Implement standard common dialogs: `DismissOK(label, data)`, `ConfirmRow(yesLabel, yesData, noLabel, noData)`, `PaginationRow(page, totalPages, prefix)`.
- **Verification**: `go test -v ./keyboard`

### Task 3: HTML Processing, Lifecycle Ephemeral & Scrub (`message/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/message/html.go`
  - `orbit/orbit-tg/message/ephemeral.go`
  - `orbit/orbit-tg/message/scrub.go`
  - `orbit/orbit-tg/message/message_test.go`
- **Responsibilities**:
  - `HTMLEscape(text string) string`: fast entity escaping (`&`, `<`, `>`).
  - `SanitizeHTML(html string) string`: strip invalid tags, retain allowed Telegram tags (`b`, `i`, `u`, `s`, `code`, `pre`, `blockquote`, `tg-spoiler`, `a`, `tg-emoji`), auto-balance open tags.
  - `ChunkHTML(html string, maxLength int) []string`: split along semantic boundaries (`\n\n` -> `\n` -> `' '` -> hard cut), balancing tags across chunks.
  - `ephemeral.go`: Lifecycle-bound auto-delete message scheduler with cleanup supervisor.
  - `scrub.go`: Bounded exponential-retry trigger message deletion.
- **Verification**: `go test -v ./message`

### Task 4: Clean HTML String Templates (`rich/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/rich/rich.go`
  - `orbit/orbit-tg/rich/rich_test.go`
- **Responsibilities**:
  - Eliminate all `telego.InputRichBlock` dependencies.
  - Build pure Telegram HTML templates:
    - `Notice(title string, sections ...Section) string`
    - `BulletList(items []string) string`
    - `NumberedList(items []string) string`
    - `CodeBlock(lang, code string) string`
    - `Blockquote(text string, expandable bool) string`
- **Verification**: `go test -v ./rich`

### Task 5: Entity Offsets & Chat Helpers (`mention/`, `chat/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/mention/mention.go`
  - `orbit/orbit-tg/mention/mention_test.go`
  - `orbit/orbit-tg/chat/chat.go`
  - `orbit/orbit-tg/chat/chat_test.go`
- **Responsibilities**:
  - `mention`: Retain `UTF16Len` and `LineUTF16Offset`. Remove `MatchLine` and `parseLeadingIndex`. Add tests.
  - `chat`: Refactor `IsGroup(c tg.Chat)`, `MessageThread(m tg.Message)`, `CallbackThread(q tg.CallbackQuery)`, `DisplayName(u tg.User)`. Add tests.
- **Verification**: `go test -v ./mention ./chat`

### Task 6: In-Memory Mock Driver (`mock/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/mock/messenger.go`
  - `orbit/orbit-tg/mock/messenger_test.go`
- **Responsibilities**:
  - `MockMessenger` implementing `tg.Messenger` fully.
  - Thread-safe tracking of `Sent`, `Edited`, `Deleted`, `CallbacksAnswered`, `Pinned`, `Unpinned`.
  - Configurable error injection (`FailWith`, `FailCallback`).
- **Verification**: `go test -v ./mock`

### Task 7: Telego Adapter (`adapter/telego/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/adapter/telego/client.go`
  - `orbit/orbit-tg/adapter/telego/mapper.go`
  - `orbit/orbit-tg/adapter/telego/client_test.go`
- **Responsibilities**:
  - Map `tg.Markup` to `telego.ReplyMarkup` (`InlineKeyboardMarkup`, `ReplyKeyboardMarkup`, `ReplyKeyboardRemove`, `ForceReply`).
  - Implement `tg.Messenger` wrapping `*telego.Bot`.
  - Wire `SendText` with auto-escaping and `SendHTML` with sanitization.
  - Handle graceful shutdown via `Close()`.
- **Verification**: `go test -v ./adapter/telego`

### Task 8: Bot Runner & Abstract Router (`bot/`)
- **Files to create/modify**:
  - `orbit/orbit-tg/bot/router.go`
  - `orbit/orbit-tg/bot/bot.go`
  - `orbit/orbit-tg/bot/bot_test.go`
- **Responsibilities**:
  - `Router`: command matching (`/cmd`), prefix callback routing (`col:action`), regex.
  - `Bot`: wraps polling/webhook with router and `tg.Messenger`.
- **Verification**: `go test -v ./...` in `orbit/orbit-tg`

---

## Phase 2: Consumer Migrations (`orbit-staff` & `orbit-cli`)

### Task 9: Migrate `orbit/orbit-staff`
- **Files to create/modify**:
  - `orbit/orbit-staff/internal/infrastructure/telegram/notify.go`
  - `orbit/orbit-staff/internal/infrastructure/telegram/notify_test.go`
  - `orbit/orbit-staff/cmd/staff/main.go`
- **Responsibilities**:
  - Create branch `refactor/staff-telegram-abstraction` from `feat/keycloak-staff-identity`.
  - Replace `type SendText func(...)` with `tg.Messenger`.
  - Update `Notifier.Activity` and `Notifier.Secrets` to call `tg.SendHTML`.
  - Update `main.go` to construct `telego.New(...)` adapter.
  - Update `notify_test.go` to use `mock.NewMessenger()`.
- **Verification**: `go test -v ./...` in `orbit/orbit-staff`

### Task 10: Modernize `orbit/orbit-cli`
- **Files to create/modify**:
  - `orbit/orbit-cli/pkg/owner/telegram.go`
  - `orbit/orbit-cli/pkg/owner/telegram_test.go`
- **Responsibilities**:
  - Create branch `refactor/cli-telegram-abstraction` from `main`.
  - Clean up Telegram HTML formatting and escaping.
  - Maintain standalone CI requirements without breaking `scripts/test-standalone-build.sh`.
- **Verification**:
  - `go test -v ./pkg/owner/...`
  - `bash scripts/test-standalone-build.sh`

---

## Phase 3: TypeScript Architecture (`clients/manova/tgbot`)

### Task 11: Reorganize `src/shared/telegram/`
- **Files to create/modify**:
  - `clients/manova/tgbot/src/shared/telegram/types.ts`
  - `clients/manova/tgbot/src/shared/telegram/html.ts`
  - `clients/manova/tgbot/src/shared/telegram/builder.ts`
  - `clients/manova/tgbot/src/shared/telegram/dialogs.ts`
  - `clients/manova/tgbot/src/shared/telegram/index.ts`
- **Responsibilities**:
  - Relocate HTML utilities (`escapeHtml`, `sanitizeTelegramHtml`, `chunkTelegramHtml`, `formatTelegramRichText`, `formatTelegramStreamText`) to `src/shared/telegram/html.ts`.
  - Create `InlineKeyboardBuilder` and `ReplyKeyboardBuilder`.
  - Create standard dialog helpers (`createConfirmKeyboard`, `createDismissKeyboard`, `createPaginationKeyboard`).

### Task 12: Backward Compatibility Re-export & Tests
- **Files to create/modify**:
  - `clients/manova/tgbot/src/shared/telegram-html.ts`
  - `clients/manova/tgbot/test/telegram-builder.test.ts`
- **Responsibilities**:
  - In `src/shared/telegram-html.ts`, re-export everything from `src/shared/telegram/html.ts`.
  - Add comprehensive unit tests in `test/telegram-builder.test.ts`.
- **Verification**: `bun test` in `clients/manova/tgbot`

---

## Phase 4: Documentation & Final Workspace Quality Gate

### Task 13: Modernize Cursor Rule
- **Files to modify**:
  - `handbook/cursor/rules/telegram-bot-architecture.mdc`
- **Responsibilities**:
  - Purge references to `aiogram 3` and Telethon Python scripts.
  - Document the zero-spam in-place state machine, 100% Telegram HTML standard, and `tg.Messenger` interface pattern.

### Task 14: Comprehensive Verification Across Repositories
- Run all test suites:
  - `git -C orbit/orbit-tg status` && `go test -v ./...`
  - `git -C orbit/orbit-staff status` && `go test -v ./...`
  - `git -C orbit/orbit-cli status` && `go test -v ./...` && `bash scripts/test-standalone-build.sh`
  - `git -C clients/manova/tgbot status` && `bun test`
