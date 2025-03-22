# Product Requirements Document (PRD) for Assistant Implementation

## Overview

The `Assistant` class is designed to facilitate conversational interactions with an OpenAI-like API. It manages conversation threads, sends requests to the API (via the `OpenAiClient`), and processes responses. The implementation is modular, leveraging interfaces for flexibility and testability.

---

## Requirements

### 1. **Message Management**

- **Requirement**: Define a structure to represent conversation messages.
- **Implementation**:
  - The `Message` struct includes:
    - `Role`: The role of the message sender (e.g., `RoleSystem`, `RoleUser`, `RoleAssistant`).
    - `Content`: The content of the message.

---

### 2. **Thread Management**

- **Requirement**: Manage threads of messages for conversation context.
- **Implementation**:
  - The `ThreadRepository` interface defines:
    - `CreateThread(tid string)`: Creates a new thread.
    - `AppendMessage(tid string, msg Message)`: Appends a message to a thread.
    - `GetMessages(tid string)`: Retrieves all messages in a thread.
    - `ThreadExists(tid string)`: Checks if a thread exists by its ID.

---

### 3. **API Communication**

- **Requirement**: Use the `OpenAiClient` to send requests and receive responses.
- **Implementation**:
  - The `Assistant` relies on the `OpenAiClient` for API communication.
  - For details on the `OpenAiClient`, refer to the [OpenAiClient PRD](client_prd.md).

---

### 4. **Assistant Behavior**

- **Requirement**: Provide a conversational assistant that uses the `HttpClient` and `ThreadRepository` to manage conversations.
- **Implementation**:

  - **Fields**:
    - `model`: The model used for OpenAI requests.
    - `system`: The system message for the assistant's behavior.
    - `client`: An instance of `HttpClient` for API communication.
    - `threads`: An instance of `ThreadRepository` for thread management.
  - **Methods**:

    1. **`Ask`**:

       - **Requirement**: Send a message to the assistant within a thread.
       - **Implementation**:
         - Uses the `OpenAiClient` to send requests and retrieve responses.
         - Stores usage statistics for later retrieval.

    2. **`GetThread`**:

       - **Requirement**: Retrieve all messages in a thread.
       - **Implementation**:
         - Calls `GetMessages` from the `ThreadRepository` interface.

    3. **`GetUsage`**:
       - **Requirement**: Retrieve the latest usage statistics.
       - **Implementation**:
         - Returns the `usage` field stored in the `Assistant` instance.

---

## Error Handling

- **Requirement**: Handle errors gracefully in all methods.
- **Implementation**:
  - Errors are checked and returned at every step (e.g., thread existence check, thread creation, message appending, API requests).
  - Helper methods reduce repetitive error handling.

---

## Summary

The `Assistant` class is a modular and extensible implementation for managing conversational interactions. It relies on the `OpenAiClient` for API communication and the `ThreadRepository` for thread management. New features for thread review and usage statistics enhance its usability and provide valuable insights for users.
