package stage7

import "time"

// Server is configured with the "functional options" pattern, widely used in
// Go libraries instead of overloaded constructors or giant config structs.
type Server struct {
	Addr     string
	Timeout  time.Duration
	MaxConns int
}

// Option mutates a Server during construction.
type Option func(*Server)

// WithTimeout returns an Option that sets Server.Timeout.
//
// TODO: Return a closure that assigns the field.
func WithTimeout(d time.Duration) Option {
	return func(s *Server) {}
}

// WithMaxConns returns an Option that sets Server.MaxConns.
//
// TODO: Return a closure that assigns the field.
func WithMaxConns(n int) Option {
	return func(s *Server) {}
}

// NewServer builds a Server with sensible defaults (Timeout 30s, MaxConns
// 100), then applies the options in order, so later options win.
//
// TODO: Start from the defaults, loop over opts, return the pointer.
// Why is this friendlier to evolve than adding a parameter to NewServer?
func NewServer(addr string, opts ...Option) *Server {
	return &Server{}
}
