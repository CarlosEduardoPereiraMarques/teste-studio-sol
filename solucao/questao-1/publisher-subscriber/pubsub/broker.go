package pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Message struct {
	Queue string      `json:"queue"`
	Data  interface{} `json:"data"`
}

type Broker struct {
    sync.RWMutex
    subscribers  map[string][]chan interface{}
    closeFunction func(string, <-chan interface{})
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan interface{}),
	}
}


func (b *Broker) Publish(queue string, data interface{}) {
	b.Lock()
	defer b.Unlock()

	channels, ok := b.subscribers[queue]
	if !ok {
		log.Printf("Queue '%s' does not exist\n", queue)
		return
	}

	for _, ch := range channels {
		select {
		case ch <- data:
		default:
			log.Println("Channel is full. Discarding message.")
		}
	}

	if !isChannel(data) {
		b.saveMessage(Message{Queue: queue, Data: data})
	}
}


func isChannel(value interface{}) bool {
	switch value.(type) {
	case chan interface{}, chan int, chan string, chan bool:
		return true
	default:
		return false
	}
}


func (b *Broker) Subscribe(queue string) <-chan interface{} {
	b.Lock()
	defer b.Unlock()

	ch := make(chan interface{})

	b.subscribers[queue] = append(b.subscribers[queue], ch)
	return ch
}

func (b *Broker) Run(queue string, data interface{}) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		subscriber := b.Subscribe(queue)
		defer b.Close(queue, subscriber)

		b.Publish(queue, data)
	}()

	return ch
}

func (b *Broker) saveMessage(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Erro ao serializar a mensagem: %v\n", err)
		return
	}

	filename := fmt.Sprintf("%s.txt", message.Queue)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("Erro ao escrever no arquivo: %v\n", err)
	}
}

func (b *Broker) Close(queue string, ch <-chan interface{}) {
    b.Lock()
    defer b.Unlock()

    subscribers, ok := b.subscribers[queue]
    if !ok {
        return
    }

    for i, c := range subscribers {
        if c == ch {
            close(c)
            b.subscribers[queue] = append(subscribers[:i], subscribers[i+1:]...)
            break
        }
    }

    if b.closeFunction != nil {
        b.closeFunction(queue, ch)
    }
}

func (b *Broker) CloseAll(queue string) {
	b.Lock()
	defer b.Unlock()

	subscribers, ok := b.subscribers[queue]
	if !ok {
		return
	}

	for _, c := range subscribers {
		close(c)
	}

	delete(b.subscribers, queue)
}

func (b *Broker) CloseAllChannels() {
	b.Lock()
	defer b.Unlock()

	for queue, subscribers := range b.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
		delete(b.subscribers, queue)
	}
}

func (b *Broker) SetCloseFunction(closeFn func(string, <-chan interface{})) {
    b.Lock()
    defer b.Unlock()
    b.closeFunction = closeFn
}