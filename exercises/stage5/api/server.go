// Package api is capstone option B: an in-memory REST API built only on net/http.
//
// Tip: run with the race detector: go test -race ./exercises/stage5/api
package api

import (
	"net/http"
	"sync"
)

// Todo is the JSON resource served by the API.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// store is an in-memory, concurrency-safe collection of todos.
//
// TODO: Add whatever fields you need (a map or slice, the next ID) and use mu
// in every method. Handlers run on many goroutines at once.
type store struct {
	mu sync.Mutex
}

// NewServer returns the HTTP handler for the API. Each call must return a
// handler with its OWN empty store (the tests create one per test).
//
// Routes (use the Go 1.22+ method-and-wildcard patterns of http.ServeMux,
// e.g. mux.HandleFunc("GET /todos/{id}", ...) and r.PathValue("id")):
//
//	GET    /todos       -> 200, JSON array of all todos ordered by ID ([] when empty)
//	POST   /todos       -> body {"title":"..."}; 201 + created Todo as JSON.
//	                       IDs start at 1 and increase by 1.
//	                       Invalid JSON or an empty title -> 400.
//	GET    /todos/{id}  -> 200 + Todo, or 404 if it does not exist.
//	                       A non-numeric id -> 400.
//	DELETE /todos/{id}  -> 204 with no body, or 404 if it does not exist.
//
// All JSON responses must set the header Content-Type: application/json.
// Unsupported methods on a known path should yield 405 (ServeMux does this for you).
//
// TODO: Build the mux, register the routes, and return it.
func NewServer() http.Handler {
	return http.NewServeMux()
}
