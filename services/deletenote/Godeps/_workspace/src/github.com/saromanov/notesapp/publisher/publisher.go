package publisher

import (
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn *amqp.Connection
}

func NewPublisher(exchangename, addr string) (*Publisher, error) {
	var err error
	pub := new(Publisher)
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangename, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return nil, err
	}

	pub.conn = conn

	return pub, nil

}

func (pub *Publisher) Send(exchange, msg string) error {
	ch, err := pub.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(msg),
		})

	if err != nil {
		return err
	}

	return nil
}

func (pub *Publisher) Close() {
	pub.conn.Close()
}
