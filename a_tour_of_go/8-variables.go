// Variables
// The var statement declares a list of variables; as in function argument lists, the type is last.
//
// A var statement can be at package or function level. We see both in this example.

package a_tour_of_go

import "fmt"

var c, py, java bool


// Variables with initializers
// A var declaration can include initializers, one per variable.

// If an initializer is present, the type can be omitted; the variable will take the type of the initializer.



var k, j int = 1, 2

func main() {

	var i int
	fmt.Println(i, c, py, java)

	var r, rust, js = true, false, "no!"
	fmt.Println(k, j, r, rust, js)
}
