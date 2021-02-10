package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"persons.com/api/infrastructure/events"
)

type RabbitMqService struct {
	host   string
	client *amqp.Connection
}

func NewRabbitMqService(host string, Qname string) (events.EventService, error) {

	connection, err := amqp.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	ch, err := connection.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	q, err := ch.QueueDeclare(Qname, false, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "events.rabbitmq.CreateQueue")
	}

	defer ch.Close()

	log.Printf("Queue created: %v", q)

	return &RabbitMqService{
		host:   host,
		client: connection,
	}, nil
}

func (r *RabbitMqService) Publish(exchange string, busName string, value interface{}) error {

	json, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.Publish")
	}

	ch, err := r.client.Channel()
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.NewRabbitMqService")
	}

	err = ch.Publish(exchange, busName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        json,
	})
	if err != nil {
		return errors.Wrap(err, "events.rabbitmq.Publish")
	}

	defer ch.Close()

	return nil
}
