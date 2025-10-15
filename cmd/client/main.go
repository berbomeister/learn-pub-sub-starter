package main

import (
	"fmt"
	"os"
	"os/signal"

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
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Error welcoming client:", err)
		return
	}
	channel, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
	)
	if err != nil {
		fmt.Println("Error declaring and binding queue:", err)
		return
	}
	defer channel.Close()

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nShutting down gracefully...")
	fmt.Println("Starting Peril client for user:", username)
}
