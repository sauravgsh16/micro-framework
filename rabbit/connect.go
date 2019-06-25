package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

type connectionPrimitives struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// Connection Rabbit Connection object
type Connection struct {
	connectionPrimitives
	server string
}

func (r *Connection) setConnectionPrimitives() {
	var err error
	r.Conn, err = amqp.Dial(r.server)
	failOnError(err, "Failed to connect to RabbitMQ")

	r.Ch, err = r.Conn.Channel()
	failOnError(err, "Failed to acquire the connection Channel primitive")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}

// New returns a connection object to Rabbit
func New(server string) Connection {
	c := &Connection{
		server: server,
	}
	c.setConnectionPrimitives()
	return *c
}
