package client_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mwazovzky/assistant"
	"github.com/mwazovzky/assistant/http/client"
)

type MockHttpDoer struct {
	mock.Mock
}

func (m *MockHttpDoer) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestRequest(t *testing.T) {
	type testCase struct {
		name           string
		mockResponse   *http.Response
		mockError      error
		expectedError  string
		expectedResult assistant.Message
		expectedUsage  assistant.Usage
	}

	tests := []testCase{
		{
			name: "Success",
			mockResponse: func() *http.Response {
				rec := httptest.NewRecorder()
				rec.WriteHeader(http.StatusOK)
				rec.Body.WriteString(`{
					"choices": [{"message": {"role": "assistant", "content": "2+2=4"}}],
					"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
				}`)
				return rec.Result()
			}(),
			expectedResult: assistant.Message{Role: "assistant", Content: "2+2=4"},
			expectedUsage:  assistant.Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		},
		{
			name:          "HTTP Error",
			mockResponse:  nil,
			mockError:     errors.New("mock network error"),
			expectedError: "mock network error",
		},
		{
			name: "Invalid JSON Response",
			mockResponse: func() *http.Response {
				rec := httptest.NewRecorder()
				rec.WriteHeader(http.StatusOK)
				rec.Body.WriteString("invalid-json")
				return rec.Result()
			}(),
			expectedError: "failed to decode response",
		},
		{
			name: "Empty Choices",
			mockResponse: func() *http.Response {
				rec := httptest.NewRecorder()
				rec.WriteHeader(http.StatusOK)
				rec.Body.WriteString(`{
					"choices": [],
					"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
				}`)
				return rec.Result()
			}(),
			expectedError: "no choices returned in the response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHttpDoer := &MockHttpDoer{}
			openAiClient := client.NewOpenAiClient("http://example.com", "test-api-key")
			openAiClient.SetHttpClient(mockHttpDoer)

			if tt.mockResponse != nil || tt.mockError != nil {
				mockHttpDoer.On("Do", mock.Anything).Return(tt.mockResponse, tt.mockError)
			}

			result, usage, err := openAiClient.Request("gpt-4", []assistant.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is 2+2?"},
			})

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
				assert.Equal(t, tt.expectedUsage, usage)
			}

			mockHttpDoer.AssertExpectations(t)
		})
	}
}
