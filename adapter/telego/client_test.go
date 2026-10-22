package telego_test

import (
	"testing"
	"time"

	telegotype "github.com/mymmrac/telego"

	tg "github.com/manovaspace/orbit-tg"
	telegoadapter "github.com/manovaspace/orbit-tg/adapter/telego"
	"github.com/manovaspace/orbit-tg/keyboard"
)

func TestMapMarkup(t *testing.T) {
	t.Run("nil markup", func(t *testing.T) {
		if res := telegoadapter.MapMarkup(nil); res != nil {
			t.Fatalf("expected nil for nil markup, got %v", res)
		}
	})

	t.Run("inline keyboard", func(t *testing.T) {
		ik := keyboard.NewInline().
			Row(
				keyboard.Callback("Btn1", "cb1"),
				keyboard.URL("Link", "https://example.com"),
			).
			Build()

		mapped := telegoadapter.MapMarkup(ik)
		tgMarkup, ok := mapped.(*telegotype.InlineKeyboardMarkup)
		if !ok || len(tgMarkup.InlineKeyboard) != 1 || len(tgMarkup.InlineKeyboard[0]) != 2 {
			t.Fatalf("unexpected inline keyboard mapping: %T %+v", mapped, mapped)
		}
		if tgMarkup.InlineKeyboard[0][0].Text != "Btn1" || tgMarkup.InlineKeyboard[0][0].CallbackData != "cb1" {
			t.Errorf("button 1 mismatch: %+v", tgMarkup.InlineKeyboard[0][0])
		}
		if tgMarkup.InlineKeyboard[0][1].Text != "Link" || tgMarkup.InlineKeyboard[0][1].URL != "https://example.com" {
			t.Errorf("button 2 mismatch: %+v", tgMarkup.InlineKeyboard[0][1])
		}
	})

	t.Run("reply keyboard", func(t *testing.T) {
		rk := keyboard.NewReply().
			Row(keyboard.ContactButton("Send Contact")).
			Resize(true).
			Build()

		mapped := telegoadapter.MapMarkup(rk)
		tgMarkup, ok := mapped.(*telegotype.ReplyKeyboardMarkup)
		if !ok || !tgMarkup.ResizeKeyboard || len(tgMarkup.Keyboard) != 1 {
			t.Fatalf("unexpected reply keyboard mapping: %T %+v", mapped, mapped)
		}
		if !tgMarkup.Keyboard[0][0].RequestContact {
			t.Errorf("expected RequestContact to be true")
		}
	})

	t.Run("remove keyboard", func(t *testing.T) {
		remk := &keyboard.RemoveKeyboard{Selective: true}
		mapped := telegoadapter.MapMarkup(remk)
		tgMarkup, ok := mapped.(*telegotype.ReplyKeyboardRemove)
		if !ok || !tgMarkup.RemoveKeyboard || !tgMarkup.Selective {
			t.Fatalf("unexpected remove keyboard mapping: %+v", mapped)
		}
	})

	t.Run("force reply", func(t *testing.T) {
		fr := &keyboard.ForceReply{Placeholder: "Reply here", Selective: true}
		mapped := telegoadapter.MapMarkup(fr)
		tgMarkup, ok := mapped.(*telegotype.ForceReply)
		if !ok || !tgMarkup.ForceReply || tgMarkup.InputFieldPlaceholder != "Reply here" {
			t.Fatalf("unexpected force reply mapping: %+v", mapped)
		}
	})
}

func TestMapUserAndChat(t *testing.T) {
	tu := &telegotype.User{
		ID:        42,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		IsBot:     false,
	}
	u := telegoadapter.MapUser(tu)
	if u.ID != 42 || u.FirstName != "John" || u.Username != "johndoe" {
		t.Errorf("User mapping mismatch: %+v", u)
	}

	tc := telegotype.Chat{
		ID:    -10012345,
		Type:  "supergroup",
		Title: "Dev Team",
	}
	c := telegoadapter.MapChat(tc)
	if c.ID != -10012345 || !c.IsGroup() || c.Title != "Dev Team" {
		t.Errorf("Chat mapping mismatch: %+v", c)
	}
}

func TestMapMessageAndUpdate(t *testing.T) {
	now := time.Now().Unix()
	tm := &telegotype.Message{
		MessageID:       101,
		Date:            now,
		MessageThreadID: 42,
		Chat: telegotype.Chat{
			ID:   -10099,
			Type: "supergroup",
		},
		From: &telegotype.User{
			ID:        999,
			FirstName: "Bot",
			IsBot:     true,
		},
		Text: "Hello update",
	}

	m := telegoadapter.MapMessage(tm)
	if m.ID != 101 || m.ThreadID != 42 || m.Text != "Hello update" || m.Chat.ID != -10099 {
		t.Errorf("Message mapping mismatch: %+v", m)
	}

	tq := &telegotype.CallbackQuery{
		ID:      "cb_query_1",
		From:    *tm.From,
		Message: tm,
		Data:    "click:action",
	}

	cb := telegoadapter.MapCallbackQuery(tq)
	if cb.ID != "cb_query_1" || cb.Data != "click:action" || cb.Message == nil || cb.Message.ID != 101 {
		t.Errorf("CallbackQuery mapping mismatch: %+v", cb)
	}

	tu := &telegotype.Update{
		UpdateID:      505,
		Message:       tm,
		CallbackQuery: tq,
	}

	u := telegoadapter.MapUpdate(tu)
	if u.ID != 505 || u.Message == nil || u.CallbackQuery == nil {
		t.Errorf("Update mapping mismatch: %+v", u)
	}
}

func TestTelegoAdapterInterface(t *testing.T) {
	var _ tg.Messenger = (*telegoadapter.TelegoAdapter)(nil)
}
