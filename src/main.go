//go run main.go -raddr="129.241.187.145:20033" -lport=20034
package main 

import (
	"fmt"
	"time"
	"strconv"
	"flag"
	"../src/network"
)



var raddr = flag.String("raddr", "127.0.0.1:21331", "the ip adress for the target connection")
var lport = flag.Int("lport", 21337, "the local port to listen on for new conns")


func main (){
	flag.Parse()
	fmt.Println("tcp_example_test.go")
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
			fmt.Println("%v Sent: %v", lport,msg)
			time.Sleep(1*time.Second)	
		}	
	}(schan)

	for {
		msg := <- rchan
		fmt.Println("%v Received: %v",lport, msg)	
	}
}
