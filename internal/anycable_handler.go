package internal

import (
	"net/http"

	"github.com/anycable/anycable-go/cli"
)

type AnyCableHandler struct {
	handler http.Handler
	next    http.Handler
}

func NewAnyCableHandler(anycable *cli.Embedded, next http.Handler) *AnyCableHandler {
	handler, _ := anycable.WebSocketHandler()

	return &AnyCableHandler{
		handler: handler,
		next:    next,
	}
}

func (h *AnyCableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/cable" {
		h.handler.ServeHTTP(w, r)
	} else {
		h.next.ServeHTTP(w, r)
	}
}
