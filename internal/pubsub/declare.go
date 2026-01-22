package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type SimpleQueueType string

const Durable SimpleQueueType = "durable"
const Transient SimpleQueueType = "transient"

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}
	queue, err := ch.QueueDeclare(queueName, queueType == Durable, queueType == Transient, queueType == Transient, false, nil)
	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatalf("Error binding queue: %v", err)
	}
	return ch, queue, nil
}
