package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

var (
	M         int64   // number of times to repeat the tests (avg)
	TICKS     int64   // length of the test
	TICK_time float64 // 1 TICK = X nanoseconds

	N_start, N_end, N_step int64   // Number of computers  connected to the LAN
	A                      float64 // Number of packets/sec
	W                      float64 // The speed of the lan in bits/sec
	L                      int64   // Length of a packet in bits

	length_of_line    float64 //in meters
	speed_over_line   float64 //in meters per sec
	jam_signal_length int64   //in bits

	Prop_ticks         int64
	Packet_trans_ticks int64
	Jam_trans_ticks    int64

	kmax              int64
	tp                int64
	medium_sense_time int64

	logger *log.Logger
)

func main() {
	fmt.Println("Values:\n\n")
	// gets and sets all global varaibles not dependant on N
	get_variables()

	fmt.Println("\t ---")
	var nextPacket int64 = getExpRandNum()
	fmt.Println("\t genPacket       ", nextPacket, "ticks")
	fmt.Println("\t maxGenPacket    ", (TICKS / nextPacket), "packets")
}

func getExpRandNum() int64 {
	a := 1 - rand.Float64()
	b := math.Log(a)
	c := (-1 / A)
	d := b * c             // sec / packet
	ans := sec_to_tickb(d) // tick / packet

	return ans
}

//convert seconds into ticks
func sec_to_tickb(s float64) int64 {
	return int64(math.Ceil((s * 1e9) / TICK_time))
}

//The exponential backoff wait time, when a collision is detected
func expBackOff(i int32) int64 {
	return int64(rand.Int31n(2^i-1)) * tp
}
