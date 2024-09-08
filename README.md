# mwazovzky/assistant

Package mwazovzky/assistant implements simple open ai api client.

## Install

```
go get github.com/mwazovzky/assistant
```

## Basic Example

```
apiKey := os.Getenv("OPENAI_API_KEY")
assistantRole := "You are a helpful assistant."
a := openapi.NewAssistant(apiKey, assistantRole)
msg, err := a.Ask("2+2=")
if err != nil {
	log.Fatal(err)
}
log.Println(msg) // 2024/09/08 19:11:53 2 + 2 = 4.
```
