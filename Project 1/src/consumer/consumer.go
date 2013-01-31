package consumer

import (
	"log"
	"stats"
)

var logger *log.Logger
var (
	next_tick      float64 = 1
	wait_time_tick float64 // in milliseconds
	queue          stats.QueueMgr
)

func Init(l *log.Logger, wait_tick float64) *stats.QueueMgr {
	logger = l
	logger.Println("[Consumer] Started")

	wait_time_tick = wait_tick

	logger.Println("[Consumer] Service time is", wait_time_tick, "ticks")

	return &queue
}

func Tick(t float64) {
	if next_tick <= t {
		consume_packet(t)
	}
}

func consume_packet(t float64) {
	//logger.Println("[Consumer] Consuming Packet for", wait_time_tick, "ticks")
	if queue.Size > 0 {
		var p stats.Packet

		p, err := queue.Pop()

		if err != nil {
			logger.Fatalf("Recieved Error From Pop: %v", err)
		}

		next_tick = t + wait_time_tick
		stats.Proportion_idle.Total++

		p.Finished = t
		stats.Avg_sojourn.AddAvg(p.SojournTime())
	} else {
		stats.Proportion_idle.AddOne()
	}
}
