package main

import (
	"csv"
	"os"
)

func get_variables() {

	/*
		Accept the variables: A, W, L
		Accept all variables need for dicrete event simulator
		Accept the Variable: N_start, N_end, N_step
	*/

	// Get Variables
	if len(os.Args) == 2 {
		get_args_params()
	} else {
		logger.Fatalln("Only a CSV file is premitted as commandline arguments")
	}
	// End of Get Variables

	// Set Variables
	Prop_ticks = (length_of_line / speed_over_line) * 1e9 / TICK_time // 100 m/s divided by propagation speed. converted to nano secounds, converted to ticks.
	Packet_trans_ticks = (L / W) * 1e9 / TICK_time                    // packet length divided by trans speed. converted to nano secounds, converted to ticks
	Jam_trans_ticks = (jam_signal_length / W) * 1e9 / TICK_time       // Jam length divided by trans speed. converted to nano secounds, converted to ticks
	// End of Set Variables

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
	get_float64_csv(rec[2], &TICK_time)

	get_int64_csv(rec[3], &N_start)
	get_int64_csv(rec[3], &N_end)
	get_int64_csv(rec[3], &N_step)

	get_int64_csv(rec[4], &A)
	get_int64_csv(rec[5], &W)
	get_int64_csv(rec[6], &L)
}

func get_float64_csv(r string, b *float64) {
	_, errI := fmt.Sscan(r, b)
	if errI != nil {
		logger.Fatalf("converted = %f\n%v\n", b, errI)
	}
}

func get_int64_csv(r string, b *int64) {
	_, errI := fmt.Sscan(r, b)
	if errI != nil {
		logger.Fatalf("converted = %d\n%v\n", b, errI)
	}
}
