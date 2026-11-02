package scrub_test

import (
	"context"
	"errors"
	"testing"

	"github.com/manovaspace/orbit-tg/internal/application/scrub"
)

func TestScrubStopsAfterSuccess(t *testing.T) {
	attempts := 0
	err := scrub.Scrub(context.Background(), 1, 2, func(context.Context, int64, int) error {
		attempts++
		if attempts < 3 {
			return errors.New("still there")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestScrubReturnsLastErrorAfterBoundedRetries(t *testing.T) {
	want := errors.New("always fails")
	attempts := 0
	err := scrub.Scrub(context.Background(), 1, 2, func(context.Context, int64, int) error {
		attempts++
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("err = %v, want %v", err, want)
	}
	if attempts != 4 {
		t.Fatalf("attempts = %d, want 4", attempts)
	}
}

func TestScrubIgnoresMissingMessage(t *testing.T) {
	called := false
	err := scrub.Scrub(context.Background(), 1, 0, func(context.Context, int64, int) error {
		called = true
		return nil
	})
	if err != nil || called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}
