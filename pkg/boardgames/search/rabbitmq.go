package search

import (
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/streadway/amqp"
)

func insertQueueItem(ch *amqp.Channel, q *amqp.Queue, item models.Price) error {
	body, err := item.ToGOB64()
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func setupQueue(topic string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, error) {
	uri := config.App.Databases["rabbitmq"]
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return conn, ch, &q, nil
}
