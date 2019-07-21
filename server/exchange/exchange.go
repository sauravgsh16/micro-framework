package exchange

import (
	mq "github.com/sauravgsh16/micro-framework/messagequeue"
)

// Exchange routes messages to desired destination
type Exchange struct {
	qtype  string
	queue  *mq.Queue
	sendCh chan mq.Value
	rcvCh  chan mq.Value
	ackCh  chan bool
}

// New returns new exchange
func New() *Exchange {
	q := mq.NewQueue()
	e := &Exchange{
		queue: q,
	}
	return e
}

// Register to be called by the publisher who wants to send
// message to the exchange
func (e *Exchange) Register(qtype string) (chan<- mq.Value, <-chan bool) {
	e.qtype = qtype
	e.sendCh = make(chan mq.Value)
	e.ackCh = make(chan bool)
	return e.sendCh, e.ackCh
}

// Subscribe to be called by the consumer which wants to receive
// message from the exchange
func (e *Exchange) Subscribe(qtype string) <-chan mq.Value {
	e.qtype = qtype
	e.rcvCh = make(chan mq.Value)
	return e.rcvCh
}

// Publish a message to the exchange
func (e *Exchange) Publish(v *mq.Value) {

}
