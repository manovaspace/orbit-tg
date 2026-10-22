package chat

import (
	tg "github.com/manovaspace/orbit-tg"
)

// IsGroup returns whether the chat is a group or supergroup.
func IsGroup(c tg.Chat) bool {
	return c.IsGroup()
}

// MessageThread extracts the thread/topic ID of a message, if any.
func MessageThread(msg tg.Message) int {
	return msg.ThreadID
}

// CallbackThread extracts the thread/topic ID of the message attached to a callback query.
func CallbackThread(q tg.CallbackQuery) int {
	if q.Message != nil {
		return q.Message.ThreadID
	}
	return 0
}

// CallbackMessageID extracts the message ID attached to a callback query.
func CallbackMessageID(q tg.CallbackQuery) int {
	if q.Message != nil {
		return q.Message.ID
	}
	return 0
}

// DisplayName formats a readable name for a user.
func DisplayName(u tg.User) string {
	return u.DisplayName()
}
