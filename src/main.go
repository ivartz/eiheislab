//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
	"strconv"
	"flag"
	"../src/driver"
	"../src/states"
	"../src/queue"
	"../src/communication"
)

var elevatorNumber int = 1;
var numberOfElevators int = 3;


func main(){

	// Initialize hardware
	if (!ElevInit()){
		fmt.Println("Unable to initialize hardware..")
		return 1
	} 
	fmt.Printf("****Successfully initialized driver on elevator nr.: %v****\n****to communicate with %v other elevators****\n", elevatorNumber, numberOfElevators)



}