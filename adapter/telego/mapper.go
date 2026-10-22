package telego

import (
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/keyboard"
)

// MapMarkup converts an abstract tg.Markup into a telego.ReplyMarkup.
func MapMarkup(markup tg.Markup) telego.ReplyMarkup {
	if markup == nil {
		return nil
	}

	switch m := markup.(type) {
	case *keyboard.InlineKeyboard:
		if m == nil || len(m.Rows) == 0 {
			return nil
		}
		rows := make([][]telego.InlineKeyboardButton, len(m.Rows))
		for rIdx, row := range m.Rows {
			btnRow := make([]telego.InlineKeyboardButton, len(row))
			for cIdx, btn := range row {
				tb := tu.InlineKeyboardButton(btn.Text)
				if btn.CallbackData != "" {
					tb = tb.WithCallbackData(btn.CallbackData)
				}
				if btn.URL != "" {
					tb = tb.WithURL(btn.URL)
				}
				if btn.WebAppURL != "" {
					tb = tb.WithWebApp(&telego.WebAppInfo{URL: btn.WebAppURL})
				}
				if btn.CopyText != "" {
					tb = tb.WithCopyText(&telego.CopyTextButton{Text: btn.CopyText})
				}
				btnRow[cIdx] = tb
			}
			rows[rIdx] = btnRow
		}
		return &telego.InlineKeyboardMarkup{InlineKeyboard: rows}

	case *keyboard.ReplyKeyboard:
		if m == nil {
			return nil
		}
		rows := make([][]telego.KeyboardButton, len(m.Rows))
		for rIdx, row := range m.Rows {
			btnRow := make([]telego.KeyboardButton, len(row))
			for cIdx, btn := range row {
				kb := tu.KeyboardButton(btn.Text)
				if btn.RequestContact {
					kb = kb.WithRequestContact()
				}
				if btn.RequestLocation {
					kb = kb.WithRequestLocation()
				}
				btnRow[cIdx] = kb
			}
			rows[rIdx] = btnRow
		}
		rk := tu.Keyboard(rows...)
		if m.Resize {
			rk = rk.WithResizeKeyboard()
		}
		if m.OneTime {
			rk = rk.WithOneTimeKeyboard()
		}
		if m.Placeholder != "" {
			rk = rk.WithInputFieldPlaceholder(m.Placeholder)
		}
		if m.Selective {
			rk = rk.WithSelective()
		}
		if m.IsPersistent {
			rk = rk.WithIsPersistent()
		}
		return rk

	case *keyboard.RemoveKeyboard:
		if m == nil {
			return nil
		}
		rem := &telego.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      m.Selective,
		}
		return rem

	case *keyboard.ForceReply:
		if m == nil {
			return nil
		}
		fr := &telego.ForceReply{
			ForceReply:            true,
			InputFieldPlaceholder: m.Placeholder,
			Selective:             m.Selective,
		}
		return fr

	default:
		return nil
	}
}

// MapUser converts a telego.User to tg.User.
func MapUser(u *telego.User) *tg.User {
	if u == nil {
		return nil
	}
	return &tg.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		IsBot:     u.IsBot,
	}
}

// MapChat converts a telego.Chat to tg.Chat.
func MapChat(c telego.Chat) tg.Chat {
	return tg.Chat{
		ID:       c.ID,
		Type:     c.Type,
		Title:    c.Title,
		Username: c.Username,
	}
}

// MapMessage converts a telego.Message to tg.Message.
func MapMessage(m *telego.Message) *tg.Message {
	if m == nil {
		return nil
	}
	var date time.Time
	if m.Date > 0 {
		date = time.Unix(m.Date, 0)
	}
	return &tg.Message{
		ID:       m.MessageID,
		Chat:     MapChat(m.Chat),
		ThreadID: m.MessageThreadID,
		From:     MapUser(m.From),
		Date:     date,
		Text:     m.Text,
	}
}

// MapCallbackQuery converts a telego.CallbackQuery to tg.CallbackQuery.
func MapCallbackQuery(q *telego.CallbackQuery) *tg.CallbackQuery {
	if q == nil {
		return nil
	}
	var user tg.User
	if mapped := MapUser(&q.From); mapped != nil {
		user = *mapped
	}

	var msg *tg.Message
	if m, ok := q.Message.(*telego.Message); ok && m != nil {
		msg = MapMessage(m)
	}

	return &tg.CallbackQuery{
		ID:      q.ID,
		From:    user,
		Message: msg,
		Data:    q.Data,
	}
}

// MapUpdate converts a telego.Update to tg.Update.
func MapUpdate(u *telego.Update) *tg.Update {
	if u == nil {
		return nil
	}
	return &tg.Update{
		ID:            u.UpdateID,
		Message:       MapMessage(u.Message),
		CallbackQuery: MapCallbackQuery(u.CallbackQuery),
	}
}
