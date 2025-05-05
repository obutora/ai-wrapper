# AI Wrapper for Go

[English](README.md) | [日本語](README_ja.md)

# AI Wrapper for Go

A unified Go library for interacting with multiple LLM providers (OpenAI, Anthropic, and Gemini) through a consistent interface.

## Features

- Single, consistent API for multiple LLM providers
- Support for the latest models from OpenAI, Anthropic, and Gemini
- Simple conversation handling with message history
- Token usage tracking
- Error handling with provider-specific details

## Installation

```bash
go get github.com/obutora/ai-wrapper
```

## Environment Variables

Set the following environment variables for each provider you want to use:

```bash
# OpenAI
export OPENAI_API_KEY=your_openai_api_key

# Anthropic
export ANTHROPIC_API_KEY=your_anthropic_api_key

# Gemini
export GEMINI_API_KEY=your_gemini_api_key
```

You can also use a `.env` file with a package like [godotenv](https://github.com/joho/godotenv) to load these variables.

## Quick Start

### Using a Single Provider

```go
package main

import (
    "fmt"
    "os"
    
    wrapper "github.com/obutora/ai-wrapper"
)

func main() {
    // Create a client for OpenAI
    client, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
    if err != nil {
        panic(err)
    }
    
    // Generate text with a single message
    text, err, tokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGPT4o,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Response: %s\nTokens used: %d\n", text, tokens)
}
```

### Using the Unified Client (Multiple Providers)

```go
package main

import (
    "fmt"
    "os"
    
    wrapper "github.com/obutora/ai-wrapper"
)

func main() {
    // Create a map of API keys for different providers
    apiKeys := map[wrapper.Provider]string{
        wrapper.ProviderOpenAI:    os.Getenv("OPENAI_API_KEY"),
        wrapper.ProviderAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
        wrapper.ProviderGemini:    os.Getenv("GEMINI_API_KEY"),
    }
    
    // Create a unified client that can use any provider
    client, err := wrapper.NewUnifiedClient(apiKeys)
    if err != nil {
        panic(err)
    }
    
    // Use an OpenAI model (automatically selects the OpenAI provider)
    openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGPT4o,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("OpenAI Response: %s\nTokens used: %d\n\n", openaiText, openaiTokens)
    
    // Use an Anthropic model (automatically selects the Anthropic provider)
    anthropicText, err, anthropicTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelClaude3Opus,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "What is the capital of Germany?"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Anthropic Response: %s\nTokens used: %d\n\n", anthropicText, anthropicTokens)
    
    // Use a Gemini model (automatically selects the Gemini provider)
    geminiText, err, geminiTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGemini20Pro,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "What is the capital of Japan?"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Gemini Response: %s\nTokens used: %d\n", geminiText, geminiTokens)
    
    // Register a custom model mapping if needed
    client.RegisterCustomModel("my-custom-model", wrapper.ProviderOpenAI)
}
```

## Supported Providers and Models

### OpenAI

- `ModelGPT4o` - GPT-4o
- `ModelGPT4` - GPT-4
- `ModelGPT35Turbo` - GPT-3.5 Turbo
- `ModelO3Mini` - O3 Mini
- `ModelO4Mini` - O4 Mini
- `Model4_1Nano` - GPT-4.1 Nano
- `ModelO3` - O3

### Anthropic

- `ModelClaude3Opus` - Claude 3 Opus
- `ModelClaude37Sonnet` - Claude 3.7 Sonnet
- `ModelClaude3Haiku` - Claude 3 Haiku

### Gemini

- `ModelGemini20Flash` - Gemini 2.0 Flash
- `ModelGemini20Pro` - Gemini 2.0 Pro
- `ModelGemini25FlashPreview` - Gemini 2.5 Flash Preview
- `ModelGemini25ProPreview` - Gemini 2.5 Pro Preview

## Detailed Usage

### Creating a Client

```go
// OpenAI client
openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))

// Anthropic client
anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"))

// Gemini client
geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"))
```

### Generating Text

```go
// Basic text generation
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
    },
})

// With conversation history
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
        {Role: wrapper.RoleAssistant, Content: "The capital of France is Paris."},
        {Role: wrapper.RoleUser, Content: "What is its population?"},
    },
})

// With system message (for supported providers)
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleSystem, Content: "You are a helpful assistant that provides concise answers."},
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
    },
})
```

### Error Handling

```go
text, err, tokens := client.GenText(params)
if err != nil {
    switch {
    case errors.Is(err, wrapper.ErrInvalidAPIKey):
        // Handle invalid API key
    case errors.Is(err, wrapper.ErrInvalidModel):
        // Handle invalid model
    case errors.Is(err, wrapper.ErrEmptyMessages):
        // Handle empty messages
    case errors.Is(err, wrapper.ErrAPIRequest):
        // Handle API request error
    case errors.Is(err, wrapper.ErrUnsupportedProvider):
        // Handle unsupported provider
    default:
        // Handle other errors
    }
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    wrapper "github.com/obutora/ai-wrapper"
)

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }
}

func main() {
    // Create OpenAI client
    openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
    if err != nil {
        panic(err)
    }

    // First text generation with specific information
    openaiText1, err, openaiTokens1 := openaiClient.GenText(wrapper.GenTextParams{
        Model: wrapper.Model4_1Nano,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "Tanaka Taro is a 42-year-old engineer living in Tokyo. His hobbies are mountain climbing and photography. Last month, he climbed Mount Fuji."},
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("OpenAI Response 1: %s\nTokens used: %d\n\n", openaiText1, openaiTokens1)

    // Second text generation with conversation history
    openaiText2, err, openaiTokens2 := openaiClient.GenText(wrapper.GenTextParams{
        Model: wrapper.Model4_1Nano,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "Tanaka Taro is a 42-year-old engineer living in Tokyo. His hobbies are mountain climbing and photography. Last month, he climbed Mount Fuji."},
            {Role: wrapper.RoleAssistant, Content: openaiText1},
            {Role: wrapper.RoleUser, Content: "Please tell me Tanaka's age, occupation, hobbies, and what he did last month."},
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("OpenAI Response 2: %s\nTokens used: %d\n", openaiText2, openaiTokens2)
}
```

## API Reference

### Types

```go
// Provider represents an LLM provider
type Provider string

const (
    ProviderOpenAI    Provider = "openai"
    ProviderAnthropic Provider = "anthropic"
    ProviderGemini    Provider = "gemini"
)

// Role represents the role of a message
type Role string

const (
    RoleUser      Role = "user"
    RoleAssistant Role = "assistant"
    RoleSystem    Role = "system"
)

// Message represents a message in a conversation
type Message struct {
    Role    Role   `json:"role"`
    Content string `json:"content"`
}

// GenTextParams represents parameters for text generation
type GenTextParams struct {
    Model        Model     `json:"model"`
    Prompt       string    `json:"prompt,omitempty"`
    CacheEnabled bool      `json:"cache_enabled"`
    Messages     []Message `json:"messages"`
}

// LLMWrapper is an interface for interacting with LLM providers
type LLMWrapper interface {
    GenText(params GenTextParams) (string, error, int)
}
```

### Functions

```go
// NewClient creates a new LLM client for the specified provider
func NewClient(provider Provider, apiKey string) (LLMWrapper, error)
```

### Error Constants

```go
var (
    ErrUnsupportedProvider = errors.New("unsupported provider")
    ErrInvalidAPIKey       = errors.New("invalid API key")
    ErrInvalidModel        = errors.New("invalid model")
    ErrEmptyMessages       = errors.New("empty messages")
    ErrAPIRequest          = errors.New("API request error")
)
```

## License

[MIT License](LICENSE)
