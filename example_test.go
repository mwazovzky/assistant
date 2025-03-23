package assistant_test

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mwazovzky/assistant"
)

type MockHttpClient struct{}

func (c MockHttpClient) Request(model string, msgs []assistant.Message) (msg assistant.Message, usage assistant.Usage, err error) {
	return assistant.Message{Role: assistant.RoleAssistant, Content: "Mock response"}, assistant.Usage{
		PromptTokens:     10,
		CompletionTokens: 5,
		TotalTokens:      15,
	}, nil
}

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

func (tr *ThreadRepository) ThreadExists(tid string) (bool, error) {
	_, ok := tr.data[tid]
	return ok, nil
}

func ExampleAssistant_Ask() {
	model := "gpt-4o-mini"
	system := "You are assistant"
	client := MockHttpClient{} // Use the mock client
	threads := NewThreadRepository()

	a := assistant.NewAssistant(model, system, client, threads)

	tid := uuid.New().String()

	msg, err := a.Ask(tid, "2+2=")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg)
	// Output: Mock response
}

func ExampleAssistant_GetMessages() {
	model := "gpt-4o-mini"
	system := "You are assistant"
	client := MockHttpClient{} // Use the mock client
	threads := NewThreadRepository()

	a := assistant.NewAssistant(model, system, client, threads)

	tid := uuid.New().String()

	// Ensure the thread contains at least one message
	_ = threads.CreateThread(tid)
	_ = threads.AppendMessage(tid, assistant.Message{Role: assistant.RoleSystem, Content: "You are assistant"})

	messages, err := a.GetMessages(tid)
	if err != nil {
		fmt.Println("expected no error, got", err)
		return
	}

	len := len(messages)
	if len != 1 {
		fmt.Println("expected one message in thread, got", len)
		return
	}

	fmt.Println(messages[0])
	// Output: {system You are assistant}
}
