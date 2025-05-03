package wrapper

import (
	"fmt"

	"github.com/obutora/ai-wrapper/internal/providers"
	"github.com/obutora/ai-wrapper/internal/types"
)

// Provider は、LLMプロバイダの種類を表す型です。
type Provider = types.Provider

// 利用可能なプロバイダの定数
const (
	ProviderOpenAI    = types.ProviderOpenAI
	ProviderAnthropic = types.ProviderAnthropic
	ProviderGemini    = types.ProviderGemini
)

// Role は、メッセージの役割を表す型です。
type Role = types.Role

// 利用可能なロールの定数
const (
	RoleUser      = types.RoleUser
	RoleAssistant = types.RoleAssistant
	RoleSystem    = types.RoleSystem
)

// Message は、LLMとのやり取りに使用するメッセージを表す構造体です。
type Message = types.Message

// GenTextParams は、テキスト生成に必要なパラメータを表す構造体です。
type GenTextParams = types.GenTextParams

// GenTextResponse は、テキスト生成の結果を表す構造体です。
type GenTextResponse = types.GenTextResponse

// LLMWrapper は、LLMプロバイダとのやり取りを抽象化するインターフェースです。
type LLMWrapper = types.LLMWrapper

// エラー定数
var (
	ErrUnsupportedProvider = types.ErrUnsupportedProvider
	ErrInvalidAPIKey       = types.ErrInvalidAPIKey
	ErrInvalidModel        = types.ErrInvalidModel
	ErrEmptyMessages       = types.ErrEmptyMessages
	ErrAPIRequest          = types.ErrAPIRequest
)

// NewClient は、指定されたプロバイダとAPIキーに基づいて新しいLLMWrapperクライアントを作成します。
func NewClient(provider Provider, apiKey string) (LLMWrapper, error) {
	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	switch provider {
	case ProviderOpenAI:
		return providers.NewOpenAIClient(apiKey), nil
	case ProviderAnthropic:
		return providers.NewAnthropicClient(apiKey), nil
	case ProviderGemini:
		client := providers.NewGeminiClient(apiKey)
		if client == nil {
			return nil, fmt.Errorf("failed to create Gemini client")
		}
		return client, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, provider)
	}
}
