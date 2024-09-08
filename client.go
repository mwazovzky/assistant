package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type choice struct {
	Index   int     `json:"index"`
	Message message `json:"message"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openAiResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []choice `json:"choices"`
	Usage             usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type openAiClient struct {
	url    string
	apiKey string
}

const url = "https://api.openai.com/v1/chat/completions"

func newOpenAiClient(apiKey string) *openAiClient {
	return &openAiClient{url, apiKey}
}

func (c *openAiClient) Request(model string, system string, msg string) (string, error) {
	request := openAiRequest{
		Model: model,
		Messages: []message{
			{
				Role:    "system",
				Content: system,
			},
			{
				Role:    "user",
				Content: msg,
			},
		},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(reqBody)

	httpReq, err := http.NewRequest(http.MethodPost, c.url, reader)
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	httpClient := &http.Client{}

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}

	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(httpRes.Body)
		log.Println(string(resBody))
		return "", fmt.Errorf("http request error, status %d", httpRes.StatusCode)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", err
	}

	var res openAiResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}
