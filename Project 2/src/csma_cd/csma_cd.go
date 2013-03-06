package csma_cd

import (
	"log"
	"math"
	"math/rand"
	"stats"
	"time"
)

type CSMA struct {
	qm          stats.QueueMgr
	logger      *log.Logger
	nextTick    float64
	curTick     float64
	lambda      float64
	time_2_tick float64 //1 tick = X milliseconds
}

func (CSMA *csma) Init(l *log.Logger, _lambda float64, TICK_time float64) {
	csma.logger = l
	csma.lambda = _lambda
	csma.qm = stats.QueueMgr
	csma.time_2_tick = TICK_time
	csma.nextTick = 1

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func (CSMA *csma) Tick(t float64) {
	csma.curTick = t
	if t == nextTick {
		producePacket()
		getExpRandNum()
		//logger.Println("[Producer] Next Tick", nextTick)
	}
}

func (CSMA *csma) producePacket(i float64) {
	var p = stats.Packet{csma.curTick, 0}
	if err := csma.qm.Push(p); err != nil {
		//logger.Printf("[Producer] Received Error %v\n", err)
	}
	//logger.Println("[Producer] Packet produced")
}

func (CSMA *csma) getExpRandNum() {
	a := 1 - rand.Float64()
	b := math.Log(a)
	c := (-1 / csma.lambda)
	d := b * c * 1000 // in milliseconds
	ans := d / float64(csma.time_2_tick)

	csma.nextTick = csma.curTick + math.Ceil(ans)
}
