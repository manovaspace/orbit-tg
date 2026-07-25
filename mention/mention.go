package mention

import (
	"strings"

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

func LineUTF16Offset(text string, lineIndex int) int {
	pos := 0
	lines := strings.Split(text, "\n")
	for i := 0; i < lineIndex && i < len(lines); i++ {
		pos += UTF16Len(lines[i]) + 1
	}
	return pos
}

func OnLine(entities []telego.MessageEntity, lineStart, lineEnd int) (*telego.User, bool) {
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

// MatchLine finds the first unconsumed message line whose stripped body equals index,
// then returns a text_mention user on that line.
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
		if u, ok := OnLine(entities, start, end); ok {
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
