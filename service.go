package notesapp

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	errRabbitMQConnect = errors.New("Failed to connect to RabbitMQ")
)

// Service provides implementation of basic Service
type Service struct {
	Title      string
	Addr       string
	Port       string
	connection *amqp.Connection
	running    chan bool
}

func (service *Service) CreateService(title string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return errRabbitMQConnect
	}
	service.connection = conn
}

// Received messages from event bus (AMQP)
func (service *Service) Received(item string) {
	//Received
}

// Publish message to event bus (AMQP)
func (service *Service) Publish(name, msg string) error {
	ch, err := service.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		id,       // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}
	err = ch.Publish(
		id,    // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(msg),
		})
	//failOnError(err, "Failed to publish a message")

	if err != nil {
		return err
	}

}

// Start set of service is alive
func (service *Service) Start() {
	go func() {
		service.running <- true
	}()
}

// Stop provides off service
func (service *Service) Stop() {
	service.connection.Close()
}
