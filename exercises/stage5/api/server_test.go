package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeTodo(t *testing.T, rec *httptest.ResponseRecorder) Todo {
	t.Helper()
	var td Todo
	if err := json.NewDecoder(rec.Body).Decode(&td); err != nil {
		t.Fatalf("response is not a valid Todo JSON: %v", err)
	}
	return td
}

func TestListEmpty(t *testing.T) {
	rec := do(t, NewServer(), "GET", "/todos", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("body = %q, want []", got)
	}
}

func TestCreateAndGet(t *testing.T) {
	h := NewServer()

	rec := do(t, h, "POST", "/todos", `{"title":"learn go"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201", rec.Code)
	}
	created := decodeTodo(t, rec)
	if created.ID != 1 || created.Title != "learn go" || created.Done {
		t.Errorf("created = %+v, want {1 learn go false}", created)
	}

	rec = do(t, h, "POST", "/todos", `{"title":"write tests"}`)
	if second := decodeTodo(t, rec); second.ID != 2 {
		t.Errorf("second ID = %d, want 2", second.ID)
	}

	rec = do(t, h, "GET", "/todos/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /todos/1 status = %d, want 200", rec.Code)
	}
	if got := decodeTodo(t, rec); got != created {
		t.Errorf("GET returned %+v, want %+v", got, created)
	}

	rec = do(t, h, "GET", "/todos", "")
	var list []Todo
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("list is not valid JSON: %v", err)
	}
	if len(list) != 2 || list[0].ID != 1 || list[1].ID != 2 {
		t.Errorf("list = %+v, want two todos ordered by ID", list)
	}
}

func TestBadRequests(t *testing.T) {
	h := NewServer()
	for _, body := range []string{`{"title":""}`, `{}`, `not json`} {
		if rec := do(t, h, "POST", "/todos", body); rec.Code != http.StatusBadRequest {
			t.Errorf("POST %q status = %d, want 400", body, rec.Code)
		}
	}
	if rec := do(t, h, "GET", "/todos/abc", ""); rec.Code != http.StatusBadRequest {
		t.Errorf("GET /todos/abc status = %d, want 400", rec.Code)
	}
}

func TestNotFound(t *testing.T) {
	h := NewServer()
	do(t, h, "POST", "/todos", `{"title":"exists"}`)
	if rec := do(t, h, "GET", "/todos/1", ""); rec.Code != http.StatusOK {
		t.Fatalf("GET existing todo status = %d, want 200", rec.Code)
	}
	if rec := do(t, h, "GET", "/todos/99", ""); rec.Code != http.StatusNotFound {
		t.Errorf("GET missing status = %d, want 404", rec.Code)
	}
	if rec := do(t, h, "DELETE", "/todos/99", ""); rec.Code != http.StatusNotFound {
		t.Errorf("DELETE missing status = %d, want 404", rec.Code)
	}
}

func TestDelete(t *testing.T) {
	h := NewServer()
	do(t, h, "POST", "/todos", `{"title":"temp"}`)

	rec := do(t, h, "DELETE", "/todos/1", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("204 response must have no body, got %q", rec.Body.String())
	}
	if rec := do(t, h, "GET", "/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Errorf("GET after DELETE status = %d, want 404", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := NewServer()
	do(t, h, "POST", "/todos", `{"title":"x"}`)
	if rec := do(t, h, "PUT", "/todos/1", `{}`); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("PUT status = %d, want 405", rec.Code)
	}
}

func TestConcurrentCreates(t *testing.T) {
	h := NewServer()
	const n = 50
	ids := make(chan int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec := do(t, h, "POST", "/todos", `{"title":"parallel"}`)
			if rec.Code != http.StatusCreated {
				t.Errorf("status = %d, want 201", rec.Code)
				return
			}
			ids <- decodeTodo(t, rec).ID
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[int]bool)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID %d: ID generation is not safe for concurrent use", id)
		}
		seen[id] = true
	}
	if len(seen) != n {
		t.Errorf("got %d unique IDs, want %d", len(seen), n)
	}
}
