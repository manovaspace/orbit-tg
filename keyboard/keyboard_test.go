package keyboard_test

import (
	"encoding/json"
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/keyboard"
)

func TestMarkupContracts(t *testing.T) {
	var _ tg.Markup = (*keyboard.InlineKeyboard)(nil)
	var _ tg.Markup = (*keyboard.ReplyKeyboard)(nil)
	var _ tg.Markup = (*keyboard.RemoveKeyboard)(nil)
	var _ tg.Markup = (*keyboard.ForceReply)(nil)

	ik := &keyboard.InlineKeyboard{}
	if ik.MarkupType() != "inline" {
		t.Fatalf("unexpected markup type %q", ik.MarkupType())
	}

	rk := &keyboard.ReplyKeyboard{}
	if rk.MarkupType() != "reply" {
		t.Fatalf("unexpected markup type %q", rk.MarkupType())
	}

	remk := &keyboard.RemoveKeyboard{}
	if remk.MarkupType() != "remove" {
		t.Fatalf("unexpected markup type %q", remk.MarkupType())
	}

	fr := &keyboard.ForceReply{}
	if fr.MarkupType() != "force_reply" {
		t.Fatalf("unexpected markup type %q", fr.MarkupType())
	}
}

func TestInlineBuilder(t *testing.T) {
	kb := keyboard.NewInline().
		Row(
			keyboard.Callback("Confirm", "action:confirm"),
			keyboard.URL("Site", "https://manova.space"),
		).
		Row(
			keyboard.WebApp("App", "https://app.manova.space"),
			keyboard.CopyText("Copy Token", "secret-token-123"),
		).
		Build()

	if len(kb.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(kb.Rows))
	}
	if len(kb.Rows[0]) != 2 || len(kb.Rows[1]) != 2 {
		t.Fatalf("expected 2 buttons per row")
	}

	b1 := kb.Rows[0][0]
	if b1.Text != "Confirm" || b1.CallbackData != "action:confirm" {
		t.Errorf("button 1 mismatch: %+v", b1)
	}

	b2 := kb.Rows[0][1]
	if b2.Text != "Site" || b2.URL != "https://manova.space" {
		t.Errorf("button 2 mismatch: %+v", b2)
	}

	b3 := kb.Rows[1][0]
	if b3.Text != "App" || b3.WebAppURL != "https://app.manova.space" {
		t.Errorf("button 3 mismatch: %+v", b3)
	}

	b4 := kb.Rows[1][1]
	if b4.Text != "Copy Token" || b4.CopyText != "secret-token-123" {
		t.Errorf("button 4 mismatch: %+v", b4)
	}
}

func TestInlineBuilder_Chunk(t *testing.T) {
	buttons := []keyboard.InlineButton{
		keyboard.Callback("1", "1"),
		keyboard.Callback("2", "2"),
		keyboard.Callback("3", "3"),
		keyboard.Callback("4", "4"),
		keyboard.Callback("5", "5"),
	}

	kb := keyboard.NewInline().Chunk(2, buttons...).Build()
	if len(kb.Rows) != 3 {
		t.Fatalf("expected 3 rows for 5 items chunked by 2, got %d", len(kb.Rows))
	}
	if len(kb.Rows[0]) != 2 || len(kb.Rows[1]) != 2 || len(kb.Rows[2]) != 1 {
		t.Errorf("unexpected row lengths: %d, %d, %d", len(kb.Rows[0]), len(kb.Rows[1]), len(kb.Rows[2]))
	}
}

func TestGrid(t *testing.T) {
	fruits := []string{"Apple", "Banana", "Cherry", "Date"}
	kb := keyboard.Grid(fruits, 2, func(f string) keyboard.InlineButton {
		return keyboard.Callback(f, "pick:"+f)
	})

	if len(kb.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(kb.Rows))
	}
	if kb.Rows[0][0].Text != "Apple" || kb.Rows[1][1].Text != "Date" {
		t.Errorf("unexpected button grid content: %+v", kb.Rows)
	}
}

func TestReplyBuilder(t *testing.T) {
	kb := keyboard.NewReply().
		Row(
			keyboard.TextButton("Start"),
			keyboard.ContactButton("Send Phone"),
			keyboard.LocationButton("Send Location"),
		).
		Resize(true).
		OneTime(true).
		Placeholder("Type command...").
		Selective(true).
		Persistent(true).
		Build()

	if len(kb.Rows) != 1 || len(kb.Rows[0]) != 3 {
		t.Fatalf("unexpected row structure: %+v", kb.Rows)
	}
	if !kb.Resize || !kb.OneTime || !kb.Selective || !kb.IsPersistent {
		t.Errorf("flags not set properly: %+v", kb)
	}
	if kb.Placeholder != "Type command..." {
		t.Errorf("unexpected placeholder: %s", kb.Placeholder)
	}
	if !kb.Rows[0][1].RequestContact {
		t.Errorf("expected contact button")
	}
	if !kb.Rows[0][2].RequestLocation {
		t.Errorf("expected location button")
	}
}

func TestCommonDialogs(t *testing.T) {
	t.Run("DismissOK default", func(t *testing.T) {
		kb := keyboard.DismissOK("")
		if len(kb.Rows) != 1 || len(kb.Rows[0]) != 1 {
			t.Fatalf("unexpected rows: %+v", kb)
		}
		btn := kb.Rows[0][0]
		if btn.Text != "Ok ✅" || btn.CallbackData != keyboard.CallbackDismiss {
			t.Errorf("unexpected dismiss button: %+v", btn)
		}
		if !keyboard.IsDismiss(btn.CallbackData) {
			t.Errorf("expected IsDismiss to be true")
		}
	})

	t.Run("DismissOK custom", func(t *testing.T) {
		kb := keyboard.DismissOK("Close", "custom:close")
		btn := kb.Rows[0][0]
		if btn.Text != "Close" || btn.CallbackData != "custom:close" {
			t.Errorf("unexpected dismiss button: %+v", btn)
		}
	})

	t.Run("ConfirmRow", func(t *testing.T) {
		kb := keyboard.ConfirmRow("Yes", "yes_payload", "No", "no_payload")
		if len(kb.Rows) != 1 || len(kb.Rows[0]) != 2 {
			t.Fatalf("unexpected rows: %+v", kb)
		}
		if kb.Rows[0][0].CallbackData != "yes_payload" || kb.Rows[0][1].CallbackData != "no_payload" {
			t.Errorf("unexpected confirm row: %+v", kb.Rows[0])
		}
	})

	t.Run("PaginationRow", func(t *testing.T) {
		// Single page returns nil
		if row := keyboard.PaginationRow(1, 1, "items"); row != nil {
			t.Errorf("expected nil for 1 page, got %v", row)
		}

		// First page
		row := keyboard.PaginationRow(1, 5, "items")
		if len(row) != 3 {
			t.Fatalf("expected 3 buttons, got %d", len(row))
		}
		if row[0].CallbackData != "noop" || row[2].CallbackData != "items:page:2" {
			t.Errorf("unexpected first page row: %+v", row)
		}

		// Middle page
		row = keyboard.PaginationRow(3, 5, "items")
		if row[0].CallbackData != "items:page:2" || row[2].CallbackData != "items:page:4" {
			t.Errorf("unexpected middle page row: %+v", row)
		}

		// Last page
		row = keyboard.PaginationRow(5, 5, "items")
		if row[0].CallbackData != "items:page:4" || row[2].CallbackData != "noop" {
			t.Errorf("unexpected last page row: %+v", row)
		}
	})
}

func TestJSONSerialization(t *testing.T) {
	kb := keyboard.NewInline().
		Row(keyboard.Callback("Hello", "greet")).
		Build()

	data, err := json.Marshal(kb)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded keyboard.InlineKeyboard
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if len(decoded.Rows) != 1 || decoded.Rows[0][0].Text != "Hello" {
		t.Errorf("decoded mismatch: %+v", decoded)
	}
}
