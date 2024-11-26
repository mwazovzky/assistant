package assistant_test

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mwazovzky/assistant"
	"github.com/mwazovzky/assistant/http/client"
)

const url = "https://api.openai.com/v1/chat/completions"

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

func ExampleAssistant_Ask() {
	model := "gpt-4o-mini"
	system := "You are assistant"
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := client.NewOpenAiClient(url, apiKey)
	threads := NewThreadRepository()

	a := assistant.NewAssistant(model, system, client, threads)

	msg, err := a.Ask("2+2=")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg)
	// Output: 2 + 2 = 4.
}

func ExampleAssistant_CreateThread() {
	model := "gpt-4o-mini"
	system := "You are assistant"
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := client.NewOpenAiClient(url, apiKey)
	threads := NewThreadRepository()

	a := assistant.NewAssistant(model, system, client, threads)

	tid := uuid.New().String()
	a.CreateThread(tid)

	messages, err := a.GetThread(tid)
	if err != nil {
		fmt.Println("expected no error, got", err)
	}

	len := len(messages)
	if len != 1 {
		fmt.Println("expected one message in thread, got", len)
	}

	fmt.Println(messages[0])
	// Output: {system You are assistant}
}

func ExampleAssistant_Post() {
	model := "gpt-4o-mini"
	system := "You are assistant"
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := client.NewOpenAiClient(url, apiKey)
	threads := NewThreadRepository()

	a := assistant.NewAssistant(model, system, client, threads)

	tid := uuid.New().String()
	a.CreateThread(tid)

	msg, err := a.Post(tid, "2+2=")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg)
	// Output: 2 + 2 = 4.
}
