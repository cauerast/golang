package stage6

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"
)

// ---------- io ----------

func TestCountBytes(t *testing.T) {
	n, err := CountBytes(strings.NewReader("hello world"))
	if err != nil || n != 11 {
		t.Errorf("CountBytes = (%d, %v), want (11, nil)", n, err)
	}
	boom := errors.New("boom")
	if _, err := CountBytes(iotest.ErrReader(boom)); !errors.Is(err, boom) {
		t.Errorf("CountBytes must return the reader's error, got %v", err)
	}
}

func TestRot13Reader(t *testing.T) {
	r := Rot13Reader{R: strings.NewReader("Hello, World!")}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "Uryyb, Jbeyq!" {
		t.Errorf("rot13 = %q, want %q", got, "Uryyb, Jbeyq!")
	}
	// Reads in tiny chunks must give the same result.
	r = Rot13Reader{R: iotest.OneByteReader(strings.NewReader("abcxyz"))}
	got, _ = io.ReadAll(r)
	if string(got) != "nopklm" {
		t.Errorf("one-byte reads = %q, want %q", got, "nopklm")
	}
}

func TestWriteFrameWireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrame(&buf, []byte("hi")); err != nil {
		t.Fatalf("WriteFrame: %v", err)
	}
	want := []byte{0, 0, 0, 2, 'h', 'i'}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %v, want %v", buf.Bytes(), want)
	}
	big := make([]byte, MaxFrameSize+1)
	var sink bytes.Buffer
	if err := WriteFrame(&sink, big); !errors.Is(err, ErrFrameTooLarge) || sink.Len() != 0 {
		t.Errorf("oversized WriteFrame: err=%v written=%d, want ErrFrameTooLarge and 0 bytes", err, sink.Len())
	}
}

func TestFrameRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	msgs := []string{"hello", "", "world"}
	for _, m := range msgs {
		if err := WriteFrame(&buf, []byte(m)); err != nil {
			t.Fatalf("WriteFrame(%q): %v", m, err)
		}
	}
	for _, want := range msgs {
		got, err := ReadFrame(&buf)
		if err != nil || string(got) != want {
			t.Fatalf("ReadFrame = (%q, %v), want (%q, nil)", got, err, want)
		}
	}
	if _, err := ReadFrame(&buf); err != io.EOF {
		t.Errorf("ReadFrame at end of stream = %v, want io.EOF", err)
	}
}

func TestReadFrameErrors(t *testing.T) {
	// Truncated payload.
	_, err := ReadFrame(bytes.NewReader([]byte{0, 0, 0, 10, 'a', 'b', 'c'}))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("truncated payload error = %v, want io.ErrUnexpectedEOF", err)
	}
	// Truncated header.
	_, err = ReadFrame(bytes.NewReader([]byte{0, 0}))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("truncated header error = %v, want io.ErrUnexpectedEOF", err)
	}
	// Declared length (2 MiB) over the limit.
	_, err = ReadFrame(bytes.NewReader([]byte{0x00, 0x20, 0x00, 0x00}))
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Errorf("oversized frame error = %v, want ErrFrameTooLarge", err)
	}
}

func TestSHA256Hex(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	want := hex.EncodeToString(sum[:])
	got, err := SHA256Hex(strings.NewReader("hello"))
	if err != nil || got != want {
		t.Errorf("SHA256Hex = (%q, %v), want (%q, nil)", got, err, want)
	}
	boom := errors.New("boom")
	if _, err := SHA256Hex(iotest.ErrReader(boom)); !errors.Is(err, boom) {
		t.Errorf("SHA256Hex must return the reader's error, got %v", err)
	}
}

// ---------- net ----------

func TestEchoServer(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- EchoServer(l) }()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg := []byte(strings.Repeat("x", i+1) + "-ping")
			got, err := SendAndReceive(l.Addr().String(), msg, time.Second)
			if err != nil || !bytes.Equal(got, msg) {
				t.Errorf("client %d: got (%q, %v), want (%q, nil)", i, got, err, msg)
			}
		}()
	}
	wg.Wait()

	l.Close()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("EchoServer after Close returned %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Error("EchoServer did not return after the listener was closed")
	}
}

func TestSendAndReceiveTimeout(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	stop := make(chan struct{})
	defer close(stop)
	go func() { // a peer that accepts but never answers
		c, err := l.Accept()
		if err != nil {
			return
		}
		<-stop
		c.Close()
	}()

	start := time.Now()
	_, err = SendAndReceive(l.Addr().String(), []byte("hello"), 100*time.Millisecond)
	var ne net.Error
	if !errors.As(err, &ne) || !ne.Timeout() {
		t.Fatalf("error = %v, want a net.Error with Timeout() == true", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("took %v, the timeout was ignored", elapsed)
	}
}

func TestSendAndReceiveRefused(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := l.Addr().String()
	l.Close() // nothing listens here anymore
	if _, err := SendAndReceive(addr, []byte("x"), time.Second); err == nil {
		t.Error("expected a connection error, got nil")
	}
}

// ---------- http client ----------

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.Error(w, "nope", http.StatusNotFound)
			return
		}
		io.WriteString(w, "hello")
	}))
	defer srv.Close()

	status, body, err := Fetch(context.Background(), srv.Client(), srv.URL+"/ok")
	if err != nil || status != 200 || body != "hello" {
		t.Errorf("Fetch ok = (%d, %q, %v), want (200, \"hello\", nil)", status, body, err)
	}
	status, body, err = Fetch(context.Background(), srv.Client(), srv.URL+"/missing")
	if err == nil || status != 404 || !strings.Contains(body, "nope") {
		t.Errorf("Fetch 404 = (%d, %q, %v), want status 404, body and a non-nil error", status, body, err)
	}
}

func TestFetchDeadline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, _, err := Fetch(ctx, srv.Client(), srv.URL)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want one satisfying context.DeadlineExceeded", err)
	}
}

func TestGetJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/err" {
			http.Error(w, "{}", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"name":"go","size":7}`)
	}))
	defer srv.Close()

	var out struct {
		Name string `json:"name"`
		Size int    `json:"size"`
	}
	if err := GetJSON(context.Background(), srv.Client(), srv.URL, &out); err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
	if out.Name != "go" || out.Size != 7 {
		t.Errorf("decoded = %+v, want {go 7}", out)
	}
	if err := GetJSON(context.Background(), srv.Client(), srv.URL+"/err", &out); err == nil {
		t.Error("expected an error for a 500 response")
	}
}

// ---------- encoding & context ----------

func TestEncodeMessage(t *testing.T) {
	got, err := EncodeMessage(Message{Type: "file", Size: 10})
	if err != nil || string(got) != `{"type":"file","size":10}` {
		t.Errorf("EncodeMessage = (%s, %v)", got, err)
	}
	got, _ = EncodeMessage(Message{Type: "chat", Body: "hi", Size: 2})
	if string(got) != `{"type":"chat","body":"hi","size":2}` {
		t.Errorf("EncodeMessage with body = %s", got)
	}
}

func TestDecodeMessage(t *testing.T) {
	m, err := DecodeMessage([]byte(`{"type":"chat","body":"hi","size":2}`))
	if err != nil || m != (Message{Type: "chat", Body: "hi", Size: 2}) {
		t.Errorf("DecodeMessage = (%+v, %v)", m, err)
	}
	if _, err := DecodeMessage([]byte(`{"type":"x","size":1,"extra":true}`)); err == nil {
		t.Error("unknown fields must be rejected")
	}
	if _, err := DecodeMessage([]byte(`not json`)); err == nil {
		t.Error("invalid JSON must return an error")
	}
}

func TestCollectUntilDone(t *testing.T) {
	// Case 1: the channel is closed.
	in := make(chan int, 3)
	in <- 1
	in <- 2
	in <- 3
	close(in)
	if got := CollectUntilDone(context.Background(), in); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("closed channel: got %v, want [1 2 3]", got)
	}

	// Case 2: the channel stays open; cancellation ends the collection.
	in2 := make(chan int, 2)
	in2 <- 1
	in2 <- 2
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan []int, 1)
	go func() { result <- CollectUntilDone(ctx, in2) }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case got := <-result:
		if !reflect.DeepEqual(got, []int{1, 2}) {
			t.Errorf("cancelled: got %v, want [1 2]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("CollectUntilDone ignored context cancellation")
	}
}

// ---------- text ----------

func TestExtractHrefs(t *testing.T) {
	html := `<p><a href="/a">A</a> <A HREF='b.html'>B</A>
<a class="x" href="https://x.com/c">C</a> <a name="anchor">no href</a></p>`
	want := []string{"/a", "b.html", "https://x.com/c"}
	if got := ExtractHrefs(html); !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractHrefs = %v, want %v", got, want)
	}
}

func TestResolveLinks(t *testing.T) {
	base := "https://example.com/dir/page.html"
	hrefs := []string{
		"/a", "b.html", "../c", "https://other.org/x#frag",
		"mailto:someone@example.com", "javascript:void(0)", "/a", "#top",
	}
	want := []string{
		"https://example.com/a",
		"https://example.com/dir/b.html",
		"https://example.com/c",
		"https://other.org/x",
		"https://example.com/dir/page.html",
	}
	got, err := ResolveLinks(base, hrefs)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveLinks = (%v, %v), want %v", got, err, want)
	}
	if _, err := ResolveLinks("://bad", []string{"/a"}); err == nil {
		t.Error("an invalid base URL must return an error")
	}
}

func TestTokenize(t *testing.T) {
	got := Tokenize("Go is FUN! 100% (really)")
	want := []string{"go", "is", "fun", "100", "really"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Tokenize = %v, want %v", got, want)
	}
}

func TestIndexAndSearch(t *testing.T) {
	docs := map[int]string{
		1: "Go is fun",
		2: "Go concurrency is fun!",
		3: "Rust is fast",
	}
	idx := BuildIndex(docs)
	if !reflect.DeepEqual(idx["go"], []int{1, 2}) {
		t.Errorf(`index["go"] = %v, want [1 2]`, idx["go"])
	}
	if !reflect.DeepEqual(idx["is"], []int{1, 2, 3}) {
		t.Errorf(`index["is"] = %v, want [1 2 3]`, idx["is"])
	}
	if got := Search(idx, "go fun"); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf(`Search("go fun") = %v, want [1 2]`, got)
	}
	if got := Search(idx, "IS fast"); !reflect.DeepEqual(got, []int{3}) {
		t.Errorf(`Search("IS fast") = %v, want [3]`, got)
	}
	if got := Search(idx, "python"); len(got) != 0 {
		t.Errorf(`Search("python") = %v, want empty`, got)
	}
	if got := Search(idx, ""); len(got) != 0 {
		t.Errorf(`Search("") = %v, want empty`, got)
	}
}

// ---------- files ----------

func TestListGoFiles(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"a.go":          "package a",
		"sub/b.go":      "package b",
		"sub/c.txt":     "text",
		".hidden/d.go":  "package d",
		"sub/deep/e.go": "package e",
	}
	for name, content := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := ListGoFiles(root)
	want := []string{"a.go", "sub/b.go", "sub/deep/e.go"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("ListGoFiles = (%v, %v), want %v", got, err, want)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("some content"), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := CopyFile(dst, src)
	if err != nil || n != 12 {
		t.Fatalf("CopyFile = (%d, %v), want (12, nil)", n, err)
	}
	got, err := os.ReadFile(dst)
	if err != nil || string(got) != "some content" {
		t.Errorf("destination = (%q, %v)", got, err)
	}
	if _, err := CopyFile(dst, filepath.Join(dir, "missing")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing source error = %v, want os.ErrNotExist", err)
	}
}
