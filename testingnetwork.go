package main

import "fmt"
//import "reflect"

var elevIpAddresses = make([]string, 3)
var elevPorts = make([]int, 3)


func main(){
	elevIpAddresses = []string{"129.241.187.143", "129.241.187.141", "129.241.187.146"}
	elevPorts = []int{2005, 2004, 2006}

	for c := range elevIpAddresses{
		fmt.Println(elevIpAddresses[c])
	}
}