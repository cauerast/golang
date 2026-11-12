package stage4

import (
	"sort"
	"sync"
	"testing"
	"time"
)

const testTimeout = 2 * time.Second

// collect drains ch, failing the test instead of hanging when ch is nil,
// never closed, or never produces values.
func collect(t *testing.T, ch <-chan int) []int {
	t.Helper()
	var got []int
	timeout := time.After(testTimeout)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-timeout:
			t.Fatalf("timed out waiting for channel to close (got so far: %v)", got)
			return nil
		}
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSumConcurrent(t *testing.T) {
	nums := make([]int, 1000)
	want := 0
	for i := range nums {
		nums[i] = i + 1
		want += nums[i]
	}
	for _, workers := range []int{-3, 0, 1, 4, 7, 2000} {
		if got := SumConcurrent(nums, workers); got != want {
			t.Errorf("SumConcurrent(workers=%d) = %d, want %d", workers, got, want)
		}
	}
	if got := SumConcurrent(nil, 4); got != 0 {
		t.Errorf("SumConcurrent(nil) = %d, want 0", got)
	}
}

func TestGenerator(t *testing.T) {
	got := collect(t, Generator(5))
	if !equalInts(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Generator(5) = %v, want [1 2 3 4 5]", got)
	}
	if got := collect(t, Generator(0)); len(got) != 0 {
		t.Errorf("Generator(0) = %v, want empty", got)
	}
}

func TestSquare(t *testing.T) {
	got := collect(t, Square(Generator(4)))
	if !equalInts(got, []int{1, 4, 9, 16}) {
		t.Errorf("Square(Generator(4)) = %v, want [1 4 9 16]", got)
	}
}

func TestMerge(t *testing.T) {
	got := collect(t, Merge(Generator(3), Generator(2), Generator(4)))
	sort.Ints(got)
	want := []int{1, 1, 1, 2, 2, 2, 3, 3, 4}
	if !equalInts(got, want) {
		t.Errorf("Merge = %v, want %v (any order)", got, want)
	}
}

func TestSafeCounter(t *testing.T) {
	var c SafeCounter
	var wg sync.WaitGroup
	const goroutines, perG = 100, 100
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
	if got := c.Value(); got != goroutines*perG {
		t.Errorf("Value() = %d, want %d", got, goroutines*perG)
	}
}

func TestWithTimeout(t *testing.T) {
	fast := func() string { return "done" }
	if got, ok := WithTimeout(fast, time.Second); !ok || got != "done" {
		t.Errorf("fast call = (%q, %v), want (\"done\", true)", got, ok)
	}

	slow := func() string {
		time.Sleep(300 * time.Millisecond)
		return "late"
	}
	start := time.Now()
	got, ok := WithTimeout(slow, 20*time.Millisecond)
	if ok || got != "" {
		t.Errorf("slow call = (%q, %v), want (\"\", false)", got, ok)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("WithTimeout waited %v; it must return when the timeout fires", elapsed)
	}
}
