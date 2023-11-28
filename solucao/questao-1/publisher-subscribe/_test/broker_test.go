package test

import (
	"sync"
	"testing"
	"time"
	"publisher-subscriber/pubsub"
)

func TestBroker(t *testing.T) {
	broker := pubsub.NewBroker()
	queue := "test_queue"

	var wg sync.WaitGroup

	// Adiciona ao WaitGroup o número total de goroutines que serão criadas
	wg.Add(10)

	// Subscribing
	for i := 0; i < 5; i++ {
		go func(index int) {
			defer wg.Done() // Indica que a goroutine terminou

			ch := broker.Subscribe(queue)
			for {
				select {
				case msg := <-ch:
					if msg == nil {
						// A goroutine será encerrada se receber uma mensagem nula
						t.Logf("Goroutine %d received nil, closing", index)
						return
					}
					t.Logf("Goroutine %d received message: %+v", index, msg)
				}
			}
		}(i)
	}

	// Publishing
	for i := 0; i < 5; i++ {
		go func(index int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				broker.Publish(queue, j)
			}
		}(i)
	}

	// Aguarda um tempo antes de encerrar a produção de mensagens
	time.Sleep(2 * time.Second)

	// Encerra a produção de mensagens
	broker.CloseAll(queue)

	// Aguarda até que todas as goroutines tenham terminado
	wg.Wait()
}
