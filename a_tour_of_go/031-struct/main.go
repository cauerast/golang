// Structs
// A struct is a collection of fields.

package main

import "fmt"

type gasEngine struct{
	mpg uint8
	gallons uint8
	// onwerInfo onwer
	onwer
}

type onwer struct{
	name string
}

func main() {
	var myEngine = gasEngine{25, 15, onwer{"Alex"}}
	myEngine.mpg = 20
	fmt.Println(myEngine)
	fmt.Println(myEngine.mpg, myEngine.gallons, myEngine.name)
}
