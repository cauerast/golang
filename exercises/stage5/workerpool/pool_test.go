package workerpool

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestInvalidWorkers(t *testing.T) {
	res, err := Process([]Job{{1, 1}}, 0, func(j Job) (int, error) { return 0, nil })
	if !errors.Is(err, ErrInvalidWorkers) || res != nil {
		t.Errorf("Process(workers=0) = (%v, %v), want (nil, ErrInvalidWorkers)", res, err)
	}
}

func TestProcessOrderAndValues(t *testing.T) {
	jobs := make([]Job, 100)
	for i := range jobs {
		jobs[i] = Job{ID: i, Payload: i}
	}
	res, err := Process(jobs, 4, func(j Job) (int, error) {
		// Later jobs finish sooner, to punish implementations that append
		// results in completion order.
		time.Sleep(time.Duration(100-j.ID) * 50 * time.Microsecond)
		return j.Payload * j.Payload, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != len(jobs) {
		t.Fatalf("got %d results, want %d", len(res), len(jobs))
	}
	for i, r := range res {
		if r.JobID != i || r.Value != i*i || r.Err != nil {
			t.Fatalf("res[%d] = %+v, want {JobID:%d Value:%d Err:nil}", i, r, i, i*i)
		}
	}
}

func TestProcessCollectsErrors(t *testing.T) {
	bad := errors.New("bad job")
	jobs := []Job{{0, 1}, {1, 2}, {2, 3}}
	res, err := Process(jobs, 2, func(j Job) (int, error) {
		if j.ID == 1 {
			return 0, bad
		}
		return j.Payload, nil
	})
	if err != nil {
		t.Fatalf("job errors must not become Process errors: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("got %d results, want 3", len(res))
	}
	if !errors.Is(res[1].Err, bad) {
		t.Errorf("res[1].Err = %v, want %v", res[1].Err, bad)
	}
	if res[0].Err != nil || res[0].Value != 1 || res[2].Err != nil || res[2].Value != 3 {
		t.Errorf("other jobs must still succeed: %+v", res)
	}
}

func TestProcessIsBoundedAndConcurrent(t *testing.T) {
	const workers, njobs = 4, 16
	jobs := make([]Job, njobs)
	for i := range jobs {
		jobs[i] = Job{ID: i}
	}
	var running, peak int32
	_, err := Process(jobs, workers, func(j Job) (int, error) {
		cur := atomic.AddInt32(&running, 1)
		for {
			p := atomic.LoadInt32(&peak)
			if cur <= p || atomic.CompareAndSwapInt32(&peak, p, cur) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&running, -1)
		return 0, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if peak < 2 {
		t.Errorf("peak concurrency = %d; jobs should run concurrently", peak)
	}
	if peak > workers {
		t.Errorf("peak concurrency = %d, exceeds %d workers", peak, workers)
	}
}
