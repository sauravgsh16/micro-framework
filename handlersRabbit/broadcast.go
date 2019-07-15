package messengers

import (
	"github.com/sauravgsh16/src/constants"
	"github.com/sauravgsh16/src/rabbit"
	"github.com/streadway/amqp"
)

// BroadCast type
type broadCast struct {
	exchange string
}

// BroadCastMessage  braodcasts a message on an exchange
func (b *broadCast) BroadCastMessage(msg amqp.Publishing, exchangeName string) error {
	r := rabbit.New(constants.RABBIT_SERVER)

	defer r.Conn.Close()
	defer r.Ch.Close()

	err := r.Ch.Publish(
		exchangeName, // exchange name
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}

// New broadcaster
func New(exchange string) broadCast {
	b := &broadCast{
		exchange: exchange,
	}
	return *b
}
