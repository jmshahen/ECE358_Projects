package producer

import (
	"log"
	"math"
	"math/rand"
	"stats"
	"time"
)

var (
	logger      *log.Logger
	nextTick    float64 = 1
	lambda      float64
	qm          *stats.QueueMgr
	time_2_tick float64 //1 tick = X milliseconds
)

func Init(l *log.Logger, _qm *stats.QueueMgr, _lambda float64, TICK_time float64) {
	logger = l
	lambda = _lambda
	qm = _qm
	time_2_tick = TICK_time

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func Tick(t float64) {
	if t == nextTick {
		producePacket(t)
		nextTick = t + getExpRandNum(lambda)
		//logger.Println("[Producer] Next Packet at", nextTick)
	}
}

var debug_getExpRandNum = false

func getExpRandNum(l float64) float64 {
	a := 1 - rand.Float64()
	b := math.Log(a)
	c := (-1 / l)
	d := b * c * 1000 // in milliseconds
	ans := d / float64(time_2_tick)
	if debug_getExpRandNum {
		logger.Println("b", b, "\nc", c, "\nd", d, "\nans", ans, "\n")
	}
	return float64(math.Ceil(ans))
}

func producePacket(i float64) {
	var p = stats.Packet{i, 0}
	if err := qm.Push(p); err != nil {
		logger.Printf("[Producer] Received Error %v\n", err)
	}
}
