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
	STATE_START               = 0
	STATE_MEDIUM_SENSING_INIT = 1
	STATE_MEDIUM_SENSING      = 2
	STATE_STANDARD_WAIT       = 3
	STATE_TRANSMIT_FRAME      = 4
	STATE_JAMMING_SIGNAL      = 5
	STATE_EXP_BACKOFF         = 6
	STATE_EXP_WAIT            = 7
	STATE_ERROR_SEND          = 8
	STATE_SUCCESS_SEND        = 9
)

type StateMachine struct {
	state int
}

type CSMA struct {
	id          int64
	qm          stats.QueueMgr
	logger      *log.Logger
	nextTick    int64
	curTick     int64
	lambda      float64
	time_2_tick float64 //1 tick = X nanoseconds
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
func (csma *CSMA) Init(id int64, lan *lan.LAN, logger *log.Logger, _lambda float64, TICK_time float64, kmax int64, tp int64, medium_sense_time int64) {
	csma.id = id
	csma.logger = logger
	csma.time_2_tick = TICK_time
	csma.lambda = _lambda
	csma.lan = lan
	csma.kmax = kmax
	csma.tp = tp
	csma.medium_sense_time = medium_sense_time

	csma.qm.Clear()
	csma.nextTick = -1
	csma.state.state = STATE_START

	logger.Println("[ CSMA /CD ] Started")

	rand.Seed(time.Now().UnixNano())
}

//The main logic for the CSMA Computer
//producers packets and then enters the CSMA_CD State Machine
func (csma *CSMA) Tick(t int64) {
	csma.curTick = t

	//produce packets
	if t >= csma.nextTick {
		if csma.nextTick != -1 {
			csma.producePacket()
		}
		csma.getExpRandNum()
		//csma.logger.Println("[ Comp", csma.id, "] Next Tick", csma.nextTick)
	}

	switch csma.state.state {
	case STATE_START:
		csma.i = 0
		if csma.qm.Size != 0 {
			//csma.logger.Println("[ Comp", csma.id, "] STATE_START: Packet Readying")
			csma.packet, _ = csma.qm.Pop()
			csma.packet.ExitQueue = t
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_MEDIUM_SENSING_INIT:
		csma.waitFor = t + csma.medium_sense_time
		csma.state.state = STATE_MEDIUM_SENSING
		fallthrough
	case STATE_MEDIUM_SENSING:
		//sense medium
		if csma.lan.Sense_line(csma.id) {
			//csma.logger.Println("[ Comp", csma.id, "] STATE_MEDIUM_SENSING: LAN Busy")
			csma.waitFor = t + csma.tp
			csma.state.state = STATE_STANDARD_WAIT
		} else if csma.waitFor == t {
			//ready to send file
			//csma.logger.Println("[ Comp", csma.id, "] STATE_MEDIUM_SENSING: Sending Packet")
			csma.waitFor = t + csma.lan.Put_packet(&csma.packet, csma.id, t)
			csma.state.state = STATE_TRANSMIT_FRAME
		}
	case STATE_STANDARD_WAIT:
		if csma.waitFor == t {
			csma.state.state = STATE_MEDIUM_SENSING_INIT
		}
	case STATE_TRANSMIT_FRAME:
		//sense medium
		if csma.lan.Sense_line(csma.id) {
			//csma.logger.Println("[ Comp", csma.id, "] STATE_TRANSMIT_FRAME: Collision Detected")
			csma.waitFor = t + csma.lan.Send_jam_signal(csma.id, t)
			csma.state.state = STATE_JAMMING_SIGNAL
		} else if csma.waitFor == t {
			//done waiting
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
	case STATE_SUCCESS_SEND:
		//record successful packet
		//csma.logger.Println(csma.id, "Correctly Sent a Packet")
		csma.state.state = STATE_START
	case STATE_ERROR_SEND:
		//record failed packet
		csma.lan.Record_lost_packet(&csma.packet, t)
		csma.logger.Println(csma.id, "Failed to send packet", csma.kmax, "times")
		csma.state.state = STATE_START
	}
}

func (csma *CSMA) producePacket() {
	var p = stats.Packet{csma.id, csma.curTick, 0, 0}
	if err := csma.qm.Push(p); err != nil {
		//logger.Printf("[Producer] Received Error %v\n", err)
	}
	//logger.Println("[Producer] Packet produced")
}

//Exponentially Random num that will set when the next packet should be produced
func (csma *CSMA) getExpRandNum() {
	a := 1 - rand.Float64()
	b := math.Log(a)
	c := (-1 / csma.lambda)
	d := b * c                              // sec / packet
	ans := sec_to_tick(d, csma.time_2_tick) // tick / packet

	csma.nextTick = csma.curTick + ans
}

//convert seconds into ticks
func sec_to_tick(s float64, TICK_time float64) int64 {
	return int64(math.Ceil((s * 1e9) / TICK_time))
}

//The exponential backoff wait time, when a collision is detected
func (csma *CSMA) expBackOff() int64 {
	return int64(rand.Int31n(2^int32(csma.i)-1)) * csma.tp
}
