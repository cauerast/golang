// Package stage7 covers patterns you meet in production Go code: generics,
// functional options, HTTP middleware, structured logging, retries, rate
// limiting, graceful shutdown, error groups, and injectable clocks.
package stage7

// Number is a type constraint. The ~ means "any type whose underlying type is
// this one", so `type Celsius float64` satisfies it too.
type Number interface {
	~int | ~int64 | ~float64
}

// Map applies f to every element of s and returns the results in order.
//
// TODO: Preallocate the result with make([]U, 0, len(s)).
func Map[T, U any](s []T, f func(T) U) []U {
	return nil
}

// Filter returns the elements of s for which keep returns true, in order.
// The result is never nil.
//
// TODO: Same idea as the non-generic version, now for any T.
func Filter[T any](s []T, keep func(T) bool) []T {
	return nil
}

// Reduce folds s into a single value, starting from init.
//
// TODO: acc = f(acc, v) for every element.
func Reduce[T, A any](s []T, init A, f func(A, T) A) A {
	return init
}

// SumOf adds up the elements of xs (the zero value if xs is empty).
//
// TODO: The Number constraint is what makes `+` legal on a type parameter.
func SumOf[T Number](xs []T) T {
	var zero T
	return zero
}

// Set is a generic set. Its zero value is NOT ready to use: construct it with NewSet.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet returns a set containing the given items.
//
// TODO: Initialize the map (a nil map panics on write) and add the items.
func NewSet[T comparable](items ...T) *Set[T] {
	return &Set[T]{}
}

// Add inserts v into the set.
//
// TODO: map[T]struct{} is the idiomatic set: struct{} takes zero bytes.
func (s *Set[T]) Add(v T) {
}

// Has reports whether v is in the set.
//
// TODO: Use the comma-ok idiom.
func (s *Set[T]) Has(v T) bool {
	return false
}

// Len returns the number of distinct elements.
//
// TODO: Implement it.
func (s *Set[T]) Len() int {
	return -1
}
