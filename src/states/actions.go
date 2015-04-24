package states

import (
	"fmt"
	"driver"
	"queue"
	"communication"
	"time"
)

// Run as goroutines
func HandleRemoteCalls(){
	for temp := range communication.ENOEQChan{
		EvNewOrderInEmptyQueue(temp.Floor)
		fmt.Println("states: EvNewOrderInEmptyQueue() called from HandleRemoteCalls()")
		queue.AddOrder(temp.Dir, temp.Floor)
		fmt.Printf("states: Motor remote started from IDLE because this elevator was best fit to take order to floor %v\n", temp.Floor)
		if temp.Dir == 0{
			driver.SetButtonLight(0, temp.Floor)
			communication.NotifyTheOthers("OU", temp.Floor, true, 0)	
		}else if temp.Dir == 1{
			driver.SetButtonLight(1, temp.Floor)
			communication.NotifyTheOthers("OD", temp.Floor, true, 0)
		}else{
			fmt.Println("states: ERROR: HandleRemoteCalls() identified invalid temp.Dir on ENOEQChan! Consequence: NotifyTheOthers() not called")
			//r := fmt.Errorf("states: ERROR: HandleRemoteCalls() identified invalid temp.Dir on ENOEQChan! Consequence: NotifyTheOthers() not called")
			//fmt.Println(r)
			//return r
		}
		time.Sleep(3 * time.Second) //Remeber this! Can be set to lower value if needed
	}
}