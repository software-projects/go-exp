package httpctx

/*

import (
	"golang.org/x/net/context"
	"net/http"
)

// A Stack is a list of middleware functions.
type Stack struct {
	middleware func(Handler) Handler
	prev       *Stack
}

func (s *Stack) Use(f ...func(h Handler) Handler) *Stack {
	stack := s

	for _, m := range f {
		if m != nil {
			stack = &Stack{
				middleware: m,
				prev:       stack,
			}
		}
	}

	return stack
}

func (s *Stack) Then(h Handler) http.Handler {
	for stack := s; stack != nil; stack = stack.prev {
		h = stack.middleware(h)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelFunc := newContext(w, r)
		defer cancelFunc()
		err := h.ServeHTTPContext(ctx, w, r)
		if err != nil {
			sendError(w, r, err)
		}
	})
}

func Use(f ...func(h Handler) Handler) *Stack {
	var stack *Stack

	for _, m := range f {
		if m != nil {
			stack = &Stack{
				middleware: m,
				prev:       stack,
			}
		}
	}

	return stack
}
*/