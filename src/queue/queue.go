package queue

import (
	"fmt"
//	"driver"
//	"communication"
)

// Unique for each elevator
const elevatorNumber int = 2

const numberOfFloors int = 4

//const numberOfElevators int = 3
const numberOfElevators int = 2

// Must be synchronized
var FloorElevator = make([]int, numberOfElevators)
var DirectionElevator = make([]int, numberOfElevators)

var OrderFloorUp = make([]bool, numberOfFloors)
var OrderFloorDown = make([]bool, numberOfFloors)

// Not synchronized
var orderCommand = make([]bool, numberOfFloors)

var task int = -1

func Initialize(){
	InitializeQueue()
	InitializeFloorsAndDirections()
}

func InitializeQueue(){
	for floor := 0; floor < numberOfFloors; floor++{
		OrderFloorUp[floor] = false
		OrderFloorDown[floor] = false
		orderCommand[floor] = false
	}
}

func InitializeFloorsAndDirections(){
	for i := range FloorElevator{
		FloorElevator[i] = -1
		DirectionElevator[i] = 0
	}
}

func AddOrder(typeOrder int, floorButton int){
	fmt.Println("queue: AddOrder() was called")
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
	fmt.Println("queue: RemoveOrder() was called")
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
	fmt.Println("queue: SetCurrentFloor() was called")
	FloorElevator[elevatorNumber - 1] = floor
}

func GetCurrentFloor() int{
	return FloorElevator[elevatorNumber - 1]
}

func SetDirectionElevator(dir int){
	fmt.Printf("queue: SetDirectionElevator(%v) was called\n", dir)
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

//NB! DENNE FUNGERER IKKE FULLSTENDIG UT!!
func IsClosest(floor int) int{
	fmt.Println("queue: IsClosest() was called")
	diff := numberOfFloors
	elev := elevatorNumber
	for index := range FloorElevator{
		if FloorElevator[index] != -1{
			temp := FloorElevator[index] - floor 
			if temp < 0{
				temp = -temp
			}
			if temp < diff{
				diff = temp
				elev = index + 1
			}		
		}
	}
	fmt.Printf("queue: IsClosest(): Elevator %v was closest\n", elev)
	return elev
}