# Go用AIラッパー

[English](README.md) | [日本語](README_ja.md)

# Go用AIラッパー

一貫したインターフェースを通じて複数のLLMプロバイダ（OpenAI、Anthropic、Gemini）と対話するための統一されたGoライブラリです。

## 特徴

- 複数のLLMプロバイダに対する単一の一貫したAPI
- OpenAI、Anthropic、Geminiの最新モデルをサポート
- メッセージ履歴による簡単な会話処理
- トークン使用量の追跡
- プロバイダ横断の推論深さ指定（`ThinkingLevel`）。各プロバイダのネイティブな thinking / reasoning パラメータへ自動変換
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
        Model: models.ModelGPT4o,
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
        Model: models.ModelGPT4o,
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
        Model: models.ModelClaudeOpus5,
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
        Model: models.ModelGeminiProLatest,
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

- `ModelClaudeOpus5` - Claude Opus 5（最新の Opus）
- `ModelClaudeSonnet5` - Claude Sonnet 5（最新の Sonnet）
- `ModelClaudeHaiku45` - Claude Haiku 4.5（最新の Haiku）
- `ModelClaudeFable51` / `ModelClaudeFable5` - Claude Fable 5.1 / 5
- `ModelClaudeOpus48` / `ModelClaudeOpus47` / `ModelClaudeOpus46` / `ModelClaudeSonnet46` - Claude 4.6〜4.8 系
- `ModelClaudeOpus45` / `ModelClaudeSonnet45` - Claude 4.5 系
- `ModelClaude3Opus`, `ModelClaude37Sonnet`, `ModelClaude3Haiku`, `ModelClaudeOpus41` - **非推奨**（提供終了モデル。互換性のために残置）

### Gemini

特定バージョンに固定する必要がなければ `*-latest` エイリアスを推奨します。

- `ModelGeminiFlashLatest` - `gemini-flash-latest`（常に最新の Flash。現在は Gemini 3.8 Flash）
- `ModelGeminiFlashLiteLatest` - `gemini-flash-lite-latest`（現在は Gemini 3.5 Flash-Lite）
- `ModelGeminiProLatest` - `gemini-pro-latest`（現在は Gemini 3.1 Pro Preview）
- `ModelGemini38Flash` / `ModelGemini35FlashLite` / `ModelGemini31ProPreview` - 上記のバージョン固定版
- `ModelGemini25Flash` / `ModelGemini25FlashLite` - Gemini 2.5 Flash / Flash-Lite
- `ModelGemini25Pro`, `ModelGemini20Flash`, `ModelGemini20Pro`, `ModelGemini25FlashPreview`, `ModelGemini25ProPreview` - **非推奨**（提供終了モデル）

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
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "フランスの首都はどこですか？"},
    },
})

// 会話履歴を含む
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "フランスの首都はどこですか？"},
        {Role: wrapper.RoleAssistant, Content: "フランスの首都はパリです。"},
        {Role: wrapper.RoleUser, Content: "その人口は？"},
    },
})

// システムメッセージを含む（サポートされているプロバイダ向け）
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model: models.ModelGPT4o,
    Messages: []wrapper.Message{
        {Role: wrapper.RoleSystem, Content: "あなたは簡潔な回答を提供する役立つアシスタントです。"},
        {Role: wrapper.RoleUser, Content: "フランスの首都はどこですか？"},
    },
})
```

### ThinkingLevel（推論の深さ）

`GenTextParams.ThinkingLevel` で、モデルが回答前に行う推論（thinking）の深さを指定できます。
同じ値が全プロバイダで使え、ラッパー内部で各プロバイダのネイティブなパラメータに変換されます。
推論をサポートしないモデルでは無視されます。

```go
text, err, tokens := client.GenText(wrapper.GenTextParams{
    Model:         models.ModelGeminiFlashLatest, // 現在は gemini-3.8-flash に解決される
    ThinkingLevel: wrapper.ThinkingLevelHigh,      // minimal / low / medium / high / max
    Messages: []wrapper.Message{
        {Role: wrapper.RoleUser, Content: "次の問題を段階的に解いてください: ..."},
    },
})
```

| `ThinkingLevel` | Gemini 3 以降 / `*-latest` (`thinkingLevel`) | Gemini 2.5 (`thinkingBudget`) | OpenAI 推論モデル (`reasoning_effort`) | Claude 4.6 以降 (`thinking: adaptive` + `output_config.effort`) | Claude 3.7〜4.5 (`thinking: enabled` + `budget_tokens`) |
|---|---|---|---|---|---|
| `minimal` | `LOW`（Flash-Lite は `MINIMAL`） | 0（Pro は 128） | `minimal`（GPT-5.1 以降は `none`、o 系は `low`） | `low` | 1024 |
| `low` | `LOW` | 1024 | `low` | `low` | 2048 |
| `medium` | `MEDIUM` | 8192 | `medium` | `medium` | 8192 |
| `high` | `HIGH` | 24576 | `high` | `high` | 16384 |
| `max` | `HIGH` | 24576（Pro は 32768） | `high`（GPT-5.2 以降は `xhigh`） | `max` | 32768 |

補足:

- 推論非対応モデル（Gemini 2.0 以前、GPT-4 系、Claude 3.x など）では `ThinkingLevel` は無視されます。
- `gemini-flash-latest` / `gemini-pro-latest` / `gemini-flash-lite-latest` のようなバージョン無しの別名は現在 Gemini 3 系に解決されるため、`thinkingLevel` を使います。
- `MINIMAL` の対応はモデルごとに異なります（2026-09 時点: 3.5 Flash-Lite は対応、3.8 Flash と 3.1 Pro は 400 で拒否）。ラッパーは Flash-Lite 系にのみ `MINIMAL` を送り、拒否された場合は `LOW` で 1 回だけ再試行します。
- Claude 3.7〜4.5 では `budget_tokens` が `MaxToken` 未満である必要があるため、`MaxToken - 1` に丸められます。`MaxToken` が 1024 以下の場合は `ErrInvalidThinkingConfig` を返します。
- 未知の値を指定すると、API を呼ぶ前に `ErrInvalidThinkingLevel` を返します。
- 返却されるトークン数には推論トークンも含まれます。

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

// ThinkingLevel は推論の深さをプロバイダ横断で指定する型です
type ThinkingLevel string

const (
    ThinkingLevelDefault ThinkingLevel = ""        // プロバイダ・モデルの既定動作
    ThinkingLevelMinimal ThinkingLevel = "minimal"
    ThinkingLevelLow     ThinkingLevel = "low"
    ThinkingLevelMedium  ThinkingLevel = "medium"
    ThinkingLevelHigh    ThinkingLevel = "high"
    ThinkingLevelMax     ThinkingLevel = "max"
)

// GenTextParams はテキスト生成のパラメータを表す構造体です
type GenTextParams struct {
    Model         Model         `json:"model"`
    Prompt        string        `json:"prompt,omitempty"`
    CacheEnabled  bool          `json:"cache_enabled"`
    Messages      []Message     `json:"messages"`
    ThinkingLevel ThinkingLevel `json:"thinking_level,omitempty"`
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
    ErrInvalidThinkingLevel  = errors.New("invalid thinking level")
    ErrInvalidThinkingConfig = errors.New("invalid thinking config")
)
```

## ライセンス

[MITライセンス](LICENSE)
