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

type Bucket struct {
	logger *log.Logger
	// Stats
	Avg_Full_Delay           Avg
	Avg_Full_Delay_per_Comp  []Avg
	Avg_Queue_Delay          Avg
	Avg_Queue_Delay_per_Comp []Avg
	Avg_CSMA_Delay           Avg
	Avg_CSMA_Delay_per_Comp  []Avg
	packet_len               int64
	Packets                  QueueMgr
}

func (Bucket *b) Init(l *log.Logger, packet_len int64, num_comps int64) {
	b.logger = l
	b.logger.Println("[Stats] Started")
	b.Avg_Delay_per_Comp = make([]Avg, N, N)

	b.packet_len = packet_len
}

func (Bucket *b) Throughput(total_ticks int64) float64 {
	return float64(b.Packets.Size) * float64(b.packet_len) / float64(total_ticks)
}

func (Bucket *b) Accept_packet(p *Packet) {
	//calculate delay
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
