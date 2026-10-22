package mention

import (
	"strings"

	tg "github.com/manovaspace/orbit-tg"
)

// UTF16Len calculates the UTF-16 code unit length of a string (standard Telegram entity measurement).
func UTF16Len(s string) int {
	n := 0
	for _, r := range s {
		if r > 0xFFFF {
			n += 2
		} else {
			n++
		}
	}
	return n
}

// LineUTF16Offset computes the UTF-16 code unit offset of a specific line in text.
func LineUTF16Offset(text string, lineIndex int) int {
	pos := 0
	lines := strings.Split(text, "\n")
	for i := 0; i < lineIndex && i < len(lines); i++ {
		pos += UTF16Len(lines[i]) + 1 // +1 for the newline character
	}
	return pos
}

// Entity represents an abstract message entity (mention, text_mention, command, etc.).
type Entity struct {
	Type   string   `json:"type"`
	Offset int      `json:"offset"`
	Length int      `json:"length"`
	User   *tg.User `json:"user,omitempty"`
}

// OnLine checks if any text_mention entity is located within the given line offset range.
func OnLine(entities []Entity, lineStart, lineEnd int) (*tg.User, bool) {
	for _, e := range entities {
		if e.Type != "text_mention" || e.User == nil {
			continue
		}
		if e.Offset >= lineStart && e.Offset < lineEnd {
			return e.User, true
		}
	}
	return nil, false
}
