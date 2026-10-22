package keyboard

// InlineButton represents an individual interactive inline button.
type InlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
	WebAppURL    string `json:"web_app_url,omitempty"`
	CopyText     string `json:"copy_text,omitempty"`
}

// InlineKeyboard represents an inline keyboard attached to a message.
type InlineKeyboard struct {
	Rows [][]InlineButton `json:"rows"`
}

func (k *InlineKeyboard) MarkupType() string {
	return "inline"
}

// ReplyButton represents a button in a custom reply keyboard.
type ReplyButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

// ReplyKeyboard represents a custom reply keyboard shown in place of system keyboard.
type ReplyKeyboard struct {
	Rows         [][]ReplyButton `json:"rows"`
	Resize       bool            `json:"resize"`
	OneTime      bool            `json:"one_time"`
	Placeholder  string          `json:"placeholder,omitempty"`
	Selective    bool            `json:"selective"`
	IsPersistent bool            `json:"is_persistent"`
}

func (k *ReplyKeyboard) MarkupType() string {
	return "reply"
}

// RemoveKeyboard instructs clients to dismiss the custom reply keyboard.
type RemoveKeyboard struct {
	Selective bool `json:"selective"`
}

func (k *RemoveKeyboard) MarkupType() string {
	return "remove"
}

// ForceReply prompts the user to reply directly to the message.
type ForceReply struct {
	Placeholder string `json:"placeholder,omitempty"`
	Selective   bool   `json:"selective"`
}

func (k *ForceReply) MarkupType() string {
	return "force_reply"
}
