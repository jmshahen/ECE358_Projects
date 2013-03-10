package main

import (
	"csma_cd"
	"fmt"
	"lan"
	"stats"
)

var (
	M         int64 // number of times to repeat the tests (avg)
	TICKS     int64 // length of the test
	TICK_time int64 // 1 TICK = X milliseconds

	N_start, N_end, N_step int64 // Number of computers  connected to the LAN
	A                      int64 // Number of packets/sec
	W                      int64 // The speed of the lan in bits/sec
	L                      int64 // Length of a packet in bits

	bit_time int64 // the bit time. - time it takes to put 1 bit on the line.

	lan         *lan.LAN
	cmsa        *csma_cd.CSMA
	bucket      *stats.Bucket
	lost_bucket *stats.Bucket

	Prop_ticks         int64
	Packet_trans_ticks int64
	Jam_trans_ticks    int64

	kmax              int64
	tp                int64
	medium_sense_time int64

	computers 		[]csma_cd.CSMA 

	logger 			*log.Logger


	Avg_Avg_Full_Delay       		stats.Avg
    Avg_Avg_Full_Delay_per_Comp  	[]stats.Avg
    Avg_Avg_Queue_Delay          	stats.Avg
    Avg_Avg_Queue_Delay_per_Comp 	[]stats.Avg
    Avg_Avg_CSMA_Delay           	stats.Avg
    Avg_Avg_CSMA_Delay_per_Comp  	[]stats.Avg

    Sum_throughput					float64
    Avg_throughput					float64

    csv_cols int = 11
)

func main() {
	//Header
	fmt.Println("ECE 358 Project 2 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// End of Header
	write_csv_header()
	// gets and sets all global varaibles not dependant on N
	get_variables()

	var test_t time.Time
	var test_b time.Time = time.Now()
	var tsince time.Duration
	for i := N_start; i < N_End; i+=N_step {

		for m := 1.0; m <= M; m++ {
			// initilize components
			bucket.Init(logger, L, i)
			lost_bucket.Init(logger, L, i)

			init_computers(i)

			lan.Init(i, Prop_ticks, Packet_trans_ticks, Jam_trans_ticks, bucket, lost_bucket)
			// end initialize components

			for t := 0; t < TICKS; t++ {
				for c := range computers {
					computers[c].Tick(t)
					lan.Complete_Tick(t)
				}
			}

			// compute stats
			Avg_Avg_Full_Delay.AddAvg(bucket.Avg_Full_Delay.GetAvg())
			Avg_Avg_Queue_Delay.AddAvg(bucket.Avg_Queue_Delay.GetAvg())
			Avg_Avg_CSMA_Delay.AddAvg(bucket.Avg_CSMA_Delay.GetAvg())
			Sum_throughput += bucket.Throughput()

			for a := range Avg_Avg_Full_Delay_per_Comp {
				Avg_Avg_Full_Delay_per_Comp[a].AddAvg(bucket.Avg_Full_Delay_per_Comp[a].GetAvg())
				Avg_Avg_Queue_Delay_per_Comp[a].AddAvg(bucket.Avg_Queue_Delay_per_Comp[a].GetAvg())
				Avg_Avg_CSMA_Delay_per_Comp[a].AddAvg(bucket.Avg_CSMA_Delay_per_Comp[a].GetAvg())
			}
			// end compute stats
		}
		Avg_throughput = (Sum_throughput / m)
		write_csv_output(i)
	}

}

func init_computers(N int64) {

	computers = make([]csma_cd.CSMA, N, N)

	for i:= 0; i < N; i++ {
		computers[].Init(i, lan, logger, A, TICK_time, kmax, tp, medium_sense_time )
	}
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
		headers[i] = "Number of Comps"
		i++
		headers[i] = "A : Packets per sec"
		i++
		headers[i] = "L : Packet length (bits)"
		i++
		headers[i] = "W : bits per secound"
		i++

		headers[i] = "Avg_Full_Delay: "
		i++
		headers[i] = "Avg_Queue_Delay"
		i++
		headers[i] = "Avg_CSMA_Delay"
		i++
		headers[i] = "Throughtput"
		i++

		writter.Write(headers)
		writter.Flush()
		file.Close()
	}
}

func write_csv_output(num_comps int64) {

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
	rec[i] = strconv.FormatFloat(num_comps, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(A, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(L, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(W, 'f', -1, 64)
	i++

	rec[i] = strconv.FormatFloat(Avg_Avg_Full_Delay.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_Queue_Delay_per_Comp.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_CSMA_Delay_per_Comp.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Throughput, 'f', -1, 64)
	i++

	writter.Write(rec)
	writter.Flush()
}