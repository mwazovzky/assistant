// Package assistant provides a conversational AI interface for applications.
//
// The assistant package allows for the creation and management of AI-powered
// conversations through an API service like OpenAI. It provides abstractions
// for managing conversations as threads of messages, tracking token usage,
// and interacting with AI models.
//
// # Core Components
//
// - Assistant: The main struct that manages conversations with AI models
// - Message: Represents a single message in a conversation
// - Usage: Tracks token consumption for billing and monitoring
// - HttpClient: Interface for making requests to AI service APIs
// - ThreadRepository: Interface for storing and retrieving conversation threads
//
// # Basic Usage
//
//	// Initialize components
//	httpClient := client.NewOpenAiClient(apiUrl, apiKey)
//	threadRepo := storage.NewInMemoryThreadRepository()
//
//	// Create a new assistant
//	assistant := assistant.NewAssistant(
//		"gpt-4",                     // Model name
//		"You are a helpful assistant", // System prompt
//		httpClient,                  // HTTP client for API requests
//		threadRepo,                  // Thread repository for conversation storage
//	)
//
//	// Start or continue a conversation
//	threadID := "user-123"
//	response, err := assistant.Ask(threadID, "What is the capital of France?")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response) // "The capital of France is Paris."
//
//	// Get conversation history
//	messages, err := assistant.GetMessages(threadID)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, msg := range messages {
//		fmt.Printf("%s: %s\n", msg.Role, msg.Content)
//	}
//
//	// Get token usage statistics
//	usage := assistant.GetUsage()
//	fmt.Printf("Tokens used: %d\n", usage.TotalTokens)
package assistant
