package communication

import(
	"fmt"
	"strconv"
	"encoding/json"
	"queue"
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
var msg struct{
	MType string
	
	ENumber int
	Floor int
	Set bool
	Dir int
}
var sendToAllOthersChan = make (chan msg)


func Initialize() bool{
	
	elevIpAddresses = []string{"129.241.187.144", "129.241.187.141", "129.241.187.146"}
	elevPorts = []int{20012, 20015, 20006}

	TCPServerInit(elevPorts[queue.GetElevatorNumber() - 1], sendChan, receiveChan)

	return true
}


func NotifyTheOthers(mtype string, floor int, set bool, dir int){
	if (ntype == "OU" || mtype == "OD" || mtype == "F" || mtype == "D"){
		temp := msg{mtype, queue.GetElevatorNumber(), floor, set, dir}
		select{
		case sendToAllChan <- temp:
		default:
			fmt.Println("communication: ERROR: Notify failed. NotifyTheOthers() can't send message because sendToAllOthersChan is busy")	
		}
	}else{
		fmt.Println("communication: ERROR: Can't NotifyTheOthers(), invalid string argument")
	}
}


// Run as goroutine
func HandleOutgoingMessages(){
	fmt.Println("communication: HandleOutgoingMessages() goroutine started")
	for temp := sendToAllChan{
		jtemp, err := json.Marshal(temp)
		if err != nil{
			fmt.Println("communication: json.Marshal() error!")
		}
		for i := range elevIpAddresses{
			if i + 1 != queue.GetElevatorNumber(){
				tcpm := Tcp_message{elevIpAddresses[i]+":"+strconv.Itoa(elevPorts[i]), jtemp, len(jtemp)}
				sendChan <- tcmp		
			}
		}	
	}
}

func HandleIncomingMessages(){
	// Updates the local elevators OU,OD,FE,DE arrays according to incoming messages 
	fmt.Println("communication: HandleIncomingMessages() goroutine started")
	var m Msg
	for temp := range receiveChan{
		err := json.Unmarshal(temp, &m)
		if err != nil{
			fmt.Println("communication: json.Unmarshal() error!")
			return err
		}else if m.MType == "OU"{
			queue.OrderFloorUp[m.Floor - 1] = m.Set
			fmt.Printf("communication: Remote floor order up-button to %v\n", m.Floor)
		}else if m.MType == "UD"{
			queue.OrderFloorDown[m.Floor - 1] = m.Set
			fmt.Printf("communication: Remote floor order down-button to %v\n", m.Floor)
		}else if m.MType == "F"{
			queue.FloorElevator[m.ENumber - 1] = m.Floor
			fmt.Printf("communication: Remote elevator floor; elevator %v set its floor to %v\n", m.ENumber, m.Floor)
		}else if m.MType == "D"{
			queue.DirectionElevator[m.ENumber - 1] = m.Dir
			fmt.Printf("communication: Remote elevator direction; elevator %v set its direction to %v\n", m.ENumber, m.Dir)
		}
		else{
			fmt.Println("communication: HandleIncomingMessages() received and unjsoned a message with unknown MType. Something is wrong with HandleIncomingMessages()")
		}
	}
}
