package telego

import (
	"github.com/mymmrac/telego"

	tg "github.com/manovaspace/orbit-tg"
	infra "github.com/manovaspace/orbit-tg/internal/infrastructure/telego"
)

// MapMarkup converts an abstract tg.Markup into a telego.ReplyMarkup.
func MapMarkup(markup tg.Markup) telego.ReplyMarkup {
	return infra.MapMarkup(markup)
}

// MapUser converts a telego.User to tg.User.
func MapUser(u *telego.User) *tg.User {
	return infra.MapUser(u)
}

// MapChat converts a telego.Chat to tg.Chat.
func MapChat(c telego.Chat) tg.Chat {
	return infra.MapChat(c)
}

// MapMessage converts a telego.Message to tg.Message.
func MapMessage(m *telego.Message) *tg.Message {
	return infra.MapMessage(m)
}

// MapCallbackQuery converts a telego.CallbackQuery to tg.CallbackQuery.
func MapCallbackQuery(q *telego.CallbackQuery) *tg.CallbackQuery {
	return infra.MapCallbackQuery(q)
}

// MapInlineQuery converts a telego.InlineQuery to tg.InlineQuery.
func MapInlineQuery(q *telego.InlineQuery) *tg.InlineQuery {
	if q == nil {
		return nil
	}
	var from tg.User
	if mapped := MapUser(&q.From); mapped != nil {
		from = *mapped
	}
	return &tg.InlineQuery{
		ID:       q.ID,
		From:     from,
		Query:    q.Query,
		Offset:   q.Offset,
		ChatType: q.ChatType,
	}
}

// MapChosenInlineResult converts a telego.ChosenInlineResult to tg.ChosenInlineResult.
func MapChosenInlineResult(r *telego.ChosenInlineResult) *tg.ChosenInlineResult {
	if r == nil {
		return nil
	}
	var from tg.User
	if mapped := MapUser(&r.From); mapped != nil {
		from = *mapped
	}
	return &tg.ChosenInlineResult{
		ResultID:        r.ResultID,
		From:            from,
		InlineMessageID: r.InlineMessageID,
		Query:           r.Query,
	}
}

// MapUpdate converts a telego.Update to tg.Update.
func MapUpdate(u *telego.Update) *tg.Update {
	mapped := infra.MapUpdate(u)
	if mapped == nil {
		return nil
	}
	mapped.InlineQuery = MapInlineQuery(u.InlineQuery)
	mapped.ChosenInlineResult = MapChosenInlineResult(u.ChosenInlineResult)
	return mapped
}
