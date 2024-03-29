package main

import (
	"bufio"
	"consumer"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"producer"
	"stats"
	"strconv"
	"strings"
	"time"
)

// Input Variables
var (
	M         float64 // number of times to repeat the tests (avg)
	TICKS     float64 // length of the test
	TICK_time float64 // 1 TICK = X milliseconds
	lambda    float64 // Average number of packets generated /arrived (packets per second)
	L         float64 // Length of a packet in bits
	C         float64 // The service time received by a packet. (Example: The transmission rate of the output link in bits per second.)
	K         float64 // The buffer size, 0 for infinite size

	//Averages over M
	Avg_Avg_packets      stats.Avg
	Avg_Avg_sojourn      stats.Avg
	Avg_Proportion_idle  stats.Avg
	Avg_Probability_loss stats.Avg
	Avg_total_packets    stats.Avg
	Avg_elapsed_time     stats.Avg

	csv_cols int = 15
)

var logger *log.Logger

var get_int_debug = false

/*
func get_int_csv(r string, b *int) {
	_, errI := fmt.Sscan(r, b)
	if errI != nil {
		logger.Fatalf("converted = %d\n%v\n", b, errI)
	}
}*/

func get_float64_csv(r string, b *float64) {
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

/*func get_int(r *bufio.Reader, b *int) {
	var trimval = get_val(r)
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	_, errI := fmt.Sscan(trimval, b)
	if errI != nil {
		logger.Fatalf("converted = %d\n%v\n", b, errI)
	}
}*/

func get_float64(r *bufio.Reader, b *float64) {
	var trimval = get_val(r)
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	_, errI := fmt.Sscan(trimval, b)
	if errI != nil {
		logger.Fatalf("converted = %f\n%v\n", b, errI)
	}
}

func main() {
	//Header
	fmt.Println("ECE 358 Project 1 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// End of Header

	// set up the logger to point to stdout
	logger = log.New(os.Stdout, "[ECE358 P1] ", log.LstdFlags)

	// Get Variables
	if len(os.Args) == 2 {
		get_args_params()
	} else {
		get_user_params()
	}
	// End of Get Variables

	// Display Variables
	fmt.Println("\nVariables being used:")
	fmt.Println("\t M          ", M)
	fmt.Println("\t TICKS      ", TICKS)
	fmt.Println("\t TICK_time  ", TICK_time, "milliseconds")
	fmt.Println("\t lambda     ", lambda)
	fmt.Println("\t L          ", L, "bits")
	fmt.Println("\t C          ", C, "bits/sec")
	if K == 0 {
		fmt.Println("\t K          ", "Infinity")
	} else {
		fmt.Println("\t K          ", K)
	}
	//End of display Variables

	// Loop for average statistics
	var test_t time.Time
	var test_b time.Time = time.Now()
	var tsince time.Duration
	for m := 1.0; m <= M; m++ {
		test_t = time.Now()

		fmt.Println("-------\n\n")

		wait_tick := get_tick_wait()
		qm := consumer.Init(logger, wait_tick)
		qm.MaxSize = K
		producer.Init(logger, qm, lambda, TICK_time)
		stats.Init(logger)

		logger.Println("---")

		for t := 1.0; t < TICKS; t++ {
			producer.Tick(t)
			consumer.Tick(t)

			// Getting the average packets in the queue
			stats.Avg_packets.AddAvg(float64(qm.Size))
		}
		logger.Println("Test Finished\n")

		tsince = time.Since(test_t)
		fmt.Println("[Stats] Elapsed Time                        =\t", tsince)
		fmt.Println("[Stats] End Queue Size                      =\t", qm.Size)
		fmt.Println("[Stats] Average Packets in Queue (E[N])     =\t", stats.Avg_packets.GetAvg())
		fmt.Println("[Stats] Average Sojourn Time (E[T]) (TICKS) =\t", stats.Avg_sojourn.GetAvg())
		fmt.Println("[Stats] Probability Packet Loss (P_LOSS)    =\t", stats.Probability_loss.GetProportion())
		fmt.Println("[Stats] Proportion Server Idle (P_IDLE)     =\t", stats.Proportion_idle.GetProportion())
		fmt.Println("[Stats] Total Packets Produced              =\t", stats.Probability_loss.Total)
		fmt.Println("[Stats] Total Packets Consumed              =\t", stats.Avg_sojourn.Num)
		fmt.Println("[Stats] Time Simulated (s)                  =\t", TICKS*TICK_time/1000)

		// Averages over M
		Avg_Avg_packets.AddAvg(stats.Avg_packets.GetAvg())
		Avg_Avg_sojourn.AddAvg(stats.Avg_sojourn.GetAvg())
		Avg_Probability_loss.AddAvg(stats.Probability_loss.GetProportion())
		Avg_Proportion_idle.AddAvg(stats.Proportion_idle.GetProportion())

		Avg_total_packets.AddAvg(stats.Probability_loss.Total)
		Avg_elapsed_time.AddAvg(tsince.Seconds()) //in seconds
	}

	tsince = time.Since(test_b)

	fmt.Println("\n---\n")
	fmt.Println("[Stats] Elapsed Time                            =\t", tsince)
	fmt.Printf("\n\tAverages over M (%d)\n", int(M))
	fmt.Println("[Stats] Avg Average Packets in Queue (E[N])     =\t", Avg_Avg_packets.GetAvg())
	fmt.Println("[Stats] Avg Average Sojourn Time (E[T]) (TICKS) =\t", Avg_Avg_sojourn.GetAvg())
	fmt.Println("[Stats] Avg Probability Packet Loss (P_LOSS)    =\t", Avg_Probability_loss.GetAvg())
	fmt.Println("[Stats] Avg Proportion Server Idle (P_IDLE)     =\t", Avg_Proportion_idle.GetAvg())
	fmt.Println("[Stats] Avg Total Packets Produced              =\t", Avg_total_packets.GetAvg())
	fmt.Println("[Stats] Avg Total Packets Consumed              =\t", Avg_Avg_sojourn.GetAvg())
	fmt.Println("[Stats] Avg Time Elapsed (s)                    =\t", Avg_elapsed_time.GetAvg())
	write_csv_header()
	write_csv_output(tsince.Seconds())
}

func get_tick_wait() float64 {
	return float64(math.Ceil(((float64(L) / C) * 1000) / float64(TICK_time)))
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

	get_float64_csv(rec[0], &M)
	get_float64_csv(rec[1], &TICKS)
	get_float64_csv(rec[2], &TICK_time)
	get_float64_csv(rec[3], &lambda)
	get_float64_csv(rec[4], &L)
	get_float64_csv(rec[5], &C)
	get_float64_csv(rec[6], &K)
}

func get_user_params() {
	var stdinR = bufio.NewReader(os.Stdin)
	fmt.Printf("M: ")
	get_float64(stdinR, &M)

	fmt.Printf("TICKS: ")
	get_float64(stdinR, &TICKS)

	fmt.Printf("TICK Time (1 TICK = X milliseconds): ")
	get_float64(stdinR, &TICK_time)

	fmt.Printf("lambda: ")
	get_float64(stdinR, &lambda)

	fmt.Printf("L (bits): ")
	get_float64(stdinR, &L)

	fmt.Printf("C (bits per sec): ")
	get_float64(stdinR, &C)
	fmt.Printf("K (zero = infinity): ")
	get_float64(stdinR, &K)
}

func write_csv_output(t_since float64) {

	file, err := os.OpenFile("test_out.csv", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Error in opening write file:", err)
	}
	defer file.Close()
	writter := csv.NewWriter(file)

	rec := make([]string, csv_cols)
	var i = 0
	rec[i] = strconv.FormatFloat(M, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(TICKS, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(TICK_time, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(lambda, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(L, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(C, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(K, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(get_tick_wait(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_total_packets.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_sojourn.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Probability_loss.GetAvg()*TICK_time/1000, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Proportion_idle.GetAvg()*TICK_time/1000, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_packets.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_elapsed_time.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(t_since, 'f', -1, 64)
	i++

	writter.Write(rec)
	writter.Flush()
}

func write_csv_header() {
	_, err := os.Open("test_out.csv")
	if os.IsNotExist(err) {
		file, _ := os.Create("test_out.csv")
		writter := csv.NewWriter(file)

		var i = 0
		headers := make([]string, csv_cols)
		headers[i] = "M"
		i++
		headers[i] = "TICKS"
		i++
		headers[i] = "TICK Time (sec)"
		i++
		headers[i] = "lambda"
		i++
		headers[i] = "L : Packet length (bits)"
		i++
		headers[i] = "C = bits per secound"
		i++
		headers[i] = "K: Buffer capacity (zero = inf)"
		i++
		headers[i] = "Service Time (Ticks)"
		i++
		headers[i] = "E[N] Avg num of packets in queue"
		i++
		headers[i] = "E[T] Avg sojourn time (sec)"
		i++
		headers[i] = "P_Loss - packet loss"
		i++
		headers[i] = "P_Idle - server idle (sec)"
		i++
		headers[i] = "Total Packets"
		i++
		headers[i] = "Average Time (sec)"
		i++
		headers[i] = "Total Time (sec)"
		i++

		writter.Write(headers)
		writter.Flush()
		file.Close()
	}
}
