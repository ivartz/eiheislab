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

var quitResetTimer = make(chan bool)

/*
func SetTimeOut(){
	timeOut = true
}
*/

func Timer(){
	
}

func ResetTimer(){
	fmt.Println("states: Timer reset/started**")
	timer := time.NewTimer(3 * time.Second)
	select{
	case <- timer.C:
		timeOut <- true
	case <- quitResetTimer:
		return
	}
}

func CheckTimeOut() bool{
	//fmt.Println("states: CheckTimeOut")
	select{
	case <- timeOut:
		fmt.Printf("states: Timeout!\n")
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
	fmt.Printf("states: Time: %v\n", time.Now())
}

func Clock(){
	for{
		tick <- true
		time.Sleep(500 * time.Millisecond)
		//tick <- false
		//time.Sleep(200 * time.Millisecond)
	}
}

func Tick() bool{
	select{
	case <- tick:
		return true
	default:
		return false
	}
}
