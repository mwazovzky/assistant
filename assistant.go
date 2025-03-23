package assistant

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type HttpClient interface {
	Request(model string, msgs []Message) (msg Message, usage Usage, err error)
}

type ThreadRepository interface {
	ThreadExists(tid string) (bool, error)
	CreateThread(tid string) error
	AppendMessage(tid string, msg Message) error
	GetMessages(tid string) ([]Message, error)
}

type Assistant struct {
	model   string
	system  string
	client  HttpClient
	threads ThreadRepository
	usage   Usage
}

func NewAssistant(model string, system string, client HttpClient, threads ThreadRepository) *Assistant {
	return &Assistant{
		model:   model,
		system:  system,
		client:  client,
		threads: threads,
		usage:   Usage{},
	}
}

func (a *Assistant) Ask(tid string, msg string) (string, error) {
	if err := a.getThread(tid); err != nil {
		return "", err
	}

	messages, err := a.threads.GetMessages(tid)
	if err != nil {
		return "", err
	}

	response, usage, err := a.client.Request(a.model, messages)
	if err != nil {
		return "", err
	}

	a.usage = usage

	if err := a.threads.AppendMessage(tid, response); err != nil {
		return "", err
	}

	return response.Content, nil
}

func (a *Assistant) GetMessages(tid string) ([]Message, error) {
	return a.threads.GetMessages(tid)
}

func (a *Assistant) GetUsage() Usage {
	return a.usage
}

func (a *Assistant) getThread(tid string) error {
	exists, err := a.threads.ThreadExists(tid)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return a.createThread(tid)
}

func (a *Assistant) createThread(tid string) error {
	if err := a.threads.CreateThread(tid); err != nil {
		return err
	}

	return a.threads.AppendMessage(tid, Message{Role: RoleSystem, Content: a.system})
}
