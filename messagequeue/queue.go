package messagequeue

import (
	"errors"
	"fmt"
	"sync"
)

// Value to be distributed
type Value interface{}

// Node - Individual components which communicate
// TODO: NEED TO PLACE THIS IN APPROPRIATE NODES MODULE
type Node struct {
	id int
	ch chan Value
}

type queue struct {
	list *List
	mux  sync.Mutex
}

// BroadcastQueue struct
type BroadcastQueue struct {
	queue
	nodes []Node
}

// TODO: NEED TO IMPLEMENT MUTEX FOR WRITING TO AND READING FROM QUEUE
// Put
func (q *queue) EnQueue(v Value) {
	// TODO: NEED TO IMPLEMENT MUTEX FOR WRITING TO AND READING FROM QUEUE
	q.list.Append(v)
}

// DeQueue from queue
func (q *queue) DeQueue() (Value, error) {
	v := q.list.Remove()
	if v == nil {
		return nil, errors.New("attempt to get message from empty queue")
	}
	return v, nil
}

// NewBroadcast pointer to BroadcastQueue
func NewBroadcast() *BroadcastQueue {
	bq := &BroadcastQueue{
		queue{},
		make([]Node, 0),
	}
	return bq
}

// AddNode to add
func (b *BroadcastQueue) AddNode(n Node) error {
	if b.contains(n) {
		return errors.New("node already present")
	}
	b.nodes = append(b.nodes, n)
	return nil
}

func (b *BroadcastQueue) contains(n Node) bool {
	for _, v := range b.nodes {
		if v.id == n.id {
			return true
		}
	}
	return false
}

// Send to all listening nodes
// MAINTAIN MANY BROADCAST QUEUES - WITH IDs OF MESSAGE - think
func (b *BroadcastQueue) Send() int {
	msg, err := b.queue.DeQueue()
	if err != nil {
		panic("Broadcast from empty Queue")
	}
	return b.broadcast(msg)
}

// SendAll messages from list to all nodes
func (b *BroadcastQueue) SendAll() int {
	msgCt := b.list.Len()
	for i := 0; i < b.list.Len(); i++ {
		b.Send()
	}
	return msgCt
}

func (b *BroadcastQueue) broadcast(msg Value) int {
	ct := 0
	var mux sync.Mutex
	for _, n := range b.nodes {
		go func(n Node) {
			select {
			case n.ch <- msg:
				mux.Lock()
				ct++
				mux.Unlock()
			default:
				fmt.Printf("Failed to send msg for node %d", n.id)
			}
		}(n)
	}
	return ct
}

// SendQueue struct
type SendQueue struct {
	queue
	nodes []Node
}

// NewSendQueue points to new send queue
func NewSendQueue() *SendQueue {
	sq := &SendQueue{
		queue{},
		make([]Node, 0),
	}
	return sq
}

// EnQueue in send queue
func (sq *SendQueue) EnQueue(id int, msg string) {
	val := map[int]string{
		id: msg,
	}
	sq.queue.EnQueue(val)
}
