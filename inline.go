package tg

import "context"

// InlineQuery represents an incoming inline query.
type InlineQuery struct {
	ID       string `json:"id"`
	From     User   `json:"from"`
	Query    string `json:"query"`
	Offset   string `json:"offset,omitempty"`
	ChatType string `json:"chat_type,omitempty"`
}

// ChosenInlineResult represents the result a user chose from an inline query.
type ChosenInlineResult struct {
	ResultID        string `json:"result_id"`
	From            User   `json:"from"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
	Query           string `json:"query,omitempty"`
}

// InlineResultType names the kind of inline result.
type InlineResultType string

const (
	InlineResultArticle  InlineResultType = "article"
	InlineResultPhoto    InlineResultType = "photo"
	InlineResultDocument InlineResultType = "document"
)

// InlineResult is one answer to an inline query. Article results use Text as
// the message body; photo and document results use URL as the file URL, Text as
// the caption and MimeType for documents.
type InlineResult struct {
	Type         InlineResultType
	ID           string
	Title        string
	Description  string
	Text         string
	ParseMode    ParseMode
	URL          string
	ThumbnailURL string
	MimeType     string
	Markup       Markup
}

// InlineQuerier answers an inline query with a bounded list of results.
type InlineQuerier interface {
	AnswerInlineQuery(ctx context.Context, queryID string, results []InlineResult, opts ...InlineOption) error
}
