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

// Syncs with floor inducator set or not set
var floorSet int

// Syncs with elevator stop button on or off
var elevatorStopButton int

func evFloorReached(f int){
	fmt.Println("**Floor reached**")
	switch state{
	case INIT:
		RemoveAllOrders()
		setFloorIndicator(f)
		Stop()
		state = IDLE
		SetCurrentFloor(f)
		break

	case MOVING:
		setFloorIndicator(f)
		SetCurrentFloor(f)
				
	}
}

