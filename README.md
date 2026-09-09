# AI Wrapper for Go

[English](README.md) | [日本語](README_ja.md)

# AI Wrapper for Go

A unified Go library for interacting with multiple LLM providers (OpenAI, Anthropic, and Gemini) through a consistent interface.

## Features

- Single, consistent API for multiple LLM providers
- Support for the latest models from OpenAI, Anthropic, and Gemini
- Simple conversation handling with message history
- Token usage tracking
- Provider-agnostic reasoning depth control (`ThinkingLevel`) mapped to each provider's native thinking / reasoning parameters
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

## Configuration

The wrapper supports configuration through the `models.Config` struct:

```go
config := models.Config{
    MaxToken: 2048,  // Maximum tokens for response generation
}
```

## Quick Start

### Using a Single Provider

```go
package main

import (
    "fmt"
    "os"
    
    wrapper "github.com/obutora/ai-wrapper"
    "github.com/obutora/ai-wrapper/models"
)

func main() {
    // Create a config with max tokens
    config := models.Config{
        MaxToken: 1000,  // Limit response to 1000 tokens
    }
    
    // Create a client for OpenAI with config
    client, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"), config)
    if err != nil {
        panic(err)
    }
    
    // Generate text with a single message
    text, err, tokens := client.GenText(wrapper.GenTextParams{
        Model: models.ModelGPT4o,
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
    "github.com/obutora/ai-wrapper/models"
)

func main() {
    // Create a map of API keys for different providers
    apiKeys := map[wrapper.Provider]string{
        wrapper.ProviderOpenAI:    os.Getenv("OPENAI_API_KEY"),
        wrapper.ProviderAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
        wrapper.ProviderGemini:    os.Getenv("GEMINI_API_KEY"),
    }
    
    // Create a config
    config := models.Config{
        MaxToken: 2048,  // Set max tokens for all providers
    }
    
    // Create a unified client with config
    client, err := wrapper.NewUnifiedClient(apiKeys, config)
    if err != nil {
        panic(err)
    }
    
    // Use an OpenAI model (automatically selects the OpenAI provider)
    openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
        Model: models.ModelGPT4o,
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
        Model: models.ModelClaudeOpus5,
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
        Model: models.ModelGeminiProLatest,
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

- `ModelClaudeOpus5` - Claude Opus 5 (latest Opus)
- `ModelClaudeSonnet5` - Claude Sonnet 5 (latest Sonnet)
- `ModelClaudeHaiku45` - Claude Haiku 4.5 (latest Haiku)
- `ModelClaudeFable51` / `ModelClaudeFable5` - Claude Fable 5.1 / 5
- `ModelClaudeOpus48` / `ModelClaudeOpus47` / `ModelClaudeOpus46` / `ModelClaudeSonnet46` - Claude 4.6-4.8 family
- `ModelClaudeOpus45` / `ModelClaudeSonnet45` - Claude 4.5 family
- `ModelClaude3Opus`, `ModelClaude37Sonnet`, `ModelClaude3Haiku`, `ModelClaudeOpus41` - **deprecated** (retired models, kept for compatibility)

### Gemini

Prefer the `*-latest` aliases unless you need to pin a specific version.

- `ModelGeminiFlashLatest` - `gemini-flash-latest` (always the latest Flash; currently Gemini 3.8 Flash)
- `ModelGeminiFlashLiteLatest` - `gemini-flash-lite-latest` (currently Gemini 3.5 Flash-Lite)
- `ModelGeminiProLatest` - `gemini-pro-latest` (currently Gemini 3.1 Pro Preview)
- `ModelGemini38Flash` / `ModelGemini35FlashLite` / `ModelGemini31ProPreview` - pinned versions of the above
- `ModelGemini25Flash` / `ModelGemini25FlashLite` - Gemini 2.5 Flash / Flash-Lite
- `ModelGemini25Pro`, `ModelGemini20Flash`, `ModelGemini20Pro`, `ModelGemini25FlashPreview`, `ModelGemini25ProPreview` - **deprecated** (retired models)

## Detailed Usage

### Creating a Client

```go
// Create config
config := models.Config{
    MaxToken: 1500,  // Maximum tokens for response
}

// OpenAI client
openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"), config)

// Anthropic client
anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"), config)

// Gemini client
geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"), config)
```

### Generating Text

```go
// Basic text generation
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
    },
})

// With conversation history
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
        {Role: wrapper.RoleAssistant, Content: "The capital of France is Paris."},
        {Role: wrapper.RoleUser, Content: "What is its population?"},
    },
})

// With system message (for supported providers)
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleSystem, Content: "You are a helpful assistant that provides concise answers."},
        {Role: wrapper.RoleUser, Content: "What is the capital of France?"},
    },
})
```

### Thinking Level (Reasoning Depth)

`GenTextParams.ThinkingLevel` controls how much reasoning ("thinking") the model performs before answering.
The same value works across all providers: the wrapper converts it to each provider's native parameter,
and silently ignores it on models that do not support reasoning.

```go
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model:         models.ModelGeminiFlashLatest, // currently resolves to gemini-3.8-flash
    ThinkingLevel: wrapper.ThinkingLevelHigh,      // minimal / low / medium / high / max
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "Solve this step by step: ..."},
    },
})
```

| `ThinkingLevel` | Gemini 3+ / `*-latest` (`thinkingLevel`) | Gemini 2.5 (`thinkingBudget`) | OpenAI reasoning models (`reasoning_effort`) | Claude 4.6+ (`thinking: adaptive` + `output_config.effort`) | Claude 3.7-4.5 (`thinking: enabled` + `budget_tokens`) |
|---|---|---|---|---|---|
| `minimal` | `LOW` (Flash-Lite: `MINIMAL`) | 0 (Pro: 128) | `minimal` (`none` on GPT-5.1+, `low` on o-series) | `low` | 1024 |
| `low` | `LOW` | 1024 | `low` | `low` | 2048 |
| `medium` | `MEDIUM` | 8192 | `medium` | `medium` | 8192 |
| `high` | `HIGH` | 24576 | `high` | `high` | 16384 |
| `max` | `HIGH` | 24576 (Pro: 32768) | `high` (`xhigh` on GPT-5.2+) | `max` | 32768 |

Notes:

- Models without reasoning support (Gemini 2.0 and earlier, GPT-4 family, Claude 3.x, ...) ignore `ThinkingLevel`.
- Version-less aliases (`gemini-flash-latest`, `gemini-pro-latest`, `gemini-flash-lite-latest`) currently resolve to Gemini 3.x and therefore use `thinkingLevel`.
- `MINIMAL` support differs per Gemini 3 model (as of 2026-09: 3.5 Flash-Lite accepts it, 3.8 Flash and 3.1 Pro reject it with 400). The wrapper sends `MINIMAL` only to Flash-Lite models, and if a model rejects it the request is retried once with `LOW`.
- On Claude 3.7-4.5, `budget_tokens` must be smaller than `MaxToken`: the budget is clamped to `MaxToken - 1`, and `ErrInvalidThinkingConfig` is returned when `MaxToken` is 1024 or less.
- An unknown value returns `ErrInvalidThinkingLevel` before any API call is made.
- Reasoning tokens are included in the returned token count.

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
    "github.com/obutora/ai-wrapper/models"
)

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }
}

func main() {
    // Create config with max tokens
    config := models.Config{
        MaxToken: 500,  // Limit responses to 500 tokens
    }
    
    // Create OpenAI client with config
    openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"), config)
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

// ThinkingLevel controls the reasoning depth across providers
type ThinkingLevel string

const (
    ThinkingLevelDefault ThinkingLevel = ""        // provider/model default
    ThinkingLevelMinimal ThinkingLevel = "minimal"
    ThinkingLevelLow     ThinkingLevel = "low"
    ThinkingLevelMedium  ThinkingLevel = "medium"
    ThinkingLevelHigh    ThinkingLevel = "high"
    ThinkingLevelMax     ThinkingLevel = "max"
)

// GenTextParams represents parameters for text generation
type GenTextParams struct {
    Model         Model         `json:"model"`
    Prompt        string        `json:"prompt,omitempty"`
    CacheEnabled  bool          `json:"cache_enabled"`
    Messages      []Message     `json:"messages"`
    ThinkingLevel ThinkingLevel `json:"thinking_level,omitempty"`
}

// Config represents configuration options for the wrapper
type Config struct {
    MaxToken int  // Maximum tokens for response generation
}

// LLMWrapper is an interface for interacting with LLM providers
type LLMWrapper interface {
    GenText(params GenTextParams) (string, error, int)
}
```

### Functions

```go
// NewClient creates a new LLM client for the specified provider with configuration
func NewClient(provider Provider, apiKey string, config Config) (LLMWrapper, error)

// NewUnifiedClient creates a unified client that can use multiple providers
func NewUnifiedClient(apiKeys map[Provider]string, config Config) (*UnifiedClient, error)
```

### Error Constants

```go
var (
    ErrUnsupportedProvider = errors.New("unsupported provider")
    ErrInvalidAPIKey       = errors.New("invalid API key")
    ErrInvalidModel        = errors.New("invalid model")
    ErrEmptyMessages       = errors.New("empty messages")
    ErrAPIRequest          = errors.New("API request error")
    ErrInvalidThinkingLevel  = errors.New("invalid thinking level")
    ErrInvalidThinkingConfig = errors.New("invalid thinking config")
)
```

## License

[MIT License](LICENSE)
