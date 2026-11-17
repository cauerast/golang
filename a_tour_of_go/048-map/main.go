// Maps
// A map maps keys to values.

// The zero value of a map is nil. A nil map has no keys, nor can keys be added.

// The make function returns a map of the given type, initialized and ready for use.

package main

import (
	"fmt"
)

func main(){
	var myMap map[string]uint8 = make(map[string]uint8)
	fmt.Println(myMap)

	var myMap2 = map[string]uint8{"Adam": 23, "Sarah": 45}
	fmt.Println(myMap2)
	fmt.Println(myMap2["Adam"])
	//delete(myMap2, "Adam")

	var age, ok = myMap2["Jason"]
	if ok {
		fmt.Printf("The age is %v", age)
	} else {
		fmt.Println("Invalid name")
	}

	// map for loop w/ key and value
	for name, age := range myMap2 {
		fmt.Printf("Name: %v, Age: %v \n", name, age)
	}

	var intArr [3]int8 = [3]int8{1, 2, 3}
	// array for loop w/ index and value
	for i, v := range intArr {
		fmt.Printf("Index: %v, Value: %v \n", i, v)
	}


}