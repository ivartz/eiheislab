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

// This struct is sent in send_ch and recived in receive_ch
type Tcp_message struct{
	Raddr string //Remote address, like "129.241.187.144:20012" is embedded in the message
	Data []byte //jsoned data
	Length int
}
var receiveChan = make (chan Tcp_message, 5)
var sendChan = make (chan Tcp_message, (11 * (queue.GetNumberOfElevators() - 1)) - 1) // (8 * (queue.GetNumberOfElevators() - 1)) - 1

// unjsoned data
type msg struct{
	MType string
	
	ENumber int
	Floor int
	Set bool
	Dir int
}
var sendToAllOthersChan = make (chan msg, 10) // 5 + 5 - 1 = 9 this is max msg simultaneous sent to sendToAllOthersChan //NB! Kan måttes gjøres større

// Struct and chan to handle remote calls
type RemoteMessage struct{
	Floor int
	Dir int
}
var RemoteChan = make (chan RemoteMessage)

func Initialize(eip []string, ep []int) bool{
	elevIpAddresses = eip
	elevPorts = ep

	if error := TCPServerInit(elevPorts[queue.GetElevatorNumber() - 1], sendChan, receiveChan); error != nil{
		fmt.Printf("communication: Initialize(): ERROR: %f", error)
		return false
	}
	
	go HandleOutgoingMessages()
	go HandleIncomingMessages()

	fmt.Printf("communication: TCP server now listening on port: %v\n", elevPorts[queue.GetElevatorNumber() - 1])	
	return true
}

func NotifyTheOthers(mtype string, floor int, set bool, dir int){
	fmt.Println("communication: NotifyTheOthers() was called")
	if (mtype == "OU" || mtype == "OD" || mtype == "F" || mtype == "D" || mtype == "ROU" || mtype == "ROD" || mtype == "T"){
		temp := msg{mtype, queue.GetElevatorNumber(), floor, set, dir}
		select{
		case sendToAllOthersChan <- temp:
			fmt.Println("communication: NotifyTheOthers(): msg was sent into sendToAllOthersChan")
		default:
			fmt.Println("communication: ***************************NotifyTheOthers(): ERROR: Can't send msg into --> sendToAllOthersChan <-- because it is BLOCKED!!")	
		}
	}else{
		fmt.Println("communication: ERROR: Can't NotifyTheOthers(), invalid string argument")
	}
}

//  Run as goroutine
func HandleOutgoingMessages() error{
	for temp := range sendToAllOthersChan{
		//fmt.Println("communication: HandleOutgoingMessages(): New message to send to the other elevators!")
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
					fmt.Println("communication: HandleOutgoingMessages(): ERROR: Can't send Tcp_message into --> sendChan <-- because it is BLOCKED!!")		
				}
				
			}
		}
		time.Sleep(100 * time.Millisecond) // sjekk mer!	fmgndfbdflgfkdfkgkdfjgbfjgdfkjgfkjdgbjkfdhgbkfdgbfdbgkjdbkgjhjdfkgbgh
	}
	r := fmt.Errorf("communication: ERROR: HandleOutgoingMessages() has quit range over sendToAllOthersChan!")
	return r
}

func HandleIncomingMessages() error{
	// Updates the local elevators OU,OD,FE,DE arrays according to incoming messages and sets/clears lights
	//fmt.Println("communication: HandleIncomingMessages() goroutine started")
	var m msg
	for temp := range receiveChan{
		err := json.Unmarshal(temp.Data[:temp.Length], &m)
		if err != nil{
			fmt.Println("communication: json.Unmarshal() error! HandleIncomingMessages() goroutine ending")
			fmt.Println(err)
			return err
		}else if m.MType == "OU"{
			queue.OrderFloorUp[m.Floor - 1] = m.Set
			fmt.Printf("communication: Remote set OrderFloorUp; order up on floor %v set to %v\n", m.Floor, m.Set)
			if m.Set == true{
				driver.SetButtonLight(0, m.Floor)
			}else if m.Set == false{
				driver.ClearButtonLight(0, m.Floor)
			}else{
				r := fmt.Errorf("communication: HandleIncomingMessages(): ERROR: Received and unjsoned a message with unknown Set. Something is wrong with HandleIncomingMessages()")
				fmt.Println(r)
				return r
			}
			
			//fmt.Printf("communication: Remote floor order up-button on floor %v set to %t\n", m.Floor, m.Set)
		}else if m.MType == "OD"{
			queue.OrderFloorDown[m.Floor - 1] = m.Set
			fmt.Printf("communication: Remote set OrderFloorDown; order down on floor %v set to %v\n", m.Floor, m.Set)
			if m.Set == true{
				driver.SetButtonLight(1, m.Floor)
			}else if m.Set == false{
				driver.ClearButtonLight(1, m.Floor)
			}else{
				r := fmt.Errorf("communication: HandleIncomingMessages(): ERROR: Received and unjsoned a message with unknown Set. Something is wrong with HandleIncomingMessages()")
				fmt.Println(r)
				return r				
			}
			//fmt.Printf("communication: Remote floor order down-button on floor %v set to %t\n", m.Floor, m.Set)
		
		}else if m.MType == "F"{
			queue.FloorElevator[m.ENumber - 1] = m.Floor
			fmt.Printf("communication: Remote elevator floor; elevator %v set its floor to %v\n", m.ENumber, m.Floor)
		}else if m.MType == "D"{
			queue.DirectionElevator[m.ENumber - 1] = m.Dir
			fmt.Printf("communication: Remote elevator direction; elevator %v set its direction to %v\n", m.ENumber, m.Dir)
		}else if (m.MType == "T"){
			queue.TaskElevator[m.ENumber - 1] = m.Floor
			fmt.Printf("communication: HandleIncomingMessages(): Elevator %v set its task to %v\n", m.ENumber, m.Floor)
		
		}else if (m.MType == "ROU" && m.Dir == queue.GetElevatorNumber()){ //In this case, m.Dir is the best fit elevator
			remotemsg := RemoteMessage{m.Floor, 0}
			select{
			case RemoteChan <- remotemsg:
			default:
				fmt.Println("communication: HandleIncomingMessages(): ERROR: Can't send RemoteMessage into --> RemoteChan <-- because it is BLOCKED!!")
			}
			fmt.Printf("communication: This best fit elevator to take order to floor %v was remote started from IDLE\n", m.Floor)		
		}else if (m.MType == "ROU" && m.Dir != queue.GetElevatorNumber()){
			fmt.Println("communication: HandleIncomingMessages(): ROU Call not for me")
		}else if (m.MType == "ROD" && m.Dir == queue.GetElevatorNumber()){ //In this case, m.Dir is the best fit elevator
			remotemsg := RemoteMessage{m.Floor, 1}
			select{
			case RemoteChan <- remotemsg:
			default:
				fmt.Println("communication: HandleIncomingMessages(): ERROR: Can't send msg into --> RemoteChan <-- because it is BLOCKED!!")
			}
			fmt.Printf("communication: HandleIncomingMessages(): This best fit elevator to take order to floor %v was remote started from IDLE\n", m.Floor)
		}else if (m.MType == "ROD" && m.Dir != queue.GetElevatorNumber()){
			fmt.Println("communication: HandleIncomingMessages(): ROD Call not for me")
		}else{
			r := fmt.Errorf("communication: HandleIncomingMessages(): ERROR: Received and unjsoned a message with unknown MType (%f). Something is wrong with HandleIncomingMessages()\n", m.MType)
			fmt.Println(r)
			return r
		}
		time.Sleep(10 * time.Millisecond)
	}
	r := fmt.Errorf("communication: HandleIncomingMessages(): ERROR: Quit range over receiveChan!")
	fmt.Println(r)
	return r
}
