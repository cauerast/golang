// Slices
// An array has a fixed size. A slice, on the other hand, is a dynamically-sized, flexible view into the elements of an array. In practice, slices are much more common than arrays.

// The type []T is a slice with elements of type T.

// A slice is formed by specifying two indices, a low and high bound, separated by a colon:

// a[low : high]

// This selects a half-open range which includes the first element, but excludes the last one.

// The following expression creates a slice which includes elements 1 through 3 of a:

// a[1:4]

package main

import "fmt"

func main() {
	intArr := [...]int32{1, 2, 3}
	fmt.Println(intArr)

	var intSlice []int32 = []int32{4, 5, 6}
	fmt.Printf("The length is %v w capacity %v", len(intSlice), cap(intSlice))

	intSlice = append(intSlice, 7)
	fmt.Printf("The length is %v w capacity %v \n", len(intSlice), cap(intSlice))

	var intSlice2 []int32 = []int32{8, 9}
	intSlice = append(intSlice, intSlice2...)
	fmt.Println(intSlice)

	// to specify the capacity use make()
	var intSlice3 []int32 = make([]int32, 3, 8)
	fmt.Println(intSlice3)
}
