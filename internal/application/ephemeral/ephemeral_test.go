package ephemeral_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/manovaspace/orbit-tg/internal/application/ephemeral"
)

func TestScheduleDeletesAfterTTL(t *testing.T) {
	var (
		mu     sync.Mutex
		gotIDs []int
	)
	m := ephemeral.NewEphemeralManager(context.Background(), func(_ context.Context, _ int64, id int) error {
		mu.Lock()
		defer mu.Unlock()
		gotIDs = append(gotIDs, id)
		return nil
	})
	m.Schedule(1, 42, 20*time.Millisecond)

	deadline := time.Now().Add(time.Second)
	for {
		mu.Lock()
		n := len(gotIDs)
		mu.Unlock()
		if n == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("scheduled message was not deleted")
		}
		time.Sleep(5 * time.Millisecond)
	}
	if gotIDs[0] != 42 {
		t.Fatalf("deleted id = %d, want 42", gotIDs[0])
	}
	m.Close()
}

func TestCloseCancelsPendingDeletion(t *testing.T) {
	called := false
	m := ephemeral.NewEphemeralManager(context.Background(), func(context.Context, int64, int) error {
		called = true
		return nil
	})
	m.Schedule(1, 42, time.Hour)
	m.Close()
	time.Sleep(10 * time.Millisecond)
	if called {
		t.Fatal("Close must cancel the pending deletion")
	}
}

func TestScheduleIgnoresInvalidInput(t *testing.T) {
	m := ephemeral.NewEphemeralManager(context.Background(), nil)
	m.Schedule(1, 0, time.Millisecond)
	m.Schedule(1, 2, time.Millisecond)
	m.Close()
}
