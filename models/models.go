package models

import (
	"fmt"
	"regexp"
	"strconv"
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

	// Anthropicモデル（2026-09 時点の最新: Opus 5 / Sonnet 5 / Haiku 4.5）
	ModelClaudeOpus5    Model = "claude-opus-5"
	ModelClaudeSonnet5  Model = "claude-sonnet-5"
	ModelClaudeHaiku45  Model = "claude-haiku-4-5"
	ModelClaudeFable51  Model = "claude-fable-5-1"
	ModelClaudeFable5   Model = "claude-fable-5"
	ModelClaudeOpus48   Model = "claude-opus-4-8"
	ModelClaudeOpus47   Model = "claude-opus-4-7"
	ModelClaudeOpus46   Model = "claude-opus-4-6"
	ModelClaudeSonnet46 Model = "claude-sonnet-4-6"
	ModelClaudeOpus45   Model = "claude-opus-4-5"
	ModelClaudeSonnet45 Model = "claude-sonnet-4-5"

	// Deprecated: 提供終了モデルです。ModelClaudeOpus5 を使用してください。
	ModelClaudeOpus41 Model = "claude-opus-4-1"
	// Deprecated: 提供終了モデルです。ModelClaudeOpus5 を使用してください。
	ModelClaude3Opus Model = "claude-3-opus"
	// Deprecated: 提供終了モデルです。ModelClaudeSonnet5 を使用してください。
	ModelClaude37Sonnet Model = "claude-3.7-sonnet"
	// Deprecated: 提供終了モデルです。ModelClaudeHaiku45 を使用してください。
	ModelClaude3Haiku Model = "claude-3-haiku"

	// Geminiモデル（2026-09 時点の最新: 3.8 Flash / 3.5 Flash-Lite / 3.1 Pro Preview）
	// *-latest エイリアスは常に最新版を指すため、特定バージョンに固定する必要がなければこちらを推奨します。
	ModelGeminiFlashLatest     Model = "gemini-flash-latest"      // 現在は gemini-3.8-flash に解決される
	ModelGeminiFlashLiteLatest Model = "gemini-flash-lite-latest" // 現在は gemini-3.5-flash-lite に解決される
	ModelGeminiProLatest       Model = "gemini-pro-latest"        // 現在は gemini-3.1-pro-preview に解決される
	ModelGemini38Flash         Model = "gemini-3.8-flash"
	ModelGemini35FlashLite     Model = "gemini-3.5-flash-lite"
	ModelGemini31ProPreview    Model = "gemini-3.1-pro-preview"
	ModelGemini25Flash         Model = "gemini-2.5-flash"
	ModelGemini25FlashLite     Model = "gemini-2.5-flash-lite"

	// Deprecated: 提供終了モデルです。ModelGeminiFlashLatest を使用してください。
	ModelGemini20Flash Model = "gemini-2.0-flash"
	// Deprecated: 提供終了モデルです。ModelGeminiProLatest を使用してください。
	ModelGemini20Pro Model = "gemini-2.0-pro"
	// Deprecated: 提供終了モデルです。ModelGeminiFlashLatest を使用してください。
	ModelGemini25FlashPreview Model = "gemini-2.5-flash-preview-04-17"
	// Deprecated: 提供終了モデルです。ModelGeminiProLatest を使用してください。
	ModelGemini25ProPreview Model = "gemini-2.5-pro-preview-03-25"
	// Deprecated: 提供終了モデルです（gemini-2.5-pro 自体も新規利用不可）。ModelGeminiProLatest を使用してください。
	ModelGemini25Pro Model = "gemini-2.5-pro-exp-03-25"
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
		return "claude-3-opus-latest"
	case ModelClaude37Sonnet:
		return "claude-3-7-sonnet-latest"
	case ModelClaude3Haiku:
		return "claude-3-5-haiku-latest"

	default:
		return string(m)
	}
}

var (
	// claude-opus-4-6 / claude-sonnet-5 / claude-fable-5-1 / claude-sonnet-4-20250514 など
	anthropicNamedVersionRe = regexp.MustCompile(`^claude-(?:opus|sonnet|haiku|fable|mythos)-(\d+)(?:-(\d+))?`)
	// claude-3-7-sonnet-latest / claude-3-5-haiku-latest / claude-3-opus-latest など
	anthropicLegacyVersionRe = regexp.MustCompile(`^claude-(\d+)(?:-(\d+))?-(?:opus|sonnet|haiku)`)
)

// AnthropicVersion は、Claude のモデル名から世代（major, minor）を取り出します。
// 例: claude-opus-4-6 → (4, 6)、claude-sonnet-5 → (5, 0)、claude-3-7-sonnet-latest → (3, 7)。
// 日付サフィックス（claude-sonnet-4-20250514 の 20250514 等）はマイナーバージョンとして扱いません。
// Claude のモデル名として解釈できない場合は ok=false を返します。
func (m Model) AnthropicVersion() (major, minor int, ok bool) {
	name := strings.ToLower(string(m.ToAnthropicModel()))
	if sub := anthropicNamedVersionRe.FindStringSubmatch(name); sub != nil {
		major, _ = strconv.Atoi(sub[1])
		if sub[2] != "" && len(sub[2]) <= 2 {
			minor, _ = strconv.Atoi(sub[2])
		}
		return major, minor, true
	}
	if sub := anthropicLegacyVersionRe.FindStringSubmatch(name); sub != nil {
		major, _ = strconv.Atoi(sub[1])
		if sub[2] != "" {
			minor, _ = strconv.Atoi(sub[2])
		}
		return major, minor, true
	}
	return 0, 0, false
}

// SupportsTemperature は、モデルがtemperatureパラメータをサポートするかを返します。
// Claude 4.5 以降のモデルでは temperature が非推奨（4.7 以降・5 系では送信すると API エラー）のため false を返します。
func (m Model) SupportsTemperature() bool {
	major, minor, ok := m.AnthropicVersion()
	if !ok {
		return true
	}
	return major < 4 || (major == 4 && minor < 5)
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

// ThinkingLevel は、推論（thinking / reasoning）の深さをプロバイダ横断で指定する型です。
// 各プロバイダのネイティブなパラメータ（Gemini: thinkingLevel / thinkingBudget、
// OpenAI: reasoning_effort、Anthropic: thinking + output_config.effort）へ内部で変換されます。
// 推論をサポートしないモデルに対して指定した場合は無視されます。
type ThinkingLevel string

const (
	// ThinkingLevelDefault は、未指定（各プロバイダ・モデルの既定動作）を表します。
	ThinkingLevelDefault ThinkingLevel = ""
	// ThinkingLevelMinimal は、推論を最小限に抑えます（応答速度・コスト優先）。
	ThinkingLevelMinimal ThinkingLevel = "minimal"
	// ThinkingLevelLow は、軽い推論を行います。
	ThinkingLevelLow ThinkingLevel = "low"
	// ThinkingLevelMedium は、中程度の推論を行います。
	ThinkingLevelMedium ThinkingLevel = "medium"
	// ThinkingLevelHigh は、深い推論を行います。
	ThinkingLevelHigh ThinkingLevel = "high"
	// ThinkingLevelMax は、そのモデルで利用可能な最大の推論を行います。
	ThinkingLevelMax ThinkingLevel = "max"
)

// Validate は、ThinkingLevel が既知の値であるかを検証します。
func (l ThinkingLevel) Validate() error {
	switch l {
	case ThinkingLevelDefault, ThinkingLevelMinimal, ThinkingLevelLow,
		ThinkingLevelMedium, ThinkingLevelHigh, ThinkingLevelMax:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrInvalidThinkingLevel, string(l))
	}
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
	// ThinkingLevel は、推論（thinking）の深さです。
	// 空（ThinkingLevelDefault）の場合は、各プロバイダ・モデルの既定動作になります。
	ThinkingLevel ThinkingLevel `json:"thinking_level,omitempty"`
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
