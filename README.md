# orbit-tg

[![CI](https://github.com/manovaspace/orbit-tg/actions/workflows/ci.yml/badge.svg)](https://github.com/manovaspace/orbit-tg/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

Go library: reusable Telegram bot shell — long-polling runner, message helpers, rich-text builders, and common inline keyboards.

Built on [telego](https://github.com/mymmrac/telego). Part of the [Manova / Orbit](https://github.com/manovaspace) open toolkit.

## Install

```bash
go get github.com/manovaspace/orbit-tg@latest
```

## Packages

| Import | Role |
| --- | --- |
| `.../orbit-tg/bot` | Long-poll lifecycle + `Messenger` accessor |
| `.../orbit-tg/message` | Send/edit/delete rich and plain messages |
| `.../orbit-tg/handler` | telegohandler predicates (`FirstLineCommand`, …) |
| `.../orbit-tg/chat` | Group/thread/user display helpers |
| `.../orbit-tg/rich` | Rich block builder + notice layout |
| `.../orbit-tg/keyboard` | Dismiss / confirm inline keyboards |
| `.../orbit-tg/mention` | UTF-16 `text_mention` resolution |

## Usage

```go
tg, err := bot.New(token, log)
if err != nil { ... }

err = tg.Run(ctx, func(bh *telegohandler.BotHandler) {
    bh.HandleMessage(onHelp, handler.FirstLineCommand("help"))
})
```

## Development

```bash
go test ./...
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md). Security reports: [SECURITY.md](./SECURITY.md).

## License

MIT — see [LICENSE](./LICENSE).
