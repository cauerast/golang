package stage1

// Increment adds 1 to the integer that p points to.
//
// TODO: Dereference the pointer. Go has no pointer arithmetic, so this is
// the only thing you can do with p besides passing it around.
func Increment(p *int) {
}

// SafeDeref returns *p, or fallback when p is nil.
//
// TODO: Check for nil before dereferencing (a nil dereference panics).
func SafeDeref(p *int, fallback int) int {
	return fallback + 1 // incorrect placeholder, replace it
}

// Point is a plain struct. Structs are values: assignment copies them.
type Point struct {
	X, Y int
}

// Add returns a NEW Point that is the sum of p and q.
//
// TODO: Implement with a value receiver (p is a copy; p must not be mutated).
func (p Point) Add(q Point) Point {
	return Point{}
}

// Scale multiplies both coordinates of p by factor, IN PLACE.
//
// TODO: This needs a pointer receiver to mutate the caller's struct.
// Change the receiver from (p Point) to (p *Point) and implement it.
func (p Point) Scale(factor int) {
}

// DeferOrder returns the order in which deferred calls run.
// It must return []int{3, 2, 1}.
//
// TODO: Defer three closures that append 1, 2, 3 (in that source order) to a
// slice captured by the function, and return the slice. Think: why does the
// function need a NAMED result for the returned slice to include the appends?
func DeferOrder() (order []int) {
	return nil
}
