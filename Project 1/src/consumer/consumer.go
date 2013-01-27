package consumer

import (
	"log"
)

var logger *log.Logger
var (
	next_tick      int = 1
	wait_time_tick int // in milliseconds
	queue          QueueMgr
)

func Init(l *log.Logger, wait_tick int) *QueueMgr {
	logger = l
	logger.Println("[Consumer] Started")

	wait_time_tick = wait_tick

	logger.Println("[Consumer] Service time is", wait_time_tick, "ticks")

	return &queue
}

func Tick(t int) {
	if next_tick == t {
		consume_packet()
		next_tick = t + wait_time_tick
	}
}

func consume_packet() {
	logger.Println("[Consumer] Consuming Packet for", wait_time_tick, "ticks")
	queue.Pop()
}
