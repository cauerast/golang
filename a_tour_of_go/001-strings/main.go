package main

import "fmt"


func main(){
	//var myString string = "résumé"
	var myString = []rune("résumé")
	var indexed = myString[1] // 'é' uses 2 bytes 
	fmt.Printf("%v, %T \n", indexed, indexed)

	for i, v := range myString {
		fmt.Println(i, v)
	}
	fmt.Printf("The length of my string is: %v", len(myString)) // 8 if string, 6 if rune
}