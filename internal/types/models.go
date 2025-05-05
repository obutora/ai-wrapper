package types

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
