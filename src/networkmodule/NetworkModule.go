package networkmodule

import (
    "fmt"
    "net"
    "time"
)
func UDPListen(){
	buffer := make([]byte, 1024)
	address,error := net.ResolveUDPAddr("udp4", ":20017")
	fmt.Println("Address: ", address, error)
	socket, error := net.ListenUDP("udp4", address)
	fmt.Println(socket, error)
	for{
		readLength, remote, error := socket.ReadFromUDP(buffer[0:])
		fmt.Println(readLength, "bytes received from", remote, error)
		fmt.Println("  ",  string(buffer[:]), "\n")
		buffer = make([]byte, 1024)
	}

}

func UDPSend(){
	address,_ := net.ResolveUDPAddr("udp4", "129.241.187.255:20017")
	fmt.Println(address)
	socket, error := net.DialUDP("udp4", nil, address)
	fmt.Println("SendSock: ", socket, error)
	for {
		time.Sleep(time.Second)
		_, error := socket.Write([]byte("yay"))

		if error != nil {
			fmt.Println("Error sending:", error)
		}
	}

}

