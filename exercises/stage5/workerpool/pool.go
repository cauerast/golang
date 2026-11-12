// Package workerpool is capstone option A: a bounded concurrent worker pool.
//
// Tip: run with the race detector: go test -race ./exercises/stage5/workerpool
package workerpool

import "errors"

// ErrInvalidWorkers is returned when the worker count is less than 1.
var ErrInvalidWorkers = errors.New("workers must be at least 1")

// Job is a unit of work.
type Job struct {
	ID      int
	Payload int
}

// Result is the outcome of processing one Job.
type Result struct {
	JobID int
	Value int
	Err   error
}

// Process runs fn on every job using exactly `workers` goroutines at most
// (never more than that many fn calls running at the same time).
//
// Requirements:
//   - If workers < 1, return (nil, ErrInvalidWorkers).
//   - The returned slice has one Result per job, in the SAME ORDER as jobs,
//     regardless of which worker finished first.
//   - If fn returns an error for a job, store it in Result.Err and keep going;
//     one failure must not stop the others.
//   - Process returns only after all jobs are done, and leaks no goroutines.
//
// TODO: Start `workers` goroutines reading job indexes from a channel, write
// each Result into its own slot of a pre-sized slice (no mutex needed: why?),
// close the channel after feeding it, and wait with a sync.WaitGroup.
func Process(jobs []Job, workers int, fn func(Job) (int, error)) ([]Result, error) {
	return nil, nil
}
