package bot

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"github.com/manovaspace/orbit-tg/message"
)

type Bot struct {
	API *telego.Bot
	Msg *message.Messenger
	log *slog.Logger
}

func New(token string, log *slog.Logger) (*Bot, error) {
	api, err := telego.NewBot(token)
	if err != nil {
		return nil, err
	}
	if log == nil {
		log = slog.Default()
	}
	return &Bot{API: api, Msg: message.New(api, log), log: log}, nil
}

// Run starts long-polling and invokes register to attach handlers.
func (b *Bot) Run(ctx context.Context, register func(bh *th.BotHandler)) error {
	updates, err := b.API.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}
	bh, err := th.NewBotHandler(b.API, updates)
	if err != nil {
		return err
	}
	defer bh.Stop()

	register(bh)
	b.log.Info("telegram long-polling started")
	bh.Start()
	<-ctx.Done()
	return nil
}
