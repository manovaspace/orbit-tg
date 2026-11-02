package telego_test

import (
	"context"
	"encoding/json"
	"testing"

	telegotype "github.com/mymmrac/telego"
	telegoapi "github.com/mymmrac/telego/telegoapi"

	tg "github.com/manovaspace/orbit-tg"
	telegoadapter "github.com/manovaspace/orbit-tg/adapter/telego"
)

func TestMapInlineQuery(t *testing.T) {
	query := &telegotype.InlineQuery{
		ID: "inline-query-1",
		From: telegotype.User{
			ID:        42,
			FirstName: "Ada",
			Username:  "ada",
		},
		Query:    "vpn",
		Offset:   "0",
		ChatType: "private",
	}

	mapped := telegoadapter.MapInlineQuery(query)
	if mapped == nil {
		t.Fatal("MapInlineQuery returned nil")
	}
	if mapped.ID != "inline-query-1" || mapped.Query != "vpn" || mapped.ChatType != "private" {
		t.Errorf("InlineQuery mapping mismatch: %+v", mapped)
	}
	if mapped.From.ID != 42 || mapped.From.Username != "ada" {
		t.Errorf("InlineQuery from mismatch: %+v", mapped.From)
	}

	update := telegoadapter.MapUpdate(&telegotype.Update{UpdateID: 9, InlineQuery: query})
	if update == nil || update.InlineQuery == nil || update.InlineQuery.ID != "inline-query-1" {
		t.Errorf("Update.InlineQuery mismatch: %+v", update)
	}
}

func TestMapChosenInlineResult(t *testing.T) {
	result := &telegotype.ChosenInlineResult{
		ResultID:        "result-9",
		From:            telegotype.User{ID: 42, FirstName: "Ada"},
		InlineMessageID: "inline-message-77",
		Query:           "vpn",
	}

	mapped := telegoadapter.MapChosenInlineResult(result)
	if mapped == nil {
		t.Fatal("MapChosenInlineResult returned nil")
	}
	if mapped.ResultID != "result-9" || mapped.InlineMessageID != "inline-message-77" {
		t.Errorf("ChosenInlineResult mapping mismatch: %+v", mapped)
	}

	update := telegoadapter.MapUpdate(&telegotype.Update{UpdateID: 10, ChosenInlineResult: result})
	if update == nil || update.ChosenInlineResult == nil || update.ChosenInlineResult.ResultID != "result-9" {
		t.Errorf("Update.ChosenInlineResult mismatch: %+v", update)
	}
}

func TestAnswerInlineQuerySendsCacheTimeZero(t *testing.T) {
	constructor := &captureConstructor{delegate: telegoapi.DefaultConstructor{}}
	caller := &recordingCaller{}
	adapter := telegoadapter.NewFromBot(newTestBot(t, caller, constructor), nil)

	results := []tg.InlineResult{{
		ID:    "article-1",
		Title: "Open VPN",
		Text:  "<b>Open VPN</b>",
	}}
	if err := adapter.AnswerInlineQuery(context.Background(), "query-1", results, tg.WithInlineCacheTime(0)); err != nil {
		t.Fatalf("AnswerInlineQuery: %v", err)
	}

	params, ok := constructor.last.(*telegotype.AnswerInlineQueryParams)
	if !ok {
		t.Fatalf("captured parameters = %T, want *telego.AnswerInlineQueryParams", constructor.last)
	}
	if params.CacheTime != 0 {
		t.Errorf("CacheTime = %d, want 0", params.CacheTime)
	}
	if params.InlineQueryID != "query-1" {
		t.Errorf("InlineQueryID = %q, want query-1", params.InlineQueryID)
	}
	if len(params.Results) != 1 {
		t.Fatalf("results = %d, want 1", len(params.Results))
	}
	if got := tg.ApplyInlineOptions(tg.WithInlineCacheTime(0)); got.CacheTime != 0 {
		t.Errorf("ApplyInlineOptions CacheTime = %d, want 0", got.CacheTime)
	}

	method, _ := caller.last(t)
	if method != "answerInlineQuery" {
		t.Errorf("method = %q, want answerInlineQuery", method)
	}
}

func TestSetMyCommandsScopes(t *testing.T) {
	cases := []struct {
		name     string
		scope    tg.CommandScope
		wantType string
		wantChat int64
	}{
		{name: "default", scope: tg.CommandScope{}, wantType: "default"},
		{name: "all private chats", scope: tg.CommandScope{Type: tg.CommandScopeAllPrivateChats}, wantType: "all_private_chats"},
		{name: "specific chat", scope: tg.CommandScope{Type: tg.CommandScopeChat, ChatID: -10042}, wantType: "chat", wantChat: -10042},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			caller := &recordingCaller{}
			adapter := telegoadapter.NewFromBot(newTestBot(t, caller, nil), nil)

			cmds := []tg.BotCommand{{Command: "start", Description: "Start the bot"}}
			if err := adapter.SetMyCommands(context.Background(), cmds, tc.scope); err != nil {
				t.Fatalf("SetMyCommands: %v", err)
			}

			method, body := caller.last(t)
			if method != "setMyCommands" {
				t.Errorf("method = %q, want setMyCommands", method)
			}
			var got struct {
				Commands []struct {
					Command     string `json:"command"`
					Description string `json:"description"`
				} `json:"commands"`
				Scope struct {
					Type   string `json:"type"`
					ChatID int64  `json:"chat_id"`
				} `json:"scope"`
			}
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("unmarshal body %s: %v", body, err)
			}
			if got.Scope.Type != tc.wantType || got.Scope.ChatID != tc.wantChat {
				t.Errorf("scope = %+v, want type %q chat %d", got.Scope, tc.wantType, tc.wantChat)
			}
			if len(got.Commands) != 1 || got.Commands[0].Command != "start" || got.Commands[0].Description != "Start the bot" {
				t.Errorf("commands = %+v", got.Commands)
			}
		})
	}
}
