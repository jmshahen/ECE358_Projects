package main

import (
	"fmt"
	"log"
	"math"
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

	logger            *log.Logger
	csv_file_name_in  string
	csv_file_name_out string
)

func main() {
	fmt.Println("Values:\n\n")
	// gets and sets all global varaibles not dependant on N
	get_variables()

	fmt.Println("\t ---")
	var nextPacket int64
	var avg int64
	for i := int64(0); i < M; i++ {
		avg += getExpRandNum()
	}
	nextPacket = int64(avg / M)
	fmt.Println("\t genPacket       ", nextPacket, "ticks")
	fmt.Println("\t maxGenPacket    ", (TICKS / nextPacket), "packets / per one comp")
	a := math.Floor(float64(nextPacket) / float64(medium_sense_time+Packet_trans_ticks))
	fmt.Println("\t maxCompsSend    ", a, "comps")
	b := math.Min(a, float64(N_end))
	fmt.Println("\t availCompsSend  ", b, "comps")
	fmt.Println("\t availGenPacket  ", math.Floor(b*float64(TICKS)/float64(nextPacket)), "packets / per", b, "comp")
	fmt.Println("\t maxGenPacket    ", math.Floor(a*float64(TICKS)/float64(nextPacket)), "packets / per", a, "comp")
}
