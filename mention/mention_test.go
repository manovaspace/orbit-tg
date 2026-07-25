package mention

import (
	"testing"

	"github.com/mymmrac/telego"
)

func TestOnLineMentionEntity(t *testing.T) {
	text := "=1 @alice"
	entities := []telego.MessageEntity{
		{Type: "mention", Offset: 3, Length: 6},
	}
	u, ok := OnLine(entities, 0, UTF16Len(text), text)
	if !ok || u.Username != "alice" {
		t.Fatalf("got %+v ok=%v", u, ok)
	}
}

func TestMatchLineTypedUsername(t *testing.T) {
	full := "=1 @alirezaopmc"
	lines := []string{full}
	u, idx, ok := MatchLine(full, lines, nil, nil, 1, func(s string) string {
		if len(s) > 0 && s[0] == '=' {
			return s[1:]
		}
		return s
	})
	if !ok || idx != 0 || u.Username != "alirezaopmc" {
		t.Fatalf("got %+v idx=%d ok=%v", u, idx, ok)
	}
}

func TestUTF16SubstringEmoji(t *testing.T) {
	text := "👋 @bob"
	got := UTF16Substring(text, 3, 4)
	if got != "@bob" {
		t.Fatalf("got %q", got)
	}
}
