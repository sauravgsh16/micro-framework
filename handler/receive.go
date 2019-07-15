package handler

import (
	mq "github.com/sauravgsh16/micro-framework/messagequeue"
)

// Receiver - struct
type Receiver struct {
	queue *mq.Queue
	// Id    int
	// rx    chan MsgObj
}

// NewReceiver - receive queue
// TO THINK IF SEPERATE RECEIVE QUEUES OR SEND QUEUES ARE REQUIRED FOR
// INDIVIDUAL NODES
func NewReceiver() *Receiver {
	q := mq.NewQueue()
	rq := &Receiver{
		queue: q,
	}
	return rq
}

// EnQueue message from queue
func (r Receiver) EnQueue(v mq.Valuer) {
	r.queue.EnQueue(v)
}

// DeQueue message from queue
func (r Receiver) DeQueue() (mq.Value, error) {
	return r.queue.DeQueue()
}

// Receive returns a channel
func (r *Receiver) Receive(done <-chan interface{}) chan<- mq.Valuer {
	rx := make(chan mq.Valuer)
	go func() {
		defer close(rx)
		for {
			select {
			case <-done:
				return
			case v := <-rx:
				r.queue.EnQueue(v)
			}
		}
	}()
	return rx
}

/*
// Process each message in the receive queue
func (rq *Receiver) Process() {
        done := make(chan interface{})
        rcv := mq.TraverseAndRemove(rq, done)
        for d := range rcv {
                switch d.(type) {
                case nil:
                        close(done)
                case mq.Value:
                        // Here is where we process the message according to
                        // the data
                        fmt.Printf("Printing %v", d)
                }
        }
}
*/
