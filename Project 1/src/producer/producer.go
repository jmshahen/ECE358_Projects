package producer

import (
	"consumer"
	"log"
	"math"
	"math/rand"
	"time"
)

var (
	logger   *log.Logger
	nextTick int = 1
	lambda   float64
	qm       *consumer.QueueMgr
)

func Init(l *log.Logger, _qm *consumer.QueueMgr, _lambda float64) {
	logger = l
	lambda = _lambda
	qm = _qm

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func Tick(t int) {
	if t == nextTick {
		producePacket(t)
		nextTick = t + getExpRandNum()
	}
}

func getExpRandNum() int {
	return int((-1 / lambda) * (math.Log(1 - rand.Float64())))
}

func producePacket(i interface{}) {
	var qi = consumer.QueueItem{nil, i}
	if err := qm.Push(qi); err != nil {
		logger.Printf("[Producer] Received Error %v\n", err)
	}
}
