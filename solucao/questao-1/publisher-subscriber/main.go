package main

import (
	"fmt"
	"time"
	"publisher-subscriber/pubsub"
)

func main() {
	broker := pubsub.NewBroker()

	// Example: Subscribe to a queue
	subscriber := broker.Subscribe("exampleQueue")

	// Example: Run the broker in a separate goroutine
	go func() {
		data := "Hello, PubSub!"
		_ = broker.Run("exampleQueue", data)
	}()

	// Example: Receive messages from the subscriber
	for i := 0; i < 5; i++ {
		select {
		case msg := <-subscriber:
			fmt.Printf("Received message: %v\n", msg)
		case <-time.After(1 * time.Second):
			fmt.Println("Timeout: No message received.")
		}
	}

	// Close the subscriber when done
	broker.CloseAll("exampleQueue")
}
