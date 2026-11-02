package message

import "github.com/manovaspace/orbit-tg/internal/application/html"

// HTMLEscape escapes raw HTML special characters (&, <, >, ", ').
func HTMLEscape(text string) string {
	return html.HTMLEscape(text)
}

// SanitizeHTML cleans and normalizes HTML for Telegram Bot API:
// - Strips unsupported HTML tags while retaining inner text.
// - Preserves allowed Telegram tags and attributes.
// - Escapes bare angle brackets and unescaped ampersands.
// - Automatically balances unclosed tags in reverse order.
func SanitizeHTML(text string) string {
	return html.SanitizeHTML(text)
}

// ChunkHTML splits HTML into chunks not exceeding maxLength characters (default 4000).
// Each chunk is balanced with open/close tags so that it is standalone valid Telegram HTML.
func ChunkHTML(htmlStr string, maxLength int) []string {
	return html.ChunkHTML(htmlStr, maxLength)
}
