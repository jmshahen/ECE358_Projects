package csma_cd

import (
	"lan"
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
	id          int64
	qm          stats.QueueMgr
	logger      *log.Logger
	nextTick    float64
	curTick     float64
	lambda      float64
	time_2_tick float64 //1 tick = X milliseconds
	lan         *lan.LAN

	state             StateMachine
	waitFor           int64
	i                 int64
	kmax              int64
	packet            stats.Packet
	tp                int64
	medium_sense_time int64
}

//Variables that are set by calling func at Init:
//	logger, lambda, time_2_tick, kmax, tp, medium_sense_time
func (CSMA *csma) Init(id int64, lan *lan.LAN, logger *log.Logger,
	_lambda float64, TICK_time float64, kmax int64,
	tp int64, medium_sense_time int64) {

	csma.id = id
	csma.logger = logger
	csma.time_2_tick = TICK_time
	csma.lambda = _lambda
	csma.lan = lan
	csma.kmax = kmax
	csma.tp = tp
	csma.medium_sense_time = medium_sense_time

	csma.qm = stats.QueueMgr
	csma.nextTick = -1
	csma.state.state = STATE_START

	logger.Println("[Producer] Started")

	rand.Seed(time.Now().UnixNano())
}

func (CSMA *csma) Tick(t int64) {
	csma.curTick = t

	//produce packets
	if t >= nextTick {
		if nextTick != -1 {
			csma.producePacket(t)
		}
		csma.getExpRandNum()
		csma.logger.Println("[Comp", csma.id, "] Next Tick", csma.nextTick)
	}

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
		if csma.lan.sense_line(csma.id) {
			csma.waitFor = t + csma.tp
			csma.state.state = STATE_STANDARD_WAIT
		} else if csma.waitFor == t {
			csma.state.state = STATE_TRANSMIT_FRAME
		}
	case STATE_STANDARD_WAIT:
		if csma.waitFor == t {
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_TRANSMIT_FRAME:
		//sense medium
		if csma.lan.sense_line(csma.id) {
			//collision
			csma.waitFor = t + csma.lan.send_jam_signal(csma.id, t)
			csma.state.state = STATE_JAMMING_SIGNAL
		} else if csma.waitFor == t {
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
		csma.logger.Println(csma.id, "Correctly Sent a Packet")
		csma.state.state = STATE_START
	case STATE_SUCCESS_SEND:
		//record successful packet
		csma.logger.Println(csma.id, "Failed to send packet", csma.kmax, "times")
		csma.state.state = STATE_START
	}
}

func (CSMA *csma) producePacket() {
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
	return rand.Int31n(2**csma.i-1) * csma.tp
}
