//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
	"time"
	"strconv"
	"flag"
	"../src/network"
//	"reflect"
)
//test


var raddr = flag.String("raddr", "127.0.0.1:21331", "the ip adress for the target connection")
var lport = flag.Int("lport", 20017, "the local port to listen on for new conns")


func main (){
	flag.Parse()
	fmt.Println("main: This is a test on sending and receiving tcp messages between two hosts using the network module")
	rchan := make (chan network.Tcp_message)
	schan := make (chan network.Tcp_message)
	network.TCPServerInit(*lport, schan, rchan)
	
	go func(ch chan network.Tcp_message){
		id := 0
		msg := network.Tcp_message{Raddr: *raddr, Data: strconv.Itoa(id), Length: 32}

		for {
			msg.Data = strconv.Itoa(id)
			schan <- msg
			id++
//			fmt.Println("%v Sent: %v", lport msg.Data)
			fmt.Printf("main: Sent data %v through port %v to %v. 5 seconds to next send. \n", msg.Data, *lport, *raddr)
			time.Sleep(5*time.Second)	
		}	
	}(schan)

	for {
		msg := <- rchan
		fmt.Printf("main: Received data %v on port %v\n", msg.Data, *lport)	
	}
}
