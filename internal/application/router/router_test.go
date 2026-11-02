package router_test

import (
	"context"
	"errors"
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/internal/application/router"
)

func TestDispatchRunsFirstMatchingRoute(t *testing.T) {
	r := router.New()
	var calls []string
	r.Handle(func(*tg.Update) bool { return true }, func(context.Context, *tg.Update) error {
		calls = append(calls, "first")
		return nil
	})
	r.Handle(func(*tg.Update) bool { return true }, func(context.Context, *tg.Update) error {
		calls = append(calls, "second")
		return nil
	})
	matched, err := r.Dispatch(context.Background(), &tg.Update{ID: 1})
	if err != nil || !matched {
		t.Fatalf("matched=%v err=%v", matched, err)
	}
	if len(calls) != 1 || calls[0] != "first" {
		t.Fatalf("calls = %v, want only first", calls)
	}
}

func TestDispatchIgnoresNilUpdate(t *testing.T) {
	matched, err := router.New().Dispatch(context.Background(), nil)
	if err != nil || matched {
		t.Fatalf("matched=%v err=%v, want no match", matched, err)
	}
}

func TestCommandMatchesFirstLineAndMentions(t *testing.T) {
	r := router.New()
	called := false
	r.Command("start", func(context.Context, *tg.Update) error {
		called = true
		return nil
	})
	for _, text := range []string{"/start", "/start@mybot\narg", "  /START  "} {
		called = false
		matched, err := r.Dispatch(context.Background(), &tg.Update{Message: &tg.Message{Text: text}})
		if err != nil || !matched || !called {
			t.Fatalf("%q: matched=%v called=%v err=%v", text, matched, called, err)
		}
	}
	matched, _ := r.Dispatch(context.Background(), &tg.Update{Message: &tg.Message{Text: "/starter"}})
	if matched {
		t.Fatal("/starter must not match /start")
	}
}

func TestCallbackPrefix(t *testing.T) {
	r := router.New()
	var data string
	r.CallbackPrefix("order:", func(_ context.Context, u *tg.Update) error {
		data = u.CallbackQuery.Data
		return nil
	})
	matched, err := r.Dispatch(context.Background(), &tg.Update{
		CallbackQuery: &tg.CallbackQuery{Data: "order:confirm"},
	})
	if err != nil || !matched || data != "order:confirm" {
		t.Fatalf("matched=%v data=%q err=%v", matched, data, err)
	}
}

func TestHandlerErrorIsReturned(t *testing.T) {
	want := errors.New("boom")
	r := router.New()
	r.Handle(func(*tg.Update) bool { return true }, func(context.Context, *tg.Update) error { return want })
	matched, err := r.Dispatch(context.Background(), &tg.Update{})
	if !matched || !errors.Is(err, want) {
		t.Fatalf("matched=%v err=%v", matched, err)
	}
}
