package consumer

import (
	"log"
	"time"
)

var logger *log.Logger
var (
	packet_len        int     // packet length in bits
	transmission_rate float32 // bits per second
	wait_time         int     // in milliseconds
)

func Init(l *log.Logger, _packet_len int, _transmission_rate float32) {
	logger = l
	logger.Println("[Consumer] Started")
	packet_len = _packet_len
	transmission_rate = _transmission_rate

	wait_time = int((float32(packet_len) / transmission_rate) * 1000) // in milliseconds
	logger.Println("[Consumer] Service time is", wait_time, "milliseconds")
}

func Tick(t int) {

}

func consume_packet() {
	logger.Println("[Consumer] Consuming Packet for", wait_time, "milliseconds")
	time.Sleep(time.Duration(wait_time) * time.Millisecond)
}
