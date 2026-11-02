package handler_test

import (
	"testing"

	apphandler "github.com/manovaspace/orbit-tg/internal/application/handler"
)

func TestFirstLineCommand(t *testing.T) {
	match := apphandler.FirstLineCommand("checklist")
	cases := map[string]bool{
		"/checklist":       true,
		"/checklist\narg":  true,
		"/checklist@bot x": true,
		"  /CHECKLIST  ":   true,
		"/checklistx":      false,
		"text /checklist":  false,
		"":                 false,
		"/other":           false,
		"/checklist@bot":   true,
	}
	for text, want := range cases {
		if got := match(text); got != want {
			t.Fatalf("%q: match = %v, want %v", text, got, want)
		}
	}
}
