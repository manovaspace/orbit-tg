package chat_test

import (
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/chat"
)

func TestChatHelpers(t *testing.T) {
	cGroup := tg.Chat{Type: "group"}
	cPrivate := tg.Chat{Type: "private"}
	if !chat.IsGroup(cGroup) {
		t.Errorf("expected IsGroup(group) to be true")
	}
	if chat.IsGroup(cPrivate) {
		t.Errorf("expected IsGroup(private) to be false")
	}

	msg := tg.Message{ID: 10, ThreadID: 42}
	if thread := chat.MessageThread(msg); thread != 42 {
		t.Errorf("expected thread 42, got %d", thread)
	}

	cb := tg.CallbackQuery{
		ID:      "cb_1",
		Message: &msg,
	}
	if thread := chat.CallbackThread(cb); thread != 42 {
		t.Errorf("expected callback thread 42, got %d", thread)
	}
	if id := chat.CallbackMessageID(cb); id != 10 {
		t.Errorf("expected callback message id 10, got %d", id)
	}

	cbNil := tg.CallbackQuery{ID: "cb_2"}
	if thread := chat.CallbackThread(cbNil); thread != 0 {
		t.Errorf("expected 0 for nil message, got %d", thread)
	}
	if id := chat.CallbackMessageID(cbNil); id != 0 {
		t.Errorf("expected 0 for nil message, got %d", id)
	}

	user := tg.User{FirstName: "Ali", LastName: "Reza"}
	if name := chat.DisplayName(user); name != "Ali Reza" {
		t.Errorf("expected 'Ali Reza', got %q", name)
	}
}
