package stage6

import (
	"net"
	"time"
)

// EchoServer accepts connections from l and echoes back every byte it receives
// on each connection, handling every connection in its own goroutine.
// It returns nil once l has been closed (that is the normal way to stop it)
// and any other Accept error otherwise.
//
// TODO: Loop on l.Accept(). Detect shutdown with errors.Is(err, net.ErrClosed).
// For each conn: go func() { defer conn.Close(); io.Copy(conn, conn) }().
func EchoServer(l net.Listener) error {
	return nil
}

// SendAndReceive dials addr (TCP), sends msg, and reads exactly len(msg) bytes
// back. The whole exchange (dial, write, read) must finish within timeout,
// otherwise an error is returned (a net.Error whose Timeout() is true when the
// peer is just silent).
//
// TODO: net.DialTimeout, then conn.SetDeadline(time.Now().Add(timeout)),
// conn.Write, and io.ReadFull into a buffer of len(msg). Always close conn.
func SendAndReceive(addr string, msg []byte, timeout time.Duration) ([]byte, error) {
	return nil, nil
}
