package assistant

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHttpClient using testify's mock
type MockHttpClient struct {
	mock.Mock
}

func (c *MockHttpClient) Request(model string, msgs []Message) (Message, Usage, error) {
	args := c.Called(model, msgs)
	return args.Get(0).(Message), args.Get(1).(Usage), args.Error(2)
}

// MockThreadRepo using testify's mock
type MockThreadRepo struct {
	mock.Mock
}

func (r *MockThreadRepo) ThreadExists(tid string) (bool, error) {
	args := r.Called(tid)
	return args.Bool(0), args.Error(1)
}

func (r *MockThreadRepo) CreateThread(tid string) error {
	return r.Called(tid).Error(0)
}

func (r *MockThreadRepo) AppendMessage(tid string, msg Message) error {
	return r.Called(tid, msg).Error(0)
}

func (r *MockThreadRepo) GetMessages(tid string) ([]Message, error) {
	args := r.Called(tid)
	return args.Get(0).([]Message), args.Error(1)
}

func TestNewAssistant(t *testing.T) {
	client := &MockHttpClient{}
	threads := &MockThreadRepo{}

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)

	assert.NotNil(t, assistant, "Expected a non-nil Assistant instance")
}

func TestAsk_Success(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"
	expectedResponse := Message{Role: RoleAssistant, Content: "Mock response"}
	expectedUsage := Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(true, nil)
	threads.On("GetMessages", tid).Return([]Message{}, nil)
	client.On("Request", "gpt-4", mock.Anything).Return(expectedResponse, expectedUsage, nil)
	threads.On("AppendMessage", tid, mock.Anything).Return(nil)

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	response, usage, err := assistant.Ask(tid, question)

	assert.NoError(t, err)
	assert.Equal(t, "Mock response", response)
	assert.Equal(t, expectedUsage, usage)
	threads.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestAsk_Success_CreateThread(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"
	expectedResponse := Message{Role: RoleAssistant, Content: "Mock response"}
	expectedUsage := Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(false, nil)
	threads.On("CreateThread", tid).Return(nil)
	threads.On("AppendMessage", tid, Message{Role: RoleSystem, Content: "You are a helpful assistant."}).Return(nil)
	threads.On("GetMessages", tid).Return([]Message{}, nil)
	client.On("Request", "gpt-4", mock.Anything).Return(expectedResponse, expectedUsage, nil)
	threads.On("AppendMessage", tid, mock.Anything).Return(nil)

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	response, usage, err := assistant.Ask(tid, question)

	assert.NoError(t, err)
	assert.Equal(t, "Mock response", response)
	assert.Equal(t, expectedUsage, usage)
	threads.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestAsk_Error_GetThread(t *testing.T) {
	tid := "error-thread"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(false, errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, _, err := assistant.Ask(tid, "What is 2+2?")

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
	threads.AssertExpectations(t)
}

func TestAsk_Error_CreateThread(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)

	threads.On("ThreadExists", tid).Return(false, nil)
	threads.On("CreateThread", tid).Return(errors.New("mock error"))

	_, _, err := assistant.Ask(tid, question)

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
	threads.AssertExpectations(t)
}

func TestAsk_Error_AppendSystemMessage(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(false, nil)
	threads.On("CreateThread", tid).Return(nil)
	threads.On("AppendMessage", tid, mock.Anything).Return(errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, _, err := assistant.Ask(tid, question)

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")

	threads.AssertExpectations(t)
}

func TestAsk_Error_GetMessages(t *testing.T) {
	tid := "thread-1"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(true, nil)
	threads.On("GetMessages", tid).Return([]Message{}, errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, _, err := assistant.Ask(tid, "What is 2+2?")

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
	threads.AssertExpectations(t)
}

func TestAsk_Error_Request(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(true, nil)
	threads.On("GetMessages", tid).Return([]Message{}, nil)
	client.On("Request", "gpt-4", mock.Anything).Return(Message{}, Usage{}, errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, _, err := assistant.Ask(tid, question)

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")

	threads.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestAsk_Error_AppendMessage(t *testing.T) {
	tid := "thread-1"
	question := "What is 2+2?"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("ThreadExists", tid).Return(true, nil)
	threads.On("GetMessages", tid).Return([]Message{{Role: RoleUser, Content: question}}, nil)
	client.On("Request", "gpt-4", mock.Anything).Return(Message{}, Usage{}, nil)
	threads.On("AppendMessage", tid, mock.Anything).Return(errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, _, err := assistant.Ask(tid, "What is 2+2?")

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
	threads.AssertExpectations(t)
}

func TestGetMessages_Success(t *testing.T) {
	tid := "thread-1"
	expectedMessages := []Message{{Role: RoleUser, Content: "Hello!"}}

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("GetMessages", tid).Return(expectedMessages, nil)

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	messages, err := assistant.GetMessages(tid)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	threads.AssertExpectations(t)
}

func TestGetMessages_Error(t *testing.T) {
	tid := "error-thread"

	client := &MockHttpClient{}
	threads := &MockThreadRepo{}
	threads.On("GetMessages", tid).Return([]Message{}, errors.New("mock error"))

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	_, err := assistant.GetMessages(tid)

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
	threads.AssertExpectations(t)
}

func TestGetUsage(t *testing.T) {
	client := &MockHttpClient{}
	threads := &MockThreadRepo{}

	assistant := NewAssistant("gpt-4", "You are a helpful assistant.", client, threads)
	expectedUsage := Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}
	assistant.usage = expectedUsage

	usage := assistant.GetUsage()

	assert.Equal(t, expectedUsage, usage, "GetUsage should return the correct usage statistics")
}
