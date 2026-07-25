package handler

import (
	"context"
	"testing"

	"github.com/mymmrac/telego"
)

func TestFirstLineCommand(t *testing.T) {
	p := FirstLineCommand("checklist")
	cases := []struct {
		text string
		want bool
	}{
		{"/checklist", true},
		{"/checklist\n", true},
		{"/checklist Camping", true},
		{"/checklist Camping\ntent\nrope", true},
		{"/checklist@MyBot Camping", true},
		{"/list", false},
		{"+checklist", false},
		{"", false},
	}
	for _, tc := range cases {
		ok := p(context.Background(), telego.Update{Message: &telego.Message{Text: tc.text}})
		if ok != tc.want {
			t.Fatalf("%q: got %v want %v", tc.text, ok, tc.want)
		}
	}
}
