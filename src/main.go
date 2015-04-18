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

// elevatorNumber, numberOfFloors and numberOfElevators constants are set in ../src/queue/queue.go

func main(){

	// Initialize hardware
	if (!driver.Initialize()){
		fmt.Println("Unable to initialize hardware..")
		return 1
	}

	// Initialize network here!
	//fmt.Printf("****Successfully initialized driver on elevator nr.: %v****\n****to communicate with %v other elevators****\n", elevatorNumber, numberOfElevators)

	// Temporary message function
	fmt.Println("****Successfully initialized driver on elevator nr.: %v.****\n****Network NOT currently initialized.****\n", queue.GetElevatorNumber())

	// Moving down as part of the initialization
	driver.MoveDown()
	/*
	var TempMoveDownChan := make(chan int, 1)
	TempMoveDownChan <- -1
	driver.MotorControl(TempMoveDownChan)
	queue.InitializeQueue()
	*/

	// -1 if a floor is not reached. If floor reached: 1-4. Belongs to HandleFloorSensor() 
	reached := -1	

	for{

		driver.MotorControl()

		HandleStopButton()
		
		HandleFloorSensor()
		
		if (!states.CheckElevatorStopButtonVariable()){
			HandleFloorButtons()
		}
		
		HandleCommandButtons()
		
		HandleTimeOut()
		
		HandleObstruction()
	}
}

func HandleFloorButtons(){
	// Checking floor buttons and adding orders, setting button lights and calling events
	for (floor := 1; floor <= driver.GetNFloors(); ++floor){
		if (floor > 1 && floor < driver.GetNFloors()){
			if (driver.CheckButton(0, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddFloorOrder(0, floor)
					driver.SetButtonLight(0, floor)
				}
				else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
				}
				else{					
					queue.AddFloorOrder(0, floor)
					driver.SetButtonLight(0, floor)
				}
			}
			if (driver.CheckButton(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)					
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
				else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()
				}
				else{				
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
			}
			else if (floor == driver.GetFloorSensorSignal()){
				states.EvNewOrderInCurrentFloor()
			}
			else{
				queue.AddFloorOrder(0, floor)
				driver.SetButtonLight(0, floor)
			}
		}
		if (floor == driver.GetNFloors()){
			if (driver.CheckButton(1, floor)){
				if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
					states.EvNewOrderInEmptyQueue(floor)
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
				else if (floor == driver.GetFloorSensorSignal()){
					states.EvNewOrderInCurrentFloor()		
				}
				else{
					queue.AddFloorOrder(1, floor)
					driver.SetButtonLight(1, floor)
				}
			}
		}
	}
}

func HandleCommandButtons(){
	for (floor := 1; floor <= driver.GetNFloors(); ++floor){
		if (driver.CheckButton(2, floor)){
			if (states.CheckElevatorStopButtonVariable()){
				states.EvStopButtonOff()
			}
			if (queue.GetAssignedTask() == -1 && floor != driver.GetFloorSensorSignal()){
				states.EvNewOrderInEmptyQueue()
				queue.AddFloorOrder(2, floor)
			}
			else if (queue.GetAssignedTask() == -1 && floor == driver.GetFloorSensorSignal()){
				states.EvNewOrderInCurrentFloor()
			}
			else{
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
		}
		else{
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
	}
	else if (!driver.CheckObstruction()){
		states.EvObstructionOff()
	}	
}