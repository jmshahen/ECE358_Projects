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
	packets_per_Comp         []int64
	packet_len               int64
	Packets                  QueueMgr
}

func (b *Bucket) Init(l *log.Logger, packet_len int64, num_comps int64) {
	N := num_comps
	b.logger = l
	b.logger.Println("[Stats] Started")

	b.Avg_Full_Delay.Clear()
	b.Avg_Full_Delay_per_Comp = make([]Avg, N, N)
	b.Avg_Queue_Delay.Clear()
	b.Avg_Queue_Delay_per_Comp = make([]Avg, N, N)
	b.Avg_CSMA_Delay.Clear()
	b.Avg_CSMA_Delay_per_Comp = make([]Avg, N, N)

	b.packets_per_Comp = make([]int64, N, N)

	b.packet_len = packet_len
	b.Packets.Clear()
}

//returned in: bits / tick (# bits per tick)
func (b *Bucket) Throughput(total_ticks int64) float64 {
	return float64(b.Packets.Size) * float64(b.packet_len) / float64(total_ticks)
}

func (b *Bucket) Throughput_per_comp(compID int64, total_ticks int64) float64 {
	return float64(b.packets_per_Comp[compID]) * float64(b.packet_len) / float64(total_ticks)
}

func (b *Bucket) Accept_packet(p *Packet) {
	//calculate delay
	id := p.CompID

	var delay float64 = float64(p.Finished - p.Generated)
	b.Avg_Full_Delay.AddAvg(delay)
	b.Avg_Full_Delay_per_Comp[id].AddAvg(delay)

	delay = float64(p.ExitQueue - p.Generated)
	b.Avg_Queue_Delay.AddAvg(delay)
	b.Avg_Queue_Delay_per_Comp[id].AddAvg(delay)

	delay = float64(p.Finished - p.ExitQueue)
	b.Avg_CSMA_Delay.AddAvg(delay)
	b.Avg_CSMA_Delay_per_Comp[id].AddAvg(delay)

	b.packets_per_Comp[id]++

	b.Packets.Push(*p)

	b.logger.Println("[ Bucket ] Accepted Packet")
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
	return float64(p.Finished - p.Generated)
}
