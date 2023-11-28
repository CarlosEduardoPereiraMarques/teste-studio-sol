package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestPubSub(t *testing.T) {
	ps := NewPubSub()

	queue := "example"
	msg := "hello"
	ps.Publish(queue, msg)

	ch := ps.Subscribe(queue)
	receivedMsg := <-ch

	if receivedMsg != msg {
		t.Errorf("Expected message '%s', but got '%s'", msg, receivedMsg)
	}

	filename := "test_output.json"
	ps.SaveToFile(filename)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	var savedMessages []Message
	err = json.Unmarshal(fileData, &savedMessages)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	if len(savedMessages) != 1 || savedMessages[0].Data != msg || savedMessages[0].Queue != queue {
		t.Errorf("Unexpected file content: %+v", savedMessages)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("Error removing test file: %v", err)
	}
}
