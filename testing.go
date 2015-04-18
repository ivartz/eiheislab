package main

import "fmt"
import "reflect"

var HELLO int

const ha int = 5

func GetHello() int{
	return HELLO
}

func testeste(){

}

func SetConst(c int){
	HELLO = c
}

const numberOfElevators int = 3

var size int = 3

var array = make([]int, size)

var floorElevator [numberOfElevators]int

func main(){
	/*
	for index := 0; index < numberOfElevators; index++{
		fmt.Println(floorElevator[index])
	}
	*/
	/*for index := 0; index < size; index++{
		fmt.Println(array[index])
	}
	testeste()
	*/
	const noe int = 7
	SetConst(noe)
	fmt.Println(HELLO)
	fmt.Println(reflect.TypeOf(HELLO))
	fmt.Println(reflect.TypeOf(ha))

}