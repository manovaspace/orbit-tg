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

// MockMessenger provides a thread-safe in-memory implementation of tg.Messenger.
type MockMessenger struct {
	mu           sync.Mutex
	Sent         []SentRecord
	Edited       []EditRecord
	Deleted      []int
	Callbacks    []CallbackRecord
	Pinned       []PinRecord
	Unpinned     []PinRecord
	NextID       int
	FailWith     error
	FailCallback error
	Closed       bool
}

// NewMessenger initializes a MockMessenger.
func NewMessenger() *MockMessenger {
	return &MockMessenger{
		NextID:    1000,
		Sent:      make([]SentRecord, 0),
		Edited:    make([]EditRecord, 0),
		Deleted:   make([]int, 0),
		Callbacks: make([]CallbackRecord, 0),
		Pinned:    make([]PinRecord, 0),
		Unpinned:  make([]PinRecord, 0),
	}
}

var _ tg.Messenger = (*MockMessenger)(nil)

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

func (m *MockMessenger) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Closed = true
	return nil
}
