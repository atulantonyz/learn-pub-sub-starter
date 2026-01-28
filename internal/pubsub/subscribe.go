package pubsub

import (
	"encoding/json"
	"fmt"

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
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		defer ch.Close()
		for m := range msgs {
			var body T
			err := json.Unmarshal(m.Body, &body)
			if err != nil {
				fmt.Printf("Could not unmarshal message: %v\n", err)
				continue
			}

			handler(body)
			err = m.Ack(false)
			if err != nil {
				fmt.Printf("Failed to Ack delivery: %v\n", err)
			}
		}
	}()

	return nil
}
