package models

import (
	"regexp"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go/shared"
)

// Model は、LLMモデルの種類を表す型です。
type Model string

// 利用可能なモデルの定数
const (
	// OpenAIモデル
	ModelGPT4o      Model = "gpt-4o"
	ModelGPT4       Model = "gpt-4"
	ModelGPT35Turbo Model = "gpt-3.5-turbo"
	ModelO3Mini     Model = "o3-mini-2025-01-31k"
	ModelO4Mini     Model = "o4-mini-2025-04-16"
	Model4_1Nano    Model = "gpt-4.1-nano-2025-04-14"
	ModelO3         Model = "o3-2025-04-16"

	// Anthropicモデル
	ModelClaude3Opus    Model = "claude-3-opus"
	ModelClaude37Sonnet Model = "claude-3.7-sonnet"
	ModelClaude3Haiku   Model = "claude-3-haiku"

	// Geminiモデル
	ModelGemini20Flash        Model = "gemini-2.0-flash"
	ModelGemini20Pro          Model = "gemini-2.0-pro"
	ModelGemini25FlashPreview Model = "gemini-2.5-flash-preview-04-17"
	ModelGemini25ProPreview   Model = "gemini-2.5-pro-preview-03-25"
	ModelGemini25Pro          Model = "gemini-2.5-pro-exp-03-25"
)

// Provider は、LLMプロバイダの種類を表す型です。
type Provider string

const (
	// ProviderOpenAI は、OpenAIプロバイダを表します。
	ProviderOpenAI Provider = "openai"
	// ProviderAnthropic は、Anthropicプロバイダを表します。
	ProviderAnthropic Provider = "anthropic"
	// ProviderGemini は、Geminiプロバイダを表します。
	ProviderGemini Provider = "gemini"
)

// ToOpenAIModel は、共通モデル型をOpenAI SDKのモデル型に変換します。
func (m Model) ToOpenAIModel() shared.ChatModel {
	switch m {
	case ModelGPT4o:
		return shared.ChatModelGPT4o
	case ModelGPT4:
		return shared.ChatModelGPT4
	case ModelO3Mini:
		return shared.ChatModelO3Mini
	case ModelO4Mini:
		return "o4-mini-2025-04-16"
	case Model4_1Nano:
		return "gpt-4.1-nano-2025-04-14"
	case ModelO3:
		return "o3-2025-04-16"
	// case ModelGPT35Turbo:
	// 	return shared.ChatModelGPT35Turbo
	default:
		return string(m)
	}
}

// ToAnthropicModel は、共通モデル型をAnthropic SDKのモデル型に変換します。
func (m Model) ToAnthropicModel() anthropic.Model {
	switch m {
	case ModelClaude3Opus:
		return anthropic.ModelClaude3OpusLatest
	case ModelClaude37Sonnet:
		return anthropic.ModelClaude3_7SonnetLatest
	case ModelClaude3Haiku:
		return anthropic.ModelClaude3_5HaikuLatest

	default:
		return string(m)
	}
}

// GetProvider はモデル名からプロバイダーを判定します
func (m Model) GetProvider() Provider {
	modelName := string(m)

	// OpenAIモデルのパターン
	// - "gpt-" で始まるモデル (例: gpt-4, gpt-3.5-turbo)
	// - "o" + 数字 + "-" で始まるモデル (例: o1-, o2-, o3-, o4-)
	if strings.HasPrefix(modelName, "gpt-") ||
		regexp.MustCompile(`^o\d+-`).MatchString(modelName) {
		return ProviderOpenAI
	}

	// Anthropicモデルのパターン
	if strings.HasPrefix(modelName, "claude-") {
		return ProviderAnthropic
	}

	// Geminiモデルのパターン
	if strings.HasPrefix(modelName, "gemini-") {
		return ProviderGemini
	}

	// 不明なモデル
	return ""
}

// Role は、メッセージの役割を表す型です。
type Role string

const (
	// RoleUser は、ユーザーからのメッセージを表します。
	RoleUser Role = "user"
	// RoleAssistant は、アシスタントからのメッセージを表します。
	RoleAssistant Role = "assistant"
	// RoleSystem は、システムからのメッセージを表します。
	RoleSystem Role = "system"
)

// Message は、LLMとのやり取りに使用するメッセージを表す構造体です。
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// GenTextParams は、テキスト生成に必要なパラメータを表す構造体です。
type GenTextParams struct {
	// Model は、使用するLLMモデルです。
	Model Model `json:"model"`
	// Prompt は、単一のプロンプトテキストです。
	Prompt string `json:"prompt,omitempty"`
	// CacheEnabled は、キャッシュを有効にするかどうかを指定します。
	CacheEnabled bool `json:"cache_enabled"`
	// Messages は、会話履歴を表すメッセージのスライスです。
	Messages []Message `json:"messages"`
}

// GenTextResponse は、テキスト生成の結果を表す構造体です。
type GenTextResponse struct {
	// Text は、生成されたテキストです。
	Text string
	// Tokens は、使用されたトークン数です。
	Tokens int
}

// LLMWrapper は、LLMプロバイダとのやり取りを抽象化するインターフェースです。
type LLMWrapper interface {
	// GenText は、指定されたパラメータに基づいてテキストを生成します。
	// 生成されたテキスト、エラー、使用されたトークン数を返します。
	GenText(params GenTextParams) (string, error, int)
}
