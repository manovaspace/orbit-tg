package bot_test

import (
	"testing"

	"log/slog"

	ogbot "github.com/manovaspace/orbit-tg/bot"
)

func TestNewRequiresToken(t *testing.T) {
	_, err := ogbot.New("", slog.Default())
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
