// Package stage3 covers implicit interfaces and error handling.
package stage3

import "io"

// Shape is satisfied implicitly by any type with these two methods.
// There is no "implements" keyword anywhere.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle is a Shape.
type Rectangle struct {
	Width, Height float64
}

// Area returns Width * Height.
//
// TODO: Implement it.
func (r Rectangle) Area() float64 {
	return 0
}

// Perimeter returns 2 * (Width + Height).
//
// TODO: Implement it.
func (r Rectangle) Perimeter() float64 {
	return 0
}

// Circle is a Shape.
type Circle struct {
	Radius float64
}

// Area returns pi * r^2.
//
// TODO: Implement it using math.Pi.
func (c Circle) Area() float64 {
	return 0
}

// Perimeter returns 2 * pi * r.
//
// TODO: Implement it using math.Pi.
func (c Circle) Perimeter() float64 {
	return 0
}

// TotalArea sums the areas of all shapes.
//
// TODO: Range over the slice and call Area() on each element.
func TotalArea(shapes []Shape) float64 {
	return -1
}

// Largest returns the shape with the greatest area, or nil if shapes is empty.
//
// TODO: Return the interface value itself, not a copy of a concrete type.
func Largest(shapes []Shape) Shape {
	return Rectangle{}
}

// CountWords reads r to the end and returns how many whitespace-separated
// words it contained. Any read error must be returned.
//
// TODO: Accept the small io.Reader interface (not a string or a file) and use
// bufio.NewScanner with bufio.ScanWords. Don't forget scanner.Err().
func CountWords(r io.Reader) (int, error) {
	return -1, nil
}
