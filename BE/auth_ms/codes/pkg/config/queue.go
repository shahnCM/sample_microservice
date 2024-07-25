package config

import (
	"fmt"
	"os"
)

type queueConfig = struct {
	Connection string
	QueueName string
}

func GetQueueConfig() *queueConfig{
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQUsername := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	amqpConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUsername, rabbitMQPassword, rabbitMQHost, rabbitMQPort)
	return &queueConfig {
		Connection: amqpConnectionString,
		QueueName: os.Getenv("RABBITMQ_QUEUE"),
	}
}