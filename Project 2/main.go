package main

import (
	"csma_cd"
	"fmt"
	"lan"
	"stats"
)

var (
	M         float64 // number of times to repeat the tests (avg)
	TICKS     float64 // length of the test
	TICK_time float64 // 1 TICK = X milliseconds

	N float64 // Number of computers  connected tot he LAN
	A float64 // Number of packets/sec
	W float64 // The speed of the lan in bits/msec
	L float64 // Length of a packet in bits

	bit_time float64 // the bit time. - time it takes to put 1 bit on the line.
)

func main() {
	fmt.Println("HELLO WORLD\n")

	lan.Init(10, 5, 1000, 1000000, 1500)

}

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
