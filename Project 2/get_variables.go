package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"time"
)

func get_variables() {

	// set up the logger to point to stdout
	logger = log.New(os.Stdout, "[ECE358 P2] ", log.LstdFlags)

	/**************************************************************************/
	// Get Variables

	if len(os.Args) == 2 {
		get_args_params()
	} else {
		logger.Fatalln("Only a CSV file is premitted as commandline arguments")
	}

	// End of Get Variables

	/**************************************************************************/
	// Set Variables
	rand.Seed(time.Now().UnixNano())
	tp = bit_time_to_ticks(tp)
	medium_sense_time = bit_time_to_ticks(medium_sense_time)

	Prop_ticks = sec_to_tick(length_of_line / speed_over_line)
	Packet_trans_ticks = sec_to_tick(float64(L) / W)
	Jam_trans_ticks = sec_to_tick(float64(jam_signal_length) / W)

	dir, file := path.Split(csv_file_name_in)
	csv_file_name_out = dir + "test_results_-_" + file
	// End of Set Variables

	// Display Variables
	fmt.Println("\nInput file:", csv_file_name_in)
	fmt.Println("\nResults saved to:", csv_file_name_out)
	fmt.Println("\nVariables being used:")
	fmt.Println("\t M               ", M)
	fmt.Println("\t TICKS           ", TICKS)
	fmt.Println("\t TICK_time       ", TICK_time, "nanoseconds / tick")
	fmt.Println("\t ---")
	fmt.Println("\t N_start         ", N_start)
	fmt.Println("\t N_end           ", N_end)
	fmt.Println("\t N_step          ", N_step)
	fmt.Println("\t ---")
	fmt.Println("\t A               ", A, "pcks/sec")
	fmt.Println("\t W               ", W, "bits/sec")
	fmt.Println("\t L               ", L, "bits")
	fmt.Println("\t ---")
	fmt.Println("\t kmax            ", kmax)
	fmt.Println("\t tp              ", tp, "ticks")
	fmt.Println("\t m_sense         ", medium_sense_time, "ticks")
	fmt.Println("\t ---")
	fmt.Println("\t jam_sig         ", jam_signal_length, "bits")
	fmt.Println("\t line            ", length_of_line, "meters")
	fmt.Println("\t e_speed         ", speed_over_line, "meters / sec")
	fmt.Println("\t ---")
	fmt.Println("\t prop_d          ", Prop_ticks, "ticks")
	fmt.Println("\t trans_d         ", Packet_trans_ticks, "ticks")
	fmt.Println("\t jam_d           ", Jam_trans_ticks, "ticks")
	fmt.Println("\t ---")
	nextPacket := getExpRandNum()
	fmt.Println("\t genPacket       ", nextPacket, "ticks / packet")
	fmt.Println("\t maxGenPacket    ", (TICKS / nextPacket), "packets")
	fmt.Println("\t ---")
	fmt.Println("\t --- ---\n")
	// End of Display Variables
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

func get_args_params() {
	csv_file_name_in = os.Args[1]
	file, err := os.Open(csv_file_name_in)
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

	i := 0
	get_int64_csv(rec[i], &M)
	i++
	get_float64_csv(rec[i], &TICK_time)
	TICKS = int64(TICK_time)
	i++
	get_float64_csv(rec[i], &TICK_time)
	i++

	get_int64_csv(rec[i], &N_start)
	i++
	get_int64_csv(rec[i], &N_end)
	i++
	get_int64_csv(rec[i], &N_step)
	i++

	get_float64_csv(rec[i], &A)
	i++
	get_float64_csv(rec[i], &W) //in bits per second
	i++
	get_int64_csv(rec[i], &L)
	i++

	get_int64_csv(rec[i], &kmax)
	i++
	get_int64_csv(rec[i], &tp) //bit times
	i++
	get_int64_csv(rec[i], &medium_sense_time) //bit times
	i++

	get_int64_csv(rec[i], &jam_signal_length) //bits
	i++
	get_float64_csv(rec[i], &length_of_line) //in meters
	i++
	get_float64_csv(rec[i], &speed_over_line) //in meters per second
	i++
}

func bit_time_to_ticks(v int64) int64 {
	return sec_to_tick(float64(v) / W)
}

func sec_to_tick(s float64) int64 {
	return int64(math.Ceil((s * 1e9) / float64(TICK_time)))
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
