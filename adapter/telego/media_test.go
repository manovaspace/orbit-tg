package telego_test

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	telegotype "github.com/mymmrac/telego"
	telegoapi "github.com/mymmrac/telego/telegoapi"

	tg "github.com/manovaspace/orbit-tg"
	telegoadapter "github.com/manovaspace/orbit-tg/adapter/telego"
)

// recordingCaller captures the request the adapter hands to telego.
type recordingCaller struct {
	mu     sync.Mutex
	method string
	body   []byte
}

func (r *recordingCaller) Call(_ context.Context, url string, data *telegoapi.RequestData) (*telegoapi.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if i := strings.LastIndex(url, "/"); i >= 0 {
		r.method = url[i+1:]
	}
	r.body = append([]byte(nil), data.BodyRaw...)
	return &telegoapi.Response{Ok: true, Result: json.RawMessage(`true`)}, nil
}

func (r *recordingCaller) last(t *testing.T) (string, []byte) {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.method, append([]byte(nil), r.body...)
}

// captureConstructor keeps the parameters before telego marshals them.
type captureConstructor struct {
	delegate telegoapi.RequestConstructor
	last     any
}

func (c *captureConstructor) JSONRequest(parameters any) (*telegoapi.RequestData, error) {
	c.last = parameters
	return c.delegate.JSONRequest(parameters)
}

func (c *captureConstructor) MultipartRequest(parameters map[string]string, files map[string]telegoapi.NamedReader) (*telegoapi.RequestData, error) {
	return c.delegate.MultipartRequest(parameters, files)
}

func newTestBot(t *testing.T, caller telegoapi.Caller, constructor telegoapi.RequestConstructor) *telegotype.Bot {
	t.Helper()
	opts := []telegotype.BotOption{
		telegotype.WithAPICaller(caller),
		telegotype.WithDiscardLogger(),
	}
	if constructor != nil {
		opts = append(opts, telegotype.WithRequestConstructor(constructor))
	}
	api, err := telegotype.NewBot("123456:"+strings.Repeat("a", 35), opts...)
	if err != nil {
		t.Fatalf("new test bot: %v", err)
	}
	return api
}

func TestSendChatActionMapsTyping(t *testing.T) {
	caller := &recordingCaller{}
	adapter := telegoadapter.NewFromBot(newTestBot(t, caller, nil), nil)

	if err := adapter.SendChatAction(context.Background(), -100123, tg.ChatActionTyping); err != nil {
		t.Fatalf("SendChatAction: %v", err)
	}

	method, body := caller.last(t)
	if method != "sendChatAction" {
		t.Errorf("method = %q, want sendChatAction", method)
	}
	var got struct {
		ChatID int64  `json:"chat_id"`
		Action string `json:"action"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal body %s: %v", body, err)
	}
	if got.ChatID != -100123 || got.Action != "typing" {
		t.Errorf("params = %+v, want chat_id -100123 and action typing", got)
	}
}
