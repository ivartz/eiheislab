package driver

import "fmt"

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

var directionMotor MotorDirection

func Initialize(nFloors int) bool{

	init := IoInit()
	if (init == 0){
		fmt.Println("driver: IO Driver not initialized")
		return false
	}
	ClearAllOrderLights(nFloors)
	ClearStopButtonLight()
	ClearDoorLight()
	SetFloorLight(1)

	return true
}

func MoveUp(){
	directionMotor = MOVE_UP
	fmt.Println("driver: MoveUp(): MoveUp")
	IoClearBit(MOTORDIR)
	IoWriteAnalog(MOTOR, 2800)
}

func MoveDown(){
	directionMotor = MOVE_DOWN
	fmt.Println("driver: MoveDown(): MoveDown")
	IoSetBit(MOTORDIR)
	IoWriteAnalog(MOTOR, 2800)
}

func Stop(){
	directionMotor = MOVE_STOP
	fmt.Println("driver: Stop(): Stop")
	IoWriteAnalog(MOTOR, 0)
}

func SetButtonLight(dir OrderType, floor int){
	//var hardware OrderType
	hardware := LocalizeHardware("light", floor, dir)
	IoSetBit(int(hardware))
}

func ClearButtonLight(dir OrderType, floor int){
	//var hardware OrderType
	hardware := LocalizeHardware("light", floor, dir)
	IoClearBit(int(hardware))
}

func CheckButton(t OrderType, floor int) bool{
	//var hardware OrderType
	/*if t == BUTTON_RELEASED{
		return false
	}*/ 
	hardware := LocalizeHardware("button", floor, t)
	if IoReadBit(int(hardware)) != 0{
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

func SetStopButtonLight(){
	IoSetBit(LIGHT_STOP)
}

func ClearStopButtonLight(){
	IoClearBit(LIGHT_STOP)
}

func CheckStopButton() bool{
	if IoReadBit(STOP) != 0{
		return true
	}else{
		return false
	}
}

func CheckObstruction()bool{
	if IoReadBit(OBSTRUCTION) != 0{
		return true
	}else{
		return false
	}
}

func SetDoorLight(){
	IoSetBit(LIGHT_DOOR_OPEN)
}

func ClearDoorLight(){
	IoClearBit(LIGHT_DOOR_OPEN)
}

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