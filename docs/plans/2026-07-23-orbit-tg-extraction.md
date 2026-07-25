# orbit-tg Extraction Implementation Plan

> **For agentic workers:** Implement task-by-task with tests after each package.

**Goal:** Extract reusable Telegram bot shell (handlers, messaging, rich text, keyboards) from `tgbot` into `orbit-tg`, then slim `tgbot` to checklist-specific adapter code.

**Architecture:** `orbit-tg` is a Go library (like `orbit-observability`) with small subpackages: `bot` (long-poll runner), `message` (send/edit/delete), `handler` (predicates), `chat` (group/thread/user helpers), `rich` (block builder), `keyboard` (inline patterns), `mention` (UTF-16 entity helpers). Product bots depend on `orbit-tg`; domain logic stays in product `internal/domain` + `internal/application`.

**Tech Stack:** Go 1.26, telego v1.11.1, telegohandler

## Global Constraints

- Module path: `github.com/manovaspace/orbit-tg`
- No product/checklist types in `orbit-tg`
- `tgbot` uses `replace` for local dev (same pattern as `orbit-observability`)
- Handbook: add `orbit-tg` to module catalog as `beta`

---

### Task 1: Scaffold orbit-tg module

**Files:** `go.mod`, `README.md`, `AGENTS.md`, `.forgejo/workflows/ci.yml`, `LICENSE`

### Task 2: chat + handler packages

**Produces:** `chat.IsGroup`, `chat.MessageThread`, `chat.DisplayName`, `handler.FirstLineCommand`

### Task 3: message package

**Produces:** `message.Messenger` with `SendText`, `SendRich`, `EditRich`, `TryDelete`, `TryDeleteMany`, `Scrub`, `ReplyEphemeral`

### Task 4: rich + keyboard packages

**Produces:** `rich.Builder`, `rich.Notice`, `rich.HTMLEscape`, `keyboard.DismissOK`, `keyboard.ConfirmRow`, `keyboard.CallbackDismiss`

### Task 5: bot runner + mention

**Produces:** `bot.Bot` with `New`, `Run`, `Messenger()`; `mention` UTF-16 helpers

### Task 6: Refactor tgbot

**Files:** split `internal/bot/` into focused files; wire `orbit-tg`; update `go.mod` replace

### Task 7: Verification

`go test ./...` in `orbit-tg` and `tgbot`
