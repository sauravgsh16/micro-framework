package handler

import (
	"errors"
	"fmt"
	"log"
	"sync"

	mq "github.com/sauravgsh16/micro-framework/messagequeue"
)

// Messenger interface
type Messenger interface {
	GetID() int
	GetChan() chan mq.Valuer
}

// Broadcaster struct
type Broadcaster struct {
	queue *mq.Queue
	nodes []Messenger
}

// NewBroadcaster pointer to Broadcaster
func NewBroadcaster() *Broadcaster {
	q := mq.NewQueue()
	bq := &Broadcaster{
		queue: q,
		nodes: make([]Messenger, 0),
	}
	fmt.Printf("%+v\n", *bq)
	return bq
}

// EnQueue a message to queue
func (b Broadcaster) EnQueue(v mq.Valuer) {
	b.queue.EnQueue(v)
}

// DeQueue a message from queue
func (b Broadcaster) DeQueue() (mq.Value, error) {
	return b.queue.DeQueue()
}

// Send to all listening nodes
// MAINTAIN MANY BROADCAST QUEUES - WITH IDs OF MESSAGE - think
func (b *Broadcaster) Send() int {
	v, err := b.DeQueue()
	if err != nil {
		log.Fatalf("Broadcast from empty Queue")
	}
	return b.broadcast(v)
}

// SendAll messages from list to all nodes
func (b *Broadcaster) SendAll() int {
	msgCt := b.queue.Len()
	for i := 0; i < b.queue.Len(); i++ {
		b.Send()
	}
	return msgCt
}

func (b *Broadcaster) broadcast(v mq.Valuer) int {
	ct := 0
	var mux sync.Mutex
	for _, m := range b.nodes {
		go func(m Messenger) {
			select {
			case m.GetChan() <- v:
				mux.Lock()
				ct++
				mux.Unlock()
			default:
				fmt.Printf("Failed to send msg for node %d", m.GetID())
			}
		}(m)
	}
	return ct
}

// AddNode to add
func (b *Broadcaster) AddNode(m Messenger) error {
	if b.contains(m) {
		return errors.New("node already present")
	}
	b.nodes = append(b.nodes, m)
	return nil
}

func (b *Broadcaster) contains(m Messenger) bool {
	for _, v := range b.nodes {
		if v.GetID() == m.GetID() {
			return true
		}
	}
	return false
}
