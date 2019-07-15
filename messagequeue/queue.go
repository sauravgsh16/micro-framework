package messagequeue

import (
	"errors"
	"fmt"
	_ "reflect"
	"sync"
)

//ErrEmptyQueue empty Queue
var ErrEmptyQueue = errors.New("attempt to get message from empty queue")

type Queuer interface {
	EnQueue(v Valuer)
	DeQueue() (v Value, err error)
}

type Valuer interface {
	getId() int
	getVal() string
}

// Value to be distributed
type Value struct {
	Id   int // Id of node which has created the value
	data string
}

func (v Value) getId() int {
	return v.Id
}

func (v Value) getVal() string {
	return fmt.Sprintf("%s", v.data)
}

type Queue struct {
	list *List
	mux  sync.Mutex
}

func NewQueue() *Queue {
	l := newList()
	q := &Queue{
		list: l,
	}
	return q
}

// TODO: NEED TO IMPLEMENT MUTEX FOR WRITING TO AND READING FROM QUEUE
// EnQueue
func (q *Queue) EnQueue(v Valuer) {
	q.list.Append(v)
}

// DeQueue from queue
func (q *Queue) DeQueue() (v Value, err error) {
	v = q.list.Remove()
	if v == (Value{}) {
		return Value{}, ErrEmptyQueue
	}
	return v, nil
}

func (q *Queue) Len() int {
	return q.list.Len()
}

// Traverse the list
func Traverse(qu Queuer) <-chan Valuer {
	switch q := qu.(type) {
	case *Queue:
		cur := q.list.Root
		rcv := make(chan Valuer)
		go func(c *msg) {
			defer close(rcv)
			for c != nil {
				rcv <- c.value
				c = c.next
			}
		}(cur)
		return rcv
	default:
		return nil
	}
}

// TraverseAndRemove from list till empty
func TraverseAndRemove(qu Queuer, done <-chan interface{}) <-chan Valuer {
	switch q := qu.(type) {
	case *Queue:
		rcv := make(chan Valuer)
		go func() {
			defer close(rcv)
			for {
				val := q.list.Remove()
				select {
				case rcv <- val:
				case <-done:
					return
				}
			}
		}()
		return rcv
	default:
		return nil
	}
}
