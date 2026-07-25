package mention

import (
	"strings"
	"unicode/utf16"

	"github.com/mymmrac/telego"
)

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

func UTF16Substring(text string, offset, length int) string {
	if offset < 0 || length <= 0 {
		return ""
	}
	u16 := utf16.Encode([]rune(text))
	end := offset + length
	if offset >= len(u16) {
		return ""
	}
	if end > len(u16) {
		end = len(u16)
	}
	return string(utf16.Decode(u16[offset:end]))
}

func LineUTF16Offset(text string, lineIndex int) int {
	pos := 0
	lines := strings.Split(text, "\n")
	for i := 0; i < lineIndex && i < len(lines); i++ {
		pos += UTF16Len(lines[i]) + 1
	}
	return pos
}

func OnLine(entities []telego.MessageEntity, lineStart, lineEnd int, fullText string) (*telego.User, bool) {
	for _, e := range entities {
		if e.Offset < lineStart || e.Offset >= lineEnd {
			continue
		}
		switch e.Type {
		case "text_mention":
			if e.User != nil {
				return e.User, true
			}
		case "mention":
			mention := UTF16Substring(fullText, e.Offset, e.Length)
			username := strings.TrimPrefix(mention, "@")
			if username != "" {
				return &telego.User{Username: username, FirstName: username}, true
			}
		}
	}
	return nil, false
}

// MatchLine finds the first unconsumed message line whose stripped body equals index,
// then returns an assignee from a text_mention, @mention entity, or typed @username.
func MatchLine(fullText string, msgLines []string, consumed []bool, entities []telego.MessageEntity, index int, strip func(string) string) (*telego.User, int, bool) {
	for lineIdx, line := range msgLines {
		if lineIdx < len(consumed) && consumed[lineIdx] {
			continue
		}
		body := strip(line)
		got, ok := parseLeadingIndex(body)
		if !ok || got != index {
			continue
		}
		start := LineUTF16Offset(fullText, lineIdx)
		end := start + UTF16Len(line)
		if u, ok := OnLine(entities, start, end, fullText); ok {
			return u, lineIdx, true
		}
		if u, ok := parseTypedUsername(body); ok {
			return u, lineIdx, true
		}
	}
	return nil, -1, false
}

func parseLeadingIndex(line string) (int, bool) {
	line = strings.TrimSpace(line)
	i := 0
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i == 0 {
		return 0, false
	}
	n := 0
	for j := 0; j < i; j++ {
		n = n*10 + int(line[j]-'0')
	}
	if n < 1 {
		return 0, false
	}
	return n, true
}

func parseTypedUsername(body string) (*telego.User, bool) {
	body = strings.TrimSpace(body)
	i := 0
	for i < len(body) && body[i] >= '0' && body[i] <= '9' {
		i++
	}
	rest := strings.TrimSpace(body[i:])
	if !strings.HasPrefix(rest, "@") {
		return nil, false
	}
	rest = rest[1:]
	end := 0
	for end < len(rest) && isUsernameChar(rest[end]) {
		end++
	}
	if end == 0 {
		return nil, false
	}
	username := rest[:end]
	return &telego.User{Username: username, FirstName: username}, true
}

func isUsernameChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}
