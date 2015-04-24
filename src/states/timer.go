package states

import(
	"fmt"
	"time"
)

var tick = make (chan bool)
var timeOut = make(chan bool)
var quitResetTimer = make(chan bool)

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
	select{
	case <- timeOut:
		fmt.Printf("states: Timeout!\n")
		return true
	default:

		return false
	}
}

func PrintCurrentTime(){
	fmt.Printf("states: Time: %v\n", time.Now())
}

func Clock(){
	for{
		tick <- true
		time.Sleep(100 * time.Millisecond)
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
