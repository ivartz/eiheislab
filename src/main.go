//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
	"driver"
	"states"
	"queue"
	"communication"
	"time"

)
// elevatorNumber, numberOfFloors and numberOfElevators constants are set in ../src/queue/queue.go



func main(){

	if (!driver.Initialize(queue.GetNumberOfFloors())){
		fmt.Println("main: Unable to initialize hardware..")
	}

	elevIpAddresses := []string{"129.241.187.146", "129.241.187.141"}
	elevPorts := []int{20011, 20004}

	if (!communication.Initialize(elevIpAddresses, elevPorts)){
		fmt.Println("main: Unable to initialize network..")
	}
	
	queue.Initialize()

	driver.MoveDown()
	queue.SetDirectionElevator(-1)

	fmt.Println("****************************************")
	fmt.Printf("main: Elevator %v successfully initialized driver and TCP listening and send server on port: %v\n\n", queue.GetElevatorNumber(), elevPorts[queue.GetElevatorNumber() - 1])	
	fmt.Printf("main: Elevator #: %v\n", queue.GetElevatorNumber())
	fmt.Printf("main: # floors: %v\n", queue.GetNumberOfFloors())
	fmt.Printf("main: # elevators: %v\n", queue.GetNumberOfElevators())
	fmt.Printf("main: Current task in initialization: %v\n", queue.GetAssignedTask())
	//fmt.Printf("********for loop Go!********\n")

	go states.Clock()
	go HandleStopButton()
	go HandleFloorSensor()
	go HandleFloorButtons()
	go HandleCommandButtons()
	go states.HandleRemoteCalls()
	go HandleTimeOut()
	go HandleObstruction()

	select{
	}
}

type previousButton struct{
	Type driver.OrderType
	Floor int
}

// Called as goroutines
func HandleFloorButtons(){
	// Checking floor buttons and adding orders, setting button lights and calling events
	previous := previousButton{driver.BUTTON_PRESSED, 0}
	for{
		if (!states.CheckElevatorStopButtonVariable() && !driver.CheckButton(previous.Type, previous.Floor)){
			for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
				if (floor > 1 && floor < queue.GetNumberOfFloors()){
					if (driver.CheckButton(0, floor) && !queue.CheckOrder(0, floor) && floor != driver.GetFloorSensorSignal()){
						previous = previousButton{driver.BUTTON_CALL_UP, floor}
						if (queue.IsEmpty() && queue.GetElevatorNumber() == queue.IsClosest(floor)){
							states.EvNewOrderInEmptyQueue(floor)					
							queue.AddOrder(0, floor)
							driver.SetButtonLight(0, floor)
							fmt.Println("main: Elevator was closest and took the order, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OU", floor, true, 0)
						}else if (queue.IsEmpty() && queue.GetElevatorNumber() != queue.IsClosest(floor)){// && states.Tick()){
							fmt.Println("main: Elevator was not closest, calling NotifyTheOthers()")
							communication.NotifyTheOthers("ENOEQU", floor, false, queue.IsClosest(floor)) // Should cause EvNewOrderInEmptyQueue(), AddOrder() SetButtonLight() and NotifyTheOthers() to be called on closest (best) remote elevator 
						}else if (states.Tick()){
							queue.AddOrder(0, floor)
							driver.SetButtonLight(0, floor)
							fmt.Println("main: Added order, queue should be non-empty, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OU", floor, true, 0)
						}			
					}else if (driver.CheckButton(0, floor) && floor == driver.GetFloorSensorSignal()){//  && states.Tick()){
						previous = previousButton{driver.BUTTON_CALL_UP, floor}
						states.EvNewOrderInCurrentFloor()
						fmt.Println("main: EvNewOrderInCurrentFloor() called")
					}
					if (driver.CheckButton(1, floor) && !queue.CheckOrder(1, floor) && floor != driver.GetFloorSensorSignal()){
						previous = previousButton{driver.BUTTON_CALL_DOWN, floor}
						if (queue.IsEmpty() && queue.GetElevatorNumber() == queue.IsClosest(floor)){
							states.EvNewOrderInEmptyQueue(floor)					
							queue.AddOrder(1, floor)
							driver.SetButtonLight(1, floor)
							fmt.Println("main: Elevator was closest and took the order, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OD", floor, true, 0)
						}else if (queue.IsEmpty() && queue.GetElevatorNumber() != queue.IsClosest(floor)){//  && states.Tick()){
							fmt.Println("main: Elevator was not closest, calling NotifyTheOthers()")
							communication.NotifyTheOthers("ENOEQD", floor, false, queue.IsClosest(floor))
						}else if (states.Tick()){				
							queue.AddOrder(1, floor)
							driver.SetButtonLight(1, floor)
							fmt.Println("main: Added order, queue should be non-empty, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OD", floor, true, 0)
						}
					}else if (driver.CheckButton(1, floor) && floor == driver.GetFloorSensorSignal()){//  && states.Tick()){
						previous = previousButton{driver.BUTTON_CALL_DOWN, floor}
						states.EvNewOrderInCurrentFloor()
						fmt.Println("main: EvNewOrderInCurrentFloor() called")
					}
				}
				// Only one direction from floor 1 and GetNFloors()
				if (floor == 1){
					if (driver.CheckButton(0, floor) && !queue.CheckOrder(0, floor) && floor != driver.GetFloorSensorSignal()){
						previous = previousButton{driver.BUTTON_CALL_UP, floor}
						if (queue.IsEmpty() && queue.GetElevatorNumber() == queue.IsClosest(floor)){//  && states.Tick()){
							states.EvNewOrderInEmptyQueue(floor)
							queue.AddOrder(0, floor)
							driver.SetButtonLight(0, floor)
							fmt.Println("main: Elevator was closest and took the order, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OU", floor, true, 0)
						}else if (queue.IsEmpty() && queue.GetElevatorNumber() != queue.IsClosest(floor)){//  && states.Tick()){ //Adding Tick() here, because there is a little delay before the remote elevator has set the queue and synced it with the others elevator's queues, which in turn will block this NotifyTheOthers() from being called 
							fmt.Println("main: Elevator was not closest, calling NotifyTheOthers()")
							communication.NotifyTheOthers("ENOEQU", floor, false, queue.IsClosest(floor))	
						}else if (states.Tick()){
							queue.AddOrder(0, floor)
							driver.SetButtonLight(0, floor)
							fmt.Println("main: Added order, queue should be non-empty, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OU", floor, true, 0)
						}
					}else if (driver.CheckButton(0, floor) && floor == driver.GetFloorSensorSignal()){//  && states.Tick()){
						previous = previousButton{driver.BUTTON_CALL_UP, floor}
						states.EvNewOrderInCurrentFloor()
						fmt.Println("main: EvNewOrderInCurrentFloor() called")
					}
				}
				if (floor == queue.GetNumberOfFloors()){
					if (driver.CheckButton(1, floor) && !queue.CheckOrder(1, floor) && floor != driver.GetFloorSensorSignal()){
						previous = previousButton{driver.BUTTON_CALL_DOWN, floor}
						if (queue.IsEmpty() && queue.GetElevatorNumber() == queue.IsClosest(floor)){//  && states.Tick()){
							states.EvNewOrderInEmptyQueue(floor)
							queue.AddOrder(1, floor)
							driver.SetButtonLight(1, floor)
							fmt.Println("main: Elevator was closest and took the order, calling NotifyTheOthers()")	
							communication.NotifyTheOthers("OD", floor, true, 0)
						}else if (queue.IsEmpty() && queue.GetElevatorNumber() != queue.IsClosest(floor)){//  && states.Tick()){
							fmt.Println("main: Elevator was not closest, calling NotifyTheOthers()")
							communication.NotifyTheOthers("ENOEQD", floor, false, queue.IsClosest(floor))
						}else if (states.Tick()){
							queue.AddOrder(1, floor)
							driver.SetButtonLight(1, floor)
							fmt.Println("main: Added order, queue should be non-empty, calling NotifyTheOthers()")
							communication.NotifyTheOthers("OD", floor, true, 0)
						}
					}else if (driver.CheckButton(1, floor) && floor == driver.GetFloorSensorSignal()){//  && states.Tick()){
						previous = previousButton{driver.BUTTON_CALL_DOWN, floor}
						states.EvNewOrderInCurrentFloor()
						fmt.Println("main: EvNewOrderInCurrentFloor() called")
					}
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func HandleCommandButtons(){
	previous := previousButton{driver.BUTTON_PRESSED, 0}
	for{
		if (!driver.CheckButton(previous.Type, previous.Floor)){
			for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
				if (driver.CheckButton(2, floor) && !queue.CheckOrder(2, floor) && floor != driver.GetFloorSensorSignal()){
					previous = previousButton{driver.BUTTON_COMMAND, floor}
					if (queue.IsEmpty()){
						if (states.CheckElevatorStopButtonVariable()){
							states.EvStopButtonOff()
						}
						states.EvNewOrderInEmptyQueue(floor)
						queue.AddOrder(2, floor)
						driver.SetButtonLight(2, floor)
						//fmt.Println("main: EvNewOrderInEmptyQueue() was called from HandleCommandButtons()")
					}else{
						queue.AddOrder(2, floor)
						driver.SetButtonLight(2, floor)
						//fmt.Println("main: ************************ANOTHER command order was added to queue!!!")
					}
				}else if (driver.CheckButton(2, floor) && floor == driver.GetFloorSensorSignal()){// && states.Tick()){
					previous = previousButton{driver.BUTTON_COMMAND, floor}
					states.EvStopButtonOff()
					states.EvNewOrderInCurrentFloor()
					//fmt.Println("main: EvNewOrderInCurrentFloor() was called from HandleCommandButtons()")
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func HandleFloorSensor(){
	// Check if floor reached and call EvFloorReached() once
	// -1 if a floor is not reached. If floor reached: 1-4. Belongs to HandleFloorSensor() 
	reached := -1	
	for{
		if (driver.GetFloorSensorSignal() != reached){
			if (reached == -1){
				reached = driver.GetFloorSensorSignal()
				states.EvFloorReached(reached)
			}else if (reached != -1){
				reached = driver.GetFloorSensorSignal()
			}
		}
		if (driver.GetFloorSensorSignal() != reached){
			fmt.Printf("\nmain: Floor sensor says: %v\n", driver.GetFloorSensorSignal())
			states.PrintState()
		}
		time.Sleep(10 * time.Millisecond)		
	}
}

func HandleStopButton(){
	// Check if stop button is pressed, if so, stop elevator and remove all orders
	for{
		if (driver.CheckStopButton() && !states.CheckElevatorStopButtonVariable()){
			states.EvStopButton()
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func HandleTimeOut(){
	// Time out signal check
	for{
		if (states.CheckTimeOut() && !driver.CheckObstruction()){
			states.EvTimerOut()
		}
		time.Sleep(1 * time.Second)
	}
}

func HandleObstruction(){
	// Universal obstruction signal
	for{
		if (driver.CheckObstruction() && !states.CheckObstructionVariable()){
			states.SetObstructionVariable()
			states.EvObstructionOn()
		}else if (!driver.CheckObstruction() && states.CheckObstructionVariable()){
			states.ClearObstructionVariable()
			states.EvObstructionOff()
		}
		time.Sleep(500 * time.Millisecond)	
	}
}