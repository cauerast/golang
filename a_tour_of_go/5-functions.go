// Functions
// A function can take zero or more arguments.

// In this example, add takes two parameters of type int.

// Notice that the type comes after the variable name.

// (For more about why types look the way they do, see the article on Go's declaration syntax.)

package a_tour_of_go

import "fmt"

func add(a int, b int) int {
	return a + b
}

func main(){
	fmt.Println(add(42, 32))	
	fmt.Println(addContinued(10, 11))	
}

// Functions continued
// When two or more consecutive named function parameters share a type, you can omit the type from all but the last.

// In this example, we shortened

// x int, y int
// to

// x, y int

func addContinued(a, b int) int {
	return a + b
}
