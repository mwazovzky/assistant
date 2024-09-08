package assistant

type Assistant struct {
	client *openAiClient
	model  string
	system string
}

func NewAssistant(apiKey string, system string) *Assistant {
	return &Assistant{
		client: newOpenAiClient(apiKey),
		model:  "gpt-4o-mini",
		system: system,
	}
}

func (a *Assistant) Ask(msg string) (string, error) {
	return a.client.Request(a.model, a.system, msg)
}
