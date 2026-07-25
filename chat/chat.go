package chat

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
)

func IsGroup(c telego.Chat) bool {
	return c.Type == telego.ChatTypeGroup || c.Type == telego.ChatTypeSupergroup
}

func MessageThread(msg telego.Message) int {
	return msg.MessageThreadID
}

func CallbackThread(msg telego.MaybeInaccessibleMessage) int {
	if m, ok := msg.(*telego.Message); ok && m != nil {
		return m.MessageThreadID
	}
	return 0
}

func CallbackMessageID(msg telego.MaybeInaccessibleMessage) int {
	if m, ok := msg.(*telego.Message); ok && m != nil {
		return m.MessageID
	}
	return 0
}

func DisplayName(u telego.User) string {
	if u.Username != "" {
		return u.Username
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	if name != "" {
		return name
	}
	return fmt.Sprintf("%d", u.ID)
}
