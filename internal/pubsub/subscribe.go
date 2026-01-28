package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType)
	if err != nil {
		return err
	}
	delivery_ch, err := ch.Consume(queueName, "", false, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		for m := range delivery_ch {
			var body T
			err := json.Unmarshal(m.Body, &body)
			if err != nil {
				log.Fatalf("Failed to marshal message body: %v", err)
			}

			handler(body)
			err = m.Ack(false)
			if err != nil {
				log.Fatalf("Failed to Ack", err)
			}
		}
	}()

	return nil
}
