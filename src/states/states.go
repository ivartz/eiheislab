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
		RemoveAllOrders()
		SetFloorIndicator(f)
		StopElevatorGood()
		state = IDLE
		SetCurrentFloor(f)
		break

	case MOVING:
		SetFloorIndicator(f)
		SetCurrentFloor(f)
		if (ShallStop(f, GetElevatorDirection())){
			StopElevatorGood()

			if (f == 2 || f == 3){
				RemoveFloorOrder(f, GetElevatorDirection())
				if (ShallRemoveOppositeFloorOrder(f), GetElevatorDirection()){
					if (GetElevatorDirection() == 1){
						RemoveFloorOrder(f, -1)
					}
					else if (GetElevatorDirection() == -1){
						RemoveFloorOrder(f, 1)
					}
				}
			}
			else if (f == 1){
				RemoveFloorOrder(f, 1)
			}
			else if (f == 4){
				RemoveFloorOrder(f, -1)
			}
			SetDoorLight()
			fmt.Println("**Door open**")
			ResetTimer()
			state = DOOR_OPEN
		}
		else if (f == 1){
			SetElevatorDirection(1)
			MoveUp()						
		} 
		else if (f == N_FLOORS){
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
	SetTimeOut(false)
	switch state{
	case DOOR_OPEN:
		SetDoorLight()
		fmt.Println("**Timeout, door closed**")
		AssignNewTask()
		if (GetAssignedTask() > GetCurrentFloor()){
			SetElevatorDirection(1)
			MoveUp()
		}
		else if (GetAssignedTask() == -1){
			return
		}
		else if (GetAssignedTask() < GetCurrentFloor()){
			SetElevatorDirection(-1)
			MoveDown()
		}
		fmt.Printf("Assigned task is: %v\n", GetAssignedTask())
		if (GetAssignedTask() != -1){
			state = MOVING
		}
		else if (GetAssignedTask() == -1){
			state = IDLE
		}
		break

	case DOOR_OBSTRUCTED:
		break

	default:
		break
	}
}

func EvStop(){
	switch state{
	case MOVING:
		StopElevatorGood()
		RemoveAllOrders()
		SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED
		break

	case STOPPED:
		break

	case OBSTRUCTION:
		RemoveAllOrders()
		SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED

	case STOPPED_OBSTRUCTION:
		break

	case IDLE:
		RemoveAllOrders()
		SetStopButtonLight()
		SetElevatorStopButtonVariable()
		SetDoorLight()
		state = STOPPED
		break

	default:
		RemoveAllOrders()
		SetStopButtonLight()
		SetElevatorStopButtonVariable()
		state = STOPPED
		break
	}
}

func EvStopOff(){
	switch state{
	case STOPPED:
		ClearStopButtoonLight()
		ClearDoorLight()
		ClearElevatorStopButtonVariable()
		state = STOPPED // Fordi evStopOff() kalles av bestilling fra COMMAND_BUTTONS
		break

	case STOPPED_OBSTRUCTION:
		ClearStopButtoonLight()
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
		SetDoorLight()
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
		if (GetAssignedTask() == -1){
			AssignNewTask()
		}
		if (GetAssignedTask() == -1){
			state = IDLE
		}
		else if (GetElevatorDirection() == 1){
			if (GetCurrentFloor() == N_FLOORS){
				SetElevatorDirection(-1)
				MoveDown()
				state = MOVING
			}
			else{
				SetElevatorDirection(1)
				MoveUp()
				state = MOVING
			}
		}
		else if (GetElevatorDirection() == -1){
			if (GetCurrentFloor() == 1){
				SetElevatorDirection(1)
				MoveUp()
				state = MOVING
			}
			else{
				SetElevatorDirection(-1)
				MoveDown()
				state = MOVING
			}
		}
		ClearDoorLight()
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
		if (buttonFloor > GetCurrentFloor()){
			SetElevatorDirection(1)
			MoveUp()
		}
		else if (buttonFloor < GetCurrentFloor()){
			SetElevatorDirection(-1)
			MoveDown()
		}
		state = MOVING
		break

	case STOPPED:
		if (buttonFloor > GetCurrentFloor()){
			SetElevatorDirection(1)
			MoveUp()
		}
		else if (buttonFloor < GetCurrentFloor()){
			SetElevatorDirection(-1)
			MoveDown()
		}
		else if (buttonFloor == GetCurrentFloor()){
			fmt.Printf("WARNING:EvNewOrderInEmptyQueue() was called when there was a \nnew order in the same floor as the elevator. \nAction Catching up by reversing motor.\n")
			if (GetElevatorDirection() == 1){
				SetElevatorDirection(-1)
				MoveDown()
			}
			else if (GetElevatorDirection() == -1){
				SetElevatorDirection(1)
				MoveUp()
			}
		}
		state = MOVING
		break

	case OBSTRUCTION:
		break

	case INIT:
		RemoveAllOrders()
		break

	default:
		break
	}
}

func EvNewOrderInCurrentFloor(f int, buttonDirection int){
	switch state{
	case IDLE:
		SetDoorLight()
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
	if (GetElevatorDirection() == 1){
		MoveDown()
		time.Sleep(10000 * time.Millisecond)
		Stop()
	}
	else if (GetElevatorDirection() == -1){
		MoveUp()
		time.Sleep(10000 * time.Millisecond)
		Stop()
	}
}

func SetFloorIndicator(floor int){
	if (floor == -1){
		return
	}
	else if (floorIndicatorSet != floor){
		SetFloorLight(floor)
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