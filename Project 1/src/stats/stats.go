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
	logger *log.Logger
	// Stats
	Avg_packets      Avg // E[N]
	Avg_sojourn      Avg // E[T]
	Proportion_idle  Pro // P_IDLE
	Probability_loss Pro // P_LOSS
)

func Init(l *log.Logger) {
	logger = l
	logger.Println("[Stats] Started")
}

// Helper functions

func (a Avg) GetAvg() float64 {
	if a.Num == 0 {
		return 0
	}
	return a.Total / float64(a.Num)
}

func (a *Avg) AddAvg(add float64) {
	a.Total += add
	a.Num++
}

func (p Pro) GetProportion() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(p.Less) / float64(p.Total)
}

func (p *Pro) AddOne() {
	p.Total++
	p.Less++
}

func (p Packet) SojournTime() float64 {
	return float64(p.Finished - p.Generated)
}
