package communication

import(
	"fmt"
	"strconv"
	"encoding/json"
	"queue"
	"driver"
	"time"
)

var elevIpAddresses = make([]string, queue.GetNumberOfElevators())
var elevPorts = make([]int, queue.GetNumberOfElevators())

// The message form is a struct. Capital letter because the struct is used outside network.go
// This struct is sent in send_ch and recived in receive_ch
type Tcp_message struct{
	Raddr string //Remote address, like "129.241.187.144:20012" is embedded in the message
	Data []byte //jsoned data
	Length int
}
var receiveChan = make (chan Tcp_message)
var sendChan = make (chan Tcp_message)

// unjsoned data
type msg struct{
	MType string
	
	ENumber int
	Floor int
	Set bool
	Dir int
}
var sendToAllOthersChan = make (chan msg)

// Struct and chan to handle remote calls when the queue is empty
type ENOEQmsg struct{
	Floor int
	Dir int
}
var ENOEQChan = make (chan ENOEQmsg)



func Initialize(eip []string, ep []int) bool{

	// Unique list of 	
//	elevIpAddresses = []string{"129.241.187.143", "129.241.187.141", "129.241.187.146"}
//	elevPorts = []int{20005, 20004, 20006}

	elevIpAddresses = eip
	elevPorts = ep

	TCPServerInit(elevPorts[queue.GetElevatorNumber() - 1], sendChan, receiveChan)

	go HandleOutgoingMessages()
	go HandleIncomingMessages()

	return true
}


func NotifyTheOthers(mtype string, floor int, set bool, dir int){
	fmt.Println("communication: NotifyTheOthers() was called")
	if (mtype == "OU" || mtype == "OD" || mtype == "F" || mtype == "D" || mtype == "ENOEQU" || mtype == "ENOEQD"){
		temp := msg{mtype, queue.GetElevatorNumber(), floor, set, dir}
		select{
		case sendToAllOthersChan <- temp:
		default:
			fmt.Println("communication: ERROR: ************************************NotifyTheOthers() can't send message because sendToAllOthersChan is BLOCKED!")	
		}
	}else{
		fmt.Println("communication: ERROR: Can't NotifyTheOthers(), invalid string argument")
	}
}


// Run as goroutine
func HandleOutgoingMessages() error{
	fmt.Println("communication: HandleOutgoingMessages() goroutine started")
	for temp := range sendToAllOthersChan{
		fmt.Println("communication: HandleOutgoingMessages(): New message to send to the other elevators!")
		jtemp, err := json.Marshal(temp)
		if err != nil{
			fmt.Println("communication: json.Marshal() error! HandleOutgoingMessages() goroutine ending")
			return err
		}
		for i := range elevIpAddresses{
			if i + 1 != queue.GetElevatorNumber(){
				tcpm := Tcp_message{elevIpAddresses[i]+":"+strconv.Itoa(elevPorts[i]), jtemp, len(jtemp)}
				//fmt.Printf("communication: HandleOutgoingMessages(): Message to elevator %v:\n               %v\n", i + 1, tcpm)
				select{
				case sendChan <- tcpm:
					fmt.Println("communication: HandleOutgoingMessages(): Tcp_message was sent into sendChan!")		
				default:
					fmt.Println("communication: ERROR: ******************************HandleOutgoingMessages() can't send Tcp_message into sendChan because sendChan is BLOCKED!")		
				}
				
			}
		}
		time.Sleep(10 * time.Millisecond)	
	}
	r := fmt.Errorf("communication: ERROR: HandleOutgoingMessages() has quit range over sendToAllOthersChan!")

	return r
}

func HandleIncomingMessages() error{
	// Updates the local elevators OU,OD,FE,DE arrays according to incoming messages and sets/clears lights
	fmt.Println("communication: HandleIncomingMessages() goroutine started")
	var m msg
	for temp := range receiveChan{
		err := json.Unmarshal(temp.Data[:temp.Length], &m)
		if err != nil{
			fmt.Println("communication: json.Unmarshal() error! HandleIncomingMessages() goroutine ending")
			fmt.Println(err)
			return err
		}else if m.MType == "OU"{
			queue.OrderFloorUp[m.Floor - 1] = m.Set
			
			if m.Set == true{
				driver.SetButtonLight(0, m.Floor)
			}else if m.Set == false{
				driver.ClearButtonLight(0, m.Floor)
			}else{
				r := fmt.Errorf("communication: HandleIncomingMessages() received and unjsoned a message with unknown Set. Something is wrong with HandleIncomingMessages()")
				return r
			}
			
			fmt.Printf("communication: Remote floor order up-button on floor %v set to %t\n", m.Floor, m.Set)
		}else if m.MType == "UD"{
			queue.OrderFloorDown[m.Floor - 1] = m.Set
			
			if m.Set == true{
				driver.SetButtonLight(1, m.Floor)
			}else if m.Set == false{
				driver.ClearButtonLight(1, m.Floor)
			}else{
				r := fmt.Errorf("communication: HandleIncomingMessages() received and unjsoned a message with unknown Set. Something is wrong with HandleIncomingMessages()")
				return r				
			}
			
			fmt.Printf("communication: Remote floor order down-button on floor %v set to %t\n", m.Floor, m.Set)
		
		}else if m.MType == "F"{
			queue.FloorElevator[m.ENumber - 1] = m.Floor
			fmt.Printf("communication: Remote elevator floor; elevator %v set its floor to %v\n", m.ENumber, m.Floor)
		}else if m.MType == "D"{
			queue.DirectionElevator[m.ENumber - 1] = m.Dir
			fmt.Printf("communication: Remote elevator direction; elevator %v set its direction to %v\n", m.ENumber, m.Dir)
		}else if (m.MType == "ENOEQU" && m.Dir == queue.GetElevatorNumber()){ //In this case, m.Dir is the best fit elevator 
			enoeqmsg := ENOEQmsg{m.Floor, 0}
			select{
			case ENOEQChan <- enoeqmsg:
			default:
				fmt.Println("communication: ENOEQChan blocked!")
			}
			//fmt.Printf("communication: This best fit elevator to take order to floor %v was remote started from IDLE\n", m.Floor)
		}else if (m.MType == "ENOEQD" && m.Dir == queue.GetElevatorNumber()){ //In this case, m.Dir is the best fit elevator
			enoeqmsg := ENOEQmsg{m.Floor, 1}
			select{
			case ENOEQChan <- enoeqmsg:
			default:
				fmt.Println("communication: ERROR: *****************************************ENOEQChan blocked!")
			}
			//fmt.Printf("communication: This best fit elevator to take order to floor %v was remote started from IDLE\n", m.Floor)
		}else{
			r := fmt.Errorf("communication: HandleIncomingMessages() received and unjsoned a message with unknown MType. Something is wrong with HandleIncomingMessages()")
			return r
		}
		time.Sleep(10 * time.Millisecond)
	}
	r := fmt.Errorf("communication: ERROR: HandleIncomingMessages() has quit range over receiveChan!")
	return r
}
