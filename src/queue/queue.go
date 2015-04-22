package queue

import (
	"fmt"
//	"driver"
//	"communication"
)

const elevatorNumber int = 1
const numberOfFloors int = 4
const numberOfElevators int = 3

var task int = -1

// Must be synchronized
var FloorElevator = make([]int, numberOfElevators)
var DirectionElevator = make([]int, numberOfElevators)

var OrderFloorUp = make([]bool, numberOfFloors)
var OrderFloorDown = make([]bool, numberOfFloors)

// Not synchronized
var orderCommand = make([]bool, numberOfFloors)

func InitializeQueue(){
	for floor := 0; floor < numberOfFloors; floor++{
		OrderFloorUp[floor] = false
		OrderFloorDown[floor] = false
		orderCommand[floor] = false
	}
}

func AddOrder(typeOrder int, floorButton int){
	if (typeOrder == 0){
		OrderFloorUp[floorButton - 1] = true
	}else if (typeOrder == 1){
		OrderFloorDown[floorButton - 1] = true
	}else if (typeOrder == 2){
		orderCommand[floorButton -1] = true
	}
}

func CheckOrder(typeOrder int, floorButton int) bool{
	if (typeOrder == 0){
		return OrderFloorUp[floorButton - 1]
	}else if (typeOrder == 1){
		return OrderFloorDown[floorButton - 1]
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
		OrderFloorUp[floorButton - 1] = false
	}else if (typeOrder == 1){
		OrderFloorDown[floorButton - 1] = false
	}
	orderCommand[floorButton - 1] = false
}

func AssignNewTask(){

	thisFloor := GetCurrentFloor()

	fmt.Println("queue: AssignNewTask()")

	if (GetDirectionElevator() == 1){
		if (thisFloor == numberOfFloors){
			for floor := numberOfFloors - 2; floor > -1; floor--{
				if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			}
		}else{
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			} 
		}
		if (thisFloor != 1){
			for floor := thisFloor - 2; floor > -1; floor--{
				if (OrderFloorDown[floor] || OrderFloorUp[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			}
		}
	}else if (GetDirectionElevator() == -1){
		if (thisFloor == 1){
			for floor := 1; floor < numberOfFloors; floor++{
				if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			}
		}else{
			for floor := thisFloor - 2; floor > -1; floor--{
				if (OrderFloorDown[floor] || OrderFloorUp[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			} 
		}
		if (thisFloor != numberOfFloors){
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
					task = floor + 1
					return
				}
			}
		}	
	}
	// No task was found, queue is empty
	for floor := 0; floor < numberOfFloors; floor++{
		if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
			task = floor + 1
			return
		}
	}
	fmt.Println("queue: AssignNewTask(), no order in queue")
	task = -1
}

func GetAssignedTask() int{
	return task
}

func ShallStop() bool{

	thisFloor := GetCurrentFloor()
	
	if (orderCommand[thisFloor - 1]){
		return true
	}
	if (GetDirectionElevator() == 1){
		if OrderFloorUp[thisFloor - 1]{
			return true
		}
		if (thisFloor == numberOfFloors){
			if (OrderFloorDown[thisFloor - 1]){
				return true
			}
		}else{
			for floor := thisFloor; floor < numberOfFloors; floor++{
				if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
					return false
				}
			}
			if (OrderFloorDown[thisFloor - 1]){
			return true
			}
		}

	}else if (GetDirectionElevator() == -1){
		if OrderFloorDown[thisFloor - 1]{
			return true
		}
		if (thisFloor == 1){
			if (OrderFloorUp[thisFloor - 1]){
				return true
			}
		}else{
			for floor := thisFloor - 2; floor > -1; floor--{
				if OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]{
					return false
				}
			}
			if (OrderFloorUp[thisFloor -1]){
				return true
			}
		}
	}
	return false
}

func ShallRemoveOppositeFloorOrder() bool{
	// Assuming this function is not called on first and last floor!

	thisFloor := FloorElevator[elevatorNumber - 1]

	if (DirectionElevator[elevatorNumber - 1] == 1){
		for floor := thisFloor; floor < numberOfFloors; floor++{
			if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
				return false
			}
		}
		if (OrderFloorDown[thisFloor - 1]){
			return true
		}

	}else if (DirectionElevator[elevatorNumber - 1] == -1){
		for floor := thisFloor - 2; floor > -1; floor--{
			if (OrderFloorUp[floor] || OrderFloorDown[floor] || orderCommand[floor]){
				return false
			} 
		}
		if (OrderFloorUp[thisFloor - 1]){
			return true
		}
	}
	return false	
}

func SetCurrentFloor(floor int){
	FloorElevator[elevatorNumber - 1] = floor
}

func GetCurrentFloor() int{
	return FloorElevator[elevatorNumber - 1]
}

func SetDirectionElevator(dir int){
	DirectionElevator[elevatorNumber - 1] = dir
}

func GetDirectionElevator() int{
	return DirectionElevator[elevatorNumber - 1]
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

func PrintQueue(){
	fmt.Println("F | C\t\t| FUP\t| FDOWN\t")
	for floor := numberOfFloors - 1; floor > -1; floor--{
		fmt.Printf("%v | %v\t| %v\t| %v\n", floor + 1, orderCommand[floor], OrderFloorUp[floor], OrderFloorDown[floor])
	}
}

func IsEmpty() bool{
	for index := 0; index < numberOfFloors; index++{
		if (OrderFloorUp[index] || OrderFloorDown[index] || orderCommand[index] ){
			return false
		}
	}
	return true
}