// Short variable declarations
// Inside a function, the := short assignment statement can be used in place of a var declaration with implicit type.

// Outside a function, every statement begins with a keyword (var, func, and so on) and so the := construct is not available.

package a_tour_of_go

import "fmt"

func main(){

	var i, j int = 1, 2
	k := 3
	c, py := "c", "py"

	fmt.Println(i, j, k, c, py)
}