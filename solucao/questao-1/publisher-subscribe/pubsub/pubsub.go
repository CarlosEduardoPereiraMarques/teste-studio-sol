package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Message struct {
	Queue string `json:"queue"`
	Data  string `json:"data"`
}

type PubSub struct {
	mu       sync.Mutex
	channels map[string]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[string]chan string),
	}
}

func (ps *PubSub) Publish(queue string, data string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.channels[queue]; !ok {
		ps.channels[queue] = make(chan string, 100)
	}

	ps.channels[queue] <- data
}

func (ps *PubSub) Subscribe(queue string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.channels[queue]; !ok {
		ps.channels[queue] = make(chan string, 100)
	}

	return ps.channels[queue]
}

func (ps *PubSub) SaveToFile(filename string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var messages []Message

	for queue, channel := range ps.channels {
		for {
			select {
			case data, ok := <-channel:
				if !ok {
					break
				}
				messages = append(messages, Message{Queue: queue, Data: data})
			default:
				break
			}
		}
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
