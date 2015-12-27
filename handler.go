package mohttp

import (
	"golang.org/x/net/context"
)

type Handler interface {
	Handle(context.Context)
}

type HandlerFunc func(context.Context)

func (h HandlerFunc) Handle(c context.Context) {
	h(c)
}

type PriorityHandler interface {
	Priority() int
	Handler
}

type priorityHandler struct {
	priority int
	handler  Handler
}

func (p *priorityHandler) Priority() int            { return p.priority }
func (p *priorityHandler) Handle(c context.Context) { p.handler.Handle(c) }

func PriorityHandlerFunc(p int, fn HandlerFunc) PriorityHandler {
	return &priorityHandler{p, fn}
}

var EmptyBodyHandler Handler = HandlerFunc(func(c context.Context) {})

type ChainedHandlers []Handler

func (ch ChainedHandlers) Handle(c context.Context) {
	Serve(c, ch...)
}
