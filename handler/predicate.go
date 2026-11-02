package handler

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	apphandler "github.com/manovaspace/orbit-tg/internal/application/handler"
)

// FirstLineCommand matches /cmd or /cmd@bot on the first line only.
// telego's CommandEqual uses a whole-message regex that breaks on trailing newlines
// (e.g. "/checklist\n" from mobile clients).
func FirstLineCommand(name string) th.Predicate {
	match := apphandler.FirstLineCommand(name)
	return func(_ context.Context, update telego.Update) bool {
		if update.Message == nil {
			return false
		}
		return match(update.Message.Text)
	}
}
