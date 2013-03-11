package main

import (
	"csma_cd"
	"encoding/csv"
	"fmt"
	"lan"
	"log"
	"os"
	"stats"
	"strconv"
	"time"
)

var (
	M         int64   // number of times to repeat the tests (avg)
	TICKS     int64   // length of the test
	TICK_time float64 // 1 TICK = X nanoseconds

	N_start, N_end, N_step int64   // Number of computers  connected to the LAN
	A                      float64 // Number of packets/sec
	W                      float64 // The speed of the lan in bits/sec
	L                      int64   // Length of a packet in bits

	lan_v       lan.LAN
	cmsa        csma_cd.CSMA
	bucket      stats.Bucket
	lost_bucket stats.Bucket

	length_of_line    float64 //in meters
	speed_over_line   float64 //in meters per sec
	jam_signal_length int64   //in bits

	Prop_ticks         int64
	Packet_trans_ticks int64
	Jam_trans_ticks    int64

	kmax              int64
	tp                int64
	medium_sense_time int64

	computers []csma_cd.CSMA

	logger *log.Logger

	Avg_Avg_Full_Delay           stats.Avg
	Avg_Avg_Full_Delay_per_Comp  []stats.Avg
	Avg_Avg_Queue_Delay          stats.Avg
	Avg_Avg_Queue_Delay_per_Comp []stats.Avg
	Avg_Avg_CSMA_Delay           stats.Avg
	Avg_Avg_CSMA_Delay_per_Comp  []stats.Avg
	Avg_throughput               stats.Avg
	Avg_throughput_per_Comp      []stats.Avg

	csv_cols int64 = 11
)

func main() {
	//Header
	fmt.Println("ECE 358 Project 2 - Written in GO (golang.org)")
	fmt.Println("Submitted by:")
	fmt.Println("\tJon Shahen    \t(20334465)")
	fmt.Println("\tKevin Carlton \t(20337152)")
	fmt.Println("---\n")
	// gets and sets all global varaibles not dependant on N
	get_variables()

	// End of Header
	write_csv_header(N_end)

	var test_t time.Time
	var test_b time.Time = time.Now()
	var tsince time.Duration
	for i := N_start; i <= N_end; i += N_step {

		// run M times for each i value.
		for m := int64(1); m <= M; m++ {
			test_t = time.Now()
			// initilize components
			bucket.Init(logger, L, i)
			lost_bucket.Init(logger, L, i)

			init_computers(i)

			lan_v.Init(i, Prop_ticks, Packet_trans_ticks, Jam_trans_ticks, &bucket, &lost_bucket)
			// end initialize components

			//must start at 1
			for t := int64(1); t <= TICKS; t++ {
				for c := range computers {
					computers[c].Tick(t)
				}
				lan_v.Complete_Tick(t)
			}

			// compute stats
			Avg_Avg_Full_Delay.AddAvg(bucket.Avg_Full_Delay.GetAvg())

			Avg_Avg_Queue_Delay.AddAvg(bucket.Avg_Queue_Delay.GetAvg())

			Avg_Avg_CSMA_Delay.AddAvg(bucket.Avg_CSMA_Delay.GetAvg())

			Avg_throughput.AddAvg(bucket.Throughput(TICKS))

			// Probably only needed for testing.
			for a := range Avg_Avg_Full_Delay_per_Comp {
				Avg_Avg_Full_Delay_per_Comp[a].AddAvg(bucket.Avg_Full_Delay_per_Comp[a].GetAvg())
				Avg_Avg_Queue_Delay_per_Comp[a].AddAvg(bucket.Avg_Queue_Delay_per_Comp[a].GetAvg())
				Avg_Avg_CSMA_Delay_per_Comp[a].AddAvg(bucket.Avg_CSMA_Delay_per_Comp[a].GetAvg())
				Avg_throughput_per_Comp[a].AddAvg(bucket.Throughput_per_comp(int64(a), TICKS))
			}
			// end compute stats

			tsince = time.Since(test_t)
			fmt.Println("[Stats] Elapsed Time                        =\t", tsince)
		}
		write_csv_output(i)

	}
	tsince = time.Since(test_b)

	fmt.Println("\n---\n")
	fmt.Println("[Stats] Total Elapsed Time                         =\t", tsince)

}

func init_computers(N int64) {

	Avg_Avg_Full_Delay_per_Comp = make([]stats.Avg, N, N)
	Avg_Avg_Queue_Delay_per_Comp = make([]stats.Avg, N, N)
	Avg_Avg_CSMA_Delay_per_Comp = make([]stats.Avg, N, N)
	Avg_throughput_per_Comp = make([]stats.Avg, N, N)

	computers = make([]csma_cd.CSMA, N, N)

	for i := int64(0); i < N; i++ {
		computers[i].Init(i, &lan_v, logger, A, TICK_time, kmax, tp, medium_sense_time)
	}
}

func write_csv_header(max_comps int64) {
	_, err := os.Open("test_out.csv")
	if os.IsNotExist(err) {
		file, _ := os.Create("test_out.csv")
		writter := csv.NewWriter(file)

		var i = 0
		headers := make([]string, csv_cols+4*max_comps)
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

		var d = i
		for c := int64(0); c < max_comps; c, d = c+1, d+4 {
			s := strconv.FormatInt(c, 10)
			headers[d+0] = "Avg_Full_Delay (" + s + ")"
			headers[d+1] = "Avg_Queue_Delay (" + s + ")"
			headers[d+2] = "Avg_CSMA_Delay (" + s + ")"
			headers[d+3] = "Throughtput (" + s + ")"
		}

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

	rec := make([]string, csv_cols+4*num_comps)
	var i = 0
	rec[i] = strconv.FormatInt(M, 10)
	i++
	rec[i] = strconv.FormatInt(TICKS, 10)
	i++
	rec[i] = strconv.FormatFloat(TICK_time, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatInt(num_comps, 10)
	i++
	rec[i] = strconv.FormatFloat(A, 'f', -1, 64)
	i++
	rec[i] = strconv.FormatInt(L, 10)
	i++
	rec[i] = strconv.FormatFloat(W, 'f', -1, 64)
	i++

	rec[i] = strconv.FormatFloat(Avg_Avg_Full_Delay.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_Queue_Delay.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_Avg_CSMA_Delay.GetAvg(), 'f', -1, 64)
	i++
	rec[i] = strconv.FormatFloat(Avg_throughput.GetAvg(), 'f', -1, 64)
	i++

	// clear stats.
	Avg_Avg_Full_Delay.Clear()
	Avg_Avg_Queue_Delay.Clear()
	Avg_Avg_CSMA_Delay.Clear()
	Avg_throughput.Clear()

	var d = i
	for c := int64(0); c < num_comps; c, d = c+1, d+4 {
		rec[d+0] = strconv.FormatFloat(Avg_Avg_Full_Delay_per_Comp[c].GetAvg(), 'f', -1, 64)
		rec[d+1] = strconv.FormatFloat(Avg_Avg_Queue_Delay_per_Comp[c].GetAvg(), 'f', -1, 64)
		rec[d+2] = strconv.FormatFloat(Avg_Avg_CSMA_Delay_per_Comp[c].GetAvg(), 'f', -1, 64)
		rec[d+3] = strconv.FormatFloat(Avg_throughput_per_Comp[c].GetAvg(), 'f', -1, 64)

		Avg_Avg_Full_Delay_per_Comp[c].Clear()
		Avg_Avg_Queue_Delay_per_Comp[c].Clear()
		Avg_Avg_CSMA_Delay_per_Comp[c].Clear()
		Avg_throughput_per_Comp[c].Clear()
	}

	writter.Write(rec)
	writter.Flush()
}
