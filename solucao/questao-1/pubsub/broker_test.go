package pubsub

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type TestBrokerWithError struct {
	*Broker
}

func (t *TestBrokerWithError) Subscribe(queue string) <-chan interface{} {
	ch := make(chan interface{})
	close(ch)
	return ch
}

func isChannelTest(value interface{}) bool {
	switch value.(type) {
	case chan interface{}, chan int, chan string, chan bool:
		return true
	default:
		return false
	}
}

func TestBroker(t *testing.T) {
	broker := NewBroker()

	t.Run("Publish_HappyPath", func(t *testing.T) {
		subscriber := broker.Subscribe("test_queue")
		defer broker.CloseAllChannels()

		done := make(chan struct{})
		defer close(done)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Goroutine panicked: %v", r)
				}
			}()

			for {
				select {
				case msg := <-subscriber:
					if msg == nil {
						t.Error("Received nil, closing")
						return
					}

					t.Logf("Received message: %+v\n", msg)
				case <-done:
					return
				}
			}
		}()

		message := "Test Message"
		broker.Publish("test_queue", message)


		time.Sleep(time.Second)
	})

	t.Run("Publish_NoSubscribers", func(t *testing.T) {
		message := "Test Message"
		broker.Publish("nonexistent_queue", message)

	})

	t.Run("Subscribe_HappyPath", func(t *testing.T) {
		subscriber := broker.Subscribe("test_queue")
		defer broker.CloseAllChannels()
		if subscriber == nil {
			t.Error("Subscribe did not return a channel")
		}
	})

	t.Run("Subscribe_ExistingSubscriber", func(t *testing.T) {
		subscriber1 := broker.Subscribe("existing_queue")
		defer broker.CloseAllChannels()

		subscriber2 := broker.Subscribe("existing_queue")

		if subscriber1 == subscriber2 {
			t.Error("Subscribe returned the same channel for an existing queue")
		}
	})

	t.Run("Run_HappyPath", func(t *testing.T) {
		subscriber := broker.Run("test_queue", "Run Message")
		if subscriber == nil {
			t.Error("Run did not return a channel")
		}

		time.Sleep(time.Second)
	})

	t.Run("Run_ErrorPath", func(t *testing.T) {
		testBroker := &TestBrokerWithError{
			Broker: NewBroker(),
		}

		subscriber := testBroker.Subscribe("test_queue")

		select {
		case msg, ok := <-subscriber:
			if ok || msg != nil {
				t.Error("Run should return an error during subscription")
			}
		case <-time.After(time.Second):
			t.Error("Run did not return an error during subscription")
		}
	})

	t.Run("saveMessage_HappyPath", func(t *testing.T) {
		message := Message{Queue: "test_queue", Data: "Test Data"}
		broker.saveMessage(message)
		filename := fmt.Sprintf("%s.txt", message.Queue)
		_, err := os.Stat(filename)
		if err != nil {
			t.Errorf("saveMessage did not create the file: %v", err)
		}
		os.Remove(filename)
	})

	t.Run("saveMessage_ErrorPath", func(t *testing.T) {
		data := struct {
			Queue string      `json:"queue"`
			Data  interface{} `json:"data"`
		}{
			Queue: "test_queue",
			Data:  "Test Data",
		}
	
		if !isChannelTest(data.Data) {
			broker.saveMessage(data)
		} else {
			t.Log("Data is a channel, skipping serialization")
		}
	})

	t.Run("Close_HappyPath", func(t *testing.T) {
		subscriber := broker.Subscribe("test_queue")
		broker.Close("test_queue", subscriber)
	})

	t.Run("Close_ErrorPath", func(t *testing.T) {
		broker.Close("nonexistent_queue", make(chan interface{}))
	})

	t.Run("CloseAll_HappyPath", func(t *testing.T) {
		broker.Subscribe("test_queue")
		broker.Subscribe("test_queue")
		broker.CloseAll("test_queue")
	})

	t.Run("CloseAll_ErrorPath", func(t *testing.T) {
		broker.CloseAll("nonexistent_queue")
	})

	t.Run("CloseAllChannels_HappyPath", func(t *testing.T) {
		broker.Subscribe("test_queue1")
		broker.Subscribe("test_queue2")
		broker.CloseAllChannels()
	})

	t.Run("CloseAllChannels_ErrorPath", func(t *testing.T) {
		subscriber1 := broker.Subscribe("test_queue1")
		subscriber2 := broker.Subscribe("test_queue2")

		simulateCloseError := func(queue string, ch <-chan interface{}) {
			t.Error("Simulated error during channel close")
			broker.Close(queue, ch)
		}

		broker.SetCloseFunction(simulateCloseError)

		broker.CloseAllChannels()

		<-subscriber1
		<-subscriber2
	})
}
