package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
)

func main() {
	AMQP_URL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(AMQP_URL)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		os.Exit(1)
	}
	fmt.Println("Connection successful")
	defer conn.Close()
	fmt.Println("Starting Peril server...")
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Interrupt received, shutting down program")

}
