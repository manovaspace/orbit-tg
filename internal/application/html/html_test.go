package html_test

import (
	"strings"
	"testing"

	"github.com/manovaspace/orbit-tg/internal/application/html"
)

func TestHTMLEscape(t *testing.T) {
	if got := html.HTMLEscape(`<b>&"'`); got != "&lt;b&gt;&amp;&quot;&#39;" {
		t.Fatalf("HTMLEscape = %q", got)
	}
}

func TestSanitizeHTMLKeepsTelegramTagsAndDropsUnknown(t *testing.T) {
	got := html.SanitizeHTML(`<div><b>ok</b> <script>x</script></div>`)
	if strings.Contains(got, "<div") || strings.Contains(got, "<script") {
		t.Fatalf("unsupported tags survived: %q", got)
	}
	if !strings.Contains(got, "<b>ok</b>") {
		t.Fatalf("allowed tag missing: %q", got)
	}
}

func TestChunkHTMLBalancesTags(t *testing.T) {
	long := "<b>" + strings.Repeat("word ", 20) + "</b>"
	chunks := html.ChunkHTML(long, 40)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for _, chunk := range chunks {
		if strings.Count(chunk, "<b>") != strings.Count(chunk, "</b>") {
			t.Fatalf("unbalanced chunk: %q", chunk)
		}
	}
}
