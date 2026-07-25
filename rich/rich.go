package rich

import (
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HTMLEscape(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	).Replace(s)
}

func MessageHTML(html string) telego.InputRichMessage {
	return tu.RichMessageHTML(html)
}

// Section is a titled bullet list used in notice messages.
type Section struct {
	Title string
	Items []string
}

// Notice builds a hierarchical rich message: main heading + optional subsections with bullet lists.
func Notice(mainTitle string, sections ...Section) telego.InputRichMessage {
	blocks := []telego.InputRichBlock{
		tu.RichBlockSectionHeading(tu.RichTextPlain(mainTitle), 2),
	}
	for _, sec := range sections {
		if len(sec.Items) == 0 {
			continue
		}
		if sec.Title != "" {
			blocks = append(blocks, tu.RichBlockSectionHeading(tu.RichTextPlain(sec.Title), 3))
		}
		blocks = append(blocks, BulletList(sec.Items))
	}
	return tu.RichMessage(blocks...)
}

func BulletList(items []string) telego.InputRichBlock {
	paras := make([]telego.RichText, len(items))
	for i, name := range items {
		paras[i] = tu.RichTextPlain(name)
	}
	return BulletListRich(paras...)
}

// BulletListRich builds a bullet list from rich-text paragraphs (emoji, bold, etc.).
func BulletListRich(items ...telego.RichText) telego.InputRichBlock {
	listItems := make([]telego.InputRichBlockListItem, len(items))
	for i, text := range items {
		listItems[i] = tu.RichBlockListItem(tu.RichBlockParagraph(text))
	}
	return tu.RichBlockList(listItems...)
}

// NumberedList builds an ordered list (1. 2. 3. …) without separators between items.
func NumberedList(items []string) telego.InputRichBlock {
	listItems := make([]telego.InputRichBlockListItem, len(items))
	for i, name := range items {
		listItems[i] = tu.RichBlockListItem(
			tu.RichBlockParagraph(tu.RichTextPlain(name)),
		).WithType(telego.OrderedListDecimal).WithValue(i + 1)
	}
	return tu.RichBlockList(listItems...)
}

// Builder accumulates rich blocks for custom layouts.
type Builder struct {
	blocks []telego.InputRichBlock
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Heading(title string, level int) *Builder {
	b.blocks = append(b.blocks, tu.RichBlockSectionHeading(tu.RichTextPlain(title), level))
	return b
}

func (b *Builder) ParagraphPlain(text string) *Builder {
	b.blocks = append(b.blocks, tu.RichBlockParagraph(tu.RichTextPlain(text)))
	return b
}

func (b *Builder) ParagraphItalic(text string) *Builder {
	b.blocks = append(b.blocks, tu.RichBlockParagraph(tu.RichTextItalic(tu.RichTextPlain(text))))
	return b
}

func (b *Builder) Divider() *Builder {
	b.blocks = append(b.blocks, tu.RichBlockDivider())
	return b
}

func (b *Builder) Bullets(items []string) *Builder {
	if len(items) > 0 {
		b.blocks = append(b.blocks, BulletList(items))
	}
	return b
}

func (b *Builder) BulletsRich(items ...telego.RichText) *Builder {
	if len(items) > 0 {
		b.blocks = append(b.blocks, BulletListRich(items...))
	}
	return b
}

func (b *Builder) Block(block telego.InputRichBlock) *Builder {
	b.blocks = append(b.blocks, block)
	return b
}

func (b *Builder) Numbered(items []string) *Builder {
	if len(items) > 0 {
		b.blocks = append(b.blocks, NumberedList(items))
	}
	return b
}

func (b *Builder) Build() telego.InputRichMessage {
	return tu.RichMessage(b.blocks...)
}
