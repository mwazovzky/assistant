![CI](https://github.com/mwazovzky/assistant/actions/workflows/go.yml/badge.svg)

# mwazovzky/assistant

Package mwazovzky/assistant implements simple open ai api client.

## Install

```
go get github.com/mwazovzky/assistant
```

## Basic Example

```
package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mwazovzky/assistant"
	"github.com/mwazovzky/assistant/http/client"
)

type ThreadRepository struct {
	data map[string][]assistant.Message
}

func NewThreadRepository() *ThreadRepository {
	data := make(map[string][]assistant.Message)
	return &ThreadRepository{data}
}

func (tr *ThreadRepository) CreateThread(tid string) error {
	tr.data[tid] = []assistant.Message{}
	return nil
}

func (tr *ThreadRepository) AppendMessage(tid string, msg assistant.Message) error {
	messages, ok := tr.data[tid]
	if !ok {
		return fmt.Errorf("thread [%s] does not exist", tid)
	}

	tr.data[tid] = append(messages, msg)
	return nil
}

func (tr *ThreadRepository) GetMessages(tid string) ([]assistant.Message, error) {
	messages, ok := tr.data[tid]
	if !ok {
		return nil, fmt.Errorf("thread [%s] does not exist", tid)
	}

	return messages, nil
}

func main () {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := client.NewOpenAiClient(url, apiKey)
	tr := NewThreadRepository()

	role := "You are assistant."
	a := assistant.NewAssistant(role, client, tr)

	tid := uuid.New().String()
	a.CreateThread(tid)

	msg, err := a.Post(tid, "2+2=")
	fmt.Println(msg)
	// Output: 2 + 2 = 4.
}
```

## Testing

```
go test
go test client_test.go -v
go test assistant_test.go -v
go test -test.run=TestCreateThread -v
go test example_test
go test ./...
```
