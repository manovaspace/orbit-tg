package tg

import "context"

// BotCommand is one entry in the bot's command menu.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// CommandScopeType selects which users see a command list.
type CommandScopeType string

const (
	CommandScopeDefault               CommandScopeType = "default"
	CommandScopeAllPrivateChats       CommandScopeType = "all_private_chats"
	CommandScopeAllGroupChats         CommandScopeType = "all_group_chats"
	CommandScopeAllChatAdministrators CommandScopeType = "all_chat_administrators"
	CommandScopeChat                  CommandScopeType = "chat"
	CommandScopeChatAdministrators    CommandScopeType = "chat_administrators"
	CommandScopeChatMember            CommandScopeType = "chat_member"
)

// CommandScope describes the users a command list applies to. The zero value
// is the default scope; ChatID is required for chat scopes and UserID for
// CommandScopeChatMember.
type CommandScope struct {
	Type   CommandScopeType
	ChatID int64
	UserID int64
}

// MenuButtonType selects the bot's private-chat menu button.
type MenuButtonType string

const (
	MenuButtonDefault  MenuButtonType = "default"
	MenuButtonCommands MenuButtonType = "commands"
	MenuButtonWebApp   MenuButtonType = "web_app"
)

// MenuButton configures the menu button shown in a private chat. The zero
// value is the default button; Text and URL are used by the web_app button.
type MenuButton struct {
	Type MenuButtonType
	Text string
	URL  string
}

// CommandSetter publishes the bot's command list and menu button.
type CommandSetter interface {
	SetMyCommands(ctx context.Context, cmds []BotCommand, scope CommandScope) error
	SetChatMenuButton(ctx context.Context, chatID int64, button MenuButton) error
}
