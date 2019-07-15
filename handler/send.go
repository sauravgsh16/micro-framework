package handler

import (
	"fmt"

	mq "github.com/sauravgsh16/micro-framework/messagequeue"
)

// Sender struct
type Sender struct {
	queue      *mq.Queue
	messengers map[int]Messenger
}

// NewSender points to new send queue
func NewSender() *Sender {
	q := mq.NewQueue()
	s := &Sender{
		queue:      q,
		messengers: make(map[int]Messenger),
	}
	return s
}

// EnQueue message to Queue
func (s Sender) EnQueue(v mq.Valuer) {
	s.queue.EnQueue(v)
}

// DeQueue message from Queue
func (s Sender) DeQueue() (mq.Value, error) {
	return s.queue.DeQueue()
}

// Addnode for sending data
func (s *Sender) Addnode(id int, m Messenger) {
	if _, ok := s.messengers[id]; ok {
		return
	}
	s.messengers[id] = m
}

func (s *Sender) send(v mq.Valuer) (int, error) {
	var id int
	switch val := v.(type) {
	case mq.Value:
		id = val.Id
	default:
		return 0, fmt.Errorf("message type incorrect: %s", val)
	}
	m, ok := s.messengers[id]
	if !ok {
		return 0, fmt.Errorf("node not present in send queue: %d", id)
	}
	go func(m Messenger, v mq.Valuer) {
		select {
		case m.GetChan() <- v:
		default:
			// Need to implement logger with io.Writer interface
			// Use Fprintf - to write error to logger
			fmt.Printf("msg not sent to: %d", id)
		}
	}(m, v)
	return 0, nil
}

// Send root msg from queue
func (s *Sender) Send() (int, error) {
	val, err := s.queue.DeQueue()
	if err != nil {
		return 0, err
	}
	return s.send(val)
}

// SendAll message present in the queue
