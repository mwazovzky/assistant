# Product Requirements Document (PRD) for OpenAiClient Implementation

## Overview

The `OpenAiClient` is a Go client designed to interact with the OpenAI API. It abstracts the HTTP communication required to send requests and receive responses, including usage statistics. The client is designed to be reusable, testable, and extensible.

---

## Requirements

### 1. **Request Management**

- **Requirement**: Send requests to the OpenAI API with a model and a list of messages.
- **Implementation**:
  - The `Request` method:
    - Accepts the model name and a list of `Message` objects.
    - Sends an HTTP POST request to the OpenAI API.
    - Includes the API key in the `Authorization` header.
    - Sets the `Content-Type` header to `application/json`.

---

### 2. **Response Handling**

- **Requirement**: Parse the OpenAI API response and extract the assistant's message and usage statistics.
- **Implementation**:
  - The `Request` method:
    - Parses the JSON response to extract:
      - The assistant's message (`Message`).
      - Usage statistics (`Usage`), including:
        - `PromptTokens`: Number of tokens in the input prompt.
        - `CompletionTokens`: Number of tokens in the generated response.
        - `TotalTokens`: Total tokens used.

---

## Summary

The `OpenAiClient` provides a clean and modular interface for interacting with the OpenAI API. It handles request creation, response parsing, and error handling, while exposing usage statistics for better insights into API usage. The design ensures flexibility for future extensions and ease of testing.
