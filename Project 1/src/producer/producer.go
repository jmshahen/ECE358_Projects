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
	nextTick    int = 1
	lambda      float64
	qm          *stats.QueueMgr
	time_2_tick int //1 tick = X milliseconds
)

func Init(l *log.Logger, _qm *stats.QueueMgr, _lambda float64, TICK_time int) {
	logger = l
	lambda = _lambda
	qm = _qm
	time_2_tick = TICK_time

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func Tick(t int) {
	if t == nextTick {
		producePacket(t)
		nextTick = t + getExpRandNum(lambda)
		//logger.Println("[Producer] Next Packet at", nextTick)
	}
}

var debug_getExpRandNum = false

func getExpRandNum(l float64) int {
	a := 1 - rand.Float64()
	b := math.Log(a)
	c := (-1 / l)
	d := b * c * 1000 // in milliseconds
	ans := d / float64(time_2_tick)
	if debug_getExpRandNum {
		logger.Println("b", b, "\nc", c, "\nd", d, "\nans", ans, "\n")
	}
	return int(math.Ceil(ans))
}

func producePacket(i interface{}) {
	var qi = stats.QueueItem{nil, i}
	if err := qm.Push(qi); err != nil {
		logger.Printf("[Producer] Received Error %v\n", err)
	}
}
