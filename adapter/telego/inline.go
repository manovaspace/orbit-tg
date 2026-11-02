package telego

import (
	"context"

	"github.com/mymmrac/telego"
	telegoapi "github.com/mymmrac/telego/telegoapi"
	tu "github.com/mymmrac/telego/telegoutil"

	tg "github.com/manovaspace/orbit-tg"
)

var (
	_ tg.InlineQuerier = (*TelegoAdapter)(nil)
	_ tg.CommandSetter = (*TelegoAdapter)(nil)
)

// AnswerInlineQuery sends inline results for an inline query. Cache time is
// always set explicitly so a zero value means "do not cache".
func (a *TelegoAdapter) AnswerInlineQuery(ctx context.Context, queryID string, results []tg.InlineResult, opts ...tg.InlineOption) error {
	cfg := tg.ApplyInlineOptions(opts...)
	return a.API().AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
		InlineQueryID: queryID,
		Results:       mapInlineResults(results),
		CacheTime:     cfg.CacheTime,
		IsPersonal:    cfg.IsPersonal,
		NextOffset:    cfg.NextOffset,
	})
}

// SetMyCommands publishes the command list for the given scope.
func (a *TelegoAdapter) SetMyCommands(ctx context.Context, cmds []tg.BotCommand, scope tg.CommandScope) error {
	commands := make([]telego.BotCommand, 0, len(cmds))
	for _, cmd := range cmds {
		commands = append(commands, telego.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		})
	}
	return a.API().SetMyCommands(ctx, &telego.SetMyCommandsParams{
		Commands: commands,
		Scope:    mapCommandScope(scope),
	})
}

// SetChatMenuButton configures the private-chat menu button. A zero chatID
// changes the bot-wide default.
func (a *TelegoAdapter) SetChatMenuButton(ctx context.Context, chatID int64, button tg.MenuButton) error {
	return a.API().SetChatMenuButton(ctx, &telego.SetChatMenuButtonParams{
		ChatID:     chatID,
		MenuButton: mapMenuButton(button),
	})
}

// inlineQueryParams mirrors telego.AnswerInlineQueryParams but always emits
// cache_time, including zero. telego tags the field omitempty, which would let
// Telegram apply its 300 second default instead of the requested value.
type inlineQueryParams struct {
	InlineQueryID string                           `json:"inline_query_id"`
	Results       []telego.InlineQueryResult       `json:"results"`
	CacheTime     int                              `json:"cache_time"`
	IsPersonal    bool                             `json:"is_personal,omitempty"`
	NextOffset    string                           `json:"next_offset,omitempty"`
	Button        *telego.InlineQueryResultsButton `json:"button,omitempty"`
}

// inlineQueryConstructor serializes answerInlineQuery with an explicit
// cache_time while delegating every other method to telego's constructor.
type inlineQueryConstructor struct{}

func (inlineQueryConstructor) JSONRequest(parameters any) (*telegoapi.RequestData, error) {
	if params, ok := parameters.(*telego.AnswerInlineQueryParams); ok && params != nil {
		return telegoapi.DefaultConstructor{}.JSONRequest(inlineQueryParams{
			InlineQueryID: params.InlineQueryID,
			Results:       params.Results,
			CacheTime:     params.CacheTime,
			IsPersonal:    params.IsPersonal,
			NextOffset:    params.NextOffset,
			Button:        params.Button,
		})
	}
	return telegoapi.DefaultConstructor{}.JSONRequest(parameters)
}

func (inlineQueryConstructor) MultipartRequest(parameters map[string]string, files map[string]telegoapi.NamedReader) (*telegoapi.RequestData, error) {
	return telegoapi.DefaultConstructor{}.MultipartRequest(parameters, files)
}

func mapInlineResults(results []tg.InlineResult) []telego.InlineQueryResult {
	mapped := make([]telego.InlineQueryResult, 0, len(results))
	for _, result := range results {
		switch result.Type {
		case tg.InlineResultPhoto:
			thumbnail := result.ThumbnailURL
			if thumbnail == "" {
				thumbnail = result.URL
			}
			mapped = append(mapped, &telego.InlineQueryResultPhoto{
				Type:         telego.ResultTypePhoto,
				ID:           result.ID,
				PhotoURL:     result.URL,
				ThumbnailURL: thumbnail,
				Title:        result.Title,
				Description:  result.Description,
				Caption:      result.Text,
				ParseMode:    inlineParseMode(result.ParseMode),
				ReplyMarkup:  mapInlineKeyboard(result.Markup),
			})
		case tg.InlineResultDocument:
			mapped = append(mapped, &telego.InlineQueryResultDocument{
				Type:         telego.ResultTypeDocument,
				ID:           result.ID,
				Title:        result.Title,
				Description:  result.Description,
				DocumentURL:  result.URL,
				MimeType:     result.MimeType,
				Caption:      result.Text,
				ParseMode:    inlineParseMode(result.ParseMode),
				ThumbnailURL: result.ThumbnailURL,
				ReplyMarkup:  mapInlineKeyboard(result.Markup),
			})
		default:
			mapped = append(mapped, &telego.InlineQueryResultArticle{
				Type:         telego.ResultTypeArticle,
				ID:           result.ID,
				Title:        result.Title,
				Description:  result.Description,
				URL:          result.URL,
				ThumbnailURL: result.ThumbnailURL,
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: result.Text,
					ParseMode:   inlineParseMode(result.ParseMode),
				},
				ReplyMarkup: mapInlineKeyboard(result.Markup),
			})
		}
	}
	return mapped
}

func inlineParseMode(mode tg.ParseMode) string {
	if mode == "" {
		return string(tg.ParseModeHTML)
	}
	return string(mode)
}

func mapInlineKeyboard(markup tg.Markup) *telego.InlineKeyboardMarkup {
	if markup == nil {
		return nil
	}
	keyboard, _ := MapMarkup(markup).(*telego.InlineKeyboardMarkup)
	return keyboard
}

func mapCommandScope(scope tg.CommandScope) telego.BotCommandScope {
	switch scope.Type {
	case tg.CommandScopeAllPrivateChats:
		return &telego.BotCommandScopeAllPrivateChats{Type: telego.ScopeTypeAllPrivateChats}
	case tg.CommandScopeAllGroupChats:
		return &telego.BotCommandScopeAllGroupChats{Type: telego.ScopeTypeAllGroupChats}
	case tg.CommandScopeAllChatAdministrators:
		return &telego.BotCommandScopeAllChatAdministrators{Type: telego.ScopeTypeAllChatAdministrators}
	case tg.CommandScopeChat:
		return &telego.BotCommandScopeChat{Type: telego.ScopeTypeChat, ChatID: tu.ID(scope.ChatID)}
	case tg.CommandScopeChatAdministrators:
		return &telego.BotCommandScopeChatAdministrators{Type: telego.ScopeTypeChatAdministrators, ChatID: tu.ID(scope.ChatID)}
	case tg.CommandScopeChatMember:
		return &telego.BotCommandScopeChatMember{Type: telego.ScopeTypeChatMember, ChatID: tu.ID(scope.ChatID), UserID: scope.UserID}
	default:
		return &telego.BotCommandScopeDefault{Type: telego.ScopeTypeDefault}
	}
}

func mapMenuButton(button tg.MenuButton) telego.MenuButton {
	switch button.Type {
	case tg.MenuButtonCommands:
		return &telego.MenuButtonCommands{Type: telego.ButtonTypeCommands}
	case tg.MenuButtonWebApp:
		return &telego.MenuButtonWebApp{
			Type:   telego.ButtonTypeWebApp,
			Text:   button.Text,
			WebApp: telego.WebAppInfo{URL: button.URL},
		}
	default:
		return &telego.MenuButtonDefault{Type: telego.ButtonTypeDefault}
	}
}
