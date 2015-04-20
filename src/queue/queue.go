package queue

import (
//	"fmt"
//	"driver"
//	"communication"
)

const elevatorNumber int = 1
const numberOfFloors int = 4
const numberOfElevators int = 3

var task int = -1

// Must be synchronized
var floorElevator = make([]int, numberOfElevators)
var directionElevator = make([]int, numberOfElevators)

var orderFloorUp = make([]bool, numberOfFloors)
var orderFloorDown = make([]bool, numberOfFloors)

// Not synchronized
var orderCommand = make([]bool, numberOfFloors)

func InitializeQueue(){
	for floor := 0; floor < numberOfFloors; floor++{
		orderFloorUp[floor] = false
		orderFloorDown[floor] = false
		orderCommand[floor] = false
	}
}

func AddOrder(typeOrder int, floorButton int){
	if (typeOrder == 0){
		orderFloorUp[floorButton] = true
	}else if (typeOrder == 1){
		orderFloorDown[floorButton] = true
	}else if (typeOrder == 2){
		orderCommand[floorButton] = true
	}
}

func CheckOrder(typeOrder int, floorButton int) bool{
	if (typeOrder == 0){
		return orderFloorUp[floorButton - 1]
	}else if (typeOrder == 1){
		return orderFloorDown[floorButton - 1]
	}else if (typeOrder == 2){
		return orderCommand[floorButton - 1]
	}else{
		return false
	}
}

func RemoveAllOrders() {
	InitializeQueue()
}

func RemoveOrder(typeOrder int, floorButton int){
	if (typeOrder == 0){
		orderFloorUp[floorButton] = false
	}else if (typeOrder == 1){
		orderFloorDown[floorButton] = false
	}
	orderCommand[floorButton] = false
}

func AssignNewTask(){

	thisFloor := floorElevator[elevatorNumber - 1]

	if (directionElevator[elevatorNumber - 1] == 1){
		if (thisFloor == numberOfFloors){
			for floor := numberOfFloors - 2; floor > -1; floor--{
				if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
					task = floor
					return
				}
			}
		}else{
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
					task = floor
					return
				}
			} 
		}
		if (thisFloor != 1){
			for floor := thisFloor - 2; floor < -1; floor--{
				if (orderFloorDown[floor] || orderFloorUp[floor] || orderCommand[floor]){
					task = floor
					return
				}
			}
		}
	}else if (directionElevator[elevatorNumber - 1] == -1){
		if (thisFloor == 1){
			for floor := 1; floor < numberOfFloors; floor++{
				if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
					task = floor
					return
				}
			}
		}else{
			for floor := thisFloor - 2; floor < -1; floor--{
				if (orderFloorDown[floor] || orderFloorUp[floor] || orderCommand[floor]){
					task = floor
					return
				}
			} 
		}
		if (thisFloor != numberOfFloors){
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
					task = floor
					return
				}
			}
		}	
	}
}

func GetAssignedTask() int{
	return task
}

func ShallStop() bool{

	thisFloor := floorElevator[elevatorNumber - 1]
	
	if (orderCommand[thisFloor - 1]){
		return true
	}
	
	if (directionElevator[elevatorNumber - 1] == 1){

		if (thisFloor == numberOfFloors){
			if (orderFloorDown[thisFloor - 1]){
				return true
			}
		}else{
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
					return false
				}
			}
			if (orderFloorDown[thisFloor - 1]){
			return true
			}
		}
	}else if (directionElevator[elevatorNumber - 1] == -1){

		if (thisFloor == 1){
			if (orderFloorUp[thisFloor - 1]){
				return true
			}
		}else{
			for floor := thisFloor - 2; floor > -1; floor--{
				if orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]{
					return false
				}
			}
			if (orderFloorUp[thisFloor -1]){
				return true
			}
		}
	}
	return false
}

func ShallRemoveOppositeFloorOrder() bool{
	// Assuming this function is not called on first and last floor!

	thisFloor := floorElevator[elevatorNumber - 1]

	if (directionElevator[elevatorNumber - 1] == 1){
		for floor := thisFloor; floor < numberOfFloors; floor++{
			if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
				return false
			}
		}
		if (orderFloorDown[thisFloor - 1]){
			return true
		}

	}else if (directionElevator[elevatorNumber - 1] == -1){
		for floor := thisFloor - 2; floor > -1; floor--{
			if (orderFloorUp[floor] || orderFloorDown[floor] || orderCommand[floor]){
				return false
			} 
		}
		if (orderFloorUp[thisFloor - 1]){
			return true
		}
	}
	return false	
}

func SetCurrentFloor(floor int){
	floorElevator[elevatorNumber - 1] = floor
}

func GetCurrentFloor() int{
	return floorElevator[elevatorNumber - 1]
}

func SetDirectionElevator(dir int){
	directionElevator[elevatorNumber - 1] = dir
}

func GetDirectionElevator() int{
	return directionElevator[elevatorNumber - 1]
}

func GetElevatorNumber() int{
	return elevatorNumber
}

func GetNumberOfFloors() int{
	return numberOfFloors
}

func GetNumberOfElevators() int{
	return numberOfElevators
}
