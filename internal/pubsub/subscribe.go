package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackReque
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
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

			retAck := handler(body)
			switch retAck {
			case Ack:
				err = m.Ack(false)

				if err != nil {
					fmt.Printf("Failed to Ack delivery: %v\n", err)
				}
			case NackReque:
				err = m.Nack(false, true)
				if err != nil {
					fmt.Printf("Failed to Nack reque: %v\n", err)
				}
			case NackDiscard:
				err = m.Nack(false, false)
				if err != nil {
					fmt.Printf("Failed to Nack discard: %v\n", err)
				}
			default:
				fmt.Printf("Invalid AckType\n")
			}
		}
	}()

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
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
			buffer := bytes.NewBuffer(m.Body)
			dec := gob.NewDecoder(buffer)
			var body T
			err := dec.Decode(&body)

			if err != nil {
				fmt.Printf("Could not unmarshal message: %v\n", err)
				continue
			}

			retAck := handler(body)
			switch retAck {
			case Ack:
				err = m.Ack(false)

				if err != nil {
					fmt.Printf("Failed to Ack delivery: %v\n", err)
				}
			case NackReque:
				err = m.Nack(false, true)
				if err != nil {
					fmt.Printf("Failed to Nack reque: %v\n", err)
				}
			case NackDiscard:
				err = m.Nack(false, false)
				if err != nil {
					fmt.Printf("Failed to Nack discard: %v\n", err)
				}
			default:
				fmt.Printf("Invalid AckType\n")
			}
		}
	}()

	return nil
}
