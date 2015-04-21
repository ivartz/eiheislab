package states

import(
	"fmt"
	"time"
)

//var timeOut bool = true
//var timerStartTime time.Time
//var threeSeconds time.Duration = 3 * time.Second

var tick = make (chan bool)

var timeOut = make(chan bool)

/*
func SetTimeOut(){
	timeOut = true
}
*/

func Timer(){
	
}

func ResetTimer(){
	timer := time.NewTimer(3 * time.Second)
	fmt.Println("**states: Timer reset/started**")
	<- timer.C
	timeOut <- true
//	timeOut = false
}

func CheckTimeOut() bool{
	//fmt.Println("states: CheckTimeOut")
	select{
	case <- timeOut:
		fmt.Printf("states: Timeout = true\n")
		return true
	default:
		//fmt.Printf("states: Timeout = false\n")
		return false
	}
	//return false
}

	/*
	if (!timeOut){
		if (time.Since(timerStartTime) == threeSeconds){
			timeOut = true
			return timeOut
		}
	}else if (timeOut){
		return timeOut
	}
	fmt.Printf("states: Timeout = %v\n", timeOut)
	return false
}
*/
/*
func ClearTimeOut(){
	timeOut = false
}
*/
func PrintCurrentTime(){
	fmt.Printf("Time: %v\n", time.Now())
}

func Clock(){
	for{
		tick <- true
		time.Sleep(500 * time.Millisecond)
		tick <- false
		time.Sleep(500 * time.Millisecond)
	}
}

func ClockTick() bool{
	tickvar := <- tick
	return tickvar
}
