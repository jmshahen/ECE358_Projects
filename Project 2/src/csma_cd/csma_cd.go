package csma_cd

import (
	"log"
	"math"
	"math/rand"
	"stats"
	"time"
)

const (
	STATE_START = 0
	STATE_MEDIUM_SENSING_INIT
	STATE_MEDIUM_SENSING
	STATE_STANDARD_WAIT
	STATE_TRANSMIT_FRAME
	STATE_JAMMING_SIGNAL
	STATE_EXP_BACKOFF
	STATE_EXP_WAIT
	STATE_ERROR_SEND
	STATE_SUCCESS_SEND
)

type StateMachine struct {
	state int
}

type CSMA struct {
	qm          stats.QueueMgr
	logger      *log.Logger
	nextTick    float64
	curTick     float64
	lambda      float64
	time_2_tick float64 //1 tick = X milliseconds

	state             StateMachine
	waitFor           int64
	i                 int
	kmax              int64
	packet            stats.Packet
	tp                int64
	medium_sense_time int64
}

func (CSMA *csma) Init(l *log.Logger, _lambda float64, TICK_time float64) {
	csma.logger = l
	csma.lambda = _lambda
	csma.qm = stats.QueueMgr
	csma.time_2_tick = TICK_time
	csma.nextTick = 1
	csma.state.state = STATE_START

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func (CSMA *csma) Tick(t int64) {
	csma.curTick = t

	switch csma.state.state {
	case STATE_START:
		i = 0
		if csma.qm.Size != 0 {
			csma.packet = csma.qm.Pop()
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_MEDIUM_SENSING_INIT:
		csma.waitFor = t + csma.medium_sense_time
		csma.state.state = STATE_MEDIUM_SENSING
		fallthrough
	case STATE_MEDIUM_SENSING:
		//sense medium
		//if busy:
		//		csma.waitFor = t + csma.tp
		//		goto STATE_STANDARD_WAIT
		//if csma.waitFor == t goto STATE_TRANSMIT_FRAME
	case STATE_STANDARD_WAIT:
		if csma.waitFor == t {
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_TRANSMIT_FRAME:
		//sense medium
		//if busy
		//	collision
		//	csma.waitFor = t + lan.sendJammingSignal()
		//	csma.state.state = STATE_JAMMING_SIGNAL
		if csma.waitFor == t {
			csma.state.state = STATE_SUCCESS_SEND
		}
	case STATE_JAMMING_SIGNAL:
		if csma.waitFor == t {
			csma.i++
			if csma.i > csma.kmax {
				csma.state.state = STATE_ERROR_SEND
			} else {
				csma.waitFor = t + csma.expBackOff()
				csma.state.state = STATE_EXP_BACKOFF
			}
		}
	case STATE_EXP_BACKOFF:
		if csma.waitFor == t {
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_ERROR_SEND:
		//record failed packet
		csma.state.state = STATE_START
	case STATE_SUCCESS_SEND:
		//record successful packet
		csma.state.state = STATE_START
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

func (CSMA *cs) expBackOff() int64 {

}
