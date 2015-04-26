//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
	"driver"
	"states"
	"queue"
	"communication"

)
// elevatorNumber, numberOfFloors and numberOfElevators constants are set in ../src/queue/queue.go
var	elevIpAddresses []string = []string{"129.241.187.159", "129.241.187.158", "129.241.187.161"}//, "129.241.187.154"}
var elevPorts []int = []int{20011, 20010, 20025}//, 20008, 20025}

func main(){

	if (!driver.Initialize(queue.GetNumberOfFloors())){
		fmt.Println("main: Unable to initialize hardware..")
	}
	if (!communication.Initialize(elevIpAddresses, elevPorts)){
		fmt.Println("main: Unable to initialize network..")
	}
	if (!queue.Initialize()){
		fmt.Println("main: Unable to initialize queue..")
	}
	if (!states.Initialize()){
		fmt.Println("main: Unable to initialize states..")
	}
	
	go states.CheckOrderButtonsAndSendToOrderChannels()
	
	go states.CheckOrderChansAndCallEvents()
	
	go states.CheckRemoteChanAndCallEvents()
	
	go states.CheckIfTimeoutCallEventAndPrintQueue()
	
	go states.CheckFloorSensorAndCallEvents()
	
	go states.CheckStopAndObstructionAndCallEvents()

	select{
	}
}