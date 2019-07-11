package messagequeue

import (
	"errors"
	"fmt"
)

type msg struct {
	next  *msg
	value Value
}

// List - linked
type List struct {
	Root *msg // Pointer to root msg
	len  int  // length of the list
	ch   chan Value
}

// Len of list
func (l *List) Len() int {
	return l.len
}

// NewList points to pointer to a new list
func NewList() *List {
	l := &List{
		ch: make(chan Value),
	}
	return l
}

func (l *List) findLast() *msg {
	cur := l.Root
	for cur.next != nil {
		cur = cur.next
	}
	return cur
}

func (l *List) append(v Value) {
	n := &msg{value: v}
	l.len++
	if l.Root == nil {
		l.Root = n
		return
	}
	last := l.findLast()
	last.next = n
}

func (l *List) remove() (Value, error) {
	if l.Root == nil {
		return "", errors.New("Cannot remove from empty list")
	}
	n := *l.Root
	l.Root = n.next
	n.next = nil // remove reference to next pointer
	l.len--
	return n.value, nil
}

// Append to end of list
func (l *List) Append(v Value) {
	l.append(v)
}

// Remove one msg
func (l *List) Remove() Value {
	d, err := l.remove()
	if err != nil {
		return nil
	}
	return d
}

func (l *List) String() string {
	data := ""
	cur := l.Root
	for cur != nil {
		val := fmt.Sprintf("%v->", cur.value)
		data = data + val
		cur = cur.next
	}
	return fmt.Sprintf("%s, of lenght %d", data, l.len)
}

// Traverse the list
func Traverse(l *List) <-chan Value {
	cur := l.Root
	rcv := make(chan Value)
	go func(c *msg) {
		defer close(rcv)
		for c != nil {
			rcv <- c.value
			c = c.next
		}
	}(cur)
	return rcv
}

// TraverseAndRemove from list till empty
func TraverseAndRemove(l *List, done <-chan interface{}) <-chan Value {
	rcv := make(chan Value)
	go func() {
		defer close(rcv)
		for {
			val := l.Remove()
			select {
			case rcv <- val:
			case <-done:
				return
			}
		}
	}()
	return rcv
}
