package mock_test

import (
	"context"
	"errors"
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/mock"
)

func TestMockMessenger_Lifecycle(t *testing.T) {
	m := mock.NewMessenger()
	ctx := context.Background()

	// SendText
	msg, err := m.SendText(ctx, 123, "hello", tg.WithThreadID(10))
	if err != nil || msg.ID == 0 {
		t.Fatalf("SendText failed: %v", err)
	}
	if len(m.Sent) != 1 || m.Sent[0].Text != "hello" || m.Sent[0].Config.ThreadID != 10 {
		t.Errorf("Sent record mismatch: %+v", m.Sent)
	}

	// SendHTML
	htmlMsg, err := m.SendHTML(ctx, 123, "<b>bold</b>")
	if err != nil || !m.Sent[1].IsHTML {
		t.Fatalf("SendHTML failed: %v", err)
	}
	if htmlMsg.Text != "<b>bold</b>" {
		t.Errorf("expected HTML text '<b>bold</b>', got %q", htmlMsg.Text)
	}

	// EditText & EditHTML
	_, _ = m.EditText(ctx, 123, 1001, "edited plain")
	_, _ = m.EditHTML(ctx, 123, 1002, "<i>edited html</i>")
	if len(m.Edited) != 2 {
		t.Fatalf("expected 2 edited records, got %d", len(m.Edited))
	}

	// Delete
	_ = m.Delete(ctx, 123, 1001, 1002)
	if len(m.Deleted) != 2 {
		t.Fatalf("expected 2 deleted message IDs, got %v", m.Deleted)
	}

	// AnswerCallback
	_ = m.AnswerCallback(ctx, "query_123", tg.WithAlert(true), tg.WithCallbackText("Done"))
	if len(m.Callbacks) != 1 || !m.Callbacks[0].Config.ShowAlert || m.Callbacks[0].Config.Text != "Done" {
		t.Errorf("Callback record mismatch: %+v", m.Callbacks)
	}

	// Pin & Unpin
	_ = m.Pin(ctx, 123, 1001, true)
	_ = m.Unpin(ctx, 123, 1001)
	if len(m.Pinned) != 1 || !m.Pinned[0].Silent {
		t.Errorf("Pin mismatch: %+v", m.Pinned)
	}
	if len(m.Unpinned) != 1 {
		t.Errorf("Unpin mismatch: %+v", m.Unpinned)
	}

	// Close
	_ = m.Close()
	if !m.Closed {
		t.Errorf("expected Closed to be true")
	}
}

func TestMockMessenger_ErrorInjection(t *testing.T) {
	m := mock.NewMessenger()
	ctx := context.Background()

	simErr := errors.New("network failure")
	m.FailWith = simErr

	_, err := m.SendText(ctx, 123, "test")
	if !errors.Is(err, simErr) {
		t.Errorf("expected simErr, got %v", err)
	}

	err = m.Delete(ctx, 123, 1)
	if !errors.Is(err, simErr) {
		t.Errorf("expected simErr on Delete, got %v", err)
	}

	m.FailWith = nil
	m.FailCallback = simErr

	err = m.AnswerCallback(ctx, "q1")
	if !errors.Is(err, simErr) {
		t.Errorf("expected simErr on AnswerCallback, got %v", err)
	}
}
