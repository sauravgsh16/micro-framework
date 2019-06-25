package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Executing the server")
	go server()
	go client()

	var a string
	fmt.Scanln(&a)
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	}

	err := ch.Publish(
		"",     // exchange name - default here
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		msg)
	failOnError(err, "Failed to publish a message")

}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue string
		"",     // consumer string
		true,   // autoAck bool - to see ACK and NACK
		false,  // exclusive bool
		false,  // noLocal bool
		false,  // no Wait bool
		nil,    // args amqp.Table
	)

	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message : %s", msg.Body)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello", // queue name
		false,   // durable bool
		false,   // autoDelete bool
		false,   // exclusive bool
		false,   // no wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare/create Queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
