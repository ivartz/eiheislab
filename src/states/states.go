package states

import (
	"fmt"
	"time"
	"../src/queue"
	"../src/communication"
	"../src/driver"

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
var floorIndicatorSet int

// Syncs with elevator stop button on or off
var elevatorStopButton bool

func EvFloorReached(f int){
	fmt.Println("**Floor reached**")
	
	switch state{
	case INIT:
		queue.RemoveAllOrders()
		SetFloorIndicator(f)
		StopElevatorGood()
		state = IDLE
		communication.SetCurrentFloor(f)
		break

	case MOVING:
		SetFloorIndicator(f)
		communication.SetCurrentFloor(f)
		if (queue.ShallStop(f, communication.GetElevatorDirection())){
			StopElevatorGood()

			if (f == 2 || f == 3){
				queue.RemoveFloorOrder(f, communication.GetElevatorDirection())
				if (queue.ShallRemoveOppositeFloorOrder(f), communication.GetElevatorDirection()){
					if (communication.GetElevatorDirection() == 1){
						queue.RemoveFloorOrder(f, -1)
					}
					else if (communication.GetElevatorDirection() == -1){
						queue.RemoveFloorOrder(f, 1)
					}
				}
			}
			else if (f == 1){
				queue.RemoveFloorOrder(f, 1)
			}
			else if (f == 4){
				queue.RemoveFloorOrder(f, -1)
			}
			driver.SetDoorLight()
			fmt.Println("**Door open**")
			ResetTimer()
			state = DOOR_OPEN
		}
		else if (f == 1){
			SetElevatorDirection(1)
			MoveUp()						
		} 
		else if (f == driver.GetNFloors()){
			SetElevatorDirection(-1)
			MoveDown()
		}
		else{
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
		driver.SetDoorLight()
		fmt.Println("**Timeout, door closed**")
		queue.AssignNewTask()
		if (queue.GetAssignedTask() > communication.GetCurrentFloor()){
			communication.SetElevatorDirection(1)
			driver.MoveUp()
		}
		else if (queue.GetAssignedTask() == -1){
			return
		}
		else if (queue.GetAssignedTask() < communication.GetCurrentFloor()){
			communication.SetElevatorDirection(-1)
			driver.MoveDown()
		}
		fmt.Printf("Assigned task is: %v\n", queue.GetAssignedTask())
		if (queue.GetAssignedTask() != -1){
			state = MOVING
		}
		else if (queue.GetAssignedTask() == -1){
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
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED
		break

	case STOPPED:
		break

	case OBSTRUCTION:
		queue.RemoveAllOrders()
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED

	case STOPPED_OBSTRUCTION:
		break

	case IDLE:
		queue.RemoveAllOrders()
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		driver.SetDoorLight()
		state = STOPPED
		break

	default:
		queue.RemoveAllOrders()
		driver.SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED
		break
	}
}

func EvStopButtonOff(){
	switch state{
	case STOPPED:
		driver.ClearStopButtoonLight()
		driver.ClearDoorLight()
		ClearElevatorStopButtonVariable()
		state = STOPPED // Fordi evStopOff() kalles av bestilling fra COMMAND_BUTTONS
		break

	case STOPPED_OBSTRUCTION:
		driver.ClearStopButtoonLight()
		ClearElevatorStopButtonVariable()
		state = OBSTRUCTIOIN
		break

	case default:
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
		}
		else if (communication.GetElevatorDirection() == 1){
			if (communication.GetCurrentFloor() == driver.GetNFloors()){
				communication.SetElevatorDirection(-1)
				driver.MoveDown()
				state = MOVING
			}
			else{
				communication.SetElevatorDirection(1)
				driver.MoveUp()
				state = MOVING
			}
		}
		else if (communication.GetElevatorDirection() == -1){
			if (communication.GetCurrentFloor() == 1){
				communication.SetElevatorDirection(1)
				driver.MoveUp()
				state = MOVING
			}
			else{
				communication.SetElevatorDirection(-1)
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

func EvNewOrderInEmptyQueue(buttonFloor int){
	switch state{
	case IDLE:
		if (buttonFloor > communcation.GetCurrentFloor()){
			communication.SetElevatorDirection(1)
			driver.MoveUp()
		}
		else if (buttonFloor < communication.GetCurrentFloor()){
			communication.SetElevatorDirection(-1)
			driver.MoveDown()
		}
		state = MOVING
		break

	case STOPPED:
		if (buttonFloor > communication.GetCurrentFloor()){
			communication.SetElevatorDirection(1)
			driver.MoveUp()
		}
		else if (buttonFloor < communication.GetCurrentFloor()){
			communication.SetElevatorDirection(-1)
			driver.MoveDown()
		}
		else if (buttonFloor == communication.GetCurrentFloor()){
			fmt.Printf("WARNING:EvNewOrderInEmptyQueue() was called when there was a \nnew order in the same floor as the elevator. \nAction Catching up by reversing motor.\n")
			if (communication.GetElevatorDirection() == 1){
				communication.SetElevatorDirection(-1)
				driver.MoveDown()
			}
			else if (communication.GetElevatorDirection() == -1){
				communication.SetElevatorDirection(1)
				driver.MoveUp()
			}
		}
		state = MOVING
		break

	case OBSTRUCTION:
		break

	case INIT:
		queue.RemoveAllOrders()
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
	if (communication.GetElevatorDirection() == 1){
		driver.MoveDown()
		time.Sleep(10000 * time.Microsecond)
		driver.Stop()
	}
	else if (communication.GetElevatorDirection() == -1){
		driver.MoveUp()
		time.Sleep(10000 * time.Microsecond)
		driver.Stop()
	}
}

func SetFloorIndicator(floor int){
	if (floor == -1){
		return
	}
	else if (floorIndicatorSet != floor){
		driver.SetFloorLight(floor)
		floorIndicatorSet = floor
	}
}

func PrintState(){
	if (state == INIT){
		fmt.Printf("Nåværende tilstand: INIT\n")
	}
	if (state == IDLE){
		printf("Nåværende tilstand: IDLE\n")
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