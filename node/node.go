package node

import (
	"fmt"

	"github.com/sauravgsh16/micro-framework/handler"
	mq "github.com/sauravgsh16/micro-framework/messagequeue"
)

const (
	_ = iota
	leader
	member
)

var (
	// ErrInvalidEcu invalid ecu type supplied
	ErrInvalidEcu = "invalid ecu type supplied: %d"
	// ErrNoQueue invalid search
	ErrNoQueue = "no %s queue"
)

// Nodes stores all the nodes which has been instantiated
type Nodes []Node

// Node information of an node
type Node struct {
	ID      int
	ecutype int
	Ch      chan mq.Valuer
	queues  map[string]mq.Queuer
}

// GetID return node's Id
func (n Node) GetID() int {
	return n.ID
}

// GetChan returns node channel
func (n Node) GetChan() chan mq.Valuer {
	return n.Ch
}

// NewNode returns instance of a new node
func NewNode(id, t int) (*Node, error) {
	n := &Node{
		ID:      id,
		ecutype: t,
		Ch:      make(chan mq.Valuer),
		queues:  make(map[string]mq.Queuer),
	}
	err := n.instantiate()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Node) instantiate() error {
	switch n.ecutype {
	case leader:
		instantiateLeader(n)
	case member:
		instantiateMember(n)
	default:
		return fmt.Errorf(ErrInvalidEcu, n.ecutype)
	}
	return nil
}

func (n *Node) appendhandlers(qs ...mq.Queuer) {
	for _, q := range qs {
		switch q.(type) {
		case *handler.Broadcaster:
			n.queues["broadcast"] = q
		case *handler.Sender:
			n.queues["send"] = q
		case *handler.Receiver:
			n.queues["receive"] = q
		}
	}
}

func instantiateLeader(n *Node) {
	// Needs one broadcast, send and receive Queue
	bq := handler.NewBroadcaster()
	rx := handler.NewReceiver()
	tx := handler.NewSender()
	n.appendhandlers(bq, rx, tx)
}

func instantiateMember(n *Node) {
	// Implementing just the receive Queue for now
	rx := handler.NewReceiver()
	n.queues["receive"] = rx
}

// GetQueue returns the queue from the handler
func (n *Node) GetQueue(name string) (mq.Queuer, error) {
	q, ok := n.queues[name]
	if !ok {
		return nil, fmt.Errorf(ErrNoQueue, name)
	}
	return q, nil
}

// AddReceiver adds receiver nodes
func (n *Node) AddReceiver(node Node, name string) error {
	qr, err := n.GetQueue(name)
	if err != nil {
		return err
	}
	q, ok := qr.(handler.Broadcaster)
	fmt.Printf("%+v, %t, %t", q, qr, ok)
	return nil
}
