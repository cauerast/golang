// Multiple results
// A function can return any number of results.
//
// The swap function returns two strings.

package a_tour_of_go

import "fmt"

func swap(x, y string) (string, string) {
	return y, x
}

func main(){
	a, b := swap("hello", "world")
	fmt.Println(swap(b, a))
}
