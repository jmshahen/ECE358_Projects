package stats

import (
	// "log"
	"fmt"
)

type QueueMgr struct {
	Size    float64    // the current size of the queue
	MaxSize float64    // the maximum size of the queue, 0 for infinite
	Head    *QueueItem // the head of the queue
	Tail    *QueueItem // the tail of the queue
}

type QueueItem struct {
	Next *QueueItem
	Data Packet
}

type Packet struct {
	Generated float64
	Finished  float64
}

func (qm *QueueMgr) Push(item Packet) error {
	var qItem = QueueItem{nil, item}

	if qm.Size == 0 {
		qm.Head = &qItem
		qm.Tail = &qItem
		qm.Size++
	} else if qm.MaxSize == 0 || qm.Size < qm.MaxSize {
		qm.Tail.Next = &qItem
		qm.Tail = &qItem
		qm.Size++
	} else {
		//add one to the packet loss
		Probability_loss.AddOne()
		return fmt.Errorf("Queue is full, max size %d packets", qm.MaxSize)
	}
	//successfully put packet in queue
	Probability_loss.Total++
	return nil
}

func (qm *QueueMgr) Pop() (Packet, error) {
	var item Packet
	if qm.Size != 0 {
		item = qm.Head.Data
		qm.Head = qm.Head.Next
		qm.Size--
		return item, nil
	}
	return item, fmt.Errorf("Queue is full, max size %d packets", qm.MaxSize)
}

func (qm *QueueMgr) Clear() {
	qm.Size = 0
	qm.Head = nil
	qm.Tail = nil
}
