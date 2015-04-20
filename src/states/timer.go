package states

import(
	"fmt"
	"time"
)

var timeOut bool = true
var timerStartTime time.Time
var threeSeconds time.Duration = 3 * time.Second
var tick = make (chan bool)
/*
func SetTimeOut(){
	timeOut = true
}
*/

func ResetTimer(){
	timerStartTime = time.Now()
	fmt.Println("**states: Timer reset/started**")
	timeOut = false
}

func CheckTimeOut() bool{
	fmt.Println("states: CheckTimeOut")
	if (!timeOut){
		if (time.Since(timerStartTime) == threeSeconds){
			timeOut = true
			return timeOut
		}
	}else if (timeOut){
		return timeOut
	}
	return false
}
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
