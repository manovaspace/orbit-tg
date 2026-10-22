package message_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manovaspace/orbit-tg/message"
)

func TestHTMLEscape(t *testing.T) {
	input := `<script>alert("hello & 'goodbye'")</script>`
	want := `&lt;script&gt;alert(&quot;hello &amp; &#39;goodbye&#39;&quot;)&lt;/script&gt;`
	if got := message.HTMLEscape(input); got != want {
		t.Errorf("HTMLEscape() = %q, want %q", got, want)
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "bare angle brackets and ampersands",
			input: "if 5 < 10 && 10 > 5 then ok",
			want:  "if 5 &lt; 10 &amp;&amp; 10 &gt; 5 then ok",
		},
		{
			name:  "preserves existing valid entity",
			input: "Copyright &copy; 2026 &amp; &lt;tag&gt;",
			want:  "Copyright &copy; 2026 &amp; &lt;tag&gt;",
		},
		{
			name:  "strips unsupported tags while preserving inner text",
			input: `<div class="container"><span style="color:red">Alert: <b>Critical</b></span></div>`,
			want:  "Alert: <b>Critical</b>",
		},
		{
			name:  "auto-closes unclosed tags",
			input: "<b>Bold <i>Italic text",
			want:  "<b>Bold <i>Italic text</i></b>",
		},
		{
			name:  "valid link with href preserved",
			input: `<a href="https://manova.space">Manova</a>`,
			want:  `<a href="https://manova.space">Manova</a>`,
		},
		{
			name:  "link without href stripped",
			input: `<a>No link</a>`,
			want:  "No link",
		},
		{
			name:  "code block with language class",
			input: `<pre><code class="language-go">fmt.Println("hi")</code></pre>`,
			want:  `<pre><code class="language-go">fmt.Println("hi")</code></pre>`,
		},
		{
			name:  "expandable blockquote",
			input: `<blockquote expandable>Long quote</blockquote>`,
			want:  `<blockquote expandable>Long quote</blockquote>`,
		},
		{
			name:  "strips comments",
			input: "Hello <!-- secret -->World",
			want:  "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := message.SanitizeHTML(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChunkHTML(t *testing.T) {
	t.Run("short string fits in single chunk", func(t *testing.T) {
		short := "<b>Hello</b> world"
		chunks := message.ChunkHTML(short, 100)
		if len(chunks) != 1 || chunks[0] != short {
			t.Errorf("expected 1 chunk %q, got %v", short, chunks)
		}
	})

	t.Run("splits on paragraph boundary and balances tags", func(t *testing.T) {
		p1 := "Paragraph 1: " + strings.Repeat("a", 40)
		p2 := "Paragraph 2: " + strings.Repeat("b", 40)
		full := "<b>" + p1 + "\n\n" + p2 + "</b>"

		// Max length 70 will force split between p1 and p2
		chunks := message.ChunkHTML(full, 75)
		if len(chunks) < 2 {
			t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
		}

		// First chunk must close </b>
		if !strings.HasPrefix(chunks[0], "<b>") || !strings.HasSuffix(chunks[0], "</b>") {
			t.Errorf("chunk 0 not balanced: %q", chunks[0])
		}

		// Second chunk must reopen <b> and close </b>
		if !strings.HasPrefix(chunks[1], "<b>") || !strings.HasSuffix(chunks[1], "</b>") {
			t.Errorf("chunk 1 not balanced: %q", chunks[1])
		}
	})

	t.Run("never splits inside a tag", func(t *testing.T) {
		full := strings.Repeat("x", 40) + `<a href="https://example.com/very/long/url">click here</a>`
		chunks := message.ChunkHTML(full, 50)
		for i, c := range chunks {
			if strings.Count(c, "<") != strings.Count(c, ">") {
				t.Errorf("chunk %d has unclosed tag bracket: %q", i, c)
			}
		}
	})
}

func TestEphemeralManager(t *testing.T) {
	var deletedChatID int64
	var deletedMsgID int
	var mu sync.Mutex
	done := make(chan struct{})

	mgr := message.NewEphemeralManager(context.Background(), func(ctx context.Context, chatID int64, messageID int) error {
		mu.Lock()
		deletedChatID = chatID
		deletedMsgID = messageID
		mu.Unlock()
		close(done)
		return nil
	})
	defer mgr.Close()

	mgr.Schedule(12345, 6789, 50*time.Millisecond)

	select {
	case <-done:
		mu.Lock()
		defer mu.Unlock()
		if deletedChatID != 12345 || deletedMsgID != 6789 {
			t.Errorf("unexpected deletion target: chat=%d, msg=%d", deletedChatID, deletedMsgID)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for ephemeral deletion")
	}
}

func TestEphemeralManager_CloseCancelsPending(t *testing.T) {
	var calls int32
	mgr := message.NewEphemeralManager(context.Background(), func(ctx context.Context, chatID int64, messageID int) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})

	// Schedule with long TTL and immediately close
	mgr.Schedule(12345, 6789, 2*time.Second)
	mgr.Close()

	time.Sleep(100 * time.Millisecond)
	if count := atomic.LoadInt32(&calls); count != 0 {
		t.Errorf("expected 0 calls after Close(), got %d", count)
	}
}

func TestScrub(t *testing.T) {
	t.Run("successful on first attempt", func(t *testing.T) {
		var attempts int
		err := message.Scrub(context.Background(), 100, 200, func(ctx context.Context, chatID int64, messageID int) error {
			attempts++
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("retries on error until success", func(t *testing.T) {
		var attempts int
		err := message.Scrub(context.Background(), 100, 200, func(ctx context.Context, chatID int64, messageID int) error {
			attempts++
			if attempts < 3 {
				return errors.New("rate limited")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("aborts on context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var attempts int
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := message.Scrub(ctx, 100, 200, func(c context.Context, chatID int64, messageID int) error {
			attempts++
			return errors.New("always fail")
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}
