//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
//	"strconv"
//	"flag"
//	"../src/driver"
//	"../src/states"
//	"../src/queue"

	"driver"
	"states"
	"queue"
//	"time"

//	"../src/communication"
)

// elevatorNumber, numberOfFloors and numberOfElevators constants are set in ../src/queue/queue.go

// -1 if a floor is not reached. If floor reached: 1-4. Belongs to HandleFloorSensor() 
var reached int = -1	

func main(){

	// Initialize hardware
	if (!driver.Initialize(queue.GetNumberOfFloors())){
		fmt.Println("main: Unable to initialize hardware..")
	}

	// Initialize network here!
	//fmt.Printf("****Successfully initialized driver on elevator nr.: %v****\n****to communicate with %v other elevators****\n", elevatorNumber, numberOfElevators)

	// Temporary message function
	fmt.Printf("main:\n****Successfully initialized driver on elevator nr. %v.****\n****Network NOT currently initialized.****\n\n", queue.GetElevatorNumber())

	queue.InitializeQueue()

	/*
	//Kan tas vekk etterpå
	driver.SetDoorLight()
	driver.SetButtonLight(1, 3)
	driver.SetButtonLight(2, 3)
	driver.SetButtonLight(0, 3)
	time.Sleep(3*time.Second)
	driver.ClearButtonLight(1, 3)
	driver.ClearButtonLight(2, 3)
	driver.ClearButtonLight(0, 3)
	*/
	//driver.ClearAllOrderLights(queue.GetNumberOfFloors())

	// Moving down as part of the initialization
	driver.MoveDown()
	queue.SetDirectionElevator(-1)
	
	/*
	var TempMoveDownChan := make(chan int, 1)
	TempMoveDownChan <- -1
	driver.MotorControl(TempMoveDownChan)
	queue.InitializeQueue()
	*/

	fmt.Printf("main: Elevator #: %v\n", queue.GetElevatorNumber())
	fmt.Printf("main: # floors: %v\n", queue.GetNumberOfFloors())
	fmt.Printf("main: # elevators: %v\n", queue.GetNumberOfElevators())
	fmt.Printf("main: Current task in initialization: %v\n", queue.GetAssignedTask())
	fmt.Printf("********for loop Go!********\n")

	go states.Clock()

	for{

		//fmt.Printf("state of clock is: %v\n", states.ClockTick())

		if (driver.GetFloorSensorSignal() != reached){
			fmt.Printf("\nmain: Floor sensor says: %v\n", driver.GetFloorSensorSignal())
			states.PrintState()
		}

		//driver.MotorControl()

		
		//fmt.Printf("main: loop part 1\n")

		HandleStopButton()
		
		
		//fmt.Printf("main: loop part 2\n")

		HandleFloorSensor()

		
		//fmt.Printf("main: loop part 3\n")
		
		if (!states.CheckElevatorStopButtonVariable()){
			HandleFloorButtons()
		}
		
		
		//fmt.Printf("main: loop part 4\n")

		HandleCommandButtons()

		
		//fmt.Printf("main: loop part 5\n")
		
		HandleTimeOut()

		
		//fmt.Printf("main: loop part 6\n")
		
		HandleObstruction()

		//fmt.Printf("GetCurrentFloor: %v\n", queue.GetCurrentFloor())
	}
}

func HandleFloorButtons(){
	// Checking floor buttons and adding orders, setting button lights and calling events
	for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
		if (floor > 1 && floor < queue.GetNumberOfFloors()){
			if (driver.CheckButton(0, floor) && !queue.CheckOrder(0, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddOrder(0, floor)
					driver.SetButtonLight(0, floor)
					//fmt.Println("states.EvNewOrderInEmptyQueue(floor)")
				}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
					//fmt.Println("states.EvNewOrderInCurrentFloor()")
				}else{					
					queue.AddOrder(0, floor)
					driver.SetButtonLight(0, floor)
					//fmt.Println("not")
				}
			}
			if (driver.CheckButton(1, floor) && !queue.CheckOrder(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
				}else{				
					queue.AddOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
			}
		}
		// Only one direction from floor 1 and GetNFloors()
		if (floor == 1){
			if (driver.CheckButton(0, floor) && !queue.CheckOrder(0, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)
					queue.AddOrder(0, floor)
					driver.SetButtonLight(0, floor)
				}else if (floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
				}else{
					queue.AddOrder(0, floor)
					driver.SetButtonLight(0, floor)
				}
			}
		}
		if (floor == queue.GetNumberOfFloors()){
			if (driver.CheckButton(1, floor) && !queue.CheckOrder(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)
					queue.AddOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}else if (floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()		
				}else{
					queue.AddOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
			}
		}
	}
}

func HandleCommandButtons(){
	for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
		if (driver.CheckButton(2, floor) && !queue.CheckOrder(2, floor)){
			if (states.CheckElevatorStopButtonVariable()){
				states.EvStopButtonOff()
			}
			if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
				states.EvNewOrderInEmptyQueue(floor)
				queue.AddOrder(2, floor)
				driver.SetButtonLight(2, floor)
			}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
				states.EvNewOrderInCurrentFloor()
			}else{
				queue.AddOrder(2, floor)
				driver.SetButtonLight(2, floor)

			}
		}
	}
}

func HandleFloorSensor(){
	// Check if floor reached and call EvFloorReached() once
	if (driver.GetFloorSensorSignal() != reached){
		if (reached == -1){
			reached = driver.GetFloorSensorSignal()
			states.EvFloorReached(reached)
		}else{
			reached = driver.GetFloorSensorSignal()
		}
	}
}

func HandleStopButton(){
	// Check if stop button is pressed, if so, stop elevator and remove all orders
	if (driver.CheckStopButton()){
		states.EvStopButton()
	}
}

func HandleTimeOut(){
	// Time out signal check
	if (states.CheckTimeOut() && !driver.CheckObstruction()){
		states.EvTimerOut()
	}
}

func HandleObstruction(){
	// Universal obstruction signal
	if (driver.CheckObstruction()){
		states.EvObstructionOn()
	}else if (!driver.CheckObstruction()){
		states.EvObstructionOff()
	}	
}