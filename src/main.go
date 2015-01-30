package main


import (
    "networkmodule"
)

 func main() {
	go networkmodule.UDPListen()
	go networkmodule.UDPSend()

	select {}
 }
