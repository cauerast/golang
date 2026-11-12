package stage7

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a manually advanced clock, safe for concurrent use.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{t: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

// ---------- generics ----------

type myInt int

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, strconv.Itoa)
	if !reflect.DeepEqual(got, []string{"1", "2", "3"}) {
		t.Errorf("Map = %v", got)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]string{"a", "bb", "ccc", "dd"}, func(s string) bool { return len(s) == 2 })
	if !reflect.DeepEqual(got, []string{"bb", "dd"}) {
		t.Errorf("Filter = %v", got)
	}
	if none := Filter([]int{1}, func(int) bool { return false }); none == nil {
		t.Error("Filter with no matches must return an empty non-nil slice")
	}
}

func TestReduce(t *testing.T) {
	got := Reduce([]string{"a", "bb", "ccc"}, 0, func(acc int, s string) int { return acc + len(s) })
	if got != 6 {
		t.Errorf("Reduce = %d, want 6", got)
	}
}

func TestSumOf(t *testing.T) {
	if got := SumOf([]int{1, 2, 3}); got != 6 {
		t.Errorf("SumOf ints = %d, want 6", got)
	}
	if got := SumOf([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("SumOf floats = %v, want 0.75", got)
	}
	if got := SumOf([]myInt{4, 5}); got != 9 {
		t.Errorf("SumOf named type = %d, want 9", got)
	}
}

func TestSet(t *testing.T) {
	s := NewSet("a", "b", "a")
	if s.Len() != 2 {
		t.Errorf("Len = %d, want 2", s.Len())
	}
	s.Add("c")
	s.Add("c")
	if s.Len() != 3 || !s.Has("c") || s.Has("z") {
		t.Errorf("Set contents wrong: len=%d has(c)=%v has(z)=%v", s.Len(), s.Has("c"), s.Has("z"))
	}
}

func BenchmarkFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	even := func(n int) bool { return n%2 == 0 }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Filter(data, even)
	}
}

// ---------- functional options ----------

func TestNewServerDefaults(t *testing.T) {
	s := NewServer(":8080")
	if s.Addr != ":8080" || s.Timeout != 30*time.Second || s.MaxConns != 100 {
		t.Errorf("defaults = %+v, want {:8080 30s 100}", *s)
	}
}

func TestNewServerOptions(t *testing.T) {
	s := NewServer(":9000", WithTimeout(5*time.Second), WithMaxConns(10), WithMaxConns(20))
	if s.Timeout != 5*time.Second || s.MaxConns != 20 {
		t.Errorf("options not applied (later wins): %+v", *s)
	}
}

// ---------- middleware ----------

func TestChainOrder(t *testing.T) {
	var order []string
	mk := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	h := Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}), mk("A"), mk("B"))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if !reflect.DeepEqual(order, []string{"A", "B", "handler"}) {
		t.Errorf("execution order = %v, want [A B handler]", order)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, nil))
	h := Logging(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/hello" {
			w.WriteHeader(http.StatusCreated)
		}
		io.WriteString(w, "ok")
	}))

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/hello", nil))
	out := buf.String()
	for _, want := range []string{"method=GET", "path=/hello", "status=201"} {
		if !strings.Contains(out, want) {
			t.Errorf("log %q is missing %q", out, want)
		}
	}

	buf.Reset()
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/other", nil))
	if !strings.Contains(buf.String(), "status=200") {
		t.Errorf("implicit status should be logged as 200, log: %q", buf.String())
	}
}

func TestRecoverMiddleware(t *testing.T) {
	h := Recover()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	var seen string
	h := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = RequestIDFromContext(r.Context())
	}))

	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, httptest.NewRequest("GET", "/", nil))
	id1 := rec1.Header().Get("X-Request-ID")
	if id1 == "" || id1 != seen {
		t.Fatalf("header %q and context %q must match and be non-empty", id1, seen)
	}

	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest("GET", "/", nil))
	if id2 := rec2.Header().Get("X-Request-ID"); id2 == id1 {
		t.Errorf("two requests got the same generated ID %q", id1)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "abc-123")
	rec3 := httptest.NewRecorder()
	h.ServeHTTP(rec3, req)
	if got := rec3.Header().Get("X-Request-ID"); got != "abc-123" || seen != "abc-123" {
		t.Errorf("incoming ID not reused: header=%q context=%q", got, seen)
	}

	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Errorf("empty context returned %q, want \"\"", got)
	}
}

// ---------- retry ----------

func TestRetrySucceedsEventually(t *testing.T) {
	calls := 0
	err := Retry(context.Background(), 5, time.Millisecond, func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("err=%v calls=%d, want nil and 3", err, calls)
	}
}

func TestRetryGivesUp(t *testing.T) {
	boom := errors.New("boom")
	calls := 0
	err := Retry(context.Background(), 3, time.Millisecond, func() error {
		calls++
		return boom
	})
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
	if !errors.Is(err, boom) {
		t.Fatalf("error %v must wrap the last error", err)
	}
	if !strings.Contains(err.Error(), "3 attempts") {
		t.Errorf("error %q should mention the number of attempts", err.Error())
	}
}

func TestRetryAtLeastOnce(t *testing.T) {
	calls := 0
	_ = Retry(context.Background(), 0, time.Millisecond, func() error {
		calls++
		return nil
	})
	if calls != 1 {
		t.Errorf("calls = %d, want 1 (attempts < 1 is treated as 1)", calls)
	}
}

func TestRetryRespectsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	err := Retry(ctx, 10, 200*time.Millisecond, func() error {
		calls++
		return errors.New("always failing")
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want one satisfying context.Canceled", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Retry took %v, it must stop waiting when ctx is cancelled", elapsed)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

// ---------- token bucket ----------

func TestTokenBucket(t *testing.T) {
	clock := newFakeClock()
	b := NewTokenBucket(2, time.Second, clock.Now)

	expect := func(step string, want ...bool) {
		t.Helper()
		for i, w := range want {
			if got := b.Allow(); got != w {
				t.Fatalf("%s: Allow #%d = %v, want %v", step, i+1, got, w)
			}
		}
	}

	expect("initial burst", true, true, false)
	clock.Advance(time.Second)
	expect("after 1s", true, false)
	clock.Advance(500 * time.Millisecond)
	expect("after 500ms", false)
	clock.Advance(500 * time.Millisecond)
	expect("after another 500ms", true, false)
	clock.Advance(time.Hour)
	expect("refill is capped at capacity", true, true, false)
}

// ---------- graceful shutdown ----------

func TestRunGracefulShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		io.WriteString(w, "done")
	})}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runErr := make(chan error, 1)
	go func() { runErr <- Run(ctx, srv, ln, 2*time.Second) }()

	status := make(chan int, 1)
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get("http://" + ln.Addr().String() + "/slow")
		if err != nil {
			status <- -1
			return
		}
		resp.Body.Close()
		status <- resp.StatusCode
	}()

	time.Sleep(100 * time.Millisecond) // the request is now in flight
	cancel()

	select {
	case err := <-runErr:
		if err != nil {
			t.Errorf("Run returned %v, want nil", err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
	if got := <-status; got != http.StatusOK {
		t.Errorf("in-flight request status = %d, want 200 (it must be allowed to finish)", got)
	}
}

// ---------- error group ----------

func TestGroupSuccess(t *testing.T) {
	g, _ := WithContext(context.Background())
	var n atomic.Int32
	for i := 0; i < 3; i++ {
		g.Go(func() error {
			n.Add(1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Errorf("Wait = %v, want nil", err)
	}
	if n.Load() != 3 {
		t.Errorf("ran %d functions, want 3", n.Load())
	}
}

func TestGroupFirstErrorCancels(t *testing.T) {
	boom := errors.New("boom")
	g, ctx := WithContext(context.Background())
	g.Go(func() error { return boom })
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			return errors.New("context was never cancelled")
		}
	})
	if err := g.Wait(); err != boom {
		t.Errorf("Wait = %v, want the FIRST error (boom)", err)
	}
	if ctx.Err() == nil {
		t.Error("the group context must be cancelled after Wait")
	}
}

// ---------- TTL cache ----------

func TestTTLCache(t *testing.T) {
	clock := newFakeClock()
	c := NewTTLCache[string, int](10*time.Second, clock.Now)

	if _, ok := c.Get("missing"); ok {
		t.Error("missing key reported as present")
	}
	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = (%d, %v), want (1, true)", v, ok)
	}
	clock.Advance(6 * time.Second)
	c.Set("a", 2) // restarts the TTL
	clock.Advance(6 * time.Second)
	if v, ok := c.Get("a"); !ok || v != 2 {
		t.Errorf("after refresh Get(a) = (%d, %v), want (2, true)", v, ok)
	}
	clock.Advance(5 * time.Second)
	if _, ok := c.Get("a"); ok {
		t.Error("entry should have expired")
	}

	c.Set("b", 3)
	c.Delete("b")
	if _, ok := c.Get("b"); ok {
		t.Error("deleted key still present")
	}
}

func TestTTLCacheConcurrent(t *testing.T) {
	c := NewTTLCache[int, int](time.Minute, time.Now)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Set(i, i*i)
			if v, ok := c.Get(i); !ok || v != i*i {
				t.Errorf("Get(%d) = (%d, %v)", i, v, ok)
			}
		}()
	}
	wg.Wait()
}
