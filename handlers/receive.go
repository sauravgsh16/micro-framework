package messengers

import (
	"github.com/sauravgsh16/src/constants"
	"github.com/sauravgsh16/src/rabbit"
	"github.com/streadway/amqp"
)

type receiveBroadCast struct {
	ch *amqp.Channel
}

type receivePeer struct{}

func (rb *receiveBroadCast) Receive(exchangeName string) {
	r := rabbit.New(constants.RABBIT_SERVER)

	defer r.Conn.Close()
	defer r.Ch.Close()

}

func (rb *receiveBroadCast) declareExchange() {

}

// NewPeerReceiver
func NewPeerReceiver() receivePeer {
	pr := &receivePeer{}
	return *pr
}

// NewBroadReceiver
func NewBroadReceiver() receiveBroadCast {
	br := &receiveBroadCast{}
	return *br
}
