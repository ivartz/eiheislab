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
			queue.ClearAssignedTask()
			communication.NotifyTheOthers("T", -1, false, 0)
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
	default:
		break
	}
}

func EvTimerOut(){
	//ClearTimeOut()
	switch state{
	case DOOR_OPEN:
		driver.ClearDoorLight()

		fmt.Println("states: EvTimerOut(): Door closed. Calling AssignNewTask()")
		
		task, buttonType, bestFitElev := queue.AssignNewTask() // COST FUNCTION

		if bestFitElev != -1{
			fmt.Printf("states: EvTimerOut(): Called AssignNewTask(), and elevator %v got task: %v\n", bestFitElev, task)
		}
		if (bestFitElev == queue.GetElevatorNumber()){//(tsk != -1){
			MoveInDirectionFloorAndNotifyTheOthers(task)
			communication.NotifyTheOthers("T", task, false, 0)
			state = MOVING
		
		}else if (bestFitElev != -1){
			// NotifyTheOthers her for å gi ordren til riktig heis
			if buttonType == 0{
				communication.NotifyTheOthers("ROU", task, false, bestFitElev)
			}else if buttonType == 1{
				communication.NotifyTheOthers("ROD", task, false, bestFitElev)
			}else{
				fmt.Printf("states: EvTimerOut(): ERROR: AssignNewTask() assigned task to remote elevator %v, ordered with this elevators (%v) command buttons. Something is wrong with AssignNewTask()\n", bestFitElev, queue.GetElevatorNumber())
			}
			state = IDLE
		}else{
			fmt.Println("states: EvTimerOut(): AssignNewTask could not find a best fit elevator, entering IDLE")
			state = IDLE
		}
		break
	default:
		break
	}
}

func EvOrder(floorButton int){
	fmt.Println("states: EvOrder() was called")
	switch state{
	case IDLE:
		MoveInDirectionFloorAndNotifyTheOthers(floorButton)
		communication.NotifyTheOthers("T", floorButton, false, 0)
		state = MOVING
		break
	case STOPPED:
		MoveInDirectionFloorAndNotifyTheOthers(floorButton)
		communication.NotifyTheOthers("T", floorButton, false, 0)
		state = MOVING
		break
	case INIT:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		break
	default:
		break
	}
}

func EvOrderInCurrentFloor(){//f int, buttonDirection int){
	fmt.Println("states: EvOrderInCurrentFloor() was called")
	//PrintState()
	switch state{
	case IDLE:
		driver.SetDoorLight()
		fmt.Println("states: EvOrderInCurrentFloor(): Calling ResetTimer() from IDLE!")
		go ResetTimer()
		state = DOOR_OPEN
		break
	case DOOR_OPEN:
		fmt.Println("states: EvOrderInCurrentFloor(): Calling ResetTimer() from DOOR_OPEN!")
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
	}else if (state == IDLE){
		fmt.Printf("states: PrintState(): Current state is: IDLE\n")
	}else if (state == DOOR_OPEN){
		fmt.Printf("states: PrintState(): Current state is: DOOR_OPEN\n")
	}else if (state == DOOR_OBSTRUCTED){
		fmt.Printf("states: PrintState(): Current state is: DOOR_OBSTRUCTED\n")
	}else if (state == MOVING){
		fmt.Printf("states: PrintState(): Current state is: MOVING\n")
	}else if (state == STOPPED){
		fmt.Printf("states: PrintState(): Current state is: STOPPED\n")
	}else if (state == OBSTRUCTION){
			fmt.Printf("states: PrintState(): Current state is: OBSTRUCTION\n")
	}else if (state == STOPPED_OBSTRUCTION){
		fmt.Printf("states: PrintState(): Current state is: STOPPED_OBSTRUCTION\n")
	}	
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

func EvStopButton(){
	switch state{
	case MOVING:
		driver.Stop()
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		elevatorStopButton = true
		state = STOPPED
		break
	case STOPPED:
		break
	case OBSTRUCTION:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		elevatorStopButton = true
		state = STOPPED_OBSTRUCTION
	case STOPPED_OBSTRUCTION:
		break
	case DOOR_OBSTRUCTED:
		break
	case IDLE:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		elevatorStopButton = true
		driver.SetDoorLight()
		state = STOPPED
		break
	default:
		queue.RemoveAllOrders()
		driver.ClearAllOrderLights(queue.GetNumberOfFloors())
		driver.SetStopButtonLight()
		elevatorStopButton = true
		state = STOPPED
		break
	}
}

func EvStopButtonOff(){
	switch state{
	case STOPPED:
		driver.ClearStopButtonLight()
		driver.ClearDoorLight()
		elevatorStopButton = false
		//state = MOVING
		//state = STOPPED // Fordi evStopOff() kalles av bestilling fra COMMAND_BUTTONS
		break
	case STOPPED_OBSTRUCTION:
		driver.ClearStopButtonLight()
		elevatorStopButton = false
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
			//_, _, _ := queue.AssignNewTask()
			queue.AssignNewTask()
		}
		task := queue.GetAssignedTask()
		if (task == -1){
			//fmt.Println("states: EvObstructionOff() returns IDLE here")
			communication.NotifyTheOthers("T", task, false, 0)
			state = IDLE
			break
		}else{
			MoveInDirectionFloorAndNotifyTheOthers(task)
			communication.NotifyTheOthers("T", task, false, 0)
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