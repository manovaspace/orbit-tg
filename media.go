package tg

import "context"

// InputFile identifies a file to send to Telegram. Exactly one of FileID, URL
// or Data selects the source: an existing Telegram file, a public URL, or an
// in-memory upload named by Filename.
type InputFile struct {
	FileID   string
	URL      string
	Filename string
	Data     []byte
}

// MediaType names the kind of media carried by a Media value.
type MediaType string

const (
	MediaPhoto    MediaType = "photo"
	MediaDocument MediaType = "document"
)

// Media describes replacement media for EditMedia.
type Media struct {
	Type      MediaType
	File      InputFile
	Caption   string
	ParseMode ParseMode
}

// ChatAction is the transient status shown while the bot prepares a response.
type ChatAction string

const (
	ChatActionTyping          ChatAction = "typing"
	ChatActionUploadPhoto     ChatAction = "upload_photo"
	ChatActionRecordVideo     ChatAction = "record_video"
	ChatActionUploadVideo     ChatAction = "upload_video"
	ChatActionRecordVoice     ChatAction = "record_voice"
	ChatActionUploadVoice     ChatAction = "upload_voice"
	ChatActionUploadDocument  ChatAction = "upload_document"
	ChatActionChooseSticker   ChatAction = "choose_sticker"
	ChatActionFindLocation    ChatAction = "find_location"
	ChatActionRecordVideoNote ChatAction = "record_video_note"
	ChatActionUploadVideoNote ChatAction = "upload_video_note"
)

// MediaMessenger sends photos and documents and edits their captions or media.
// It is additive to Messenger; a client that implements both can be used
// wherever either interface is required.
type MediaMessenger interface {
	SendPhoto(ctx context.Context, chatID int64, photo InputFile, caption string, opts ...SendOption) (*Message, error)
	SendDocument(ctx context.Context, chatID int64, doc InputFile, caption string, opts ...SendOption) (*Message, error)
	EditCaption(ctx context.Context, chatID int64, messageID int, caption string, opts ...EditOption) error
	EditMedia(ctx context.Context, chatID int64, messageID int, media Media, opts ...EditOption) error
}

// ChatActioner broadcasts a transient chat action such as "typing".
type ChatActioner interface {
	SendChatAction(ctx context.Context, chatID int64, action ChatAction) error
}
