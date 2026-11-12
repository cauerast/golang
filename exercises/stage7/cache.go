package stage7

import (
	"sync"
	"time"
)

// TTLCache is a concurrency-safe in-memory cache whose entries expire after a
// fixed time-to-live. It combines generics, a mutex, and an injected clock.
type TTLCache[K comparable, V any] struct {
	mu sync.Mutex
	// TODO: add fields: ttl, now func, and a map from K to an entry struct
	// holding the value and its expiry time.
}

// NewTTLCache returns an empty cache. now is the clock used for expiry
// (pass time.Now in production, a fake in tests).
//
// TODO: Initialize the map.
func NewTTLCache[K comparable, V any](ttl time.Duration, now func() time.Time) *TTLCache[K, V] {
	return &TTLCache[K, V]{}
}

// Set stores value under key; it expires ttl after this call.
// Setting an existing key replaces the value and restarts its TTL.
//
// TODO: Implement it (lock the mutex).
func (c *TTLCache[K, V]) Set(key K, value V) {
}

// Get returns the value for key. Missing or expired entries return the zero
// value and false; expired entries should also be removed.
//
// TODO: Implement it.
func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	var zero V
	return zero, false
}

// Delete removes key if present.
//
// TODO: Implement it.
func (c *TTLCache[K, V]) Delete(key K) {
}
