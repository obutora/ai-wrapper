# Go用AIラッパー

[English](README.md) | [日本語](README_ja.md)

# Go用AIラッパー

一貫したインターフェースを通じて複数のLLMプロバイダ（OpenAI、Anthropic、Gemini）と対話するための統一されたGoライブラリです。

## 特徴

- 複数のLLMプロバイダに対する単一の一貫したAPI
- OpenAI、Anthropic、Geminiの最新モデルをサポート
- メッセージ履歴による簡単な会話処理
- トークン使用量の追跡
- プロバイダ固有の詳細を含むエラー処理

## インストール

```bash
go get github.com/obutora/ai-wrapper
```

## 環境変数の設定

使用したいプロバイダごとに以下の環境変数を設定してください：

```bash
# OpenAI
export OPENAI_API_KEY=your_openai_api_key

# Anthropic
export ANTHROPIC_API_KEY=your_anthropic_api_key

# Gemini
export GEMINI_API_KEY=your_gemini_api_key
```

[godotenv](https://github.com/joho/godotenv)などのパッケージを使用して、`.env`ファイルからこれらの変数を読み込むこともできます。

## クイックスタート

### 単一のプロバイダを使用する場合

```go
package main

import (
    "fmt"
    "os"
    
    wrapper "github.com/obutora/ai-wrapper"
)

func main() {
    // OpenAIのクライアントを作成
    client, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
    if err != nil {
        panic(err)
    }
    
    // 単一のメッセージでテキストを生成
    text, err, tokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGPT4o,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("応答: %s\n使用トークン数: %d\n", text, tokens)
}
```

### 統合クライアントを使用する場合（複数のプロバイダ）

```go
package main

import (
    "fmt"
    "os"
    
    wrapper "github.com/obutora/ai-wrapper"
)

func main() {
    // 各プロバイダのAPIキーをマップで作成
    apiKeys := map[wrapper.Provider]string{
        wrapper.ProviderOpenAI:    os.Getenv("OPENAI_API_KEY"),
        wrapper.ProviderAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
        wrapper.ProviderGemini:    os.Getenv("GEMINI_API_KEY"),
    }
    
    // 任意のプロバイダを使用できる統合クライアントを作成
    client, err := wrapper.NewUnifiedClient(apiKeys)
    if err != nil {
        panic(err)
    }
    
    // OpenAIモデルを使用（自動的にOpenAIプロバイダが選択される）
    openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGPT4o,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("OpenAI応答: %s\n使用トークン数: %d\n\n", openaiText, openaiTokens)
    
    // Anthropicモデルを使用（自動的にAnthropicプロバイダが選択される）
    anthropicText, err, anthropicTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelClaude3Opus,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "ドイツの首都は何ですか？"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Anthropic応答: %s\n使用トークン数: %d\n\n", anthropicText, anthropicTokens)
    
    // Geminiモデルを使用（自動的にGeminiプロバイダが選択される）
    geminiText, err, geminiTokens := client.GenText(wrapper.GenTextParams{
        Model: wrapper.ModelGemini20Pro,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "日本の首都は何ですか？"},
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Gemini応答: %s\n使用トークン数: %d\n", geminiText, geminiTokens)
    
    // 必要に応じてカスタムモデルのマッピングを登録
    client.RegisterCustomModel("my-custom-model", wrapper.ProviderOpenAI)
}
```

## サポートされているプロバイダとモデル

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

## 詳細な使用方法

### クライアントの作成

```go
// OpenAIクライアント
openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))

// Anthropicクライアント
anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"))

// Geminiクライアント
geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"))
```

### テキスト生成

```go
// 基本的なテキスト生成
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
    },
})

// 会話履歴を含む
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
        {Role: wrapper.RoleAssistant, Content: "フランスの首都はパリです。"},
        {Role: wrapper.RoleUser, Content: "その人口は？"},
    },
})

// システムメッセージを含む（サポートされているプロバイダ向け）
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: wrapper.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleSystem, Content: "あなたは簡潔な回答を提供する役立つアシスタントです。"},
        {Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
    },
})
```

### エラー処理

```go
text, err, tokens := client.GenText(params)
if err != nil {
    switch {
    case errors.Is(err, wrapper.ErrInvalidAPIKey):
        // 無効なAPIキーを処理
    case errors.Is(err, wrapper.ErrInvalidModel):
        // 無効なモデルを処理
    case errors.Is(err, wrapper.ErrEmptyMessages):
        // 空のメッセージを処理
    case errors.Is(err, wrapper.ErrAPIRequest):
        // APIリクエストエラーを処理
    case errors.Is(err, wrapper.ErrUnsupportedProvider):
        // サポートされていないプロバイダを処理
    default:
        // その他のエラーを処理
    }
}
```

## 完全な例

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
    // .envファイルを読み込む
    err := godotenv.Load()
    if err != nil {
        log.Printf("警告: .envファイルの読み込みエラー: %v", err)
    }
}

func main() {
    // OpenAIクライアントを作成
    openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
    if err != nil {
        panic(err)
    }

    // 具体的な情報を含む1回目のテキスト生成
    openaiText1, err, openaiTokens1 := openaiClient.GenText(wrapper.GenTextParams{
        Model: wrapper.Model4_1Nano,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("OpenAI応答1: %s\n使用トークン数: %d\n\n", openaiText1, openaiTokens1)

    // 会話履歴を含む2回目のテキスト生成
    openaiText2, err, openaiTokens2 := openaiClient.GenText(wrapper.GenTextParams{
        Model: wrapper.Model4_1Nano,
        Messages: []wrapper.Message{
            {Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
            {Role: wrapper.RoleAssistant, Content: openaiText1},
            {Role: wrapper.RoleUser, Content: "田中さんの年齢、職業、趣味、そして先月何をしたか教えてください。"},
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("OpenAI応答2: %s\n使用トークン数: %d\n", openaiText2, openaiTokens2)
}
```

## API リファレンス

### 型

```go
// Provider はLLMプロバイダの種類を表す型です
type Provider string

const (
    ProviderOpenAI    Provider = "openai"
    ProviderAnthropic Provider = "anthropic"
    ProviderGemini    Provider = "gemini"
)

// Role はメッセージの役割を表す型です
type Role string

const (
    RoleUser      Role = "user"
    RoleAssistant Role = "assistant"
    RoleSystem    Role = "system"
)

// Message は会話内のメッセージを表す構造体です
type Message struct {
    Role    Role   `json:"role"`
    Content string `json:"content"`
}

// GenTextParams はテキスト生成のパラメータを表す構造体です
type GenTextParams struct {
    Model        Model     `json:"model"`
    Prompt       string    `json:"prompt,omitempty"`
    CacheEnabled bool      `json:"cache_enabled"`
    Messages     []Message `json:"messages"`
}

// LLMWrapper はLLMプロバイダとの対話のためのインターフェースです
type LLMWrapper interface {
    GenText(params GenTextParams) (string, error, int)
}
```

### 関数

```go
// NewClient は指定されたプロバイダの新しいLLMクライアントを作成します
func NewClient(provider Provider, apiKey string) (LLMWrapper, error)
```

### エラー定数

```go
var (
    ErrUnsupportedProvider = errors.New("unsupported provider")
    ErrInvalidAPIKey       = errors.New("invalid API key")
    ErrInvalidModel        = errors.New("invalid model")
    ErrEmptyMessages       = errors.New("empty messages")
    ErrAPIRequest          = errors.New("API request error")
)
```

## ライセンス

[MITライセンス](LICENSE)
