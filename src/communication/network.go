package communication

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"errors"
	"sync"
	"queue"
)

var BroadcastIP string = "129.241.255.255"
var BroadcastPort int = 30500

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

	// Making conn_list instance
	conn_list = make(map[string]*net.TCPConn)

	
	// Using udp to figure out its address on the local network
	baddr, err := net.ResolveUDPAddr("udp4", BroadcastIP+":"+strconv.Itoa(BroadcastPort)) // Itoa for string representation of integer
	if err != nil {
		fmt.Println("communication: TCPServerInit: ERROR: Could not resolve baddr")
		return err
	}
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	if err != nil {
		fmt.Println("communication: TCPServerInit: Failed to dial baddr for tempLAddr generation")
		return err
	}
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
		//fmt.Println("communication: tcp_transmit_server: New message to send")
		_ , ok := conn_list[msg.Raddr]
		if (ok != true ){
			// NB! This blocks s_ch (sendChan) because of reconnect for-loop. Not any more, removed for-loop!
			if !new_tcp_conn(msg.Raddr){
				fmt.Println("communication: tcp_transmit_server could not establish new tcp connection right now, after two tries")
				fmt.Println("               Remote elevator(s) must be disconnected and will not receive the sent messages until they are turned on")
				fmt.Println("--> Listening for next outgoing message to arrive on sendChan")
				continue
			}	
		}
		conn_list_mutex.Lock()
		sendConn, ok  := conn_list[msg.Raddr]
		if (ok != true) {
			conn_list_mutex.Unlock()
			err := errors.New("communication: tcp_transmit_server: Failed to add newConn to conn_list map")
			panic(err)
		}else{
			n, err := sendConn.Write(msg.Data)
			//fmt.Println("communication: tcp_transmit_server(): Message sent.")	
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
			//fmt.Println("communication: tcp_handle_server: Starting new tcp read server because a new connection was accepted")
			buf := make([]byte, 1024)
			for {
					n, err :=	conn.Read(buf)
					if err != nil || n < 0 {
						fmt.Printf("communication: tcp_handle_server: Read error: %s \n",err)
						conn_list_mutex.Lock()
						conn.Close()
						delete(conn_list, raddr)
						conn_list_mutex.Unlock()
						return 
					}else{
						//fmt.Println("communication: tcp_handle_server: New byte array read from local port")
						select{
						case r_ch <- Tcp_message{Raddr: raddr, Data: buf, Length: n}:
							fmt.Println("communication: tcp_handle_server: Sent received Tcp_message into receiveChan")
						default:
							fmt.Println("communication: tcp_handle_server: ERROR: Can't send Tcp_message into --> receiveChan (r_ch) <-- because it is BLOCKED!!\n")
						}
						//r_ch <- Tcp_message{Raddr: raddr, Data: buf, Length: n}
					}
			}		
		}(raddr.String(), newConn, r_ch)
	}
}
 
func new_tcp_conn(raddr string) bool{
	//fmt.Println("communication: new_tcp_conn: Trying to establish new tcp connection")
	//create address
	addr, err := net.ResolveTCPAddr("tcp4", raddr)
	if err != nil {
		fmt.Println("communication: new_tcp_conn: ERROR: could not resolve address")
		return false
	}
	/*
	for {
		newConn, err := net.DialTCP("tcp4", nil,  addr)
		
		if err != nil {
			fmt.Printf("communication: new_tcp_conn: DialTCP to %v failed\n", raddr)
				time.Sleep(500*time.Millisecond)
		} else {
			fmt.Printf("communication: new_tcp_conn: DialTCP to %v succeeded\n", raddr)
			conn_list_mutex.Lock()
			conn_list[raddr] = newConn
			conn_list_mutex.Unlock()
			return true//got it BREAK!
		}
	}
	*/
	newConn, err := net.DialTCP("tcp4", nil,  addr)
	if err != nil {
		fmt.Printf("communication: new_tcp_conn: DialTCP to %v failed. Trying DialTCP one more time\n", raddr)
		time.Sleep(500*time.Millisecond)
		newConn, err := net.DialTCP("tcp4", nil,  addr)
		
		if err != nil{
			//fmt.Printf("communication: new_tcp_conn: DialTCP to %v failed again!. Giving up\n", raddr)
			return false
		}else{
			//fmt.Printf("communication: new_tcp_conn: DialTCP to %v succeeded after a re-try\n", raddr)
			conn_list_mutex.Lock()
			conn_list[raddr] = newConn
			conn_list_mutex.Unlock()
			return true//got it BREAK!
		}
		
	}else{
		//fmt.Printf("communication: new_tcp_conn: DialTCP to %v succeeded\n", raddr)
		conn_list_mutex.Lock()
		conn_list[raddr] = newConn
		conn_list_mutex.Unlock()
		return true//got it BREAK!
	}
}

