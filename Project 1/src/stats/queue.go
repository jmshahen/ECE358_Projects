package stats

import (
	// "log"
	"fmt"
)

type QueueMgr struct {
	Size    int        // the current size of the queue
	MaxSize int        // the maximum size of the queue, 0 for infinite
	Head    *QueueItem // the head of the queue
	Tail    *QueueItem // the tail of the queue
}

type QueueItem struct {
	Next *QueueItem
	Data interface{}
}

func (qm *QueueMgr) Push(item interface{}) error {
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
		return fmt.Errorf("Queue is full, max size %d packets", qm.MaxSize)
	}
	return nil
}

func (qm *QueueMgr) Pop() (interface{}, error) {
	if qm.Size != 0 {
		var item interface{}
		item = qm.Head.Data
		qm.Head = qm.Head.Next
		qm.Size--
		return item, nil
	}
	return nil, fmt.Errorf("Queue is full, max size %d packets", qm.MaxSize)
}
