package consumer

import (
	"log"
	"stats"
)

var logger *log.Logger
var (
	next_tick      int = 1
	wait_time_tick int // in milliseconds
	queue          stats.QueueMgr
)

func Init(l *log.Logger, wait_tick int) *stats.QueueMgr {
	logger = l
	logger.Println("[Consumer] Started")

	wait_time_tick = wait_tick

	logger.Println("[Consumer] Service time is", wait_time_tick, "ticks")

	return &queue
}

func Tick(t int) {
	if next_tick <= t {
		consume_packet(t)
	}
}

func consume_packet(t int) {
	//logger.Println("[Consumer] Consuming Packet for", wait_time_tick, "ticks")
	if queue.Size > 0 {
		queue.Pop()
		next_tick = t + wait_time_tick
		stats.Proportion_idle.Total++
	} else {
		stats.Proportion_idle.AddOne()
	}
}
