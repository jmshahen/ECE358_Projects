package main

import (
	"fmt"
	"lan"
)

var (
	M         float64 // number of times to repeat the tests (avg)
	TICKS     float64 // length of the test
	TICK_time float64 // 1 TICK = X milliseconds

	N float64 // Number of computers  connected tot he LAN
	A float64 // Number of packets/sec
	W float64 // The speed of the lan in bits/msec
	L float64 // Length of a packet in bits

	bit_time float64 // the bit time. - time it takes to put 1 bit on the line.
)

func main() {
	fmt.Println("HELLO WORLD\n")

	lan.Init(10, 5, 1000, 1000000, 1500)

}
