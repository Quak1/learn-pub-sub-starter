package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
		log.Fatal("Error creating channel: ", err)
	}

	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerGameLog,
	)
	if err != nil {
		log.Fatal("could not create and bind queue: ", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("Sending pause message...")
			err := sendPause(pubCh, true)
			if err != nil {
				fmt.Println("Error sending pause message: ", err)
			}
		case "resume":
			fmt.Println("Sending resume message...")
			err := sendPause(pubCh, false)
			if err != nil {
				fmt.Println("Error sending resume message: ", err)
			}
		case "quit":
			fmt.Println("Server is shutting down")
			return
		default:
			fmt.Println("I don't understand this command")
		}
	}
}

func sendPause(ch *amqp.Channel, isPaused bool) error {
	return pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: isPaused,
		})
}
