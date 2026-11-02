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

// MapUpdate converts a telego.Update to tg.Update.
func MapUpdate(u *telego.Update) *tg.Update {
	return infra.MapUpdate(u)
}
