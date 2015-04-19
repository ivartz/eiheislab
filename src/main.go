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

	// Moving down as part of the initialization
	driver.MoveDown()
	/*
	var TempMoveDownChan := make(chan int, 1)
	TempMoveDownChan <- -1
	driver.MotorControl(TempMoveDownChan)
	queue.InitializeQueue()
	*/

	fmt.Printf("main: Elevator #: %v\n", queue.GetElevatorNumber())
	fmt.Printf("main: # floors: %v\n", queue.GetNumberOfFloors())
	fmt.Printf("main: # elevators: %v\n", queue.GetNumberOfElevators())
	fmt.Printf("main: Current task: %v\n", queue.GetAssignedTask())

	for{

		fmt.Printf("main: loop part 1\n")

		fmt.Printf("main: Floor sensor says: %v\n", driver.GetFloorSensorSignal())

		driver.MotorControl()

		
		fmt.Printf("main: loop part 2\n")

		HandleStopButton()
		
		
		fmt.Printf("main: loop part 3\n")

		HandleFloorSensor()

		
		fmt.Printf("main: loop part 4\n")
		
		if (!states.CheckElevatorStopButtonVariable()){
			HandleFloorButtons()
		}
		
		
		fmt.Printf("main: loop part 5\n")

		HandleCommandButtons()

		
		fmt.Printf("main: loop part 6\n")
		
		HandleTimeOut()

		
		fmt.Printf("main: loop part 7\n")
		
		HandleObstruction()
	}
}

func HandleFloorButtons(){
	// Checking floor buttons and adding orders, setting button lights and calling events
	for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
		if (floor > 1 && floor < queue.GetNumberOfFloors()){
			if (driver.CheckButton(0, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddFloorOrder(0, floor)
					driver.SetButtonLight(0, floor)
					//fmt.Println("states.EvNewOrderInEmptyQueue(floor)")
				}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
					//fmt.Println("states.EvNewOrderInCurrentFloor()")
				}else{					
					queue.AddFloorOrder(0, floor)
					driver.SetButtonLight(0, floor)
					//fmt.Println("not")
				}
			}
			if (driver.CheckButton(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
				}else{				
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
			}
		}
		// Only one direction from floor 1 and GetNFloors()
		if (floor == 1){
			if (driver.CheckButton(0, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)
					queue.AddFloorOrder(0, floor)
					driver.SetButtonLight(0, floor)
				}
			}else if (floor == driver.GetFloorSensorSignal()){
				states.EvNewOrderInCurrentFloor()
			}else{
				queue.AddFloorOrder(0, floor)
				driver.SetButtonLight(0, floor)
			}
		}
		if (floor == queue.GetNumberOfFloors()){
			if (driver.CheckButton(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}else if (floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()		
				}else{
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
			}
		}
	}
}

func HandleCommandButtons(){
	for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
		if (driver.CheckButton(2, floor)){
			if (states.CheckElevatorStopButtonVariable()){
				states.EvStopButtonOff()
			}
			if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
				states.EvNewOrderInEmptyQueue(floor)
				queue.AddFloorOrder(2, floor)
			}else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
				states.EvNewOrderInCurrentFloor()
			}else{
				queue.AddFloorOrder(2, floor)
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