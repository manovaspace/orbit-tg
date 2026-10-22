package keyboard

import "fmt"

// CallbackDismiss is the default callback payload for transient message dismissals.
const CallbackDismiss = "ogtg:dismiss"

// DismissOK creates a 1-button inline keyboard for acknowledging notifications.
func DismissOK(label string, customData ...string) *InlineKeyboard {
	if label == "" {
		label = "Ok ✅"
	}
	data := CallbackDismiss
	if len(customData) > 0 && customData[0] != "" {
		data = customData[0]
	}
	return NewInline().Row(Callback(label, data)).Build()
}

// ConfirmRow creates a 2-button inline confirmation keyboard (e.g. Yes / No).
func ConfirmRow(yesLabel, yesData, noLabel, noData string) *InlineKeyboard {
	if yesLabel == "" {
		yesLabel = "Confirm ✅"
	}
	if noLabel == "" {
		noLabel = "Cancel ❌"
	}
	return NewInline().Row(
		Callback(yesLabel, yesData),
		Callback(noLabel, noData),
	).Build()
}

// PaginationRow creates standard previous/next pagination controls.
// Prefix format used: fmt.Sprintf("%s:page:%d", prefix, targetPage)
func PaginationRow(page, totalPages int, prefix string) []InlineButton {
	if totalPages <= 1 {
		return nil
	}

	var row []InlineButton
	if page > 1 {
		row = append(row, Callback("◀ Prev", fmt.Sprintf("%s:page:%d", prefix, page-1)))
	} else {
		row = append(row, Callback("·", "noop"))
	}

	row = append(row, Callback(fmt.Sprintf("%d / %d", page, totalPages), "noop"))

	if page < totalPages {
		row = append(row, Callback("Next ▶", fmt.Sprintf("%s:page:%d", prefix, page+1)))
	} else {
		row = append(row, Callback("·", "noop"))
	}

	return row
}

// IsDismiss checks if a callback data string represents a dismiss action.
func IsDismiss(data string) bool {
	return data == CallbackDismiss
}
