package messagebus

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type MessageBus struct {
	// Addr is address to AMQP
	Addr   string
	logger *log.Logger
}

func (mb *MessageBus) createQueue(ch *amqp.Channel, name string) {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"notesapp", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" Queue: %s, [x] %s", name, d.Body)
		}
	}()

	log.Printf(fmt.Sprintf(" [*] Waiting for %s.", name))
	<-forever
}

func (mb *MessageBus) Start() {
	conn, err := amqp.Dial(mb.Addr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"notesapp", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	go mb.createQueue(ch, "db2")
	mb.createQueue(ch, "notesview2")
}
