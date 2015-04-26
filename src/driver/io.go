package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#cgo CFLAGS: -Wall
#include "io.h"
*/
import "C"

//return Non-zero on success, and 0 on failure
func IoInit() int{
	retVal := int(C.io_init())
	return retVal
}
func IoSetBit(channel int){
	//fmt.Printf("driver: IoSetBit on channel: %v\n", channel)
	C.io_set_bit(C.int(channel))
}
func IoClearBit(channel int){
	C.io_clear_bit(C.int(channel))
}
func IoWriteAnalog(channel int, value int){
	C.io_write_analog(C.int(channel), C.int(value))
}
func IoReadBit(channel int) int{
	return int(C.io_read_bit(C.int(channel)))
}
func IoReadAnalog(channel int) int{
	return int(C.io_read_analog(C.int(channel)))
}
