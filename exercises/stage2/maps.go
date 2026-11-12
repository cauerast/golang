package stage2

// WordFrequency counts how many times each word appears in text.
// Words are separated by whitespace and compared case-insensitively
// (lowercase them). An empty text yields an empty, non-nil map.
//
// TODO: Use strings.Fields, strings.ToLower, and the zero-value behavior of
// map reads (m[k]++ works on missing keys).
func WordFrequency(text string) map[string]int {
	return nil
}

// GroupByLength groups words by their length in bytes.
// Example: ["a", "bb", "c"] -> {1: ["a", "c"], 2: ["bb"]}. Order inside each
// group follows the input order.
//
// TODO: Append to m[len(w)] directly; a missing key yields a nil slice.
func GroupByLength(words []string) map[int][]string {
	return nil
}

// SortedKeys returns the keys of m in ascending order.
//
// TODO: Map iteration order is random in Go. Collect the keys, then sort
// them (sort.Strings or slices.Sorted(maps.Keys(m))).
func SortedKeys(m map[string]int) []string {
	return nil
}

// Lookup returns the value for key and whether the key was present,
// using the "comma ok" idiom. Missing keys must return (-1, false).
//
// TODO: v, ok := m[key]. Why can't you just compare v to 0 to detect absence?
func Lookup(m map[string]int, key string) (int, bool) {
	return 0, true
}
