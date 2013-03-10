package main

import (
	"fmt"
	"lan"
	"stats"
	"csma"
)

var (
	M         int64 // number of times to repeat the tests (avg)
	TICKS     int64 // length of the test
	TICK_time int64 // 1 TICK = X milliseconds

	num_comps int64 // Number of computers  connected to the LAN
	A int64 // Number of packets/sec
	W int64 // The speed of the lan in bits/sec
	L int64 // Length of a packet in bits

	bit_time int64 // the bit time. - time it takes to put 1 bit on the line.



	lan  		*lan.LAN
	cmsa		*csma_cd.CSMA
	bucket		*stats.Bucket 
	lost_bucket	*stats.Bucket

	Prop_ticks 				int64
	Packet_trans_ticks      int64
	Jam_trans_ticks			int64

	kmax 					int64
	tp 						int64
	medium_sense_time 		int64

)

func main() {
	//Header
	fmt.Println("ECE 358 Project 2 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// End of Header

}


	// Get Variables
	if len(os.Args) == 2 {
		get_args_params()
	} else {
		get_user_params()
	}
	// End of Get Variables

/*
Accept the variables: A, W, L
Accept all variables need for dicrete event simulator
Accept the Variable: N_start, N_end, N_step
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
	 Prop_ticks = (length_of_line / speed_over_line)*1e9 / TICK_time  // 100 m/s divided by propagation speed. converted to nano secounds, converted to ticks.
	 Packet_trans_ticks = (L/W)*1e9/TICK_time  // packet length divided by trans speed. converted to nano secounds, converted to ticks
	 Jam_trans_ticks = (jam_signal_length/W)*1e9/TICK_time  // Jam length divided by trans speed. converted to nano secounds, converted to ticks

	lan.Init(num_comps int64, Prop_ticks int64, Packet_trans_ticks int64, Jam_trans_ticks int64 64, bucket *stats.Bucket, lost_bucket *stats.Bucket)
	// end lan init()


	// csma Init()

		 kmax = 512
		 tp = 1 //?
		 medium_sense_time = 96

		   //(id int64, lan *lan.LAN, logger *log.Logger, _lambda float64, TICK_time float64, kmax int64, tp int64, medium_sense_time int64)
	csma.Init(id int64, lan *lan.LAN, logger *log.Logger, _lambda float64, TICK_time float64, kmax int64, tp int64, medium_sense_time int64)
	// end csma Init()




}

func get_args_params() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		logger.Fatalf("Error:", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	rec, err := reader.Read()
	if err == io.EOF {
		logger.Fatalf("Error: No Headers")
	} else if err != nil {
		logger.Fatalf("Error:", err)
	}
	// throwaway header

	rec, err = reader.Read()
	if err == io.EOF {
		logger.Fatalf("Error: No Data")
	} else if err != nil {
		logger.Fatalf("Error:", err)
	}

	get_int64_csv(rec[0], &M)
	get_int64_csv(rec[1], &TICKS)
	get_int64_csv(rec[2], &TICK_time)
	get_int64_csv(rec[3], &N)
	get_int64_csv(rec[4], &A)
	get_int64_csv(rec[5], &W)
	get_int64_csv(rec[6], &L)
}

func get_user_params() {
	var stdinR = bufio.NewReader(os.Stdin)
	fmt.Printf("M: ")
	get_int64(stdinR, &M)

	fmt.Printf("TICKS: ")
	get_int64(stdinR, &TICKS)

	fmt.Printf("TICK Time (1 TICK = X milliseconds): ")
	get_int64(stdinR, &TICK_time)

	fmt.Printf("N: ")
	get_int64(stdinR, &N)

	fmt.Printf("A: (packets/sec) ")
	get_int64(stdinR, &A)

	fmt.Printf("W: (bits/sec) ")
	get_int64(stdinR, &W)
	fmt.Printf("L: ")
	get_int64(stdinR, &L)
}

func get_int64_csv(r string, b *int64) {
	_, errI := fmt.Sscan(r, b)
	if errI != nil {
		logger.Fatalf("converted = %f\n%v\n", b, errI)
	}
}
func get_val(r *bufio.Reader) string {
	var val, err = r.ReadString('\n')
	if err != nil {
		logger.Fatalln("| err:", err)
	}

	var trimval = strings.TrimRight(val, "\n")
	if get_int_debug {
		logger.Println("Val =", val)
		logger.Println("trimmed =", trimval)
	}

	return trimval
}

func get_int64(r *bufio.Reader, b *int6) {
	var trimval = get_val(r)
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	_, errI := fmt.Sscan(trimval, b)
	if errI != nil {
		logger.Fatalf("converted = %f\n%v\n", b, errI)
	}
}