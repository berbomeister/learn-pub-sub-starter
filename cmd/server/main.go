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
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, _ := amqp.Dial(connectionString)
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ successfully.")
	channel, _ := conn.Channel()
	defer channel.Close()
	gamelogic.PrintServerHelp()
	pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DurableQueue,
	)
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		if input[0] == "pause" {
			log.Println("Sending a pause message to the queue")
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			continue
		}
		if input[0] == "resume" {
			log.Println("Sending a resume message to the queue")
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			continue
		}
		if input[0] == "quit" {
			log.Println("Exiting the game")
			break
		}
		log.Printf("Don't understand command %v", input[0])

	}
	fmt.Println("\nShutting down gracefully...")
}
