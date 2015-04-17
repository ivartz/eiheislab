package states

import(
	"fmt"
	"time"
)

var timeOut bool = true
type timerStartTime time.Time
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
		if (time.Now() - timerStartTime == 3){
			timeOut = true
			return timeOut
		}
	}
	else if (timeOut){
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
