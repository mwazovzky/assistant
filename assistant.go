package assistant

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type HttpClient interface {
	Request(model string, msgs []Message) (msg Message, err error)
}

type TreadRepository interface {
	CreateThread(tid string) error
	AppendMessage(tid string, msg Message) error
	GetMessages(tid string) ([]Message, error)
}

type Assistant struct {
	model   string
	system  string
	client  HttpClient
	threads TreadRepository
}

func NewAssistant(system string, client HttpClient, threads TreadRepository) *Assistant {
	return &Assistant{
		model:   "gpt-4o-mini",
		system:  system,
		client:  client,
		threads: threads,
	}
}

// Ask creates am OpenAi request without context
func (a *Assistant) Ask(msg string) (string, error) {
	messages := []Message{
		{"system", a.system},
		{"user", msg},
	}

	res, err := a.client.Request(a.model, messages)
	if err != nil {
		return "", err
	}

	return res.Content, nil
}

func (a *Assistant) GetThread(tid string) ([]Message, error) {
	return a.threads.GetMessages(tid)
}

// CreateThread creates a thread of messages used as a conversation context
func (a *Assistant) CreateThread(tid string) error {
	err := a.threads.CreateThread(tid)
	if err != nil {
		return err
	}

	err = a.threads.AppendMessage(tid, Message{"system", a.system})
	if err != nil {
		return err
	}

	return nil
}

// Post creates am OpenAi request with a thread of messages as a conversation context
func (a *Assistant) Post(tid string, txt string) (string, error) {
	err := a.threads.AppendMessage(tid, Message{"user", txt})
	if err != nil {
		return "", err
	}

	messages, err := a.threads.GetMessages(tid)
	if err != nil {
		return "", err
	}

	msg, err := a.client.Request(a.model, messages)
	if err != nil {
		return "", err
	}

	err = a.threads.AppendMessage(tid, msg)
	if err != nil {
		return "", err
	}

	return msg.Content, nil
}
