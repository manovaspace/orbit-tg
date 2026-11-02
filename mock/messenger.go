package mock

import (
	"context"
	"sync"
	"time"

	tg "github.com/manovaspace/orbit-tg"
)

// SentRecord records parameters from a SendText or SendHTML invocation.
type SentRecord struct {
	ChatID int64
	Text   string
	IsHTML bool
	Config tg.SendConfig
}

// EditRecord records parameters from an EditText or EditHTML invocation.
type EditRecord struct {
	ChatID    int64
	MessageID int
	Text      string
	IsHTML    bool
	Config    tg.EditConfig
}

// CallbackRecord records parameters from an AnswerCallback invocation.
type CallbackRecord struct {
	QueryID string
	Config  tg.CallbackConfig
}

// PinRecord records parameters from Pin and Unpin invocations.
type PinRecord struct {
	ChatID    int64
	MessageID int
	Silent    bool
}

// PhotoRecord records parameters from a SendPhoto invocation.
type PhotoRecord struct {
	ChatID  int64
	Photo   tg.InputFile
	Caption string
	Config  tg.SendConfig
}

// DocumentRecord records parameters from a SendDocument invocation.
type DocumentRecord struct {
	ChatID   int64
	Document tg.InputFile
	Caption  string
	Config   tg.SendConfig
}

// InlineAnswerRecord records parameters from an AnswerInlineQuery invocation.
type InlineAnswerRecord struct {
	QueryID string
	Results []tg.InlineResult
	Config  tg.InlineConfig
}

// ChatActionRecord records parameters from a SendChatAction invocation.
type ChatActionRecord struct {
	ChatID int64
	Action tg.ChatAction
}

// CommandRecord records parameters from a SetMyCommands invocation.
type CommandRecord struct {
	Commands []tg.BotCommand
	Scope    tg.CommandScope
}

// MenuButtonRecord records parameters from a SetChatMenuButton invocation.
type MenuButtonRecord struct {
	ChatID int64
	Button tg.MenuButton
}

// MockMessenger provides a thread-safe in-memory implementation of tg.Messenger.
type MockMessenger struct {
	mu            sync.Mutex
	Sent          []SentRecord
	Edited        []EditRecord
	Deleted       []int
	Callbacks     []CallbackRecord
	Pinned        []PinRecord
	Unpinned      []PinRecord
	Photos        []PhotoRecord
	Documents     []DocumentRecord
	InlineAnswers []InlineAnswerRecord
	ChatActions   []ChatActionRecord
	Commands      []CommandRecord
	MenuButtons   []MenuButtonRecord
	NextID        int
	FailWith      error
	FailCallback  error
	Closed        bool
}

// NewMessenger initializes a MockMessenger.
func NewMessenger() *MockMessenger {
	return &MockMessenger{
		NextID:        1000,
		Sent:          make([]SentRecord, 0),
		Edited:        make([]EditRecord, 0),
		Deleted:       make([]int, 0),
		Callbacks:     make([]CallbackRecord, 0),
		Pinned:        make([]PinRecord, 0),
		Unpinned:      make([]PinRecord, 0),
		Photos:        make([]PhotoRecord, 0),
		Documents:     make([]DocumentRecord, 0),
		InlineAnswers: make([]InlineAnswerRecord, 0),
		ChatActions:   make([]ChatActionRecord, 0),
		Commands:      make([]CommandRecord, 0),
		MenuButtons:   make([]MenuButtonRecord, 0),
	}
}

var (
	_ tg.Messenger      = (*MockMessenger)(nil)
	_ tg.MediaMessenger = (*MockMessenger)(nil)
	_ tg.InlineQuerier  = (*MockMessenger)(nil)
	_ tg.ChatActioner   = (*MockMessenger)(nil)
	_ tg.CommandSetter  = (*MockMessenger)(nil)
)

func (m *MockMessenger) SendText(ctx context.Context, chatID int64, text string, opts ...tg.SendOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplySendOptions(opts...)
	m.Sent = append(m.Sent, SentRecord{
		ChatID: chatID,
		Text:   text,
		IsHTML: false,
		Config: cfg,
	})

	m.NextID++
	return &tg.Message{
		ID:       m.NextID,
		Chat:     tg.Chat{ID: chatID},
		ThreadID: cfg.ThreadID,
		Date:     time.Now(),
		Text:     text,
	}, nil
}

func (m *MockMessenger) SendHTML(ctx context.Context, chatID int64, html string, opts ...tg.SendOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplySendOptions(opts...)
	m.Sent = append(m.Sent, SentRecord{
		ChatID: chatID,
		Text:   html,
		IsHTML: true,
		Config: cfg,
	})

	m.NextID++
	return &tg.Message{
		ID:       m.NextID,
		Chat:     tg.Chat{ID: chatID},
		ThreadID: cfg.ThreadID,
		Date:     time.Now(),
		Text:     html,
	}, nil
}

func (m *MockMessenger) EditText(ctx context.Context, chatID int64, messageID int, text string, opts ...tg.EditOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplyEditOptions(opts...)
	m.Edited = append(m.Edited, EditRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
		IsHTML:    false,
		Config:    cfg,
	})

	return &tg.Message{
		ID:   messageID,
		Chat: tg.Chat{ID: chatID},
		Date: time.Now(),
		Text: text,
	}, nil
}

func (m *MockMessenger) EditHTML(ctx context.Context, chatID int64, messageID int, html string, opts ...tg.EditOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplyEditOptions(opts...)
	m.Edited = append(m.Edited, EditRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      html,
		IsHTML:    true,
		Config:    cfg,
	})

	return &tg.Message{
		ID:   messageID,
		Chat: tg.Chat{ID: chatID},
		Date: time.Now(),
		Text: html,
	}, nil
}

func (m *MockMessenger) EditKeyboard(ctx context.Context, chatID int64, messageID int, markup tg.Markup) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Edited = append(m.Edited, EditRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Config:    tg.EditConfig{Markup: markup},
	})
	return nil
}

func (m *MockMessenger) Delete(ctx context.Context, chatID int64, messageIDs ...int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Deleted = append(m.Deleted, messageIDs...)
	return nil
}

func (m *MockMessenger) AnswerCallback(ctx context.Context, queryID string, opts ...tg.CallbackOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailCallback != nil {
		return m.FailCallback
	}

	cfg := tg.ApplyCallbackOptions(opts...)
	m.Callbacks = append(m.Callbacks, CallbackRecord{
		QueryID: queryID,
		Config:  cfg,
	})
	return nil
}

func (m *MockMessenger) Pin(ctx context.Context, chatID int64, messageID int, silent bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Pinned = append(m.Pinned, PinRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Silent:    silent,
	})
	return nil
}

func (m *MockMessenger) Unpin(ctx context.Context, chatID int64, messageID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Unpinned = append(m.Unpinned, PinRecord{
		ChatID:    chatID,
		MessageID: messageID,
	})
	return nil
}

func (m *MockMessenger) SendPhoto(ctx context.Context, chatID int64, photo tg.InputFile, caption string, opts ...tg.SendOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplySendOptions(opts...)
	m.Photos = append(m.Photos, PhotoRecord{
		ChatID:  chatID,
		Photo:   photo,
		Caption: caption,
		Config:  cfg,
	})

	m.NextID++
	return &tg.Message{
		ID:       m.NextID,
		Chat:     tg.Chat{ID: chatID},
		ThreadID: cfg.ThreadID,
		Date:     time.Now(),
		Text:     caption,
	}, nil
}

func (m *MockMessenger) SendDocument(ctx context.Context, chatID int64, doc tg.InputFile, caption string, opts ...tg.SendOption) (*tg.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return nil, m.FailWith
	}

	cfg := tg.ApplySendOptions(opts...)
	m.Documents = append(m.Documents, DocumentRecord{
		ChatID:   chatID,
		Document: doc,
		Caption:  caption,
		Config:   cfg,
	})

	m.NextID++
	return &tg.Message{
		ID:       m.NextID,
		Chat:     tg.Chat{ID: chatID},
		ThreadID: cfg.ThreadID,
		Date:     time.Now(),
		Text:     caption,
	}, nil
}

func (m *MockMessenger) EditCaption(ctx context.Context, chatID int64, messageID int, caption string, opts ...tg.EditOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Edited = append(m.Edited, EditRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      caption,
		Config:    tg.ApplyEditOptions(opts...),
	})
	return nil
}

func (m *MockMessenger) EditMedia(ctx context.Context, chatID int64, messageID int, media tg.Media, opts ...tg.EditOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Edited = append(m.Edited, EditRecord{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      media.Caption,
		Config:    tg.ApplyEditOptions(opts...),
	})
	return nil
}

func (m *MockMessenger) SendChatAction(ctx context.Context, chatID int64, action tg.ChatAction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.ChatActions = append(m.ChatActions, ChatActionRecord{
		ChatID: chatID,
		Action: action,
	})
	return nil
}

func (m *MockMessenger) AnswerInlineQuery(ctx context.Context, queryID string, results []tg.InlineResult, opts ...tg.InlineOption) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.InlineAnswers = append(m.InlineAnswers, InlineAnswerRecord{
		QueryID: queryID,
		Results: append([]tg.InlineResult(nil), results...),
		Config:  tg.ApplyInlineOptions(opts...),
	})
	return nil
}

func (m *MockMessenger) SetMyCommands(ctx context.Context, cmds []tg.BotCommand, scope tg.CommandScope) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.Commands = append(m.Commands, CommandRecord{
		Commands: append([]tg.BotCommand(nil), cmds...),
		Scope:    scope,
	})
	return nil
}

func (m *MockMessenger) SetChatMenuButton(ctx context.Context, chatID int64, button tg.MenuButton) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.FailWith != nil {
		return m.FailWith
	}

	m.MenuButtons = append(m.MenuButtons, MenuButtonRecord{
		ChatID: chatID,
		Button: button,
	})
	return nil
}

func (m *MockMessenger) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Closed = true
	return nil
}
