package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"persons.com/api/infrastructure/events"
)

type RabbitMqService struct {
	client *amqp.Connection
}

func NewRabbitMqService(host string, Qnames []string, exchanges []string) (events.PersonEventService, error) {

	connection, err := amqp.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	ch, err := connection.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	for _, qname := range Qnames {
		q, err := ch.QueueDeclare(qname, false, false, false, false, nil)
		log.Printf("Queue created: {%v} => | messages: %v | consumers: %v", q.Name, q.Messages, q.Consumers)
		if err != nil {
			return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
		}
	}

	for _, exchange := range exchanges {
		err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
		log.Printf("Exchange created: %v", exchange)
		if err != nil {
			return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
		}
	}

	defer ch.Close()

	return &RabbitMqService{
		client: connection,
	}, nil
}

func (r *RabbitMqService) Publish(exchange string, Qname string, value interface{}) error {

	json, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.Publish")
	}

	ch, err := r.client.Channel()
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	err = ch.Publish(exchange, Qname, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        json,
	})
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.Publish")
	}

	defer ch.Close()

	return nil
}
