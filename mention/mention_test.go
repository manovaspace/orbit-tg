package mention_test

import (
	"testing"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/mention"
)

func TestUTF16Len(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"", 0},
		// Emoji 😀 requires surrogate pair (2 code units)
		{"😀", 2},
		{"hello 😀 world", 14},
	}

	for _, tt := range tests {
		if got := mention.UTF16Len(tt.input); got != tt.want {
			t.Errorf("UTF16Len(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestLineUTF16Offset(t *testing.T) {
	text := "line 1 😀\nline 2\nline 3"
	// line 0: 0
	// line 1: len("line 1 ") + 2 (emoji) + 1 (\n) = 7 + 2 + 1 = 10
	if offset := mention.LineUTF16Offset(text, 0); offset != 0 {
		t.Errorf("LineUTF16Offset(0) = %d, want 0", offset)
	}
	if offset := mention.LineUTF16Offset(text, 1); offset != 10 {
		t.Errorf("LineUTF16Offset(1) = %d, want 10", offset)
	}
}

func TestOnLine(t *testing.T) {
	user := &tg.User{ID: 100, FirstName: "Alice"}
	entities := []mention.Entity{
		{
			Type:   "text_mention",
			Offset: 5,
			Length: 5,
			User:   user,
		},
		{
			Type:   "bold",
			Offset: 12,
			Length: 4,
		},
	}

	// Range covering offset 5
	foundUser, ok := mention.OnLine(entities, 0, 10)
	if !ok || foundUser.ID != 100 {
		t.Errorf("expected user 100, got %v (ok=%v)", foundUser, ok)
	}

	// Range not covering offset 5
	_, ok = mention.OnLine(entities, 11, 20)
	if ok {
		t.Errorf("expected false for range 11-20")
	}
}
