package main

import (
	"fmt"

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

	gameState := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gameState),
	)
loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			gameState.CommandSpawn(input)
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println("Successfully moved unit")
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming is not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break loop
		default:
			fmt.Println("Unkown commands. Type help to see possible commands!")
		}
	}

	fmt.Println("\nShutting down gracefully...")
}
