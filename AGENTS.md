---
read_with:
  - handbook/docs/orbit/architecture/contracts.md
  - handbook/docs/orbit/architecture/module-catalog.md
---

# orbit-tg

Telegram bot shell (telego). Status: `beta`. MIT — `github.com/manovaspace/orbit-tg`.

## Commands

```bash
go test ./...
```

## Architecture

- `tg.go`, `options.go` — the domain: pure Telegram types, `Markup`, `Messenger` and the functional options. No Telegram client, HTTP framework or adapter.
- `internal/application/router` — route predicates and dispatch over domain types.
- `internal/application/handler` — the first-line command matcher.
- `internal/application/html` — `HTMLEscape`, `SanitizeHTML` and `ChunkHTML`.
- `internal/application/ephemeral` — TTL deletion scheduling.
- `internal/application/scrub` — bounded exponential-backoff deletion.
- `internal/infrastructure/telego` — the telego client and update mapper. This is the only package that constructs a Telegram client besides the public `bot` facade, whose `API` field is part of the exported contract.
- `adapter/telego`, `bot`, `message`, `handler` — public facades; every exported name and signature is unchanged.
- `keyboard`, `rich`, `mention`, `chat` — pure public helpers, already free of Telegram client imports.
- `.golangci.yml` `depguard` rules enforce the domain and application boundaries. Do not weaken them to land a change.

The kit has no database, no HTTP server of its own and publishes no events. It is a Telegram client library, so the Fiber v3 server rule does not apply.

## Docs

- Public README in this repo
- Staff handbook (private): architecture/contracts#orbit-tg
