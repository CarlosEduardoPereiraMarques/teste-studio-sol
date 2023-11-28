package main

import (
	"fmt"
	"publisher-subscribe/pubsub"
	"time"
)

func main() {
	ps := pubsub.NewPubSub()

	queue := "example"
	msg := "hello"

	ps.Publish(queue, msg)

	ch := ps.Subscribe(queue)

	select {
	case receivedMsg := <-ch:
		fmt.Printf("Received message: %s\n", receivedMsg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: No message received.")
	}

	err := ps.SaveToFile("output.json")
	if err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}
}
