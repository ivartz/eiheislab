package communication

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"errors"
	"sync"
	"queue"
	//"encoding/json"
)
//test

var BroadcastIP string = "129.241.255.255"
var BroadcastPort int = 30500

// The message form is a struct. Capital letter because the struct is used outside network.go
// This struct is ent in send_ch and receive_ch
type Tcp_message struct{
	Raddr string
	Data string 
	Length int
}

//type tcp_conn struct {
//	conn *net.TCPConn
//	receive_ch chan Tcp_message
//}

// Map (dictionary) that keeps track of the existing tcp connections
var conn_list map[string]*net.TCPConn
// Mechanism that locks and unlocks conn_list to critical sections 
var conn_list_mutex = &sync.Mutex{}

func TCPServerInit(localListenPort int, send_ch, receive_ch chan Tcp_message) error{
	//var buffer [1024]byte

 	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("communication: TCPServerInit: ERROR recovering in tcp_init: %s \n", r)
		}
	}()

	//fmt.Println("TCPServerInit: called")

	// Making conn_list instance
	conn_list = make(map[string]*net.TCPConn)

	
	// Using udp to figure out its address on the local network
	baddr, err := net.ResolveUDPAddr("udp4", BroadcastIP+":"+strconv.Itoa(BroadcastPort)) // Itoa for string representation of integer
	
	if err != nil {
		fmt.Println("communication: TCPServerInit: Could not resolve baddr")
		return err
	}

	tempConn, err := net.DialUDP("udp4", nil, baddr)
	if err != nil {
		fmt.Println("communication: TCPServerInit: Failed to dial baddr for tempLAddr generation")
		return err
	}

	fmt.Printf("communication: TCPServerInit: UDP-dialed %v to resolve its local address\n", baddr) // midlertidig
	//fmt.Println(baddr)

	tempLAddr := tempConn.LocalAddr()

	

	// Using tempLAddr to make tcp connection to listen to
 	laddr, err := net.ResolveTCPAddr("tcp4", tempLAddr.String())
 	if err != nil{
 		fmt.Println("communication: TCPServerInit: Failed to ResolveTCPAddr")
 		return err
 	}
	
	laddr.Port = localListenPort
	tempConn.Close()

	listen, err := net.ListenTCP("tcp4",laddr)
	if err != nil {
		fmt.Println("communication: TCPServerInit: Failed to initialize ListenTCP")
		return err
	}
	fmt.Printf("communication: TCPServerInit: Elevator %v with is now listening on %v, its local ip and listening port\n", queue.GetElevatorNumber(), laddr)

	// Goroutine to handle all incoming messages from send_ch. Establishing additional tcp connection if necessary (adding to conn_list) 
	go tcp_transmit_server(send_ch)

	// Goroutine to receive all incoming connections from the listener. Adding new connections to conn_list
	go tcp_handle_server(listen, receive_ch)

	return err
 }

func tcp_transmit_server (s_ch chan Tcp_message){
	for {
		msg := <- s_ch
		fmt.Println("communication: tcp_transmit_server: New message to send")
		_ , ok := conn_list[msg.Raddr]
		if (ok != true ){
			new_tcp_conn(msg.Raddr)	
		}
		conn_list_mutex.Lock()
		sendConn, ok  := conn_list[msg.Raddr]
		if (ok != true) {
			conn_list_mutex.Unlock()
			err := errors.New("communication: tcp_transmit_server: Failed to add newConn to conn_list map")
			panic(err)
		} else {
			n, err := sendConn.Write([]byte(msg.Data))	
			conn_list_mutex.Unlock()
			if err != nil || n < 0 {
				fmt.Printf("communication: tcp_transmit_server: Write error (deleting remote address): %s\n",err)
				conn_list_mutex.Lock()
				sendConn.Close()
				delete(conn_list, msg.Raddr)
				conn_list_mutex.Unlock()
			}
		}
	}
}

func tcp_handle_server (listener *net.TCPListener, r_ch chan Tcp_message){
	for {
		newConn, err := listener.AcceptTCP()
		fmt.Println("communication: tcp_handle_server: Received new request for connection")
		if err != nil {
			fmt.Printf("communication: tcp_handle_server: Error accepting tcp conn \n")
			panic(err)
		}

		// assume from here the connection was accepted

		raddr := newConn.RemoteAddr()
		
		conn_list_mutex.Lock()
		conn_list[raddr.String()] = newConn
		conn_list_mutex.Unlock()

		//reading server to read from the accepted connection newConn
		go func (raddr string, conn *net.TCPConn, ch chan Tcp_message){ 
			fmt.Println("communication: tcp_handle_server: Starting new tcp read server because a new connection was accepted")
			buf := make([]byte,1024)
			for {
					n, err :=	conn.Read(buf)
					if err != nil || n < 0 {
						fmt.Printf("communication: tcp_handle_server: Read error: %s \n",err)
						conn_list_mutex.Lock()
						conn.Close()
						delete(conn_list, raddr)
						conn_list_mutex.Unlock()
						return 
					} else {
						r_ch <- Tcp_message{Raddr: raddr, Data: string(buf), Length: n}
					}
			}		
		}(raddr.String(), newConn, r_ch)
		
	}
}
 
func new_tcp_conn(raddr string) bool{
	fmt.Println("communication: new_tcp_conn: Trying to establish new tcp connection")
	//create address
	addr, err := net.ResolveTCPAddr("tcp4", raddr)
	if err != nil {
		fmt.Println("communication: new_tcp_conn: ERROR: could not resolve address")
		return false
	}
	for {
		newConn, err := net.DialTCP("tcp4", nil,  addr)
		
		if err != nil {
			fmt.Printf("communication: new_tcp_conn: DialTCP to %v failed \n", raddr)
				time.Sleep(500*time.Millisecond)
		} else {
			conn_list_mutex.Lock()
			conn_list[raddr] = newConn
			conn_list_mutex.Unlock()
			return true//got it BREAK!
		}
	}
}

