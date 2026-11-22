---
see_also:
  - handbook/docs/orbit/architecture/contracts.md
  - handbook/docs/orbit/architecture/module-catalog.md
---

# orbit-tg

Telegram bot shell (telego). Status: `beta`. MIT — `github.com/manovaspace/orbit-tg`.

## Commands

```bash
# from the workspace root
go -C orbit/orbit-tg test ./...
```

## Architecture

- `tg.go`, `options.go`, `media.go`, `inline.go`, `commands.go` — the domain: pure Telegram types, `Markup`, the ports and the functional options. `Messenger` stays the base send/edit port; `MediaMessenger`, `InlineQuerier`, `ChatActioner` and `CommandSetter` add photos, documents, captions, inline queries, chat actions and command/menu publishing. No Telegram client, HTTP framework or adapter.
- `internal/application/router` — route predicates and dispatch over domain types.
- `internal/application/handler` — the first-line command matcher.
- `internal/application/html` — `HTMLEscape`, `SanitizeHTML` and `ChunkHTML`.
- `internal/application/ephemeral` — TTL deletion scheduling.
- `internal/application/scrub` — bounded exponential-backoff deletion.
- `internal/infrastructure/telego` — the telego client and update mapper. This is the only package that constructs a Telegram client besides the public `bot` facade, whose `API` field is part of the exported contract.
- `adapter/telego`, `bot`, `message`, `handler` — public facades; every exported name and signature is unchanged. `adapter/telego` implements the four extension ports and maps `inline_query` / `chosen_inline_result` updates.
- `keyboard`, `rich`, `mention`, `chat` — pure public helpers, already free of Telegram client imports.
- `.golangci.yml` `depguard` rules enforce the domain and application boundaries. Do not weaken them to land a change.

The kit has no database, no HTTP server of its own and publishes no events. It is a Telegram client library, so the Fiber v3 server rule does not apply.

## Docs

- Public README in this repo
- Staff handbook (private): architecture/contracts#orbit-tg

<!-- manova:conventions v1 — generated from manovaspace/ci conventions/core.md; edit it there -->
## Shared conventions

These rules apply to every Manova repository. They are copied from `manovaspace/ci` (`conventions/core.md`); change them there, not here.

### Delivery

1. Never push to `main`. Every change goes on a topic branch (`feat/<name>`, `fix/<name>`, `refactor/<name>`, `docs/<name>`, `chore/<name>`) and reaches `main` through a pull request.
2. One pull request does one thing. An unrelated problem you find on the way becomes an issue, not part of this change.
3. Pull request titles and commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`; `!` for a breaking change). Releases are computed from them.
4. Merge with squash once the checks are green. Do not use `--no-verify`; if you ever must, say so in the pull request.

### Checks

5. CI runs `manovaspace/ci@v1`. Run the same gate locally before you push; a red check is fixed, not skipped.
6. Tool versions (Go, Node, Bun, linters) come from `ci/versions.env` in `manovaspace/ci`. Do not pin a different version in this repository.
7. A behaviour change comes with a test that fails without it. Do not delete or weaken a test to make a change pass.

### Secrets

8. Never commit a secret, token, password, private key or `.env` file. Commit `.env.example` with placeholders instead.
9. Never print a secret in logs, test output, error messages or pull requests. Read secrets from the environment or a file at run time.

### Runtime

10. Containers talk to each other by name on a Docker network. Do not publish a port on the host unless a person or a host process needs it, and then bind `127.0.0.1`, never `0.0.0.0`.
11. Configuration comes from the environment. A service with a missing required setting refuses to start; it never falls back to a built-in secret or a development default in production.
12. Port numbers come from the workspace port registry. Do not invent a port.

### Documentation

13. This repository explains itself: `README.md` says how to set up, run, test and deploy it; `AGENTS.md` holds its binding rules; repository decisions live in `docs/decisions/`.
14. Write plain, short sentences. Keep exact names, paths, commands and versions. Do not document what the code already says.
<!-- /manova:conventions -->
