package states

import (
	"fmt"
	"driver"
	"queue"
	"communication"
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
	fmt.Printf("states: EvFloorReached(): Floor %v reached\n", f)
	
	switch state{
	case INIT:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		SetFloorIndicator(f)
		driver.Stop()
		state = IDLE
		queue.SetCurrentFloor(f)
		communication.NotifyTheOthers("F", f, false, 0)
		break

	case MOVING:
		SetFloorIndicator(f)
		queue.SetCurrentFloor(f)
		communication.NotifyTheOthers("F", f, false, 0)
		if (queue.ShallStop()){
			fmt.Println("states: EvFloorReached(): ShallStop() returned true")
			driver.Stop()

			RemoveCorrectOrdersClearLightsSetDirectionAndNotifyTheOthers(f)
			/*
			if (f > 1 && f < queue.GetNumberOfFloors()){
				if (queue.GetDirectionElevator() == 1 && queue.CheckOrder(0, f) || queue.CheckOrder(2, f)){
					queue.RemoveOrder(0, f)
					driver.ClearButtonLight(0, f)
					communication.NotifyTheOthers("OU", f, false, 0)
				}else if (queue.GetDirectionElevator() == -1) && queue.CheckOrder(1, f) || queue.CheckOrder(2, f){
					queue.RemoveOrder(1, f)
					driver.ClearButtonLight(1, f)
					communication.NotifyTheOthers("OD", f, false, 0)
				}
				driver.ClearButtonLight(2,f)

				if (queue.ShallRemoveOppositeFloorOrder()){
					if (queue.GetDirectionElevator() == 1 && queue.CheckOrder(1, f) || queue.CheckOrder(2, f)){
						queue.RemoveOrder(1, f)
						driver.ClearButtonLight(1, f)
						communication.NotifyTheOthers("OD", f, false, 0)
						driver.ClearButtonLight(2, f)
					}else if (queue.GetDirectionElevator() == -1 && queue.CheckOrder(0, f) || queue.CheckOrder(2, f)){
						queue.RemoveOrder(0, f)
						driver.ClearButtonLight(0, f)
						communication.NotifyTheOthers("OU", f, false, 0)
						driver.ClearButtonLight(2, f)						
					}
				}
			}else if (f == 1){
				queue.RemoveOrder(0, f)
				driver.ClearButtonLight(0, f)
				communication.NotifyTheOthers("OU", f, false, 0)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				// Also so that ShallStop() will d
				queue.SetDirectionElevator(1)
				communication.NotifyTheOthers("D", 0, false, 1)
			}else if (f == queue.GetNumberOfFloors()){
				queue.RemoveOrder(1, f)
				driver.ClearButtonLight(1, f)
				communication.NotifyTheOthers("OD", f, false, 0)
				driver.ClearButtonLight(2, f)
				// Changing direction so that AssignNewTask() quickly can find the best task for the this elevator
				queue.SetDirectionElevator(-1)
				communication.NotifyTheOthers("D", 0, false, -1)
			}
			*/

			driver.SetDoorLight()
			fmt.Println("states: EvFloorReached(): Door opening")
			fmt.Println("states: EvFloorReached(): ResetTimer() called")
			go ResetTimer()
			state = DOOR_OPEN
		}else if (f == 1){
			// Changing direction
			queue.SetDirectionElevator(1)
			communication.NotifyTheOthers("D", 0, false, 1)
			driver.MoveUp()						
		}else if (f == queue.GetNumberOfFloors()){
			// Changing direction
			queue.SetDirectionElevator(-1)
			communication.NotifyTheOthers("D", 0, false, -1)
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
		fmt.Println("states: EvTimerOut(): Door closed. Calling AssignNewTask()")
		
		queue.AssignNewTask()
		tsk := queue.GetAssignedTask()

		fmt.Printf("states: EvTimerOut(): Called AssignNewTask() and got task: %v\n", tsk)
		//fmt.Printf("states: EvTimerOut(): Current floor is: %v\n", queue.GetCurrentFloor())
		

		if (tsk != -1){
			//state = MOVING
			MoveInDirectionFloorAndNotifyTheOthers(tsk)
			communication.NotifyTheOthers("T", tsk, false, 0)
			state = MOVING
		}else if (tsk == -1){
			communication.NotifyTheOthers("T", tsk, false, 0)
			state = IDLE
		}
		break

	case DOOR_OBSTRUCTED:
		break

	default:
		break
	}
}

func EvNewOrderInEmptyQueue(floorButton int){
	fmt.Println("states: EvNewOrderInEmptyQueue() was called")
	switch state{
	case IDLE:
		MoveInDirectionFloorAndNotifyTheOthers(floorButton)
		state = MOVING
		break

	case STOPPED:
		MoveInDirectionFloorAndNotifyTheOthers(floorButton)
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
	//PrintState()
	switch state{
	case IDLE:
		driver.SetDoorLight()
		fmt.Println("states: EvNewOrderInCurrentFloor(): Calling ResetTimer() from IDLE!")
		go ResetTimer()
		state = DOOR_OPEN
		break

	case DOOR_OPEN:
		fmt.Println("states: EvNewOrderInCurrentFloor(): Calling ResetTimer() from DOOR_OPEN!")
		
		select{
		case quitResetTimer <- true:
		default: 
		}
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



func EvStopButton(){
	switch state{
	case MOVING:
		driver.Stop()
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
		driver.Stop()
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
		tsk := queue.GetAssignedTask()
		if (tsk == -1){
			//fmt.Println("states: EvObstructionOff() returns IDLE here")
			communication.NotifyTheOthers("T", tsk, false, 0)
			state = IDLE
			break
		}else{
			MoveInDirectionFloorAndNotifyTheOthers(tsk)
			communication.NotifyTheOthers("T", tsk, false, 0)
		}
		driver.ClearDoorLight()
		state = MOVING
		break

	case DOOR_OBSTRUCTED:
		go ResetTimer()
		fmt.Println("states: EvObstructionOff(): ResetTimer() called")
		state = DOOR_OPEN
		break

	case STOPPED_OBSTRUCTION:
		state = STOPPED
		break

	default:
		break
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
		fmt.Printf("states: PrintState(): Current state is: INIT\n")
	}
	if (state == IDLE){
		fmt.Printf("states: PrintState(): Current state is: IDLE\n")
	}
	if (state == DOOR_OPEN){
		fmt.Printf("states: PrintState(): Current state is: DOOR_OPEN\n")
	}
	if (state == DOOR_OBSTRUCTED){
		fmt.Printf("states: PrintState(): Current state is: DOOR_OBSTRUCTED\n")
	}
	if (state == MOVING){
		fmt.Printf("states: PrintState(): Current state is: MOVING\n")
	}
	if (state == STOPPED){
		fmt.Printf("states: PrintState(): Current state is: STOPPED\n")
	}
	if (state == OBSTRUCTION){
		fmt.Printf("states: PrintState(): Current state is: OBSTRUCTION\n")
	}
	if (state == STOPPED_OBSTRUCTION){
		fmt.Printf("states: PrintState(): Current state is: STOPPED_OBSTRUCTION\n")
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
