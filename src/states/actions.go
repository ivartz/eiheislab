package states

import (
	"fmt"
	"driver"
	"queue"
	"communication"
	"time"
)

var orderFloorUpChan = make (chan int)
var orderFloorDownChan = make (chan int)
var orderCommandChan = make (chan int)

var previousFloorUpButtonState = make ([]bool, queue.GetNumberOfFloors())
var previousFloorDownButtonState = make ([]bool, queue.GetNumberOfFloors())
var previousCommandButtonState = make ([]bool, queue.GetNumberOfFloors())

func Initialize() bool{
	driver.MoveDown()
	queue.SetDirectionElevator(-1)
	InitializePreviousButtonStateSlices()	

	return true
}

func InitializePreviousButtonStateSlices(){
	for index := range previousFloorUpButtonState{
		previousFloorUpButtonState[index] = false
		previousFloorDownButtonState[index] = false
		previousCommandButtonState[index] = false
	}
}

func updatePreviousButtonState(t driver.OrderType, f int, s bool){
	if t == 0{
		previousFloorUpButtonState[f - 1] = s
	}else if t == 1{
		previousFloorDownButtonState[f - 1] = s
	}else if t == 2{
		previousCommandButtonState[f - 1] = s
	}else{
		fmt.Println("states: UpdatePreviousButtonState(): ERROR: invalid OrderType argument!")
	}
}

func getPreviousButtonState(t driver.OrderType, f int) bool{
	if t == 0{
		return previousFloorUpButtonState[f - 1]
	}else if t == 1{
		return previousFloorDownButtonState[f - 1]
	}else if t == 2{
		return previousCommandButtonState[f - 1]
	}else{
		fmt.Println("states: getPreviousButtonState(): ERROR: invalid OrderType argument!")
		return false
	}
}

func MoveInDirectionFloorAndNotifyTheOthers(floorButton int){
	if (floorButton > queue.GetCurrentFloor()){
		queue.SetDirectionElevator(1)
		communication.NotifyTheOthers("D", 0, false, 1)
		fmt.Println("states: MoveInDirectionFloorAndNotifyTheOthers(): Calling MoveUp()")
		driver.MoveUp()
	}else if (floorButton < queue.GetCurrentFloor()){
		queue.SetDirectionElevator(-1)
		communication.NotifyTheOthers("D", 0, false, -1)
		fmt.Println("states: MoveInDirectionFloorAndNotifyTheOthers(): Calling MoveDown()")
		driver.MoveDown()
	}else if (floorButton == queue.GetCurrentFloor()){ // To fix that the elevator can go back to most recently passed floor after a sudden stop (between two floors)
		if (queue.GetDirectionElevator() == 1){
			queue.SetDirectionElevator(-1)
			communication.NotifyTheOthers("D", 0, false, -1)
			fmt.Println("states: MoveInDirectionFloorAndNotifyTheOthers(): Calling MoveDown() from top floor")
			driver.MoveDown()
		}else if (queue.GetDirectionElevator() == -1){
			queue.SetDirectionElevator(1)
			communication.NotifyTheOthers("D", 0, false, 1)
			fmt.Println("states: MoveInDirectionFloorAndNotifyTheOthers(): Calling MoveUp() from bottom floor")
			driver.MoveUp()
		}
	}
}

func RemoveCorrectOrdersClearLightsSetDirectionAndNotifyTheOthers(f int){
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
}

// Run as goroutines
func CheckOrderChansAndCallEvents(){
	for{
		select{
		case floor := <- orderFloorUpChan:
			//fmt.Println("states: CheckOrderChansAndCallEvents(): RECEIVED FROM orderFloorUpChan")
			if !elevatorStopButton{
				if (!queue.CheckOrder(0, floor) && floor != driver.GetFloorSensorSignal()){
					queueIsEmpty := queue.IsEmpty()
					closestElev := queue.IsClosest(floor)
					if (queueIsEmpty && queue.GetElevatorNumber() == closestElev){
						EvNewOrderInEmptyQueue(floor)					
						queue.AddOrder(0, floor)
						driver.SetButtonLight(0, floor)
						fmt.Printf("states: CheckOrderChansAndCallEvents(): This elevator (%v) was closest and took the order, calling NotifyTheOthers()\n", queue.GetElevatorNumber())
						communication.NotifyTheOthers("OU", floor, true, 0)
					}else if (queueIsEmpty && queue.GetElevatorNumber() != closestElev){
						fmt.Println("states: CheckOrderChansAndCallEvents(): Elevator was not closest, calling NotifyTheOthers()")
						communication.NotifyTheOthers("ROU", floor, false, closestElev) // Should cause EvNewOrderInEmptyQueue(), AddOrder() SetButtonLight() and NotifyTheOthers() to be called on closest (best) remote elevator 
					}else{
						queue.AddOrder(0, floor)
						driver.SetButtonLight(0, floor)
						fmt.Println("states: CheckOrderChansAndCallEvents(): Added order, queue should be non-empty, calling NotifyTheOthers()")
						communication.NotifyTheOthers("OU", floor, true, 0)
					}
				}else if floor == driver.GetFloorSensorSignal(){
					EvNewOrderInCurrentFloor()
					//fmt.Println("states: EvNewOrderInCurrentFloor() called")
				} 
			}
		case floor := <- orderFloorDownChan:
			//fmt.Println("states: CheckOrderChansAndCallEvents(): RECEIVED FROM orderFloorDownChan")
			if !elevatorStopButton{
				if (!queue.CheckOrder(1, floor) && floor != driver.GetFloorSensorSignal()){
					queueIsEmpty := queue.IsEmpty()
					closestElev := queue.IsClosest(floor)
					if (queueIsEmpty && queue.GetElevatorNumber() == closestElev){
						EvNewOrderInEmptyQueue(floor)					
						queue.AddOrder(1, floor)
						driver.SetButtonLight(1, floor)
						fmt.Println("states: CheckOrderChansAndCallEvents(): Elevator was closest and took the order, calling NotifyTheOthers()")
						communication.NotifyTheOthers("OD", floor, true, 0)
					}else if (queueIsEmpty && queue.GetElevatorNumber() != closestElev){
						fmt.Println("states: CheckOrderChansAndCallEvents(): Elevator was not closest, calling NotifyTheOthers()")
						communication.NotifyTheOthers("ROD", floor, false, closestElev)
					}else{
						queue.AddOrder(1, floor)
						driver.SetButtonLight(1, floor)
						fmt.Println("states: CheckOrderChansAndCallEvents(): Added order, queue should be non-empty, calling NotifyTheOthers()")
						communication.NotifyTheOthers("OD", floor, true, 0)
					}
				}else if floor == driver.GetFloorSensorSignal(){
					EvNewOrderInCurrentFloor()
					//fmt.Println("states: EvNewOrderInCurrentFloor() called")
				}
			}			
		case floor := <- orderCommandChan:
			//fmt.Println("states: CheckOrderChansAndCallEvents(): RECEIVED FROM orderCommandChan")
			if (!queue.CheckOrder(2, floor) && floor != driver.GetFloorSensorSignal()){
				if queue.IsEmpty(){
					if (elevatorStopButton){
						EvStopButtonOff()
					}
					EvNewOrderInEmptyQueue(floor)
					queue.AddOrder(2, floor)
					driver.SetButtonLight(2, floor)	
				}else{
					queue.AddOrder(2, floor)
					driver.SetButtonLight(2, floor)
				}
			}else if floor == driver.GetFloorSensorSignal(){
				EvStopButtonOff()
				EvNewOrderInCurrentFloor()
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func CheckOrderButtonsAndSendToOrderChannels(){
	for{
		for floor := 1; floor <= queue.GetNumberOfFloors(); floor++{
			if (floor > 1 && floor < queue.GetNumberOfFloors()){
				if (driver.CheckButton(0, floor) != getPreviousButtonState(0, floor)){
					if driver.CheckButton(0, floor){
						select{
							case orderFloorUpChan <- floor:
						default:
							fmt.Println("states: CheckOrderButtonsAndSendToOrderChannels(): ERROR: orderFloorUpChan is BLOCKED!!")
						}						
					}
					updatePreviousButtonState(0, floor, driver.CheckButton(0, floor))

				}else if (driver.CheckButton(1, floor) != getPreviousButtonState(1, floor)){
					if driver.CheckButton(1, floor){
						select{
							case orderFloorDownChan <- floor:
						default:
							fmt.Println("states: CheckOrderButtonsAndSendToOrderChannels(): ERROR: orderFloorDownChan is BLOCKED!!")
						}					
					}
					updatePreviousButtonState(1, floor, driver.CheckButton(1, floor))
				}
			}
			if (floor == 1){
				if (driver.CheckButton(0, floor) != getPreviousButtonState(0, floor)){
					if driver.CheckButton(0, floor){
						select{
							case orderFloorUpChan <- floor:
						default:
							fmt.Println("states: CheckOrderButtonsAndSendToOrderChannels(): ERROR: orderFloorUpChan is BLOCKED!!")
						}	
					}
					updatePreviousButtonState(0, floor, driver.CheckButton(0, floor))
				}
			}
			if (floor == queue.GetNumberOfFloors()){
				if (driver.CheckButton(1, floor) != getPreviousButtonState(1, floor)){
					if driver.CheckButton(1, floor){
						select{
							case orderFloorDownChan <- floor:
						default:
							fmt.Println("states: CheckOrderButtonsAndSendToOrderChannels(): ERROR: orderFloorDownChan is BLOCKED!!")
						}						
					}
					updatePreviousButtonState(1, floor, driver.CheckButton(1, floor))
				}
			}
			if (driver.CheckButton(2,floor) != getPreviousButtonState(2, floor)){
				if driver.CheckButton(2, floor){
					select{
						case orderCommandChan <- floor:
					default:
						fmt.Println("states: CheckOrderButtonsAndSendToOrderChannels(): ERROR: orderCommandChan is BLOCKED!!")
					}		
				}
				updatePreviousButtonState(2, floor, driver.CheckButton(2, floor))
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func CheckFloorSensorAndCallEvents(){
	// Check if floor reached and call EvFloorReached() once
	// -1 if a floor is not reached. If floor reached: 1-4. Belongs to CheckFloorSensorAndCallEvents() 
	reached := -1	
	for{
		if (driver.GetFloorSensorSignal() != reached){
			if (reached == -1){
				reached = driver.GetFloorSensorSignal()
				EvFloorReached(reached)
			}else if (reached != -1){
				reached = driver.GetFloorSensorSignal()
			}
		}
		if (driver.GetFloorSensorSignal() != reached){
			fmt.Printf("\nstates: Floor sensor says: %v\n", driver.GetFloorSensorSignal())
			PrintState()
		}
		time.Sleep(10 * time.Millisecond)		
	}
}

func CheckIfTimeoutCallEventAndPrintQueue(){
	// Time out signal check
	for{
		if (CheckTimeOut() && !driver.CheckObstruction()){
			EvTimerOut()
			//queue.PrintQueue()
		}
		time.Sleep(800 * time.Millisecond)
	}
}

func CheckRemoteChanAndCallEvents(){
	for temp := range communication.RemoteChan{
		if (temp.Floor != driver.GetFloorSensorSignal()){
			EvNewOrderInEmptyQueue(temp.Floor)
			fmt.Printf("states: CheckRemoteChanAndCallEvents(): EvNewOrderInEmptyQueue() called with floor %v\n", temp.Floor)
			queue.AddOrder(temp.Dir, temp.Floor)
			fmt.Printf("states: CheckRemoteChanAndCallEvents(): Motor remote started from IDLE because this elevator was best fit to take order to floor %v\n", temp.Floor)
			if temp.Dir == 0{
				driver.SetButtonLight(0, temp.Floor)
				communication.NotifyTheOthers("OU", temp.Floor, true, 0)	
			}else if temp.Dir == 1{
				driver.SetButtonLight(1, temp.Floor)
				communication.NotifyTheOthers("OD", temp.Floor, true, 0)
			}else{
				fmt.Println("states: ERROR: CheckRemoteChanAndCallEvents() identified invalid temp.Dir on RemoteChan! Consequence: NotifyTheOthers() not called")
				//r := fmt.Errorf("states: ERROR: CheckRemoteChanAndCallEvents() identified invalid temp.Dir on RemoteChan! Consequence: NotifyTheOthers() not called")
				//fmt.Println(r)
				//return r
			}
		}else if (temp.Floor == driver.GetFloorSensorSignal()){
			EvNewOrderInCurrentFloor()
		}

		//time.Sleep(3 * time.Second) //Remeber this! Can be set to lower value if needed
	}
}

func CheckStopAndObstructionAndCallEvents(){
	for{
		if (driver.CheckStopButton() && !elevatorStopButton){
			EvStopButton()
		}
		if (driver.CheckObstruction() && !obstruction){
			SetObstructionVariable()
			EvObstructionOn()
		}else if (!driver.CheckObstruction() && obstruction){
			ClearObstructionVariable()
			EvObstructionOff()
		}
		time.Sleep(100 * time.Millisecond)
	}
}