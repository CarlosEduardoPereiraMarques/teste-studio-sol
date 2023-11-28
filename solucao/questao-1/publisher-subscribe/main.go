package main

import (
	"fmt"
	"time"
	"publisher-subscriber/pubsub"
)

func main() {
	broker := pubsub.NewBroker()

	subscriber := broker.Subscribe("example_queue")
	defer broker.CloseAll("example_queue")

	go func() {
		for {
			select {
			case msg := <-subscriber:
				fmt.Printf("Received message: %+v\n", msg)
			}
		}
	}()

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Message %d", i)
		broker.Publish("example_queue", message)
		time.Sleep(time.Second)
	}
}
