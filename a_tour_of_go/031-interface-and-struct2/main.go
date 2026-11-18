// Structs
// A struct is a collection of fields.

package main

import "fmt"

type gasEngine struct {
	mpg     uint8
	gallons uint8
}

type eletricEngine struct {
	mpkwh uint8
	kwh   uint8
}

type engine interface{
	milesLeft() uint8
}


func (e gasEngine) milesLeft() uint8 {
	return e.gallons * e.mpg
}

func (e eletricEngine) milesLeft() uint8 {
	return e.mpkwh * e.kwh
}


func canMakeIt(e engine, miles uint8) {
	if miles <= e.milesLeft() {
		fmt.Println("Yuu can make it there")
	} else {
		fmt.Println("Need to fuel up first")
	}
}

func main() {
	var myGasEngine gasEngine = gasEngine{25, 15}
	var myEletricEngine eletricEngine = eletricEngine{25, 15}
	canMakeIt(myGasEngine, 50)
	canMakeIt(myEletricEngine, 50)
}
