package main

import (
	"bufio"
	"fmt"
	"os"
	// "strconv"
	"log"
	"strings"
)

// Input Variables
var (
	M      int //number of times to repeat the tests (avg)
	TICKS  int //
	lambda int //
	L      int //
	C      int //
	K      int //
)

var logger *log.Logger

func get_int(r *bufio.Reader) (int) {
	var val, err = r.ReadString('\n')
	if err != nil {
		fmt.Println("| err:", err)
		os.Exit(1)
	}

	logger.Println("Val =", val)
	var trimval = strings.TrimRight(val, "\n")
	logger.Println("trimmed =", trimval)
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
	fmt.Println("ECE 358 Project 1 - Written in GO")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen \t(20334465)")
	fmt.Println("\tKevin Carlton \t(203)")
	fmt.Println("---\n")
	// End of Header

	logger = log.New(os.Stdout, "[ECE358 P1]", log.LstdFlags)

	// Get Variables
	var stdinR = bufio.NewReader(os.Stdin)
	fmt.Printf("M: ")
	M = get_int(stdinR)

	fmt.Printf("TICKS: ")
	TICKS = get_int(stdinR)

	fmt.Printf("lambda: ")
	lambda = get_int(stdinR)

	fmt.Printf("L: ")
	L = get_int(stdinR)

	fmt.Printf("C: ")
	C = get_int(stdinR)

	fmt.Printf("K (zero = infinity): ")
	K = get_int(stdinR)
	// End of Get Variables

	// Display Variables
	fmt.Println("Variables being used:")
	fmt.Println("\t M     ", M)
	fmt.Println("\t TICKS ", TICKS)
	fmt.Println("\t lambda", lambda)
	fmt.Println("\t L     ", L)
	fmt.Println("\t C     ", C)
	if K == 0 {
		fmt.Println("\t K      Infinity")
	} else {
		fmt.Println("\t K     ",K)
	}
	//End of display Variables

	// Loop for average statistics
	for m := 1; m <= M; m++ {
		for t := 1; t < TICKS; t++ {
			consumer(t)
			producer(t)

		}
	}
}
