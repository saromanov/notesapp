package messagebus

import (
	"fmt"
	"log"
	"errors"

	"github.com/saromanov/notesapp/client"
	"github.com/saromanov/notesapp/logging"
	"github.com/streadway/amqp"
)

var (
	errEmptyConfig = errors.New("Config is empty")
)

type MessageBus struct {
	// Addr is address to AMQP
	Addr   string
	Exchange  string
	cna client.ClientNotesapp
	logger *logging.Logger
}

// CreateMessageBus provides creates of Messagebus object
func CreateMessageBus(config *Config)(*MessageBus, error) {
	if config == nil {
		return nil, errEmptyConfig
	}
	mb := new(MessageBus)
	mb.Addr = config.Addr
	if config.ExchangeName == "" {
		mb.Exchange = "notesapp"
	} else {
		mb.Exchange = config.ExchangeName
	}
	mb.cna = client.ClientNotesapp{Addr: "http://127.0.0.1:8085/api/incgets"}
	mb.logger = logging.NewLogger(nil)
	return mb, nil

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
	mb.failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		mb.Exchange, // exchange
		false,
		nil)
	mb.failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	mb.failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			mb.processMessages(d.Body)
		}
	}()

	log.Printf(fmt.Sprintf(" [*] Waiting for %s.", name))
	<-forever
}

func (mb *MessageBus) processMessages(msg []byte) {
	if string(msg) == "newget" {
		err := mb.cna.IncGets()
		if err != nil {
			mb.logger.Error(fmt.Sprintf("%v", err))
		} else {
			mb.logger.Info(fmt.Sprintf("RabbitMQ: Event %s is complete", string(msg)))
		}
	}
}

func (mb *MessageBus) Start() {
	conn, err := amqp.Dial(mb.Addr)
	mb.failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mb.failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		mb.Exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	mb.failOnError(err, "Failed to declare an exchange")

	mb.createQueue(ch, "notesview2")
}

func (mb *MessageBus) failOnError(err error, msg string){
	if err != nil {
		mb.logger.Error(fmt.Sprintf("%s: %s", msg, err))
	}
}
