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

func EvFloorReached(f int){
	fmt.Printf("**Floor %v reached**\n", f)
	
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
				//queue.RemoveFloorOrder(queue.GetDirectionElevator(), f)

				if (queue.GetDirectionElevator() == 1){
					queue.RemoveFloorOrder(0, f)
					driver.ClearButtonLight(0, f)
				}else if (queue.GetDirectionElevator() == -1){
					queue.RemoveFloorOrder(1, f)
					driver.ClearButtonLight(1, f)
				}
				driver.ClearButtonLight(2,f)

				if (queue.ShallRemoveOppositeFloorOrder()){
					if (queue.GetDirectionElevator() == 1){
						queue.RemoveFloorOrder(1, f)
						driver.ClearButtonLight(1, f)
						driver.ClearButtonLight(2, f)
					}else if (queue.GetDirectionElevator() == -1){
						queue.RemoveFloorOrder(0, f)
						driver.ClearButtonLight(0, f)
						driver.ClearButtonLight(2, f)						
					}
				}
			}else if (f == 1){
				queue.RemoveFloorOrder(0, f)
				driver.ClearButtonLight(0, f)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				// Also so that ShallStop() will d
				queue.SetDirectionElevator(1)
			}else if (f == queue.GetNumberOfFloors()){
				queue.RemoveFloorOrder(1, f)
				driver.ClearButtonLight(1, f)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				queue.SetDirectionElevator(-1)
			}
			driver.SetDoorLight()
			fmt.Println("**Door open**")
			ResetTimer()
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
		fmt.Println("Illegal state when evFloorReached()")
	}
}

func EvTimerOut(){
	//ClearTimeOut()
	switch state{
	case DOOR_OPEN:
		driver.ClearDoorLight()
		fmt.Println("**Timeout, door closed. Assigning new task**")
		queue.AssignNewTask()
		if (queue.GetAssignedTask() > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
		}else if (queue.GetAssignedTask() < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
		}
		fmt.Printf("Assigned task is: %v\n", queue.GetAssignedTask())
		if (queue.GetAssignedTask() != -1){
			state = MOVING
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
		state = STOPPED
		break

	case STOPPED:
		break

	case OBSTRUCTION:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED

	case STOPPED_OBSTRUCTION:
		break

	case IDLE:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		driver.SetDoorLight()
		state = STOPPED
		break

	default:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
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
		state = STOPPED // Fordi evStopOff() kalles av bestilling fra COMMAND_BUTTONS
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
	switch state{
	case OBSTRUCTION:
		if (queue.GetAssignedTask() == -1){
			queue.AssignNewTask()
		}
		if (queue.GetAssignedTask() == -1){
			state = IDLE
		}else if (queue.GetDirectionElevator() == 1){
			if (queue.GetCurrentFloor() == queue.GetNumberOfFloors()){
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
				state = MOVING
			}else{
				queue.SetDirectionElevator(1)
				driver.MoveUp()
				state = MOVING
			}
		}else if (queue.GetDirectionElevator() == -1){
			if (queue.GetCurrentFloor() == 1){
				queue.SetDirectionElevator(1)
				driver.MoveUp()
				state = MOVING
			}else{
				queue.SetDirectionElevator(-1)
				driver.MoveDown()
				state = MOVING
			}
		}
		driver.ClearDoorLight()
		break

	case DOOR_OBSTRUCTED:
		ResetTimer()
		state = DOOR_OPEN
		break

	case STOPPED_OBSTRUCTION:
		state = STOPPED
		break

	default:
		break
	}
}

func EvNewOrderInEmptyQueue(floorButton int){
	switch state{
	case IDLE:
		if (floorButton > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
		}else if (floorButton < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
		}
		state = MOVING
		break

	case STOPPED:
		if (floorButton > queue.GetCurrentFloor()){
			queue.SetDirectionElevator(1)
			driver.MoveUp()
		}else if (floorButton < queue.GetCurrentFloor()){
			queue.SetDirectionElevator(-1)
			driver.MoveDown()
		}else if (floorButton == queue.GetCurrentFloor()){
			fmt.Printf("WARNING:EvNewOrderInEmptyQueue() was called when there was a \nnew order in the same floor as the elevator. \nAction Catching up by reversing motor.\n")
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
	switch state{
	case IDLE:
		driver.SetDoorLight()
		ResetTimer()
		state = DOOR_OPEN
		break

	case DOOR_OPEN:
		ResetTimer()
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
	if (state == INIT){
		fmt.Printf("Nåværende tilstand: INIT\n")
	}
	if (state == IDLE){
		fmt.Printf("Nåværende tilstand: IDLE\n")
	}
	if (state == DOOR_OPEN){
		fmt.Printf("Nåværende tilstand: DOOR_OPEN\n")
	}
	if (state == DOOR_OBSTRUCTED){
		fmt.Printf("Nåværende tilstand: DOOR_OBSTRUCTED\n")
	}
	if (state == MOVING){
		fmt.Printf("Nåværende tilstand: MOVING\n")
	}
	if (state == STOPPED){
		fmt.Printf("Nåværende tilstand: STOPPED\n")
	}
	if (state == OBSTRUCTION){
		fmt.Printf("Nåværende tilstand: OBSTRUCTION\n")
	}
	if (state == STOPPED_OBSTRUCTION){
		fmt.Printf("Nåværende tilstand: STOPPED_OBSTRUCTION\n")
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