// Package stage4 covers goroutines, channels, select, sync, and races.
//
// Tip: run these tests with the race detector: go test -race ./exercises/stage4
package stage4

import (
	"sync"
	"time"
)

// SumConcurrent sums nums by splitting it into `workers` chunks, summing each
// chunk in its own goroutine, and combining the partial sums.
// workers < 1 must be treated as 1.
//
// TODO: Use sync.WaitGroup and a results channel (or one slot per worker in a
// pre-sized slice, so goroutines never write to the same memory).
func SumConcurrent(nums []int, workers int) int {
	return -1
}

// Generator returns a channel that emits 1, 2, ..., n and is then closed.
//
// TODO: Start a goroutine that sends and then closes the channel.
// Who should close a channel: the sender or the receiver?
func Generator(n int) <-chan int {
	return nil
}

// Square reads ints from in and emits their squares on the returned channel.
// The output channel is closed once in is closed and drained.
//
// TODO: Classic pipeline stage: for v := range in { out <- v * v }.
func Square(in <-chan int) <-chan int {
	return nil
}

// Merge fans several input channels into a single output channel. The output
// is closed after ALL inputs are closed. Order is not defined.
//
// TODO: One goroutine per input, a sync.WaitGroup, and one more goroutine that
// waits and then closes the output.
func Merge(inputs ...<-chan int) <-chan int {
	return nil
}

// SafeCounter is a counter that is safe for concurrent use.
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

// Inc increments the counter by one.
//
// TODO: Guard access with c.mu (Lock + defer Unlock).
func (c *SafeCounter) Inc() {
}

// Value returns the current count.
//
// TODO: Guard the read as well; a race can be a read racing a write.
func (c *SafeCounter) Value() int {
	return -1
}

// WithTimeout runs f in a goroutine and waits for its result for at most d.
// It returns (result, true) if f finished in time, or ("", false) on timeout.
//
// TODO: Use select with a result channel and time.After. Make the result
// channel buffered (size 1) so the goroutine can finish even after a timeout
// instead of leaking forever. Why does that matter?
func WithTimeout(f func() string, d time.Duration) (string, bool) {
	return "unimplemented", true
}
