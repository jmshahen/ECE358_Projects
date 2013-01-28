package stats

import (
	"log"
)

type Avg struct {
	Total float64
	Num   int
}

type Pro struct {
	Less  int
	Total int
}

var (
	queue  *QueueMgr
	logger *log.Logger
	// Stats
	avg_packets      Avg // E[N]
	avg_sojourn      Avg // E[T]
	proportion_idle  Pro // P_IDLE
	probability_loss Pro // P_LOSS
)

func Init(l *log.Logger, qm *QueueMgr) {
	logger = l
	queue = qm
}

func (a Avg) GetAvg() float64 {
	if a.Num == 0 {
		return 0
	}
	return a.Total / float64(a.Num)
}

func (p Pro) GetProportion() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(p.Less) / float64(p.Total)
}
