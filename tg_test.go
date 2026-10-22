package tg_test

import (
	"testing"

	tg "github.com/manovaspace/orbit-tg"
)

type dummyMarkup struct{}

func (d dummyMarkup) MarkupType() string { return "dummy" }

func TestUser_DisplayName(t *testing.T) {
	tests := []struct {
		user     tg.User
		expected string
	}{
		{
			user:     tg.User{ID: 101, Username: "alireza"},
			expected: "alireza",
		},
		{
			user:     tg.User{ID: 102, FirstName: "Sara", LastName: "Rad"},
			expected: "Sara Rad",
		},
		{
			user:     tg.User{ID: 103, FirstName: "Solo"},
			expected: "Solo",
		},
		{
			user:     tg.User{ID: 104},
			expected: "104",
		},
	}

	for _, tt := range tests {
		if got := tt.user.DisplayName(); got != tt.expected {
			t.Errorf("DisplayName() = %q, want %q", got, tt.expected)
		}
	}
}

func TestChat_IsGroup(t *testing.T) {
	tests := []struct {
		chatType string
		isGroup  bool
	}{
		{"private", false},
		{"channel", false},
		{"group", true},
		{"supergroup", true},
	}

	for _, tt := range tests {
		chat := tg.Chat{Type: tt.chatType}
		if got := chat.IsGroup(); got != tt.isGroup {
			t.Errorf("IsGroup(%q) = %v, want %v", tt.chatType, got, tt.isGroup)
		}
	}
}

func TestApplySendOptions(t *testing.T) {
	markup := dummyMarkup{}
	cfg := tg.ApplySendOptions(
		tg.WithThreadID(42),
		tg.WithMarkup(markup),
		tg.WithSilent(true),
		tg.WithProtectContent(true),
		tg.WithReplyTo(999),
		tg.WithDisableLinkPreview(true),
	)

	if cfg.ThreadID != 42 {
		t.Errorf("expected ThreadID 42, got %d", cfg.ThreadID)
	}
	if cfg.Markup == nil || cfg.Markup.MarkupType() != "dummy" {
		t.Errorf("expected dummy markup, got %v", cfg.Markup)
	}
	if !cfg.Silent {
		t.Errorf("expected Silent true")
	}
	if !cfg.ProtectContent {
		t.Errorf("expected ProtectContent true")
	}
	if cfg.ReplyToMessageID != 999 {
		t.Errorf("expected ReplyToMessageID 999, got %d", cfg.ReplyToMessageID)
	}
	if !cfg.DisableLinkPreview {
		t.Errorf("expected DisableLinkPreview true")
	}
}

func TestApplyEditOptions(t *testing.T) {
	markup := dummyMarkup{}
	cfg := tg.ApplyEditOptions(
		tg.WithEditMarkup(markup),
		tg.WithEditDisableLinkPreview(true),
	)

	if cfg.Markup == nil || cfg.Markup.MarkupType() != "dummy" {
		t.Errorf("expected dummy markup, got %v", cfg.Markup)
	}
	if !cfg.DisableLinkPreview {
		t.Errorf("expected DisableLinkPreview true")
	}
}

func TestApplyCallbackOptions(t *testing.T) {
	cfg := tg.ApplyCallbackOptions(
		tg.WithAlert(true),
		tg.WithCallbackText("Action processed"),
		tg.WithCallbackURL("https://example.com"),
		tg.WithCacheTime(300),
	)

	if !cfg.ShowAlert {
		t.Errorf("expected ShowAlert true")
	}
	if cfg.Text != "Action processed" {
		t.Errorf("expected Text 'Action processed', got %q", cfg.Text)
	}
	if cfg.URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got %q", cfg.URL)
	}
	if cfg.CacheTime != 300 {
		t.Errorf("expected CacheTime 300, got %d", cfg.CacheTime)
	}
}
