package keyboard

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const CallbackDismiss = "ogtg:dismiss"

func DismissOK(label string) *telego.InlineKeyboardMarkup {
	if label == "" {
		label = "Ok ✅"
	}
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(label).WithCallbackData(CallbackDismiss),
		),
	)
}

func ConfirmRow(yesLabel, yesData, noLabel, noData string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(yesLabel).WithCallbackData(yesData),
			tu.InlineKeyboardButton(noLabel).WithCallbackData(noData),
		),
	)
}

func IsDismiss(data string) bool {
	return data == CallbackDismiss
}
