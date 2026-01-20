package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		log.Fatalf("Error in marshalling: %v", err)
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "application/json", Body: jsonData})
	if err != nil {
		log.Fatal("Error in publishing message to exchange", err)
	}

	return nil
}
