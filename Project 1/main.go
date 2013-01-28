package main

import (
	"bufio"
	"strings"
	"fmt"
	"encoding/csv"
	"io"
	"os"
	// "strconv"
	"consumer"
	"log"
	"math"
	"producer"
	"time"
)

// Input Variables
var (
	M         int     // number of times to repeat the tests (avg)
	TICKS     int     // length of the test
	TICK_time int     // 1 TICK = X milliseconds
	lambda    float64 // Average number of packets generated /arrived (packets per second)
	L         int     // Length of a packet in bits
	C         float64 // The service time received by a packet. (Example: The transmission rate of the output link in bits per second.)
	K         int     // The buffer size, 0 for infinite size
)

var logger *log.Logger

var get_int_debug = false


func get_int_csv(r string, b *int) {
    _, errI := fmt.Sscan(r, b)
    if errI != nil {
        fmt.Printf("converted = %d\n%v\n", b, errI)
        os.Exit(1)
    }
}

func get_float64_csv(r string, b *float64) {
    _, errI := fmt.Sscan(r, b)
    if errI != nil {
        fmt.Printf("converted = %f\n%v\n", b, errI)
        os.Exit(1)
    }
}
func get_val(r *bufio.Reader) string {
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

	return trimval
}

func get_int(r *bufio.Reader, b *int) {
	var trimval = get_val(r)
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	_, errI := fmt.Sscan(trimval, b)
	if errI != nil {
		logger.Printf("converted = %d\n%v\n", b, errI)
		os.Exit(1)
	}
}

func get_float64(r *bufio.Reader, b *float64) {
	var trimval = get_val(r)
	// valI, errI := strconv.ParseInt(trimval, 10, 32)
	_, errI := fmt.Sscan(trimval, b)
	if errI != nil {
		logger.Printf("converted = %f\n%v\n", b, errI)
		os.Exit(1)
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

	logger = log.New(os.Stdout, "[ECE358 P1] ", log.LstdFlags)

	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
	    if err != nil {
	        fmt.Println("Error:", err)
	        return
	    }
	    defer file.Close()
		reader := csv.NewReader(file)

	    rec, err := reader.Read()
	    if err == io.EOF {
	        fmt.Println("Error: No Headers")
	    } else if err != nil {
	        fmt.Println("Error:", err)
	        return
	    }
	    // throwaway header

	    rec, err = reader.Read()
	    if err == io.EOF {
	        fmt.Println("Error: No Data")
	    } else if err != nil {
	        fmt.Println("Error:", err)
	        return
	    }

	    fmt.Printf("M: ")
	    get_int_csv(rec[0], &M)
	    fmt.Println(M)

	    fmt.Printf("TICKS: ")
	    get_int_csv(rec[1], &TICKS)
	    fmt.Println(TICKS)

	    fmt.Printf("TICK Time (1 TICK = X milliseconds): ")
	    get_int_csv(rec[2], &TICK_time)
	    fmt.Println(TICK_time)

	    fmt.Printf("lambda: ")
	    get_float64_csv(rec[3], &lambda)
	    fmt.Println(lambda)

	    fmt.Printf("L (bits): ")
	    get_int_csv(rec[4], &L)
	    fmt.Println(L)

	    fmt.Printf("C (bits per sec): ")
	    get_float64_csv(rec[5], &C)
	    fmt.Println(C)

	    fmt.Printf("K (zero = infinity): ")
	    get_int_csv(rec[6], &K)
	    fmt.Println(K)

	} else {
		// Get Variables
		var stdinR = bufio.NewReader(os.Stdin)
		fmt.Printf("M: ")
		get_int(stdinR, &M)

		fmt.Printf("TICKS: ")
		get_int(stdinR, &TICKS)

		fmt.Printf("TICK Time (1 TICK = X milliseconds): ")
		get_int(stdinR, &TICK_time)

		fmt.Printf("lambda: ")
		get_float64(stdinR, &lambda)

		fmt.Printf("L (bits): ")
		get_int(stdinR, &L)

		fmt.Printf("C (bits per sec): ")
		get_float64(stdinR, &C)

		fmt.Printf("K (zero = infinity): ")
		get_int(stdinR, &K)
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
	}

	// Get Variables
	
	//End of display Variables

	// Loop for average statistics
	var test_t time.Time
	for m := 1; m <= M; m++ {
		test_t = time.Now()

		var wait_tick = get_tick_wait()
		var qm = consumer.Init(logger, wait_tick)
		producer.Init(logger, qm, lambda, TICK_time)

		for t := 1; t < TICKS; t++ {
			producer.Tick(t)
			consumer.Tick(t)
		}
		logger.Println("[Info] Finished running Test #", m, "elapsed time:", time.Since(test_t))
		logger.Println("[Info] Queue Size", qm.Size)
	}
}

func get_tick_wait() int {
	return int(math.Ceil(((float64(L) / C) * 1000) / float64(TICK_time)))
}
