package communication

import(
	"queue"
//	"../src/network"
//	"queue"
)

var elevIpAddresses = make([]string, queue.GetNumberOfElevators())
var elevPorts = make([]int, queue.GetNumberOfElevators())

var receiveChan = make (chan Tcp_message)
var sendChan = make (chan Tcp_message)

var orderUpChan = make(chan map[int]bool)
var orderDownChan = make(chan map[int]bool)
var floorChan = make(chan map[int]bool)
var dirChan = make(chan map[int]int)

func Initialize() bool{
	
	elevIpAddresses = []string{"129.241.187.143", "129.241.187.141", "129.241.187.146"}
	elevPorts = []int{20005, 20004, 20006}



	TCPServerInit(elevPorts[queue.GetElevatorNumber() - 1], sendChan, receiveChan)

	return true
}

// Run as goroutines

/*
func SynchronizeOrders(){
	thisElevator := queue.GetElevatorNumber()
	for{
		select{
		case floor := <- orderUpChan:
			queue.OrderFloorUp[floor - 1] = true
		case floor := <- orderDownChan:
			queue.OrderFloorDown[floor - 1] = floor
		}
	}
}

func SynchronizeElevatorStatuses(){
	thisElevator := queue.GetElevatorNumber()
	for{
		select{
		case floor := <- floorChan:
			queue.FloorElevator[thisElevator - 1] = floor
		case direction := <- dirChan:
			queue.DirectionElevator[thisElevator - 1] = dir
		}
	}
}
*/