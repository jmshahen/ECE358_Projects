package main

import (
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
		get_user_params()
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
