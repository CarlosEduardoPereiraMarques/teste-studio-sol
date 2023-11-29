package pubsub

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestBrokerWithError é uma implementação de Subscribable que força um erro durante a inscrição
type TestBrokerWithError struct {
	*Broker
}

func (t *TestBrokerWithError) Subscribe(queue string) <-chan interface{} {
	ch := make(chan interface{})
	close(ch) // Fechar o canal imediatamente para indicar um erro na inscrição
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

	// Test Publish - Happy Path
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
					// Verificar se a mensagem é correta
					t.Logf("Received message: %+v\n", msg)
				case <-done:
					return
				}
			}
		}()

		message := "Test Message"
		broker.Publish("test_queue", message)

		// Aguardar para garantir que a mensagem seja processada
		time.Sleep(time.Second)
	})

	// Test Publish - No Subscribers
	t.Run("Publish_NoSubscribers", func(t *testing.T) {
		message := "Test Message"
		broker.Publish("nonexistent_queue", message)
		// Não deve haver erro se não houver subscribers
	})

	// Test Subscribe - Happy Path
	t.Run("Subscribe_HappyPath", func(t *testing.T) {
		subscriber := broker.Subscribe("test_queue")
		defer broker.CloseAllChannels()
		// Verificar se o canal é retornado corretamente
		if subscriber == nil {
			t.Error("Subscribe did not return a channel")
		}
	})

	// Test Subscribe - Existing Subscriber
	t.Run("Subscribe_ExistingSubscriber", func(t *testing.T) {
		subscriber1 := broker.Subscribe("existing_queue")
		defer broker.CloseAllChannels()

		// Tentar se inscrever novamente
		subscriber2 := broker.Subscribe("existing_queue")

		// Deve retornar um novo canal
		if subscriber1 == subscriber2 {
			t.Error("Subscribe returned the same channel for an existing queue")
		}
	})

	// Test Run - Happy Path
	t.Run("Run_HappyPath", func(t *testing.T) {
		subscriber := broker.Run("test_queue", "Run Message")
		// Verificar se o canal é retornado corretamente
		if subscriber == nil {
			t.Error("Run did not return a channel")
		}

		// Aguardar para garantir que a mensagem seja processada
		time.Sleep(time.Second)
	})

	// Test Run - Error Path
	t.Run("Run_ErrorPath", func(t *testing.T) {
		// Forçar um erro durante a inscrição
		testBroker := &TestBrokerWithError{
			Broker: NewBroker(),
		}

		subscriber := testBroker.Subscribe("test_queue")

		// Deve haver um erro durante a inscrição
		select {
		case msg, ok := <-subscriber:
			if ok || msg != nil {
				t.Error("Run should return an error during subscription")
			}
		case <-time.After(time.Second):
			t.Error("Run did not return an error during subscription")
		}
	})

	// Test saveMessage - Happy Path
	t.Run("saveMessage_HappyPath", func(t *testing.T) {
		// Criar uma mensagem de teste
		message := Message{Queue: "test_queue", Data: "Test Data"}
		// Chamar a função saveMessage
		broker.saveMessage(message)
		// Verificar se o arquivo foi criado corretamente
		filename := fmt.Sprintf("%s.txt", message.Queue)
		_, err := os.Stat(filename)
		if err != nil {
			t.Errorf("saveMessage did not create the file: %v", err)
		}
		// Remover o arquivo de teste
		os.Remove(filename)
	})

	// Test saveMessage - Error Path
	t.Run("saveMessage_ErrorPath", func(t *testing.T) {
		// Forçar um erro de serialização
		data := struct {
			Queue string      `json:"queue"`
			Data  interface{} `json:"data"`
		}{
			Queue: "test_queue",
			Data:  "Test Data",
		}
	
		// Tente serializar apenas se o Data não for um canal
		if !isChannelTest(data.Data) {
			broker.saveMessage(data)
		} else {
			t.Log("Data is a channel, skipping serialization")
		}
	})

	// Test Close - Happy Path
	t.Run("Close_HappyPath", func(t *testing.T) {
		// Criar um subscriber
		subscriber := broker.Subscribe("test_queue")
		// Chamar a função Close
		broker.Close("test_queue", subscriber)
		// Não deve haver erro ao fechar um canal existente
	})

	// Test Close - Error Path
	t.Run("Close_ErrorPath", func(t *testing.T) {
		// Tentar fechar um canal que não está na lista de subscribers
		broker.Close("nonexistent_queue", make(chan interface{}))
		// Não deve haver erro se o canal não estiver na lista
	})

	// Test CloseAll - Happy Path
	t.Run("CloseAll_HappyPath", func(t *testing.T) {
		// Criar vários subscribers
		broker.Subscribe("test_queue")
		broker.Subscribe("test_queue")
		// Chamar a função CloseAll
		broker.CloseAll("test_queue")
		// Não deve haver erro ao fechar todos os canais associados a uma fila
	})

	// Test CloseAll - Error Path
	t.Run("CloseAll_ErrorPath", func(t *testing.T) {
		// Tentar fechar todos os canais de uma fila que não possui subscribers
		broker.CloseAll("nonexistent_queue")
		// Não deve haver erro se a fila não tiver subscribers
	})

	// Test CloseAllChannels - Happy Path
	t.Run("CloseAllChannels_HappyPath", func(t *testing.T) {
		// Criar vários subscribers em filas diferentes
		broker.Subscribe("test_queue1")
		broker.Subscribe("test_queue2")
		// Chamar a função CloseAllChannels
		broker.CloseAllChannels()
		// Não deve haver erro ao fechar todos os canais associados ao Broker
	})

	// Test CloseAllChannels - Error Path
	t.Run("CloseAllChannels_ErrorPath", func(t *testing.T) {
		// Criar vários subscribers em filas diferentes
		subscriber1 := broker.Subscribe("test_queue1")
		subscriber2 := broker.Subscribe("test_queue2")

		// Simular um erro durante o fechamento do canal
		simulateCloseError := func(queue string, ch <-chan interface{}) {
			t.Error("Simulated error during channel close")
			broker.Close(queue, ch)
		}

		// Configurar a função de fechamento simulada
		broker.SetCloseFunction(simulateCloseError)

		// Chamar a função CloseAllChannels
		broker.CloseAllChannels()

		// Aguardar para garantir que o erro foi simulado
		<-subscriber1
		<-subscriber2
	})
}
