package stage7

import (
	"context"
	"log/slog"
	"net/http"
)

// Middleware wraps an http.Handler with extra behavior. This signature is the
// de-facto standard in the Go ecosystem.
type Middleware func(http.Handler) http.Handler

// Chain wraps h with the given middleware. The FIRST middleware in the list is
// the OUTERMOST one: Chain(h, A, B) runs A, then B, then h.
//
// TODO: Wrap from the last middleware to the first.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	return h
}

// Logging logs one structured line per request using log.Info with the
// message "request" and the attributes method, path, and status (the status
// code the handler wrote; 200 if it never called WriteHeader).
//
// TODO: Wrap http.ResponseWriter in a small struct that records the status in
// WriteHeader, then log after next.ServeHTTP returns.
func Logging(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
}

// Recover turns a panic in the wrapped handler into a 500 response instead of
// crashing the request goroutine.
//
// TODO: defer + recover(). The net/http server already recovers per
// connection, but real services recover explicitly to log and respond properly.
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
}

// RequestID makes sure every request carries an ID. It reuses the incoming
// X-Request-ID header when present, otherwise generates a random one (16 hex
// characters from crypto/rand). The ID is stored in the request context
// (see RequestIDFromContext) and echoed in the X-Request-ID response header.
//
// TODO: context.WithValue with an UNEXPORTED key type, then r.WithContext.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
}

// RequestIDFromContext returns the request ID stored by RequestID, or "".
//
// TODO: ctx.Value(key).(string) with the comma-ok form. Why an unexported key
// type instead of a plain string?
func RequestIDFromContext(ctx context.Context) string {
	return ""
}
