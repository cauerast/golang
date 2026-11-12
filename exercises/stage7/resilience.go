package stage7

import (
	"context"
	"time"
)

// Retry calls fn until it succeeds or attempts is exhausted.
//
//   - attempts < 1 is treated as 1.
//   - Between failures wait base, then 2*base, 4*base, ... (exponential backoff).
//     Do NOT wait after the last failed attempt.
//   - While waiting, stop immediately if ctx is done and return an error that
//     satisfies errors.Is(err, ctx.Err()).
//   - After the final failure return an error that wraps the last error from fn
//     and mentions the number of attempts, e.g. "after 3 attempts: boom".
//
// TODO: Wait with select on time.After(delay) and ctx.Done(). Production code
// usually adds random jitter to the delay: why?
func Retry(ctx context.Context, attempts int, base time.Duration, fn func() error) error {
	return nil
}

// TokenBucket is a rate limiter: it holds up to `capacity` tokens, refills one
// token every `refillEvery`, and each allowed event consumes one token.
// The clock is injected so tests (and you) can control time.
type TokenBucket struct {
	// TODO: add fields: a mutex, capacity, refill interval, current tokens,
	// the time of the last refill, and the now func.
}

// NewTokenBucket returns a bucket that starts FULL.
//
// TODO: Store the arguments and initialize the fields.
func NewTokenBucket(capacity int, refillEvery time.Duration, now func() time.Time) *TokenBucket {
	return &TokenBucket{}
}

// Allow reports whether one event may proceed right now, consuming a token if so.
//
// Refill rule: tokens += elapsed / refillEvery (integer division), capped at
// capacity. Advance the "last refill" time by whole intervals only, so partial
// progress is not lost (500ms + 500ms must yield one token for a 1s interval).
// When the bucket is full, simply reset the last refill time to now.
//
// TODO: Safe for concurrent use: lock the mutex.
func (b *TokenBucket) Allow() bool {
	return false
}
