// Package handler holds the pure command-matching rules. It imports no
// Telegram client or HTTP package.
package handler

import (
	"regexp"
	"strings"
)

// FirstLineCommand returns a matcher for /cmd or /cmd@bot on the first line
// only. telego's CommandEqual uses a whole-message regex that breaks on trailing
// newlines (e.g. "/checklist\n" from mobile clients).
func FirstLineCommand(name string) func(text string) bool {
	re := regexp.MustCompile(`(?i)^/` + regexp.QuoteMeta(name) + `(?:@\w+)?(?:\s|$)`)
	return func(text string) bool {
		if text == "" {
			return false
		}
		first, _, _ := strings.Cut(text, "\n")
		return re.MatchString(strings.TrimSpace(first))
	}
}
