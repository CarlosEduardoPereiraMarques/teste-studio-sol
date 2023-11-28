package pubsub

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Message representa uma mensagem genérica.
type Message struct {
	Queue string      `json:"queue"`
	Data  interface{} `json:"data"`
}

// Broker gerencia a comunicação entre publishers e subscribers.
type Broker struct {
	sync.RWMutex
	subscribers map[string][]chan interface{}
}

// NewBroker cria um novo Broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Publish publica uma mensagem na fila especificada.
func (b *Broker) Publish(queue string, data interface{}) {
	b.RLock()
	defer b.RUnlock()

	message := Message{Queue: queue, Data: data}

	// Salvar a mensagem em um arquivo
	b.saveMessage(message)

	// Enviar a mensagem para os subscribers
	if subscribers, ok := b.subscribers[queue]; ok {
		for _, ch := range subscribers {
			go func(c chan interface{}) {
				c <- message
			}(ch)
		}
	}
}

// Subscribe se inscreve em uma fila específica.
func (b *Broker) Subscribe(queue string) <-chan interface{} {
	b.Lock()
	defer b.Unlock()

	ch := make(chan interface{})
	b.subscribers[queue] = append(b.subscribers[queue], ch)

	return ch
}

// saveMessage salva a mensagem em um arquivo.
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

// Close fecha o canal para uma fila específica.
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
}

// CloseAll fecha todos os canais para uma fila específica.
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

// CloseAllChannels fecha todos os canais associados ao Broker.
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
