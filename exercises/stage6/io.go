// Package stage6 covers the standard-library building blocks you need for
// networked and I/O-heavy programs: io, framing, TCP, HTTP clients, JSON,
// context, URLs, text indexing, and the filesystem.
package stage6

import (
	"errors"
	"io"
)

// MaxFrameSize is the largest payload WriteFrame and ReadFrame accept (1 MiB).
const MaxFrameSize = 1 << 20

// ErrFrameTooLarge is returned when a frame exceeds MaxFrameSize.
var ErrFrameTooLarge = errors.New("frame too large")

// CountBytes reads r to the end and returns how many bytes it produced.
//
// TODO: Use io.Copy(io.Discard, r). It streams in small chunks, so it works
// for inputs far bigger than memory. Why is io.ReadAll the wrong tool here?
func CountBytes(r io.Reader) (int64, error) {
	return -1, nil
}

// Rot13Reader wraps another reader and applies the ROT13 cipher to every
// letter it reads ("Hello" -> "Uryyb"). Non-letters pass through unchanged.
//
// This is the decorator pattern, the heart of the io package: gzip, TLS,
// bufio, and hashing all wrap a Reader in a Reader.
type Rot13Reader struct {
	R io.Reader
}

// Read implements io.Reader.
//
// TODO: Call r.R.Read(p), then transform p[:n] IN PLACE, and return n and the
// error exactly as the inner reader reported them. Remember: Read may return
// n > 0 AND a non-nil error (including io.EOF) in the same call.
func (r Rot13Reader) Read(p []byte) (int, error) {
	return 0, io.EOF
}

// WriteFrame writes payload to w using a length-prefixed frame:
// a 4-byte big-endian unsigned length, followed by the payload bytes.
// Payloads larger than MaxFrameSize return ErrFrameTooLarge and write nothing.
//
// TCP is a byte STREAM with no message boundaries; framing is how you send
// discrete messages over it.
//
// TODO: Use encoding/binary (binary.BigEndian.PutUint32) and build the header
// and payload so that you issue as few Write calls as you reasonably can.
func WriteFrame(w io.Writer, payload []byte) error {
	return nil
}

// ReadFrame reads one frame written by WriteFrame and returns its payload.
//
//   - Clean end of stream before any header byte -> io.EOF.
//   - Stream ends in the middle of a header or payload -> io.ErrUnexpectedEOF.
//   - Declared length > MaxFrameSize -> ErrFrameTooLarge, BEFORE allocating.
//
// TODO: io.ReadFull(r, header) returns io.EOF only if zero bytes were read and
// io.ErrUnexpectedEOF for a partial read, which is exactly what you need.
// Never trust a length received from the network when allocating memory.
func ReadFrame(r io.Reader) ([]byte, error) {
	return nil, nil
}

// SHA256Hex returns the lowercase hex SHA-256 digest of everything in r.
//
// TODO: crypto/sha256's New() returns a hash.Hash, which is an io.Writer.
// Stream r into it with io.Copy, then hex-encode h.Sum(nil).
func SHA256Hex(r io.Reader) (string, error) {
	return "", nil
}
