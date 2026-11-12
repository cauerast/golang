package stage7

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Run serves srv on ln until ctx is cancelled, then shuts the server down
// gracefully: no new connections are accepted, in-flight requests are allowed
// to finish, and Run returns once they are done.
//
//   - A clean shutdown returns nil (http.ErrServerClosed is not a failure).
//   - If in-flight requests take longer than shutdownTimeout, return the error
//     from srv.Shutdown (it satisfies errors.Is(err, context.DeadlineExceeded)).
//   - If Serve fails for any other reason, return that error.
//
// In a real program ctx comes from signal.NotifyContext(context.Background(),
// os.Interrupt, syscall.SIGTERM), so Ctrl+C and container stops are graceful.
//
// TODO: Run srv.Serve(ln) in a goroutine, wait on ctx.Done() or a Serve error,
// then call srv.Shutdown with a context created by context.WithTimeout.
func Run(ctx context.Context, srv *http.Server, ln net.Listener, shutdownTimeout time.Duration) error {
	return nil
}
