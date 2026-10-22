package keyboard

// InlineBuilder fluently constructs InlineKeyboard structures.
type InlineBuilder struct {
	rows [][]InlineButton
}

// NewInline initializes an InlineBuilder.
func NewInline() *InlineBuilder {
	return &InlineBuilder{
		rows: make([][]InlineButton, 0),
	}
}

// Row appends a complete row of buttons.
func (b *InlineBuilder) Row(buttons ...InlineButton) *InlineBuilder {
	if len(buttons) > 0 {
		row := make([]InlineButton, len(buttons))
		copy(row, buttons)
		b.rows = append(b.rows, row)
	}
	return b
}

// Chunk distributes buttons into rows of max columns.
func (b *InlineBuilder) Chunk(cols int, buttons ...InlineButton) *InlineBuilder {
	if cols <= 0 {
		cols = 1
	}
	for i := 0; i < len(buttons); i += cols {
		end := i + cols
		if end > len(buttons) {
			end = len(buttons)
		}
		b.Row(buttons[i:end]...)
	}
	return b
}

// Build finalizes and returns the InlineKeyboard.
func (b *InlineBuilder) Build() *InlineKeyboard {
	return &InlineKeyboard{
		Rows: b.rows,
	}
}

// Callback creates an InlineButton with callback data.
func Callback(text, data string) InlineButton {
	return InlineButton{
		Text:         text,
		CallbackData: data,
	}
}

// URL creates an InlineButton with an external link URL.
func URL(text, url string) InlineButton {
	return InlineButton{
		Text: text,
		URL:  url,
	}
}

// WebApp creates an InlineButton that opens a Telegram Web App.
func WebApp(text, webAppURL string) InlineButton {
	return InlineButton{
		Text:      text,
		WebAppURL: webAppURL,
	}
}

// CopyText creates an InlineButton that copies text to clipboard on click.
func CopyText(text, copyText string) InlineButton {
	return InlineButton{
		Text:     text,
		CopyText: copyText,
	}
}

// Grid distributes generic items across rows using a column count and mapper.
func Grid[T any](items []T, cols int, mapper func(item T) InlineButton) *InlineKeyboard {
	b := NewInline()
	if cols <= 0 {
		cols = 1
	}
	buttons := make([]InlineButton, len(items))
	for i, item := range items {
		buttons[i] = mapper(item)
	}
	return b.Chunk(cols, buttons...).Build()
}

// ReplyBuilder fluently constructs ReplyKeyboard structures.
type ReplyBuilder struct {
	kb ReplyKeyboard
}

// NewReply initializes a ReplyBuilder.
func NewReply() *ReplyBuilder {
	return &ReplyBuilder{
		kb: ReplyKeyboard{
			Rows:   make([][]ReplyButton, 0),
			Resize: true,
		},
	}
}

// Row appends a row of reply buttons.
func (b *ReplyBuilder) Row(buttons ...ReplyButton) *ReplyBuilder {
	if len(buttons) > 0 {
		row := make([]ReplyButton, len(buttons))
		copy(row, buttons)
		b.kb.Rows = append(b.kb.Rows, row)
	}
	return b
}

// Resize toggles client auto-resizing of keyboard height (default true).
func (b *ReplyBuilder) Resize(resize bool) *ReplyBuilder {
	b.kb.Resize = resize
	return b
}

// OneTime instructs the client to hide the keyboard after one use.
func (b *ReplyBuilder) OneTime(oneTime bool) *ReplyBuilder {
	b.kb.OneTime = oneTime
	return b
}

// Placeholder sets placeholder hint text in the input field.
func (b *ReplyBuilder) Placeholder(placeholder string) *ReplyBuilder {
	b.kb.Placeholder = placeholder
	return b
}

// Selective shows keyboard only to specific users (e.g. mentioned or replied).
func (b *ReplyBuilder) Selective(selective bool) *ReplyBuilder {
	b.kb.Selective = selective
	return b
}

// Persistent keeps the custom keyboard visible when regular keyboard is opened.
func (b *ReplyBuilder) Persistent(persistent bool) *ReplyBuilder {
	b.kb.IsPersistent = persistent
	return b
}

// Build returns the configured ReplyKeyboard.
func (b *ReplyBuilder) Build() *ReplyKeyboard {
	return &b.kb
}

// TextButton creates a basic text ReplyButton.
func TextButton(text string) ReplyButton {
	return ReplyButton{Text: text}
}

// ContactButton creates a ReplyButton requesting phone number contact.
func ContactButton(text string) ReplyButton {
	return ReplyButton{
		Text:           text,
		RequestContact: true,
	}
}

// LocationButton creates a ReplyButton requesting user geo-location.
func LocationButton(text string) ReplyButton {
	return ReplyButton{
		Text:            text,
		RequestLocation: true,
	}
}
