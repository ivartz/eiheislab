package driver

import "fmt"

const N_FLOORS = 4
const N_BUTTONS = 3

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

var buttonChan chan Button
var floorChan chan int
var motorChan chan MotorDirection
var stopChan chan bool
var obstrChan chan bool

func ElevInit() bool{

	// Init hardware
	init = Io_init()
	if (!init){
		fmt.Printl("Driver not initialized")
		return false
	}

	// ToDo:
	// Zero all floor button lamps
	for i := 1; i <= N_FLOORS; ++i {
		if i != 1{
			ClearButtonLight(i, BUTTON_CALL_DOWN)
		}
		if i != N_FLOORS{
			ClearButtonLight(i, BUTTON_CALL_UP)
		}
		ClearButtonLight(i, BUTTON_COMMAND)
	}

	// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
	ClearStopButton()
	ClearDoorLight()
	SetFloorLight(1)

	return true
}


func MoveUp(){
	motorChan <- MOVE_UP
}

func MoveDown(){
	motorChan <- MOVE_DOWN
}

func Stop(){
	motorChan <- MOVE_STOP
}

func MotorControl(chanDir chan){
	dir := <- chanDir
	if dir == MOVE_UP{
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
		//MoveUp()
	}
	if dir == MOVE_DOWN{
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
		//MoveDown()
	}
	if dir == MOVE_STOP{
		Io_write_analog(MOTOR, 0)
		//Stop()
	}

}

func GetMotorChan() chan MotorDirection{
	return motorChan
}


func SetButtonLight(floor int, dir OrderType){
	hardware = LocalizeHardware("light", floor, dir)
	Io_set_bit(hardware)

/*
	if (floor == 1 && dir == BUTTON_COMMAND)
		Io_set_bit(LIGHT_COMMAND1)
	if (floor == 2 && dir == BUTTON_COMMAND)
		Io_set_bit(LIGHT_COMMAND2)
	if (floor == 3 && dir == BUTTON_COMMAND)
		Io_set_bit(LIGHT_COMMAND3)
	if (floor == 4 && dir == BUTTON_COMMAND)
		Io_set_bit(LIGHT_COMMAND4)

	if (floor == 1 && dir == BUTTON_CALL_UP)
		Io_set_bit(LIGHT_UP1)
	if (floor == 2 && dir == BUTTON_CALL_UP)
		Io_set_bit(LIGHT_UP2)
	if (floor == 3 && dir == BUTTON_CALL_UP)
		Io_set_bit(LIGHT_UP3)
	
	if (floor == 2 && dir == BUTTON_CALL_DOWN)
		Io_set_bit(LIGHT_DOWN2)
	if (floor == 3 && dir == BUTTON_CALL_DOWN)
		Io_set_bit(LIGHT_DOWN3)
	if (floor == 4 && dir == BUTTON_CALL_DOWN)
		Io_set_bit(LIGHT_DOWN4)
		*/
}

func ClearButtonLight(floor int, dir OrderType){
	hardware = LocalizeHardware("light", floor, dir)
	Io_clear_bit(hardware)
/*
	if (floor == 1 && dir == BUTTON_COMMAND)
		Io_clear_bit(LIGHT_COMMAND1)
	if (floor == 2 && dir == BUTTON_COMMAND)
		Io_clear_bit(LIGHT_COMMAND2)
	if (floor == 3 && dir == BUTTON_COMMAND)
		Io_clear_bit(LIGHT_COMMAND3)
	if (floor == 4 && dir == BUTTON_COMMAND)
		Io_clear_bit(LIGHT_COMMAND4)

	if (floor == 1 && dir == BUTTON_CALL_UP)
		Io_clear_bit(LIGHT_UP1)
	if (floor == 2 && dir == BUTTON_CALL_UP)
		Io_clear_bit(LIGHT_UP2)
	if (floor == 3 && dir == BUTTON_CALL_UP)
		Io_clear_bit(LIGHT_UP3)
	
	if (floor == 2 && dir == BUTTON_CALL_DOWN)
		Io_clear_bit(LIGHT_DOWN2)
	if (floor == 3 && dir == BUTTON_CALL_DOWN)
		Io_clear_bit(LIGHT_DOWN3)
	if (floor == 4 && dir == BUTTON_CALL_DOWN)
		Io_clear_bit(LIGHT_DOWN4)
		*/
}

func CheckButton(floor int, dir OrderType) bool{
	hardware = LocalizeHardware("button", floor, dir)
	if Io_read_bit(hardware){
		var action Button
		action.Floor = floor
		action.Dir = dir
		buttonChan <- action
		return true
	}
	else{
		return false
	}
}

func LocalizeHardware(typeof string, floor int, dir OrderType) const int{
	if typeof == "button"{
		if (floor == 1 && dir == BUTTON_COMMAND)
			hardware = BUTTON_COMMAND1
		if (floor == 2 && dir == BUTTON_COMMAND)
			hardware = BUTTON_COMMAND2
		if (floor == 3 && dir == BUTTON_COMMAND)
			hardware = BUTTON_COMMAND3
		if (floor == 4 && dir == BUTTON_COMMAND)
			hardware = BUTTON_COMMAND4

		if (floor == 1 && dir == BUTTON_CALL_UP)
			hardware = BUTTON_UP1
		if (floor == 2 && dir == BUTTON_CALL_UP)
			hardware = BUTTON_UP2
		if (floor == 3 && dir == BUTTON_CALL_UP)
			hardware = BUTTON_UP3
		
		if (floor == 2 && dir == BUTTON_CALL_DOWN)
			hardware = BUTTON_DOWN2
		if (floor == 3 && dir == BUTTON_CALL_DOWN)
			hardware = BUTTON_DOWN3
		if (floor == 4 && dir == BUTTON_CALL_DOWN)
			hardware = BUTTON_DOWN4
	}

	else if typeof == "light"{
		if (floor == 1 && dir == BUTTON_COMMAND)
			hardware = LIGHT_COMMAND1
		if (floor == 2 && dir == BUTTON_COMMAND)
			hardware = LIGHT_COMMAND2
		if (floor == 3 && dir == BUTTON_COMMAND)
			hardware = LIGHT_COMMAND3
		if (floor == 4 && dir == BUTTON_COMMAND)
			hardware = LIGHT_COMMAND4

		if (floor == 1 && dir == BUTTON_CALL_UP)
			hardware = LIGHT_UP1
		if (floor == 2 && dir == BUTTON_CALL_UP)
			hardware = LIGHT_UP2
		if (floor == 3 && dir == BUTTON_CALL_UP)
			hardware = LIGHT_UP3
		
		if (floor == 2 && dir == BUTTON_CALL_DOWN)
			hardware = LIGHT_DOWN2
		if (floor == 3 && dir == BUTTON_CALL_DOWN)
			hardware = LIGHT_DOWN3
		if (floor == 4 && dir == BUTTON_CALL_DOWN)
			hardware = LIGHT_DOWN4
	}
	return hardware
}

func GetButton() (int, OrderType){
	button := <- buttonChan
	return button.Floor, button.Dir
}

func GetButtonChan() chan Button{
	return buttonChan
}

func GetFloorSensorSignal() int{
	if Io_read_bit(SENSOR_FLOOR1){
		return 1
	}
	else if Io_read_bit(SENSOR_FLOOR2){
		return 2
	}
	else if Io_read_bit(SENSOR_FLOOR3){
		return 3
	}
	else if Io_read_bit(SENSOR_FLOOR4){
		return 4
	}
	else{
		return -1
	}
}

func SetFloorLight(floor int){
	switch floor{
	case 1:
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
	case 2:
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
	case 3:
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
	case 4:
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
	}
}

func GetFloorLight() int {
	floor := <-floorChan
	return floor
}

func GetFloorChan() chan int{
	return floorChan
}

func SetStopButtonLight(){
	Io_set_bit(LIGHT_STOP)
}

func ClearStopButtonLight(){
	Io_clear_bit(LIGHT_STOP)
}

func CheckStopButton() bool{
	if Io_read_bit(STOP){
		stopChan <- true
		return true
	}
	else{
		return false
	}
}

func GetStopButton() bool{
	stopButton := <- stopChan
	return stopButton
}

func GetStopChan() chan bool{
	return stopChan
}

func CheckObstruction()bool{
	if Io_read_bit(OBSTRUCTION){
		obstrChan <- true
		return true
	}
	else{
		return false
	}
}

func GetObstruction(){
	obstr := <- obstrChan
	return obstr
}

func GetObstructionChan() chan bool{
	return obstrChan
}

func SetDoorLight(){
	Io_set_bit(LIGHT_DOOR_OPEN)
}

func ClearDoorLight(){
	Io_clear_bit(LIGHT_DOOR_OPEN)
}

func CheckDoorLight() bool{
	return Io_read_bit(LIGHT_DOOR_OPEN)
}