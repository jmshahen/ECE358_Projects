package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		lambda  = 250
		incr    = 50
		incr_to = 11
		K       = []int{10, 25, 50}
		header  = "M,TICKS,TICK Time (1 TICK = X milliseconds),Lambda,L (bits),C (bits per sec),K (zero = infinity)\n"
		s       string
		bob     int
	)

	for i := 0; i < incr_to; i++ {
		for _, k := range K {
			bob++
			s = fmt.Sprintf("5,6e9,1e-3,%d,2000,1.00E+06,%d\n", (lambda + i*incr), k)
			write_to_file(bob, header, s)
		}
	}
}

func write_to_file(b int, header string, s string) {
	file, _ := os.Create(fmt.Sprintf("args_%0d.csv", b))
	file.WriteString(header)
	file.WriteString(s)
	file.Close()
}
