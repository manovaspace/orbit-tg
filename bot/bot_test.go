package bot_test

import (
	"context"
	"log/slog"
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	ogbot "github.com/manovaspace/orbit-tg/bot"
)

func TestNewRequiresToken(t *testing.T) {
	_, err := ogbot.New("", slog.Default())
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestRouter_Command(t *testing.T) {
	r := ogbot.NewRouter()

	var startCalled, helpCalled bool
	r.Command("start", func(ctx context.Context, u *tg.Update) error {
		startCalled = true
		return nil
	})
	r.Command("/help", func(ctx context.Context, u *tg.Update) error {
		helpCalled = true
		return nil
	})

	ctx := context.Background()

	// Test /start matching
	matched, err := r.Dispatch(ctx, &tg.Update{
		ID: 1,
		Message: &tg.Message{
			Text: "/start",
		},
	})
	if err != nil || !matched || !startCalled {
		t.Fatalf("expected /start to match, matched=%v, called=%v, err=%v", matched, startCalled, err)
	}

	// Test /start@mybot matching with mobile newline
	startCalled = false
	matched, err = r.Dispatch(ctx, &tg.Update{
		ID: 2,
		Message: &tg.Message{
			Text: "/start@mybot\narg1",
		},
	})
	if err != nil || !matched || !startCalled {
		t.Fatalf("expected /start@mybot to match, matched=%v, called=%v, err=%v", matched, startCalled, err)
	}

	// Test /help matching
	matched, err = r.Dispatch(ctx, &tg.Update{
		ID: 3,
		Message: &tg.Message{
			Text: "/help",
		},
	})
	if err != nil || !matched || !helpCalled {
		t.Fatalf("expected /help to match, matched=%v, called=%v, err=%v", matched, helpCalled, err)
	}

	// Test unmatched command
	matched, err = r.Dispatch(ctx, &tg.Update{
		ID: 4,
		Message: &tg.Message{
			Text: "/unknown",
		},
	})
	if err != nil || matched {
		t.Fatalf("expected no match, matched=%v", matched)
	}
}

func TestRouter_CallbackPrefix(t *testing.T) {
	r := ogbot.NewRouter()

	var orderAction string
	r.CallbackPrefix("order:", func(ctx context.Context, u *tg.Update) error {
		orderAction = u.CallbackQuery.Data
		return nil
	})

	ctx := context.Background()

	matched, err := r.Dispatch(ctx, &tg.Update{
		ID: 1,
		CallbackQuery: &tg.CallbackQuery{
			ID:   "cb_1",
			Data: "order:confirm:123",
		},
	})
	if err != nil || !matched || orderAction != "order:confirm:123" {
		t.Fatalf("expected callback prefix match, matched=%v, action=%q", matched, orderAction)
	}

	matched, err = r.Dispatch(ctx, &tg.Update{
		ID: 2,
		CallbackQuery: &tg.CallbackQuery{
			ID:   "cb_2",
			Data: "user:profile",
		},
	})
	if err != nil || matched {
		t.Fatalf("expected no match for user:profile, matched=%v", matched)
	}
}
