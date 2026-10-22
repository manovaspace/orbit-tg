package rich_test

import (
	"strings"
	"testing"

	"github.com/manovaspace/orbit-tg/rich"
)

func TestNotice(t *testing.T) {
	out := rich.Notice("Deployment Status",
		rich.Section{Title: "Services", Items: []string{"auth <v1>", "gateway"}},
		rich.Section{Title: "Notes", Items: []string{"all healthy"}},
	)

	if !strings.HasPrefix(out, "<b>Deployment Status</b>") {
		t.Errorf("expected heading bold prefix, got %q", out)
	}
	if !strings.Contains(out, "• auth &lt;v1&gt;") {
		t.Errorf("expected escaped item in bullet list, got %q", out)
	}
	if !strings.Contains(out, "<b>Services</b>\n• auth &lt;v1&gt;\n• gateway") {
		t.Errorf("section structure mismatch: %q", out)
	}
}

func TestNumberedList(t *testing.T) {
	out := rich.NumberedList([]string{"first", "second <item>", "third"})
	want := "1. first\n2. second &lt;item&gt;\n3. third"
	if out != want {
		t.Errorf("NumberedList() = %q, want %q", out, want)
	}
}

func TestCodeBlock(t *testing.T) {
	t.Run("with language", func(t *testing.T) {
		out := rich.CodeBlock("go", `fmt.Println("hello <world>")`)
		want := `<pre><code class="language-go">fmt.Println(&quot;hello &lt;world&gt;&quot;)</code></pre>`
		if out != want {
			t.Errorf("CodeBlock() = %q, want %q", out, want)
		}
	})

	t.Run("without language", func(t *testing.T) {
		out := rich.CodeBlock("", "plain code")
		want := "<pre>plain code</pre>"
		if out != want {
			t.Errorf("CodeBlock() = %q, want %q", out, want)
		}
	})
}

func TestBlockquote(t *testing.T) {
	t.Run("standard blockquote", func(t *testing.T) {
		out := rich.Blockquote("important note", false)
		want := "<blockquote>important note</blockquote>"
		if out != want {
			t.Errorf("Blockquote() = %q, want %q", out, want)
		}
	})

	t.Run("expandable blockquote", func(t *testing.T) {
		out := rich.Blockquote("long note", true)
		want := "<blockquote expandable>long note</blockquote>"
		if out != want {
			t.Errorf("Blockquote() = %q, want %q", out, want)
		}
	})
}

func TestInlineHelpers(t *testing.T) {
	if got := rich.Bold("bold & strong"); got != "<b>bold &amp; strong</b>" {
		t.Errorf("Bold() = %q", got)
	}
	if got := rich.Italic("slanted"); got != "<i>slanted</i>" {
		t.Errorf("Italic() = %q", got)
	}
	if got := rich.Code("x := 10"); got != "<code>x := 10</code>" {
		t.Errorf("Code() = %q", got)
	}
	if got := rich.Link("home", "https://manova.space?a=1&b=2"); got != `<a href="https://manova.space?a=1&amp;b=2">home</a>` {
		t.Errorf("Link() = %q", got)
	}
	if got := rich.Spoiler("secret"); got != "<tg-spoiler>secret</tg-spoiler>" {
		t.Errorf("Spoiler() = %q", got)
	}
}

func TestBuilder(t *testing.T) {
	b := rich.NewBuilder().
		Heading("Server Status").
		Paragraph("Everything running smoothly.").
		BulletList("API: 200 OK", "DB: Connected").
		Divider().
		CodeBlock("bash", "orbit status").
		Blockquote("Audit log recorded", false)

	result := b.Build()
	if !strings.Contains(result, "<b>Server Status</b>") {
		t.Errorf("missing heading: %q", result)
	}
	if !strings.Contains(result, "———") {
		t.Errorf("missing divider: %q", result)
	}
	if !strings.Contains(result, "<pre><code class=\"language-bash\">orbit status</code></pre>") {
		t.Errorf("missing code block: %q", result)
	}
	if !strings.Contains(result, "<blockquote>Audit log recorded</blockquote>") {
		t.Errorf("missing blockquote: %q", result)
	}
}
