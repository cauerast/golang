// Package stage2 covers arrays, slices, maps, structs, and composition.
package stage2

// ModifyArray returns a copy of a with the first element set to 100.
// The caller's array must stay untouched.
//
// TODO: Arrays are VALUES in Go: passing one copies it. Contrast this with
// slices, which share a backing array.
func ModifyArray(a [3]int) [3]int {
	return [3]int{}
}

// CloneAndModify returns a new slice with the same elements as s, but with
// the first element replaced by 99. The original slice must NOT change.
// If s is empty, return an empty (non-nil or nil) slice.
//
// TODO: Use make + copy (or slices.Clone) so the result has its own backing array.
func CloneAndModify(s []int) []int {
	return s
}

// RemoveAt removes the element at index i, preserving order, and returns the
// shortened slice. If i is out of range, return s unchanged.
//
// TODO: Use append(s[:i], s[i+1:]...). Then think about what happens to the
// backing array that the caller still holds.
func RemoveAt(s []int, i int) []int {
	return nil
}

// Reverse reverses s in place.
//
// TODO: Swap elements from both ends with a two-index loop; no new slice.
func Reverse(s []int) {
}

// Filter returns a new slice containing only the elements for which keep
// returns true. Order is preserved. Never return nil for an empty result.
//
// TODO: Preallocate with make([]int, 0, len(s)) and append.
func Filter(s []int, keep func(int) bool) []int {
	return nil
}

// GrowthCaps appends n integers one by one to an empty slice and returns the
// capacity observed after each append (len(result) == n).
//
// TODO: Start from `var s []int`, append in a loop, record cap(s) each time.
// Then print the result in a scratch main and observe the growth strategy.
func GrowthCaps(n int) []int {
	return nil
}

// SharedTail demonstrates aliasing. Given base = [1 2 3 4 5], it takes
// sub := base[1:3] and appends 100 to sub WITHOUT reallocating
// (cap(sub) is 4). It returns base after the append.
// Expected result: [1 2 3 100 5].
//
// TODO: Build base, slice it, append to sub, return base.
// Think: why did base change even though you only appended to sub?
func SharedTail() []int {
	return nil
}
