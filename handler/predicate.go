package handler

import (
	"context"
	"regexp"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// FirstLineCommand matches /cmd or /cmd@bot on the first line only.
// telego's CommandEqual uses a whole-message regex that breaks on trailing newlines
// (e.g. "/checklist\n" from mobile clients).
func FirstLineCommand(name string) th.Predicate {
	re := regexp.MustCompile(`(?i)^/` + regexp.QuoteMeta(name) + `(?:@\w+)?(?:\s|$)`)
	return func(_ context.Context, update telego.Update) bool {
		if update.Message == nil || update.Message.Text == "" {
			return false
		}
		first, _, _ := strings.Cut(update.Message.Text, "\n")
		return re.MatchString(strings.TrimSpace(first))
	}
}
