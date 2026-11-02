package bot

import (
	"context"

	tg "github.com/manovaspace/orbit-tg"
	"github.com/manovaspace/orbit-tg/internal/application/router"
)

// RoutePredicate evaluates whether an Update matches a route.
type RoutePredicate func(u *tg.Update) bool

// Route pairs a RoutePredicate with a HandlerFunc.
type Route struct {
	Predicate RoutePredicate
	Handler   tg.HandlerFunc
}

// Router dispatches incoming Telegram updates to matching handlers.
type Router struct {
	inner *router.Router
}

// NewRouter initializes an empty Router.
func NewRouter() *Router {
	return &Router{inner: router.New()}
}

// Handle registers a route with a custom predicate.
func (r *Router) Handle(predicate RoutePredicate, handler tg.HandlerFunc) {
	r.inner.Handle(router.Predicate(predicate), handler)
}

// Command registers a handler for a bot command (e.g. "/start", "help").
// Matches commands on the first line, handling mobile newline inputs and bot mentions (/cmd@bot).
func (r *Router) Command(name string, handler tg.HandlerFunc) {
	r.inner.Command(name, handler)
}

// CallbackPrefix registers a handler for callback queries whose data begins with prefix.
func (r *Router) CallbackPrefix(prefix string, handler tg.HandlerFunc) {
	r.inner.CallbackPrefix(prefix, handler)
}

// Dispatch executes the first route whose predicate evaluates to true.
func (r *Router) Dispatch(ctx context.Context, u *tg.Update) (bool, error) {
	return r.inner.Dispatch(ctx, u)
}
