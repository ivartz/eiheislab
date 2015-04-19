package states

import(
	"fmt"
	"time"
)

var timeOut bool = true
var timerStartTime time.Time
var threeSeconds time.Duration = 3 * time.Second
/*
func SetTimeOut(){
	timeOut = true
}
*/

func ResetTimer(){
	timerStartTime = time.Now()
	fmt.Println("**Timer reset/started**")
	timeOut = false
}

func CheckTimeOut() bool{
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
