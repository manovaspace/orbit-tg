package message

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	allowedTags = map[string]bool{
		"b":          true,
		"strong":     true,
		"i":          true,
		"em":         true,
		"u":          true,
		"ins":        true,
		"s":          true,
		"strike":     true,
		"del":        true,
		"code":       true,
		"pre":        true,
		"blockquote": true,
		"tg-spoiler":  true,
		"a":          true,
		"tg-emoji":   true,
	}

	allowedTagAttrs = map[string]map[string]bool{
		"a":          {"href": true},
		"code":       {"class": true},
		"pre":        {"class": true},
		"blockquote": {"expandable": true},
		"tg-emoji":   {"emoji-id": true},
	}

	tokenRegex       = regexp.MustCompile(`<!--[\s\S]*?-->|</([a-zA-Z0-9_-]+)\s*>|<([a-zA-Z0-9_-]+)((?:\s+[^>]*)?)/?>`)
	attrRegex        = regexp.MustCompile(`([a-zA-Z0-9_-]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'=<>` + "`" + `]+)))?`)
	validEntityRegex = regexp.MustCompile(`^&([a-zA-Z0-9]+|#[0-9]{1,6}|#[xX][0-9a-fA-F]{1,6});`)
	tagRe            = regexp.MustCompile(`<(/)?([a-zA-Z0-9_-]+)(?:\s+[^>]*)?>`)
)

// HTMLEscape escapes raw HTML special characters (&, <, >, ", ').
func HTMLEscape(text string) string {
	if text == "" {
		return ""
	}
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return r.Replace(text)
}

// escapeTextNode escapes bare &, <, > in text nodes while preserving existing valid entities.
func escapeTextNode(text string) string {
	if text == "" {
		return ""
	}

	var b strings.Builder
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '&':
			if entity := validEntityRegex.FindString(text[i:]); entity != "" {
				b.WriteString(entity)
				i += len(entity) - 1
			} else {
				b.WriteString("&amp;")
			}
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String()
}

// SanitizeHTML cleans and normalizes HTML for Telegram Bot API:
// - Strips unsupported HTML tags while retaining inner text.
// - Preserves allowed Telegram tags and attributes.
// - Escapes bare angle brackets and unescaped ampersands.
// - Automatically balances unclosed tags in reverse order.
func SanitizeHTML(text string) string {
	if text == "" {
		return ""
	}

	var output strings.Builder
	var tagStack []string

	matches := tokenRegex.FindAllStringSubmatchIndex(text, -1)
	lastIndex := 0

	for _, matchIdx := range matches {
		matchStart, matchEnd := matchIdx[0], matchIdx[1]
		if matchStart > lastIndex {
			output.WriteString(escapeTextNode(text[lastIndex:matchStart]))
		}
		lastIndex = matchEnd

		fullMatch := text[matchStart:matchEnd]

		// Strip comments <!-- ... -->
		if strings.HasPrefix(fullMatch, "<!--") {
			continue
		}

		// Closing tag: </tag>
		if matchIdx[2] != -1 && matchIdx[3] != -1 {
			closingTag := strings.ToLower(text[matchIdx[2]:matchIdx[3]])
			if allowedTags[closingTag] && contains(tagStack, closingTag) {
				for len(tagStack) > 0 {
					top := tagStack[len(tagStack)-1]
					tagStack = tagStack[:len(tagStack)-1]
					output.WriteString("</" + top + ">")
					if top == closingTag {
						break
					}
				}
			}
			continue
		}

		// Opening tag: <tag attrs>
		if matchIdx[4] != -1 && matchIdx[5] != -1 {
			openTag := strings.ToLower(text[matchIdx[4]:matchIdx[5]])
			if allowedTags[openTag] {
				rawAttrs := ""
				if matchIdx[6] != -1 && matchIdx[7] != -1 {
					rawAttrs = text[matchIdx[6]:matchIdx[7]]
				}

				allowedAttrs := allowedTagAttrs[openTag]
				var attrParts []string
				hasHref := false

				attrMatches := attrRegex.FindAllStringSubmatch(rawAttrs, -1)
				for _, am := range attrMatches {
					attrName := strings.ToLower(am[1])
					if allowedAttrs != nil && allowedAttrs[attrName] {
						hasVal := false
						attrVal := ""
						if am[2] != "" {
							attrVal = am[2]
							hasVal = true
						} else if am[3] != "" {
							attrVal = am[3]
							hasVal = true
						} else if len(am) > 4 && am[4] != "" {
							attrVal = am[4]
							hasVal = true
						}

						if attrName == "href" && strings.TrimSpace(attrVal) != "" {
							hasHref = true
						}

						if !hasVal {
							attrParts = append(attrParts, fmt.Sprintf(" %s", attrName))
						} else {
							attrParts = append(attrParts, fmt.Sprintf(` %s="%s"`, attrName, HTMLEscape(attrVal)))
						}
					}
				}

				// <a> tag must have non-empty href
				if openTag == "a" && !hasHref {
					continue
				}

				output.WriteString("<" + openTag + strings.Join(attrParts, "") + ">")
				tagStack = append(tagStack, openTag)
			}
		}
	}

	if lastIndex < len(text) {
		output.WriteString(escapeTextNode(text[lastIndex:]))
	}

	// Balance any unclosed tags
	for len(tagStack) > 0 {
		top := tagStack[len(tagStack)-1]
		tagStack = tagStack[:len(tagStack)-1]
		output.WriteString("</" + top + ">")
	}

	return output.String()
}

func contains(stack []string, s string) bool {
	for _, item := range stack {
		if item == s {
			return true
		}
	}
	return false
}

type openTagInfo struct {
	name string
	full string
}

func getOpenTags(initialTags []openTagInfo, text string) []openTagInfo {
	stack := make([]openTagInfo, len(initialTags))
	copy(stack, initialTags)

	matches := tagRe.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		isClosing := m[1] == "/"
		tagName := strings.ToLower(m[2])
		if !allowedTags[tagName] {
			continue
		}

		if isClosing {
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i].name == tagName {
					stack = append(stack[:i], stack[i+1:]...)
					break
				}
			}
		} else {
			fullTag := m[0]
			if tagName == "a" && !strings.Contains(strings.ToLower(fullTag), "href") {
				continue
			}
			stack = append(stack, openTagInfo{name: tagName, full: fullTag})
		}
	}
	return stack
}

func isMidTag(s string, pos int) int {
	if pos <= 0 || pos > len(s) {
		return -1
	}
	openBracket := strings.LastIndex(s[:pos], "<")
	if openBracket != -1 {
		closeBracket := strings.Index(s[openBracket:], ">")
		if closeBracket == -1 || openBracket+closeBracket >= pos {
			return openBracket
		}
	}
	return -1
}

// ChunkHTML splits HTML into chunks not exceeding maxLength characters (default 4000).
// Each chunk is balanced with open/close tags so that it is standalone valid Telegram HTML.
func ChunkHTML(htmlStr string, maxLength int) []string {
	if htmlStr == "" {
		return nil
	}
	if maxLength <= 0 {
		maxLength = 4000
	}

	clean := SanitizeHTML(htmlStr)
	if len(clean) <= maxLength {
		return []string{clean}
	}

	var chunks []string
	remaining := clean
	var openTags []openTagInfo

	for len(remaining) > 0 {
		var prefixBuilder strings.Builder
		for _, t := range openTags {
			prefixBuilder.WriteString(t.full)
		}
		prefix := prefixBuilder.String()

		tagsIfAll := getOpenTags(openTags, remaining)
		var suffixBuilder strings.Builder
		for i := len(tagsIfAll) - 1; i >= 0; i-- {
			suffixBuilder.WriteString("</" + tagsIfAll[i].name + ">")
		}
		suffixIfAll := suffixBuilder.String()

		if len(prefix)+len(remaining)+len(suffixIfAll) <= maxLength {
			chunks = append(chunks, prefix+remaining+suffixIfAll)
			break
		}

		overheadEstimate := len(prefix) + len(openTags)*10 + 20
		available := maxLength - overheadEstimate
		if available <= 0 {
			available = maxLength / 2
		}

		searchWindow := remaining
		if len(searchWindow) > available {
			searchWindow = remaining[:available]
		}

		cutPos := -1

		// Priority 1: Paragraph boundary \n\n
		idx := strings.LastIndex(searchWindow, "\n\n")
		if idx > available/4 {
			cutPos = idx + 2
		}

		// Priority 2: Line break \n
		if cutPos == -1 {
			idx = strings.LastIndex(searchWindow, "\n")
			if idx > available/4 {
				cutPos = idx + 1
			}
		}

		// Priority 3: Word boundary ' '
		if cutPos == -1 {
			idx = strings.LastIndex(searchWindow, " ")
			if idx > available/4 {
				cutPos = idx + 1
			}
		}

		// Priority 4: Hard cut
		if cutPos == -1 {
			cutPos = len(searchWindow)
		}

		// Avoid splitting inside a tag
		tagStart := isMidTag(remaining, cutPos)
		if tagStart != -1 {
			if tagStart > 0 {
				cutPos = tagStart
			} else {
				// Tag starts at 0 and extends past cutPos; if it fits in maxLength, include full opening tag
				closeBracket := strings.Index(remaining, ">")
				if closeBracket != -1 && closeBracket+1 <= maxLength {
					cutPos = closeBracket + 1
				}
			}
		}

		chunkBody := remaining[:cutPos]
		currentOpenTags := getOpenTags(openTags, chunkBody)

		var suffixClose strings.Builder
		for i := len(currentOpenTags) - 1; i >= 0; i-- {
			suffixClose.WriteString("</" + currentOpenTags[i].name + ">")
		}

		chunkText := prefix + chunkBody + suffixClose.String()
		chunks = append(chunks, chunkText)

		remaining = remaining[cutPos:]
		openTags = currentOpenTags
	}

	return chunks
}
