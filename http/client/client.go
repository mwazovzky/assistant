package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mwazovzky/assistant"
)

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
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []choice `json:"choices"`
	Usage             usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type OpenAiClient struct {
	url    string
	apiKey string
}

func NewOpenAiClient(url string, apiKey string) *OpenAiClient {
	return &OpenAiClient{url, apiKey}
}

func (c *OpenAiClient) Request(model string, messages []assistant.Message) (msg assistant.Message, err error) {
	request := openAiRequest{
		Model:    model,
		Messages: messages,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return msg, err
	}

	reader := bytes.NewReader(reqBody)

	httpReq, err := http.NewRequest(http.MethodPost, c.url, reader)
	if err != nil {
		return msg, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	httpClient := &http.Client{}

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return msg, err
	}

	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(httpRes.Body)
		log.Println(string(resBody))
		return msg, fmt.Errorf("http request error, status %d", httpRes.StatusCode)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return msg, err
	}

	// log.Println(string(resBody))

	var res openAiResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return msg, err
	}

	return res.Choices[0].Message, nil
}
