package messagequeue

import (
	"errors"
	"fmt"
)

type msg struct {
	next  *msg
	value Valuer
}

// List - linked
type List struct {
	Root *msg
	len  int
	ch   chan Valuer
}

// Len of list
func (l *List) Len() int {
	return l.len
}

// NewList points to pointer to a new list
func newList() *List {
	l := &List{
		ch: make(chan Valuer),
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

func (l *List) append(v Valuer) {
	n := &msg{value: v}
	if l.Root == nil {
		l.Root = n
		return
	}
	last := l.findLast()
	last.next = n
}

func (l *List) remove() (Value, error) {
	if l.Root == nil {
		return Value{}, errors.New("Cannot remove from empty list")
	}
	n := *l.Root
	l.Root = n.next
	n.next = nil // remove reference to next pointer
	l.len--
	return n.value.(Value), nil // Type assertion to check n.value is type Value
}

// Append to end of list
func (l *List) Append(v Valuer) {
	l.append(v)
}

// Remove one msg
func (l *List) Remove() Value {
	d, err := l.remove()
	if err != nil {
		return Value{}
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
