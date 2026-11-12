package stage7

import (
	"context"
	"sync"
)

// Group runs a collection of goroutines that return errors. It is a small
// re-implementation of golang.org/x/sync/errgroup: the first error cancels
// the shared context, and Wait returns that first error.
type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	cancel  context.CancelFunc
}

// WithContext returns a Group and a derived context that is cancelled the
// first time a function passed to Go returns a non-nil error, or when Wait returns.
//
// TODO: context.WithCancel(ctx); keep the cancel func in the Group.
func WithContext(ctx context.Context) (*Group, context.Context) {
	return &Group{}, ctx
}

// Go runs f in a new goroutine.
//
// TODO: wg.Add(1) before starting the goroutine; on a non-nil error record
// ONLY the first one (errOnce.Do) and call cancel.
func (g *Group) Go(f func() error) {
}

// Wait blocks until all functions started with Go have returned, then returns
// the first non-nil error (or nil).
//
// TODO: wg.Wait(), then cancel (if set) to release the context's resources.
func (g *Group) Wait() error {
	return nil
}
