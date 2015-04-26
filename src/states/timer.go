package states

import(
	"fmt"
	"time"
)

var tick = make (chan bool)
var timeOut = make(chan bool)
var quitResetTimer = make(chan bool)

func ResetTimer(){
	fmt.Println("states: ResetTimer(): Timer reset/started")
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
		fmt.Printf("states: CheckTimeOut(): Timeout!\n")
		return true
	default:

		return false
	}
}