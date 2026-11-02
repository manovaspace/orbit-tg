// Package router dispatches abstract Telegram updates to matching handlers. It
// imports no Telegram client or HTTP package.
package router

import (
	"context"
	"regexp"
	"strings"

	tg "github.com/manovaspace/orbit-tg"
)

// Predicate evaluates whether an Update matches a route.
type Predicate func(u *tg.Update) bool

// Route pairs a Predicate with a HandlerFunc.
type Route struct {
	Predicate Predicate
	Handler   tg.HandlerFunc
}

// Router dispatches incoming Telegram updates to matching handlers.
type Router struct {
	routes []Route
}

// New initializes an empty Router.
func New() *Router {
	return &Router{routes: make([]Route, 0)}
}

// Handle registers a route with a custom predicate.
func (r *Router) Handle(predicate Predicate, handler tg.HandlerFunc) {
	r.routes = append(r.routes, Route{Predicate: predicate, Handler: handler})
}

// Command registers a handler for a bot command (e.g. "/start", "help").
// Matches commands on the first line, handling mobile newline inputs and bot mentions (/cmd@bot).
func (r *Router) Command(name string, handler tg.HandlerFunc) {
	cleanName := strings.TrimPrefix(name, "/")
	re := regexp.MustCompile(`(?i)^/` + regexp.QuoteMeta(cleanName) + `(?:@\w+)?(?:\s|$)`)
	r.Handle(func(u *tg.Update) bool {
		if u == nil || u.Message == nil || u.Message.Text == "" {
			return false
		}
		first, _, _ := strings.Cut(u.Message.Text, "\n")
		return re.MatchString(strings.TrimSpace(first))
	}, handler)
}

// CallbackPrefix registers a handler for callback queries whose data begins with prefix.
func (r *Router) CallbackPrefix(prefix string, handler tg.HandlerFunc) {
	r.Handle(func(u *tg.Update) bool {
		if u == nil || u.CallbackQuery == nil || u.CallbackQuery.Data == "" {
			return false
		}
		return strings.HasPrefix(u.CallbackQuery.Data, prefix)
	}, handler)
}

// Dispatch executes the first route whose predicate evaluates to true.
func (r *Router) Dispatch(ctx context.Context, u *tg.Update) (bool, error) {
	if u == nil {
		return false, nil
	}
	for _, route := range r.routes {
		if route.Predicate != nil && route.Predicate(u) {
			if route.Handler != nil {
				return true, route.Handler(ctx, u)
			}
			return true, nil
		}
	}
	return false, nil
}
