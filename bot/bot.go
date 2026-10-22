package bot

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	tg "github.com/manovaspace/orbit-tg"
	telegoadapter "github.com/manovaspace/orbit-tg/adapter/telego"
)

// Bot provides high-level bot execution, routing, and messaging.
type Bot struct {
	API     *telego.Bot
	Msg     tg.Messenger
	Adapter *telegoadapter.TelegoAdapter
	Router  *Router
	log     *slog.Logger
}

// New initializes a Bot instance configured with the Telego driver.
func New(token string, log *slog.Logger) (*Bot, error) {
	if log == nil {
		log = slog.Default()
	}
	adapter, err := telegoadapter.New(token, log)
	if err != nil {
		return nil, err
	}
	return &Bot{
		API:     adapter.API(),
		Msg:     adapter,
		Adapter: adapter,
		Router:  NewRouter(),
		log:     log,
	}, nil
}

// HandlerRegistrar is a callback function for attaching custom low-level handlers if needed.
type HandlerRegistrar func(b *Bot, bh *th.BotHandler)

// Run starts long-polling and routes updates through the Router.
func (b *Bot) Run(ctx context.Context, registrars ...HandlerRegistrar) error {
	defer func() {
		_ = b.Adapter.Close()
	}()

	updates, err := b.API.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}
	bh, err := th.NewBotHandler(b.API, updates)
	if err != nil {
		return err
	}
	defer func() { _ = bh.Stop() }()

	// Attach Router dispatcher
	bh.Handle(func(thCtx *th.Context, update telego.Update) error {
		abstractUpdate := telegoadapter.MapUpdate(&update)
		if abstractUpdate != nil {
			_, err := b.Router.Dispatch(thCtx, abstractUpdate)
			if err != nil {
				b.log.Error("telegram router dispatch error", "error", err, "update_id", update.UpdateID)
				return err
			}
		}
		return nil
	})

	for _, reg := range registrars {
		if reg != nil {
			reg(b, bh)
		}
	}

	b.log.Info("telegram long-polling started")
	if err := bh.Start(); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
