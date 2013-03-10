package main

import (
	"csma"
	"fmt"
	"lan"
	"stats"
)

var (
	M         int64 // number of times to repeat the tests (avg)
	TICKS     int64 // length of the test
	TICK_time int64 // 1 TICK = X milliseconds

	N_start, N_end, N_step int64 // Number of computers  connected to the LAN
	A                      int64 // Number of packets/sec
	W                      int64 // The speed of the lan in bits/sec
	L                      int64 // Length of a packet in bits

	bit_time int64 // the bit time. - time it takes to put 1 bit on the line.

	lan         *lan.LAN
	cmsa        *csma_cd.CSMA
	bucket      *stats.Bucket
	lost_bucket *stats.Bucket

	Prop_ticks         int64
	Packet_trans_ticks int64
	Jam_trans_ticks    int64

	kmax              int64
	tp                int64
	medium_sense_time int64
)

func main() {
	//Header
	fmt.Println("ECE 358 Project 2 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// End of Header

	// gets and sets all global varaibles not dependant on N
	get_variables()

	/*
		Loop over N (i=N_start; i<= N_end; i += N_step)
			setup the buket and the computers, etc
			loop over TICK=1->TICKS
				loop over comp_i=1->COMPS
					call comp_i->Tick(TICK)
					call lan->wait_and_send(TICK)
			Save results for this N
		Generate graphs
	*/

	// Lan init()

	//lan.Init(num_comps int64, Prop_ticks int64, Packet_trans_ticks int64, Jam_trans_ticks int64 64, bucket *stats.Bucket, lost_bucket *stats.Bucket)
	// end lan init()

	// csma Init()

	kmax = 512
	tp = 1 //?
	medium_sense_time = 96

	//(id int64, lan *lan.LAN, logger *log.Logger, _lambda float64, TICK_time float64, kmax int64, tp int64, medium_sense_time int64)
	//csma.Init(id int64, lan *lan.LAN, logger *log.Logger, _lambda float64, TICK_time float64, kmax int64, tp int64, medium_sense_time int64)
	// end csma Init()

}
