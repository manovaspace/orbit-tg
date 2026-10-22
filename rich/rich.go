package rich

import (
	"fmt"
	"strings"

	"github.com/manovaspace/orbit-tg/message"
)

// Section is a titled bullet list used in structured notice messages.
type Section struct {
	Title string
	Items []string
}

// Notice builds a hierarchical notice: main heading + optional subsections with bullet lists.
func Notice(mainTitle string, sections ...Section) string {
	var b strings.Builder
	b.WriteString("<b>" + message.HTMLEscape(mainTitle) + "</b>")

	for _, sec := range sections {
		if len(sec.Items) == 0 {
			continue
		}
		b.WriteString("\n\n")
		if sec.Title != "" {
			b.WriteString("<b>" + message.HTMLEscape(sec.Title) + "</b>\n")
		}
		b.WriteString(BulletList(sec.Items))
	}

	return b.String()
}

// BulletList builds a bulleted list string (• item).
func BulletList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var lines []string
	for _, item := range items {
		lines = append(lines, "• "+message.HTMLEscape(item))
	}
	return strings.Join(lines, "\n")
}

// NumberedList builds an ordered list string (1. item).
func NumberedList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var lines []string
	for i, item := range items {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, message.HTMLEscape(item)))
	}
	return strings.Join(lines, "\n")
}

// CodeBlock formats code with an optional language syntax class.
func CodeBlock(lang, code string) string {
	escaped := message.HTMLEscape(code)
	if lang != "" {
		return fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, message.HTMLEscape(lang), escaped)
	}
	return fmt.Sprintf("<pre>%s</pre>", escaped)
}

// Blockquote formats text as a quote, optionally expandable for long messages.
func Blockquote(text string, expandable bool) string {
	escaped := message.HTMLEscape(text)
	if expandable {
		return fmt.Sprintf("<blockquote expandable>%s</blockquote>", escaped)
	}
	return fmt.Sprintf("<blockquote>%s</blockquote>", escaped)
}

// Bold returns HTML bold text.
func Bold(text string) string {
	return "<b>" + message.HTMLEscape(text) + "</b>"
}

// Italic returns HTML italic text.
func Italic(text string) string {
	return "<i>" + message.HTMLEscape(text) + "</i>"
}

// Code returns inline code.
func Code(text string) string {
	return "<code>" + message.HTMLEscape(text) + "</code>"
}

// Link returns an HTML link.
func Link(text, url string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, message.HTMLEscape(url), message.HTMLEscape(text))
}

// Spoiler returns spoiler text.
func Spoiler(text string) string {
	return "<tg-spoiler>" + message.HTMLEscape(text) + "</tg-spoiler>"
}

// Builder accumulates HTML blocks into a structured message.
type Builder struct {
	parts []string
}

// NewBuilder initializes an HTML message Builder.
func NewBuilder() *Builder {
	return &Builder{parts: make([]string, 0)}
}

// Heading adds a bold heading line.
func (b *Builder) Heading(title string) *Builder {
	b.parts = append(b.parts, Bold(title))
	return b
}

// Paragraph adds a plain text paragraph.
func (b *Builder) Paragraph(text string) *Builder {
	b.parts = append(b.parts, message.HTMLEscape(text))
	return b
}

// ParagraphBold adds a bold text paragraph.
func (b *Builder) ParagraphBold(text string) *Builder {
	b.parts = append(b.parts, Bold(text))
	return b
}

// ParagraphItalic adds an italic text paragraph.
func (b *Builder) ParagraphItalic(text string) *Builder {
	b.parts = append(b.parts, Italic(text))
	return b
}

// BulletList adds a bulleted list.
func (b *Builder) BulletList(items ...string) *Builder {
	if len(items) > 0 {
		b.parts = append(b.parts, BulletList(items))
	}
	return b
}

// NumberedList adds a numbered list.
func (b *Builder) NumberedList(items ...string) *Builder {
	if len(items) > 0 {
		b.parts = append(b.parts, NumberedList(items))
	}
	return b
}

// CodeBlock adds a preformatted code block.
func (b *Builder) CodeBlock(lang, code string) *Builder {
	b.parts = append(b.parts, CodeBlock(lang, code))
	return b
}

// Blockquote adds a blockquote.
func (b *Builder) Blockquote(text string, expandable bool) *Builder {
	b.parts = append(b.parts, Blockquote(text, expandable))
	return b
}

// Divider adds a visual separator.
func (b *Builder) Divider() *Builder {
	b.parts = append(b.parts, "———")
	return b
}

// Build combines all blocks separated by double newlines into Telegram HTML.
func (b *Builder) Build() string {
	return strings.Join(b.parts, "\n\n")
}

// String implements fmt.Stringer by calling Build().
func (b *Builder) String() string {
	return b.Build()
}
