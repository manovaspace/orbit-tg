package rich

import (
	"testing"

	"github.com/mymmrac/telego"
)

func TestNoticeHierarchy(t *testing.T) {
	msg := Notice("Checklist updated",
		Section{Title: "Added", Items: []string{"tent"}},
		Section{Title: "Already on the list", Items: []string{"water"}},
	)
	if len(msg.Blocks) < 5 {
		t.Fatalf("blocks: got %d want >= 5", len(msg.Blocks))
	}
}

func TestHTMLEscape(t *testing.T) {
	if got := HTMLEscape("<b>"); got != "&lt;b&gt;" {
		t.Fatalf("escape: %q", got)
	}
}

func TestNumberedList(t *testing.T) {
	block := NumberedList([]string{"tent", "rope"})
	list, ok := block.(*telego.InputRichBlockList)
	if !ok || len(list.Items) != 2 {
		t.Fatalf("expected list block, got %T", block)
	}
	if list.Items[0].Type != telego.OrderedListDecimal || list.Items[0].Value != 1 {
		t.Fatalf("first item: type=%q value=%d", list.Items[0].Type, list.Items[0].Value)
	}
	if list.Items[1].Value != 2 {
		t.Fatalf("second value: %d", list.Items[1].Value)
	}
}

func TestBuilderListWithDividers(t *testing.T) {
	msg := NewBuilder().
		Heading("Title", 2).
		ParagraphPlain("one").
		Divider().
		ParagraphPlain("two").
		Build()
	if len(msg.Blocks) != 4 {
		t.Fatalf("blocks: got %d want 4", len(msg.Blocks))
	}
}
