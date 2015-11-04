// Package log provides diagnostic logging.
//
// There are many good logging packages available, and it is worth asking
// why the world needs another one.
//
// Here are some differentiators for this package. Not all of them are
// unique, but this is the only package (to date) that has all of them.
//
// 1. Log messages are not formatted using a printf style interface. Each
// log message should have a constant message, which makes it easier to
// filter and search for messages. Any variable information is passed as
// properties in the message (see the WithValue function).
//
// 2. Uses an api that allows for multiple options and parameters to be
// logged in a single call. (See "Functional options for friendly APIs"
// by Dave Cheney http://goo.gl/l2KaW3).
//
// 3. When a message is logged, a non-nil *Message value is returned, which
// can be returned as an error value.
//
// 4. This package is context aware (golang.org/x/net/context). Contexts
// can be created with information that will be logged with the message.
//
// 5. Messages can be logged as text messages, or structured (JSON) messages.
package log

import (
	"os"

	"golang.org/x/net/context"
)

var handlers []Handler

func AddHandler(h Handler) {
	handlers = append(handlers, h)
}

type Handler interface {
	Handle(m *Message) error
}

// Print prints the log Text Message to standard output.
var Print func(m *Message) = func(m *Message) {
	m.WriteTo(os.Stdout)
}

func doPrint(m *Message) {
	if Print != nil && m.Level >= MinLevel {
		Print(m)
		for _, h := range handlers {
			// TODO: write to stderr if cannot send message
			_ = h.Handle(m)
		}
	}
}

// Debug logs a debug level message.
func Debug(text string, opts ...Option) *Message {
	m := newMessage(LevelDebug, text)
	m.applyOpts(opts)
	doPrint(m)
	return m
}

// DebugC logs an info level message with a context. Calling DebugC
// is identical to calling Debug with a WithContext option.
func DebugC(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(LevelDebug, text)
	m.applyOpts(opts)
	WithContext(ctx)(m)
	doPrint(m)
	return m
}

// Info logs an info level message.
func Info(text string, opts ...Option) *Message {
	m := newMessage(LevelInfo, text)
	m.applyOpts(opts)
	doPrint(m)
	return m
}

// InfoC logs an info level message with a context. Calling InfoC
// is identical to calling Info with a WithContext option.
func InfoC(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(LevelInfo, text)
	m.applyOpts(opts)
	WithContext(ctx)(m)
	doPrint(m)
	return m
}

// Warn logs a warning level message.
func Warn(text string, opts ...Option) *Message {
	m := newMessage(LevelWarning, text)
	m.applyOpts(opts)
	doPrint(m)
	return m
}

// WarnC logs an info level message with a context. Calling WarnC
// is identical to calling Warn with a WithContext option.
func WarnC(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(LevelWarning, text)
	m.applyOpts(opts)
	WithContext(ctx)(m)
	doPrint(m)
	return m
}

// Error logs an error level message.
func Error(text string, opts ...Option) *Message {
	m := newMessage(LevelError, text)
	m.applyOpts(opts)
	doPrint(m)
	return m
}

// ErrorC logs an info level message with a context. Calling ErrorC
// is identical to calling Error with a WithContext option.
func ErrorC(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(LevelError, text)
	m.applyOpts(opts)
	WithContext(ctx)(m)
	doPrint(m)
	return m
}

// ErrorCE handles the common case where an error is logged with a context and
// an error. A call to
//
//  log.ErrorCE(ctx, err, "some text")
//
// is identical to calling
//
//  log.ErrorCE(ctx, err, "some text",
//          log.WithContext(ctx),
//          log.WithError(err))
//
// and a call to
//
//  log.ErrorCE(ctx, err, "some text",
//          log.WithValue("key1", "val1"),
//          log.WithValue("key2", "val2"))
//
// is identical to calling
//
//  log.ErrorCE(ctx, err, "some text",
//          log.WithValue("key1", "val1"),
//          log.WithValue("key2", "val2"),
//          log.WithContext(ctx),
//          log.WithError(err))
//
// This function was introduced to handle a common usage pattern succinctly.
func ErrorCE(ctx context.Context, err error, text string, opts ...Option) *Message {
	m := newMessage(LevelError, text)
	m.applyOpts(opts)
	WithContext(ctx)(m)
	WithError(err)(m)
	doPrint(m)
	return m
}
