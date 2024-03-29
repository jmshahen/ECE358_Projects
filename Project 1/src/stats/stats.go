package stats

import (
	"log"
)

type Avg struct {
	Total float64
	Num   float64
}

type Pro struct {
	Less  float64
	Total float64
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

	//clear values
	Avg_packets.Clear()
	Avg_sojourn.Clear()
	Proportion_idle.Clear()
	Probability_loss.Clear()
}

// Helper functions

func (a Avg) GetAvg() float64 {
	if a.Num == 0 {
		return 0
	}
	return a.Total / a.Num
}

func (a *Avg) AddAvg(add float64) {
	a.Total += add
	a.Num++
}

func (a *Avg) Clear() {
	a.Total = 0
	a.Num = 0
}

func (p Pro) GetProportion() float64 {
	if p.Total == 0 {
		return 0
	}
	return p.Less / p.Total
}

func (p *Pro) AddOne() {
	p.Total++
	p.Less++
}

func (p *Pro) Clear() {
	p.Total = 0
	p.Less = 0
}

func (p Packet) SojournTime() float64 {
	return p.Finished - p.Generated
}
