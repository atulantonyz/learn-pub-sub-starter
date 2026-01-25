package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	gamelogic.PrintServerHelp()
	for {
		userinput := gamelogic.GetInput()
		if len(userinput) == 0 {
			continue
		}
		switch userinput[0] {
		case "pause":
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("Error publishing json: %v", err)
			}
			fmt.Println("Pause message sent!")
		case "resume":
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("Error publishing json: %v", err)
			}
			fmt.Println("Resume message sent!")
		case "quit":
			fmt.Println("Exiting!")
			return
		default:
			fmt.Println("unknown command")

		}

	}

}
