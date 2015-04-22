package states

import (
	"fmt"
	"time"
	"driver"
	"queue"
)

type ElevatorState int
const (
	INIT ElevatorState = iota
	IDLE
	MOVING
	STOPPED
	
	DOOR_OPEN
	DOOR_OBSTRUCTED
	OBSTRUCTION
	STOPPED_OBSTRUCTION
)

var state ElevatorState = INIT

// Syncs with floor inducator set or not set. set {1,2,3,4}, not set {-1}
var floorIndicatorSet int = -1

// Syncs with elevator stop button on or off
var elevatorStopButton bool = false

var obstruction bool = false

func EvFloorReached(f int){
	fmt.Printf("states: Floor %v reached\n", f)
	
	switch state{
	case INIT:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		SetFloorIndicator(f)
		StopElevatorGood()
		state = IDLE
		queue.SetCurrentFloor(f)
		break

	case MOVING:
		SetFloorIndicator(f)
		queue.SetCurrentFloor(f)
		if (queue.ShallStop()){
			StopElevatorGood()

			if (f > 1 && f < queue.GetNumberOfFloors()){
				//queue.RemoveOrder(queue.GetDirectionElevator(), f)

				if (queue.GetDirectionElevator() == 1){
					queue.RemoveOrder(0, f)
					driver.ClearButtonLight(0, f)
				}else if (queue.GetDirectionElevator() == -1){
					queue.RemoveOrder(1, f)
					driver.ClearButtonLight(1, f)
				}
				driver.ClearButtonLight(2,f)

				if (queue.ShallRemoveOppositeFloorOrder()){
					if (queue.GetDirectionElevator() == 1){
						queue.RemoveOrder(1, f)
						driver.ClearButtonLight(1, f)
						driver.ClearButtonLight(2, f)
					}else if (queue.GetDirectionElevator() == -1){
						queue.RemoveOrder(0, f)
						driver.ClearButtonLight(0, f)
						driver.ClearButtonLight(2, f)						
					}
				}
			}else if (f == 1){
				queue.RemoveOrder(0, f)
				driver.ClearButtonLight(0, f)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				// Also so that ShallStop() will d
				queue.SetDirectionElevator(1)
			}else if (f == queue.GetNumberOfFloors()){
				queue.RemoveOrder(1, f)
				driver.ClearButtonLight(1, f)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				queue.SetDirectionElevator(-1)
			}
			driver.SetDoorLight()
			fmt.Println("states: Door open")
			fmt.Println("states: Calling ResetTimer() now from EvFloorReached()!")
			go ResetTimer()
			state = DOOR_OPEN
		}else if (f == 1){
			// Changing direction
			queue.SetDirectionElevator(1)
			driver.MoveUp()						
		}else if (f == queue.GetNumberOfFloors()){
			// Changing direction
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
		}else{
			state = MOVING
		}
		break

	case IDLE:
		state = IDLE
		break

	case DOOR_OPEN:
		break

	case STOPPED:
		break

	case DOOR_OBSTRUCTED:
		break

	case STOPPED_OBSTRUCTION:
		break

	case OBSTRUCTION:
		break

	default:
		fmt.Println("states: Illegal state when evFloorReached()")
	}
}

func EvTimerOut(){
	//ClearTimeOut()
	switch state{
	case DOOR_OPEN:
		driver.ClearDoorLight()
		fmt.Println("states: Timeout, door closed. Assigning new task")
		fmt.Printf("states: Direction for correct AssignNewTask(): %v\n", queue.GetDirectionElevator())
		
		queue.AssignNewTask()
		
		fmt.Printf("states: Assigned task is: %v\n", queue.GetAssignedTask())
		fmt.Printf("states: Current floor is: %v\n", queue.GetCurrentFloor())
		queue.PrintQueue()
		
		if (queue.GetAssignedTask() != -1){
			state = MOVING
			if (queue.GetAssignedTask() > queue.GetCurrentFloor()){
				queue.SetDirectionElevator(1)
				driver.MoveUp()
				fmt.Println("states: MoveUp() called from EvTimerOut()")
			}else if (queue.GetAssignedTask() < queue.GetCurrentFloor()){
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
				fmt.Println("states: MoveDown() called from EvTimerOut()")
			}else if (queue.GetAssignedTask() == queue.GetCurrentFloor()){ // To fix that the elevator can go back to most recently passed floor after a sudden stop (between two floors)
				if (queue.GetDirectionElevator() == 1){
					queue.SetDirectionElevator(-1)
					driver.MoveDown()
					fmt.Println("states: Special MoveDown() called from EvTimerOut()")
				}else if (queue.GetDirectionElevator() == -1){
					queue.SetDirectionElevator(1)
					driver.MoveUp()
					fmt.Println("states: Special MoveUp() called from EvTimerOut()")
				}
			}
		}else if (queue.GetAssignedTask() == -1){
			state = IDLE
		}
		break

	case DOOR_OBSTRUCTED:
		break

	default:
		break
	}
}

func EvStopButton(){
	switch state{
	case MOVING:
		StopElevatorGood()
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		//fmt.Println("state = STOPPED HER 1 ****************************************")
		state = STOPPED
		break

	case STOPPED:
		break

	case OBSTRUCTION:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED_OBSTRUCTION

	case STOPPED_OBSTRUCTION:
		break

	case DOOR_OBSTRUCTED:
		break

	case IDLE:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		driver.SetDoorLight()
		//fmt.Println("state = STOPPED HER 2 ****************************************")
		state = STOPPED
		break

	default:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		//fmt.Println("state = STOPPED HER 3 ****************************************")
		state = STOPPED
		break
	}
}

func EvStopButtonOff(){
	switch state{
	case STOPPED:
		driver.ClearStopButtonLight()
		driver.ClearDoorLight()
		ClearElevatorStopButtonVariable()
		//state = MOVING
		//state = STOPPED // Fordi evStopOff() kalles av bestilling fra COMMAND_BUTTONS
		break

	case STOPPED_OBSTRUCTION:
		driver.ClearStopButtonLight()
		ClearElevatorStopButtonVariable()
		state = OBSTRUCTION
		break

	default:
		break
	}
}

func EvObstructionOn(){
	fmt.Println("states: EvObstructionOn()")
	switch state{
	case MOVING:
		StopElevatorGood()
		state = OBSTRUCTION
		break

	case DOOR_OPEN:
		state = DOOR_OBSTRUCTED
		break

	case DOOR_OBSTRUCTED:
		driver.SetDoorLight()
		break

	case STOPPED:
		state = STOPPED_OBSTRUCTION
		break

	case IDLE:
		state = OBSTRUCTION
		break

	default:
		break
	}
}

func EvObstructionOff(){
	fmt.Println("states: EvObstructionOff()")
	switch state{
	case OBSTRUCTION:
		if (queue.GetAssignedTask() == -1){
			queue.AssignNewTask()
		}
		if (queue.GetAssignedTask() == -1){
			//fmt.Println("states: EvObstructionOff() returns IDLE here")
			state = IDLE
		}else if (queue.GetAssignedTask() > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
			fmt.Println("states: MoveUp() called from EvObstructionOff()")
		}else if (queue.GetAssignedTask() < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
			fmt.Println("states: MoveDown() called from EvObstructionOff()")
		}else if (queue.GetAssignedTask() == queue.GetCurrentFloor()){ // To fix that the elevator can go back to most recently passed floor after a sudden stop (between two floors)
			if (queue.GetDirectionElevator() == 1){
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
			}else if (queue.GetDirectionElevator() == -1){
				queue.SetDirectionElevator(1)
				driver.MoveUp()
			}
		}
		driver.ClearDoorLight()
		state = MOVING
		break

	case DOOR_OBSTRUCTED:
		fmt.Println("states: Calling ResetTimer() now from EvObstructionOff()!")
		go ResetTimer()
		state = DOOR_OPEN
		break

	case STOPPED_OBSTRUCTION:
		//queue.AssignNewTask()
		fmt.Printf("states: Assigned task shoud be -1: %v\n", queue.GetAssignedTask())
		//fmt.Println("state = STOPPED HER 4 ****************************************")
		state = STOPPED
		break

	default:
		break
	}
}

func EvNewOrderInEmptyQueue(floorButton int){
	fmt.Println("states: EvNewOrderInEmptyQueue() was called")
	switch state{
	case IDLE:
		if (floorButton > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
			fmt.Println("states: MoveUp() called from EvNewOrderInEmptyQueue()")
		}else if (floorButton < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
			fmt.Println("states: MoveDown() called from EvNewOrderInEmptyQueue()")
		}else if (floorButton == queue.GetCurrentFloor()){ // To fix that the elevator can go back to most recently passed floor after a sudden stop (between two floors)
			if (queue.GetDirectionElevator() == 1){
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
			}else if (queue.GetDirectionElevator() == -1){
				queue.SetDirectionElevator(1)
				driver.MoveUp()
			}
		}
		state = MOVING
		break

	case STOPPED:
		if (floorButton > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
			fmt.Println("states: MoveUp() called from EvNewOrderInEmptyQueue()")
		}else if (floorButton < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
			fmt.Println("states: MoveDown() called from EvNewOrderInEmptyQueue()")
		}else if (floorButton == queue.GetCurrentFloor()){
			fmt.Printf("states: WARNING:EvNewOrderInEmptyQueue() was called when there was a \nnew order in the same floor as the elevator. \nAction Catching up by reversing motor.\n")
			if (queue.GetDirectionElevator() == 1){
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
			}else if (queue.GetDirectionElevator() == -1){
				queue.SetDirectionElevator(1)
				driver.MoveUp()
			}
		}
		state = MOVING
		break

	case OBSTRUCTION:
		break

	case INIT:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		break

	default:
		break
	}
}

func EvNewOrderInCurrentFloor(){//f int, buttonDirection int){
	fmt.Println("states: EvNewOrderInCurrentFloor() was called")
	switch state{
	case IDLE:
		driver.SetDoorLight()
		fmt.Println("states: Calling ResetTimer() now from EvNewOrderInCurrentFloor()!")
		go ResetTimer()
		state = DOOR_OPEN
		break

	case DOOR_OPEN:
		fmt.Println("states: Calling ResetTimer() now from EvNewOrderInCurrentFloor()!")
		quitResetTimer <- true
		go ResetTimer()
		state = DOOR_OPEN			
		break

	case STOPPED:
		state = IDLE
		break

	default:
		break
	}
}

func StopElevatorGood(){
	fmt.Println("states: Calling StopElevatorGood(). Uses MoveUp/MoveDown and Stop")
	if (queue.GetDirectionElevator() == 1){
		driver.MoveDown()
		time.Sleep(10000 * time.Microsecond)
		driver.Stop()
	}else if (queue.GetDirectionElevator() == -1){
		driver.MoveUp()
		time.Sleep(10000 * time.Microsecond)
		driver.Stop()
	}
}

func SetFloorIndicator(floor int){
	if (floor == -1){
		return
	}else if (floorIndicatorSet != floor){
		driver.SetFloorLight(floor)
		floorIndicatorSet = floor
	}
}

func PrintState(){
	fmt.Println("states: PrintState")
	if (state == INIT){
		fmt.Printf("states: Current state is: INIT\n")
	}
	if (state == IDLE){
		fmt.Printf("states: Current state is: IDLE\n")
	}
	if (state == DOOR_OPEN){
		fmt.Printf("states: Current state is: DOOR_OPEN\n")
	}
	if (state == DOOR_OBSTRUCTED){
		fmt.Printf("states: Current state is: DOOR_OBSTRUCTED\n")
	}
	if (state == MOVING){
		fmt.Printf("states: Current state is: MOVING\n")
	}
	if (state == STOPPED){
		fmt.Printf("states: Current state is: STOPPED\n")
	}
	if (state == OBSTRUCTION){
		fmt.Printf("states: Current state is: OBSTRUCTION\n")
	}
	if (state == STOPPED_OBSTRUCTION){
		fmt.Printf("states: Current state is: STOPPED_OBSTRUCTION\n")
	}	
}

func SetElevatorStopButtonVariable(){
	elevatorStopButton = true
}

func ClearElevatorStopButtonVariable(){
	elevatorStopButton = false
}

func CheckElevatorStopButtonVariable() bool{
	return elevatorStopButton
}

func SetObstructionVariable(){
	obstruction = true
}

func ClearObstructionVariable(){
	obstruction = false
}

func CheckObstructionVariable() bool{
	return obstruction
}