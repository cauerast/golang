// Struct Fields
// Struct fields are accessed using a dot.

package main

import "fmt"

type Vertex1 struct {
	X int
	Y int
}

func main() {
	v := Vertex1{1, 2}
	v.X = 4
	fmt.Println(v.X)
}
