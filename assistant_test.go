package assistant_test

import (
	"testing"

	"github.com/mwazovzky/assistant"
)

const success = "\u2713"
const failure = "\u2717"

type MockClient struct{}

func (c MockClient) Request(model string, msgs []assistant.Message) (msg assistant.Message, err error) {
	return assistant.Message{Role: "assistant", Content: "2+2=4"}, nil
}

type MockTreadRepository struct{}

func (t MockTreadRepository) CreateThread(tid string) error {
	return nil
}

func (t MockTreadRepository) GetMessages(tid string) ([]assistant.Message, error) {
	messages := []assistant.Message{
		{"system", "You are assistant "},
		{"user", "2+2="},
	}
	return messages, nil
}

func (t MockTreadRepository) AppendMessage(tid string, msg assistant.Message) error {
	return nil
}

func TestAsk(t *testing.T) {
	model := "gpt-4o-mini"
	system := "You are assistant"
	client := MockClient{}
	threads := MockTreadRepository{}
	question := "2+2="

	a := assistant.NewAssistant(model, system, client, threads)

	t.Logf("When Assistant calls Ask(\"%s\")", question)
	{
		msg, err := a.Ask(question)
		if err == nil {
			t.Logf("\t%s It shoud not get error.", success)
		} else {
			t.Fatalf("\t%s It shoud not get error, got:%v", failure, err)
		}

		expectedMsg := "2+2=4"

		if msg == expectedMsg {
			t.Logf("\t%s It should get message with content [%s].", success, expectedMsg)
		} else {
			t.Errorf("\t%s It should get message with content [%s], got [%s]", failure, expectedMsg, msg)
		}
	}
}

func TestCreateThread(t *testing.T) {
	model := "gpt-4o-mini"
	system := "You are assistant"
	client := MockClient{}
	threads := MockTreadRepository{}
	tid := "thread-one"

	a := assistant.NewAssistant(model, system, client, threads)

	t.Logf("When Assistant calls CreateThread(\"thread-one\")")
	a.CreateThread(tid)

	messages, err := a.GetThread(tid)

	if err == nil {
		t.Logf("\t%s It should not get err.", success)
	} else {
		t.Errorf("\t%s It should not get err, got [%s].", failure, err)
	}

	expectedMsg := assistant.Message{"system", "You are assistant "}
	msg := messages[0]
	if msg.Role == expectedMsg.Role {
		t.Logf("\t%s It should have message with contentrole [%s].", success, expectedMsg.Role)
	} else {
		t.Errorf("\t%s It should have message with content role [%s], got [%s]", failure, expectedMsg.Role, msg.Role)
	}
	if msg.Content == expectedMsg.Content {
		t.Logf("\t%s It should have message with content [%s].", success, expectedMsg.Content)
	} else {
		t.Errorf("\t%s It should have message with content [%s], got [%s]", failure, expectedMsg.Content, msg.Content)
	}
}

func TestPost(t *testing.T) {
	model := "gpt-4o-mini"
	system := "You are assistant"
	client := MockClient{}
	threads := MockTreadRepository{}
	tid := "thread-one"

	a := assistant.NewAssistant(model, system, client, threads)

	a.CreateThread(tid)

	t.Logf("When Assistant calls Post()")
	msg, err := a.Post(tid, "2+2=")

	if err == nil {
		t.Logf("\t%s It should not get error.", success)
	} else {
		t.Errorf("\t%s It should not get error, got [%s].", failure, err)
	}

	expectedMsg := "2+2=4"

	if msg == expectedMsg {
		t.Logf("\t%s It should get message with content [%s].", success, expectedMsg)
	} else {
		t.Errorf("\t%s It should get message with content [%s], got [%s]", failure, expectedMsg, msg)
	}
}
