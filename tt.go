package main


import (
	"fmt"
)



func main(){

	slice := make ([]bool, 5)
	slice[1] = true
	slice[4] = true


	for index := 0; index < 5 && slice[index] == true; index++{
		fmt.Printf("hei: %v\n", index)
	}
}