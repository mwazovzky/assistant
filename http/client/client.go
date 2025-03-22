package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mwazovzky/assistant"
)

// HttpDoer interface to abstract the Do method
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type openAiRequest struct {
	Model    string              `json:"model"`
	Messages []assistant.Message `json:"messages"`
}

type choice struct {
	Index   int               `json:"index"`
	Message assistant.Message `json:"message"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openAiResponse struct {
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

type OpenAiClient struct {
	url        string
	apiKey     string
	httpClient HttpDoer
}

func NewOpenAiClient(url string, apiKey string) *OpenAiClient {
	return &OpenAiClient{
		url:        url,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *OpenAiClient) SetHttpClient(httpClient HttpDoer) {
	c.httpClient = httpClient
}

func (c *OpenAiClient) Request(model string, messages []assistant.Message) (assistant.Message, assistant.Usage, error) {
	reqBody, err := json.Marshal(openAiRequest{Model: model, Messages: messages})
	if err != nil {
		return assistant.Message{}, assistant.Usage{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := c.createRequest(context.Background(), reqBody)
	if err != nil {
		return assistant.Message{}, assistant.Usage{}, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return assistant.Message{}, assistant.Usage{}, fmt.Errorf("http request failed: %w", err)
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return assistant.Message{}, assistant.Usage{}, fmt.Errorf("http request error, status %d", httpRes.StatusCode)
	}

	var res openAiResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		return assistant.Message{}, assistant.Usage{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(res.Choices) == 0 {
		return assistant.Message{}, assistant.Usage{}, fmt.Errorf("no choices returned in the response")
	}

	usage := assistant.Usage{
		PromptTokens:     res.Usage.PromptTokens,
		CompletionTokens: res.Usage.CompletionTokens,
		TotalTokens:      res.Usage.TotalTokens,
	}

	return res.Choices[0].Message, usage, nil
}

func (c *OpenAiClient) createRequest(ctx context.Context, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return req, nil
}
