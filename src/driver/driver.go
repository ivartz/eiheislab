package driver

import "fmt"

//const N_FLOORS int = 4
//const N_BUTTONS = 3

type MotorDirection int
type OrderType int

const MOVE_UP MotorDirection = 1
const MOVE_DOWN MotorDirection = -1
const MOVE_STOP MotorDirection = 0


type Button struct{
	Floor int
	Dir OrderType
}

const BUTTON_CALL_UP OrderType = 0
const BUTTON_CALL_DOWN OrderType = 1
const BUTTON_COMMAND OrderType = 2

//var buttonChan chan Button // Not needed
//var floorChan chan int
var motorChan = make(chan MotorDirection, 2) //Not currently functionality for sending 
										// and recieving on this channel on a time instant
										// Could use threads (goroutines)?
										// NOTE: Can possibly need a larger buffer! 
//var stopChan chan bool
//var obstrChan chan bool

func Initialize(nFloors int) bool{

	// Init hardware
	init := IoInit()
	if (init == 0){
		fmt.Println("Driver not initialized")
		return false
	}

	ClearAllOrderLights(nFloors)
	
	// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
	ClearStopButtonLight()
	ClearDoorLight()
	SetFloorLight(1)

	return true
}


func MoveUp(){
	motorChan <- MOVE_UP
	fmt.Println("driver: MoveUp")
}

func MoveDown(){
	motorChan <- MOVE_DOWN
	fmt.Println("driver: MoveDown")
}

func Stop(){
	motorChan <- MOVE_STOP
	fmt.Println("driver: Stop")
}

func MotorControl(){
	var dir MotorDirection
	if (dir <- motorChan){
		if dir == MOVE_UP{
			IoClearBit(MOTORDIR)
			IoWriteAnalog(MOTOR, 2800)
			//MoveUp()
		}else if dir == MOVE_DOWN{
			IoSetBit(MOTORDIR)
			IoWriteAnalog(MOTOR, 2800)
			//MoveDown()
		}else if dir == MOVE_STOP{
			IoWriteAnalog(MOTOR, 0)
			//Stop()
		}
	}
	else{
		fmt.Println("driver: Nothing on motorChan")
	}

}

func GetMotorChan() chan MotorDirection{
	return motorChan
}


func SetButtonLight(dir OrderType, floor int){
	//var hardware OrderType
	hardware := LocalizeHardware("light", floor, dir)
	IoSetBit(int(hardware))

/*
	if (floor == 1 && dir == BUTTON_COMMAND)
		IoSetBit(LIGHT_COMMAND1)
	if (floor == 2 && dir == BUTTON_COMMAND)
		IoSetBit(LIGHT_COMMAND2)
	if (floor == 3 && dir == BUTTON_COMMAND)
		IoSetBit(LIGHT_COMMAND3)
	if (floor == 4 && dir == BUTTON_COMMAND)
		IoSetBit(LIGHT_COMMAND4)

	if (floor == 1 && dir == BUTTON_CALL_UP)
		IoSetBit(LIGHT_UP1)
	if (floor == 2 && dir == BUTTON_CALL_UP)
		IoSetBit(LIGHT_UP2)
	if (floor == 3 && dir == BUTTON_CALL_UP)
		IoSetBit(LIGHT_UP3)
	
	if (floor == 2 && dir == BUTTON_CALL_DOWN)
		IoSetBit(LIGHT_DOWN2)
	if (floor == 3 && dir == BUTTON_CALL_DOWN)
		IoSetBit(LIGHT_DOWN3)
	if (floor == 4 && dir == BUTTON_CALL_DOWN)
		IoSetBit(LIGHT_DOWN4)
		*/
}

func ClearButtonLight(dir OrderType, floor int){
	//var hardware OrderType
	hardware := LocalizeHardware("light", floor, dir)
	IoClearBit(int(hardware))
/*
	if (floor == 1 && dir == BUTTON_COMMAND)
		IoClearBit(LIGHT_COMMAND1)
	if (floor == 2 && dir == BUTTON_COMMAND)
		IoClearBit(LIGHT_COMMAND2)
	if (floor == 3 && dir == BUTTON_COMMAND)
		IoClearBit(LIGHT_COMMAND3)
	if (floor == 4 && dir == BUTTON_COMMAND)
		IoClearBit(LIGHT_COMMAND4)

	if (floor == 1 && dir == BUTTON_CALL_UP)
		IoClearBit(LIGHT_UP1)
	if (floor == 2 && dir == BUTTON_CALL_UP)
		IoClearBit(LIGHT_UP2)
	if (floor == 3 && dir == BUTTON_CALL_UP)
		IoClearBit(LIGHT_UP3)
	
	if (floor == 2 && dir == BUTTON_CALL_DOWN)
		IoClearBit(LIGHT_DOWN2)
	if (floor == 3 && dir == BUTTON_CALL_DOWN)
		IoClearBit(LIGHT_DOWN3)
	if (floor == 4 && dir == BUTTON_CALL_DOWN)
		IoClearBit(LIGHT_DOWN4)
		*/
}

func CheckButton(dir OrderType, floor int) bool{
	//var hardware OrderType
	hardware := LocalizeHardware("button", floor, dir)
	if IoReadBit(int(hardware)) != 0{
		/*var action Button
		action.Floor = floor
		action.Dir = dir
		buttonChan <- action*/
		return true
	}else{
		return false
	}
}

func LocalizeHardware(typeof string, floor int, dir OrderType) OrderType{
	// Limited to max 4 floors
	var hardware OrderType
	if typeof == "button"{
		if (floor == 1 && dir == BUTTON_COMMAND){
			hardware = BUTTON_COMMAND1
		}
		if (floor == 2 && dir == BUTTON_COMMAND){
			hardware = BUTTON_COMMAND2
		}
		if (floor == 3 && dir == BUTTON_COMMAND){
			hardware = BUTTON_COMMAND3
		}
		if (floor == 4 && dir == BUTTON_COMMAND){
			hardware = BUTTON_COMMAND4
		}
		if (floor == 1 && dir == BUTTON_CALL_UP){
			hardware = BUTTON_UP1
		}
		if (floor == 2 && dir == BUTTON_CALL_UP){
			hardware = BUTTON_UP2
		}
		if (floor == 3 && dir == BUTTON_CALL_UP){
			hardware = BUTTON_UP3
		}
		if (floor == 2 && dir == BUTTON_CALL_DOWN){
			hardware = BUTTON_DOWN2
		}
		if (floor == 3 && dir == BUTTON_CALL_DOWN){
			hardware = BUTTON_DOWN3
		}
		if (floor == 4 && dir == BUTTON_CALL_DOWN){
			hardware = BUTTON_DOWN4
		}
	}else if typeof == "light"{
		if (floor == 1 && dir == BUTTON_COMMAND){
			hardware = LIGHT_COMMAND1
		}
		if (floor == 2 && dir == BUTTON_COMMAND){
			hardware = LIGHT_COMMAND2
		}
		if (floor == 3 && dir == BUTTON_COMMAND){
			hardware = LIGHT_COMMAND3
		}
		if (floor == 4 && dir == BUTTON_COMMAND){
			hardware = LIGHT_COMMAND4
		}

		if (floor == 1 && dir == BUTTON_CALL_UP){
			hardware = LIGHT_UP1
		}
		if (floor == 2 && dir == BUTTON_CALL_UP){
			hardware = LIGHT_UP2
		}
		if (floor == 3 && dir == BUTTON_CALL_UP){
			hardware = LIGHT_UP3
		}
		
		if (floor == 2 && dir == BUTTON_CALL_DOWN){
			hardware = LIGHT_DOWN2
		}
		if (floor == 3 && dir == BUTTON_CALL_DOWN){
			hardware = LIGHT_DOWN3
		}
		if (floor == 4 && dir == BUTTON_CALL_DOWN){
			hardware = LIGHT_DOWN4
		}
	}
	return hardware
}
/*
func GetButtonFromChan() (int, OrderType){
	button := <- buttonChan
	return button.Floor, button.Dir
}

func GetButtonChan() chan Button{
	return buttonChan
}
*/
func GetFloorSensorSignal() int{
	if IoReadBit(SENSOR_FLOOR1) != 0{
		return 1
	}else if IoReadBit(SENSOR_FLOOR2) != 0{
		return 2
	}else if IoReadBit(SENSOR_FLOOR3) != 0{
		return 3
	}else if IoReadBit(SENSOR_FLOOR4) != 0{
		return 4
	}else{
		return -1
	}
}

func SetFloorLight(floor int){
	switch floor{
	case 1:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 2:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	case 3:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 4:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	}
}
/*
func GetFloorLight() int {
	floor := <-floorChan
	return floor
}

func GetFloorChan() chan int{
	return floorChan
}
*/
func SetStopButtonLight(){
	IoSetBit(LIGHT_STOP)
}

func ClearStopButtonLight(){
	IoClearBit(LIGHT_STOP)
}

func CheckStopButton() bool{
	if IoReadBit(STOP) != 0{
		//stopChan <- true
		return true
	}else{
		return false
	}
}
/*
func GetStopButtonFromChan() bool{
	stopButton := <- stopChan
	return stopButton
}

func GetStopChan() chan bool{
	return stopChan
}
*/
func CheckObstruction()bool{
	if IoReadBit(OBSTRUCTION) != 0{
		//obstrChan <- true
		return true
	}else{
		return false
	}
}
/*
func GetObstructionFromChan(){
	obstr := <- obstrChan
	return obstrChan
}

func GetObstructionChan() chan bool{
	return obstrChan
}
*/
func SetDoorLight(){
	IoSetBit(LIGHT_DOOR_OPEN)
}

func ClearDoorLight(){
	IoClearBit(LIGHT_DOOR_OPEN)
}
/*
func CheckDoorLight() bool{
	if IoReadBit(LIGHT_DOOR_OPEN) != 0{
		return true
	}else{
		return false
}

func SetNFloors(floors int){
	N_FLOORS = floors
}

func GetNFloors() int{
	return N_FLOORS
}
*/
func ClearAllOrderLights(nFloors int){
	// Zero all floor button lamps
	for floor := 1; floor <= nFloors; floor++ {
		if floor != 1{
			ClearButtonLight(BUTTON_CALL_DOWN, floor)
		}
		if floor != nFloors{
			ClearButtonLight(BUTTON_CALL_UP, floor)
		}
		ClearButtonLight(BUTTON_COMMAND, floor)
	}
}