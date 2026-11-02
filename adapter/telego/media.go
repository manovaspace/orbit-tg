package telego

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/message"
)

var (
	_ tg.MediaMessenger = (*TelegoAdapter)(nil)
	_ tg.ChatActioner   = (*TelegoAdapter)(nil)
)

// SendPhoto sends a photo with an optional HTML caption.
func (a *TelegoAdapter) SendPhoto(ctx context.Context, chatID int64, photo tg.InputFile, caption string, opts ...tg.SendOption) (*tg.Message, error) {
	cfg := tg.ApplySendOptions(opts...)
	params := &telego.SendPhotoParams{
		ChatID: tu.ID(chatID),
		Photo:  mapInputFile(photo),
	}
	if caption != "" {
		params.Caption = message.SanitizeHTML(caption)
		params.ParseMode = string(tg.ParseModeHTML)
	}
	if cfg.ThreadID > 0 {
		params.MessageThreadID = cfg.ThreadID
	}
	if kb := mapInlineKeyboard(cfg.Markup); kb != nil {
		params.ReplyMarkup = kb
	}
	if cfg.Silent {
		params.DisableNotification = true
	}
	if cfg.ProtectContent {
		params.ProtectContent = true
	}
	if cfg.ReplyToMessageID > 0 {
		params.ReplyParameters = &telego.ReplyParameters{MessageID: cfg.ReplyToMessageID}
	}

	sent, err := a.API().SendPhoto(ctx, params)
	if err != nil {
		return nil, err
	}
	return MapMessage(sent), nil
}

// SendDocument sends a document with an optional HTML caption.
func (a *TelegoAdapter) SendDocument(ctx context.Context, chatID int64, doc tg.InputFile, caption string, opts ...tg.SendOption) (*tg.Message, error) {
	cfg := tg.ApplySendOptions(opts...)
	params := &telego.SendDocumentParams{
		ChatID:   tu.ID(chatID),
		Document: mapInputFile(doc),
	}
	if caption != "" {
		params.Caption = message.SanitizeHTML(caption)
		params.ParseMode = string(tg.ParseModeHTML)
	}
	if cfg.ThreadID > 0 {
		params.MessageThreadID = cfg.ThreadID
	}
	if kb := mapInlineKeyboard(cfg.Markup); kb != nil {
		params.ReplyMarkup = kb
	}
	if cfg.Silent {
		params.DisableNotification = true
	}
	if cfg.ProtectContent {
		params.ProtectContent = true
	}
	if cfg.ReplyToMessageID > 0 {
		params.ReplyParameters = &telego.ReplyParameters{MessageID: cfg.ReplyToMessageID}
	}

	sent, err := a.API().SendDocument(ctx, params)
	if err != nil {
		return nil, err
	}
	return MapMessage(sent), nil
}

// EditCaption replaces the HTML caption of a media message.
func (a *TelegoAdapter) EditCaption(ctx context.Context, chatID int64, messageID int, caption string, opts ...tg.EditOption) error {
	cfg := tg.ApplyEditOptions(opts...)
	params := &telego.EditMessageCaptionParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
	}
	if caption != "" {
		params.Caption = message.SanitizeHTML(caption)
		params.ParseMode = string(tg.ParseModeHTML)
	}
	if kb := mapInlineKeyboard(cfg.Markup); kb != nil {
		params.ReplyMarkup = kb
	}
	_, err := a.API().EditMessageCaption(ctx, params)
	return err
}

// EditMedia replaces the media of an existing message.
func (a *TelegoAdapter) EditMedia(ctx context.Context, chatID int64, messageID int, media tg.Media, opts ...tg.EditOption) error {
	mapped, err := mapMedia(media)
	if err != nil {
		return err
	}
	cfg := tg.ApplyEditOptions(opts...)
	params := &telego.EditMessageMediaParams{
		ChatID:    tu.ID(chatID),
		MessageID: messageID,
		Media:     mapped,
	}
	if kb := mapInlineKeyboard(cfg.Markup); kb != nil {
		params.ReplyMarkup = kb
	}
	_, err = a.API().EditMessageMedia(ctx, params)
	return err
}

// SendChatAction broadcasts a transient chat action such as "typing".
func (a *TelegoAdapter) SendChatAction(ctx context.Context, chatID int64, action tg.ChatAction) error {
	return a.API().SendChatAction(ctx, &telego.SendChatActionParams{
		ChatID: tu.ID(chatID),
		Action: string(action),
	})
}

func mapInputFile(file tg.InputFile) telego.InputFile {
	switch {
	case file.FileID != "":
		return tu.FileFromID(file.FileID)
	case file.URL != "":
		return tu.FileFromURL(file.URL)
	case len(file.Data) > 0:
		name := file.Filename
		if name == "" {
			name = "file"
		}
		return tu.FileFromBytes(file.Data, name)
	default:
		return telego.InputFile{}
	}
}

func mapMedia(media tg.Media) (telego.InputMedia, error) {
	parseMode := media.ParseMode
	if parseMode == "" {
		parseMode = tg.ParseModeHTML
	}
	caption := ""
	if media.Caption != "" {
		caption = message.SanitizeHTML(media.Caption)
	}

	switch media.Type {
	case "", tg.MediaPhoto:
		return &telego.InputMediaPhoto{
			Type:      telego.MediaTypePhoto,
			Media:     mapInputFile(media.File),
			Caption:   caption,
			ParseMode: string(parseMode),
		}, nil
	case tg.MediaDocument:
		return &telego.InputMediaDocument{
			Type:      telego.MediaTypeDocument,
			Media:     mapInputFile(media.File),
			Caption:   caption,
			ParseMode: string(parseMode),
		}, nil
	default:
		return nil, fmt.Errorf("tg: unsupported media type %q", media.Type)
	}
}
