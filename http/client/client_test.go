package client_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mwazovzky/assistant"
	"github.com/mwazovzky/assistant/http/client"
)

const success = "\u2713"
const failure = "\u2717"

// method params
var model = "gpt-4o-mini"
var system = "Assistant"
var question = "2+2="

// request params
var expectedMethod = http.MethodPost
var expectedContentType = "application/json"
var expectedAuthorization = "Bearer secret"
var expectedBody = `{"model":"gpt-4o-mini","messages":[{"role":"system","content":"Assistant"},{"role":"user","content":"2+2="}]}`

// response
var resBody = `{"id": "chatcmpl-A6KkJ8rHFQGCaTO54yZYYTie809wr","object": "chat.completion","created": 1726073079,"model": "gpt-4o-mini-2024-07-18",
"choices":[{"index": 0,"message": {"role": "assistant","content": "2+2=4","refusal": null},"logprobs": null,"finish_reason": "stop"}],
"usage": {"prompt_tokens": 21,"completion_tokens": 8,"total_tokens": 29},"system_fingerprint": "fp_483d39d857"}`

func TestRequest(t *testing.T) {
	t.Logf("When OpenAiClient calls Request(\"%s\", `{{\"system\", \"%s\"},{\"user\", \"%s\"}}`)", model, system, question)

	f := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == expectedMethod {
			t.Logf("\t%s Request method should be [%s].", success, method)
		} else {
			t.Errorf("\t%s Request method should be [%s], got [%s]", failure, expectedMethod, method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType == expectedContentType {
			t.Logf("\t%s Request header 'Content-Type' should be [%s].", success, expectedContentType)
		} else {
			t.Errorf("\t%s Request header 'Content-Type' should be [%s], got [%s]", failure, expectedContentType, contentType)
		}

		authorization := r.Header.Get("Authorization")
		if authorization == expectedAuthorization {
			t.Logf("\t%s Request header 'Authorization' should be [%s].", success, expectedAuthorization)
		} else {
			t.Errorf("\t%s Request header 'Authorization' should be [%s], got [%s]", failure, expectedAuthorization, authorization)
		}

		content, _ := io.ReadAll(r.Body)
		body := string(content)
		if body == expectedBody {
			t.Logf("\t%s Request body should be [%s].", success, expectedBody)
		} else {
			t.Errorf("\t%s Request body should be [%s], got [%s]", failure, expectedBody, body)
		}

		w.WriteHeader(200)
		w.Header().Set("Conten-Type", "application/json")
		fmt.Fprintln(w, resBody)
	}

	server := httptest.NewServer(http.HandlerFunc(f))
	url := server.URL
	defer server.Close()

	client := client.NewOpenAiClient(url, "secret")

	messages := []assistant.Message{
		{Role: "system", Content: system},
		{Role: "user", Content: question},
	}
	msg, err := client.Request(model, messages)
	if err == nil {
		t.Logf("\t%s Response should not be an error.", success)
	} else {
		t.Fatalf("\t%s Response should not be an error, got [%s]", failure, err)
	}

	expectedMsg := assistant.Message{Role: "assistant", Content: "2+2=4"}
	if msg.Role == expectedMsg.Role {
		t.Logf("\t%s Response message role should be [%s].", success, expectedMsg.Role)
	} else {
		t.Errorf("\t%s Response message role should be [%s], got [%s]", failure, expectedMsg.Role, msg.Role)
	}
	if msg.Content == expectedMsg.Content {
		t.Logf("\t%s Response message content should be [%s].", success, expectedMsg.Content)
	} else {
		t.Errorf("\t%s Response message content should be [%s], got [%s]", failure, expectedMsg.Content, msg.Content)
	}
}
