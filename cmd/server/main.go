package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const connString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatal("couldn't connect to RabbitMQ server: ", err)
	}
	defer conn.Close()
	fmt.Println("Connection to RabbitMQ was successful")

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatal("could not create channel: ", err)
	}

	err = pubsub.PublishJSON(
		pubCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		})
	if err != nil {
		fmt.Println("could not publish data: ", err)
	}
	fmt.Println("Message sent")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	fmt.Println("Server is shutting down")
}
