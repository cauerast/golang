package main

import (
	"fmt"
	"strings"
)


func main(){
	//var myString string = "résumé"
	var myString = []rune("résumé")
	var indexed = myString[1] // 'é' uses 2 bytes 
	fmt.Printf("%v, %T \n", indexed, indexed)

	for i, v := range myString {
		fmt.Println(i, v)
	}
	fmt.Printf("The length of my string is: %v", len(myString)) // 8 if string, 6 if rune

	var myRune = 'a'
	fmt.Printf("\n myRune: %v", myRune)

	var strSlice = []string{"y", "o", "l", "o"}

	/// inefficient /// 
	// var catStr = ""
	// for i := range strSlice{
	// 	catStr += strSlice[i]
	// }

	var strBuilder strings.Builder
	for i := range strSlice{
		strBuilder.WriteString(strSlice[i])
	}
	var catStr = strBuilder.String()

	fmt.Printf("\n %v", catStr)

}