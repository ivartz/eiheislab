package queue

import (
	"fmt"
)

// Unique for each elevator
const elevatorNumber int = 1

//

const numberOfElevators int = 3

const numberOfFloors int = 4

// Must be synchronized
var FloorElevator = make([]int, numberOfElevators)
var DirectionElevator = make([]int, numberOfElevators)
var TaskElevator = make([]int, numberOfElevators) 

var OrderFloorUp = make([]bool, numberOfFloors)
var OrderFloorDown = make([]bool, numberOfFloors)

// Not synchronized
var orderCommand = make([]bool, numberOfFloors)

//var task int = -1

func Initialize() bool{
	InitializeQueue()
	InitializeFloorsAndDirectionsAndTasks()

	fmt.Println("******************************************************************************************")
	fmt.Printf("queue: Elevator #: %v\n", GetElevatorNumber())
	fmt.Printf("queue: # floors: %v\n", GetNumberOfFloors())
	fmt.Printf("queue: # elevators: %v\n\n", GetNumberOfElevators())
	return true
}

func InitializeQueue(){
	for floor := 0; floor < numberOfFloors; floor++{
		OrderFloorUp[floor] = false
		OrderFloorDown[floor] = false
		orderCommand[floor] = false
	}
}

func InitializeFloorsAndDirectionsAndTasks(){
	for i := range FloorElevator{
		FloorElevator[i] = -1
		DirectionElevator[i] = 0
		TaskElevator[i] = -1
	}
}

func AddOrder(typeOrder int, floorButton int){
	fmt.Printf("queue: AddOrder(): Called with typeOrder = %v and floorButton = %v\n", typeOrder, floorButton)
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
	fmt.Println("queue: RemoveOrder(): Called")
	if (typeOrder == 0){
		OrderFloorUp[floorButton - 1] = false
	}else if (typeOrder == 1){
		OrderFloorDown[floorButton - 1] = false
	}
	orderCommand[floorButton - 1] = false
}

func AssignNewTask() (int, int, int){
	// Assigns new task to an elevator by updating TaskElevator and returning task, buttonType, bestFitElevator
	task := -1
	buttonType := 2
	bestFitElevator := -1

	closestDistanceFromElevatorToOrder := numberOfFloors
	
	for index := range orderCommand{
		if (orderCommand[index]){
			//closest := IsClosest(index + 1)

			dist := FloorElevator[GetElevatorNumber() - 1] - (index + 1)
			if dist < 0{
				dist = -dist
			}
			if dist < closestDistanceFromElevatorToOrder{
				closestDistanceFromElevatorToOrder = dist
				task = index + 1
				buttonType = 2
				bestFitElevator = GetElevatorNumber()
			}
		}
	}

	if bestFitElevator != -1{
		TaskElevator[bestFitElevator - 1] = task
		return task, buttonType, bestFitElevator
	}

	for index := range OrderFloorUp{

		if (OrderFloorUp[index]){
			closest := IsClosest(index + 1)
			if FloorElevator[closest - 1] != index + 1{
				dist := FloorElevator[closest - 1] - (index + 1)
				if dist < 0{
					dist = -dist
				}
				if dist < closestDistanceFromElevatorToOrder{
					closestDistanceFromElevatorToOrder = dist
					task = index + 1
					buttonType = 0
					bestFitElevator = closest
				}
			}
		}else if (OrderFloorDown[index]){
			closest := IsClosest(index + 1)
			if FloorElevator[closest - 1] != index + 1{
				dist := FloorElevator[closest - 1] - (index + 1)
				if dist < 0{
					dist = -dist
				}
				if dist < closestDistanceFromElevatorToOrder{
					closestDistanceFromElevatorToOrder = dist
					task = index + 1
					buttonType = 1
					bestFitElevator = closest
				}	
			}	
		}

	}

	if bestFitElevator != -1{
		TaskElevator[bestFitElevator - 1] = task
	}
	

	return task, buttonType, bestFitElevator
}

func GetAssignedTask() int{
	return TaskElevator[elevatorNumber - 1]
}

func ClearAssignedTask(){
	TaskElevator[elevatorNumber - 1] = -1
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

func IsEmpty() bool{
	for index := 0; index < numberOfFloors; index++{
		if (OrderFloorUp[index] || OrderFloorDown[index] || orderCommand[index] ){
			return false
		}
	}
	return true
}

func IsClosest(floor int) int{
	//fmt.Println("queue: IsClosest(): Called")
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
	//fmt.Printf("queue: IsClosest(): Elevator %v is closest to floor %v\n", elev, floor)
	return elev
}

func PrintQueue(){
	fmt.Println("****************************************")
	fmt.Println("F | C\t\t| FUP\t| FDOWN\t")
	for floor := numberOfFloors - 1; floor > -1; floor--{
		fmt.Printf("%v | %v\t| %v\t| %v\n", floor + 1, orderCommand[floor], OrderFloorUp[floor], OrderFloorDown[floor])
	}
	fmt.Println("****************************************")
}