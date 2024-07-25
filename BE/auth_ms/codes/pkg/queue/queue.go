package queue

import (
	"auth_ms/pkg/config"

	"github.com/streadway/amqp"
)

var rabbitMQ RabbitMQ

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
}

func Init() error {
	queueName := config.GetQueueConfig().QueueName
	connection, err := amqp.Dial(config.GetQueueConfig().Connection)

	if err != nil {
		return err
	}

	channel, err := connection.Channel()

	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	rabbitMQ = RabbitMQ{
		Connection: connection,
		Channel:    channel,
		QueueName:  queueName,
	}

	return nil
}

func Publish(messageBody string) error {
	return rabbitMQ.Channel.Publish(
		"",                 // exchange
		rabbitMQ.QueueName, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(messageBody),
			DeliveryMode: amqp.Persistent,
		})
}
