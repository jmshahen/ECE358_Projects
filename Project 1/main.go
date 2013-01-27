package main

import (
	"bufio"
	"fmt"
	"os"
	// "strconv"
	"consumer"
	"log"
	"producer"
	"strings"
	"time"
)

// Input Variables
var (
	M         int     // number of times to repeat the tests (avg)
	TICKS     int     // length of the test
	TICK_time int     // 1 TICK = X milliseconds
	lambda    int     // Average number of packets generated /arrived (packets per second)
	L         int     // Length of a packet in bits
	C         float32 // The service time received by a packet. (Example: The transmission rate of the output link in bits per second.)
	K         int     // The buffer size, 0 for infinite size
)

var logger *log.Logger

var get_int_debug = false

func get_int(r *bufio.Reader) int {
	var val, err = r.ReadString('\n')
	if err != nil {
		fmt.Println("| err:", err)
		os.Exit(1)
	}

	if get_int_debug {
		logger.Println("Val =", val)
	}
	var trimval = strings.TrimRight(val, "\n")
	if get_int_debug {
		logger.Println("trimmed =", trimval)
	}
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	var b int
	_, errI := fmt.Sscan(trimval, &b)
	if errI != nil {
		logger.Printf("converted = %d\n%v\n", b, errI)
		os.Exit(1)
	}

	return b
}

func main() {
	//Header
	fmt.Println("ECE 358 Project 1 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// End of Header

	logger = log.New(os.Stdout, "[ECE358 P1]", log.LstdFlags)

	// Get Variables
	var stdinR = bufio.NewReader(os.Stdin)
	fmt.Printf("M: ")
	M = get_int(stdinR)

	fmt.Printf("TICKS: ")
	TICKS = get_int(stdinR)

	fmt.Printf("TICK Time (1 TICK = X milliseconds): ")
	TICK_time = get_int(stdinR)

	fmt.Printf("lambda: ")
	lambda = get_int(stdinR)

	fmt.Printf("L (bits): ")
	L = get_int(stdinR)

	fmt.Printf("C (bits per sec): ")
	C = float32(get_int(stdinR))

	fmt.Printf("K (zero = infinity): ")
	K = get_int(stdinR)
	// End of Get Variables

	// Display Variables
	fmt.Println("Variables being used:")
	fmt.Println("\t M     ", M)
	fmt.Println("\t TICKS ", TICKS)
	fmt.Println("\t lambda", lambda)
	fmt.Println("\t L     ", L, "bits")
	fmt.Println("\t C     ", C, "bits/sec")
	if K == 0 {
		fmt.Println("\t K      Infinity")
	} else {
		fmt.Println("\t K     ", K)
	}
	//End of display Variables

	// Loop for average statistics
	var test_t time.Time
	for m := 1; m <= M; m++ {
		test_t = time.Now()

		consumer.Init(logger, L, C)
		producer.Init(logger)

		for t := 1; t < TICKS; t++ {
			producer.Tick(t)
			consumer.Tick(t)
		}
		logger.Println("[Info] Finished running Test #", m, "elapsed time:", time.Since(test_t))
	}
}
