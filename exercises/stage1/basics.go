// Package stage1 covers Go syntax, types, and pointers.
package stage1

// Swap returns its two arguments in reversed order.
//
// TODO: Use Go's multiple return values. No temporary variable needed.
func Swap(a, b int) (int, int) {
	return 0, 0
}

// Sum returns the sum of all given numbers (0 if none).
//
// TODO: Iterate over the variadic parameter with range. Inside the function,
// nums is just a []int; callers can also pass a slice with Sum(xs...).
func Sum(nums ...int) int {
	return 0
}

// NewCounter returns a function that returns 1, 2, 3, ... on each call.
// Each counter returned by NewCounter has its own independent state.
//
// TODO: Implement it with a closure capturing a local variable.
func NewCounter() func() int {
	return func() int { return 0 }
}

// Weekday is an enumeration built with typed constants.
type Weekday int

// TODO: Rewrite this block using iota so that Sunday == 0 ... Saturday == 6.
// (In the starter code every constant is zero, which is incorrect.)
const (
	Sunday    Weekday = 0
	Monday    Weekday = 0
	Tuesday   Weekday = 0
	Wednesday Weekday = 0
	Thursday  Weekday = 0
	Friday    Weekday = 0
	Saturday  Weekday = 0
)

// String implements fmt.Stringer: "Sunday", "Monday", ..., "Saturday".
// Out-of-range values must return "Unknown".
//
// TODO: Implement it (an array or a switch both work).
func (d Weekday) String() string {
	return ""
}

// Describe returns a description of v based on its dynamic type:
//
//	int    -> "int:42"
//	string -> "string:hi"
//	bool   -> "bool:true"
//	other  -> "unknown"
//
// TODO: Use a type switch (switch x := v.(type)).
func Describe(v any) string {
	return ""
}
