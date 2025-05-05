package wrapper

import (
	"fmt"

	"github.com/obutora/ai-wrapper/internal/providers"
	"github.com/obutora/ai-wrapper/models"
)

// Provider は、LLMプロバイダの種類を表す型です。
type Provider = models.Provider

// 利用可能なプロバイダの定数
const (
	ProviderOpenAI    = models.ProviderOpenAI
	ProviderAnthropic = models.ProviderAnthropic
	ProviderGemini    = models.ProviderGemini
)

// Model は、LLMモデルの種類を表す型です。
type Model = models.Model

// Role は、メッセージの役割を表す型です。
type Role = models.Role

// 利用可能なロールの定数
const (
	RoleUser      = models.RoleUser
	RoleAssistant = models.RoleAssistant
	RoleSystem    = models.RoleSystem
)

// Message は、LLMとのやり取りに使用するメッセージを表す構造体です。
type Message = models.Message

// GenTextParams は、テキスト生成に必要なパラメータを表す構造体です。
type GenTextParams = models.GenTextParams

// GenTextResponse は、テキスト生成の結果を表す構造体です。
type GenTextResponse = models.GenTextResponse

// LLMWrapper は、LLMプロバイダとのやり取りを抽象化するインターフェースです。
type LLMWrapper = models.LLMWrapper

// エラー定数
var (
	ErrUnsupportedProvider = models.ErrUnsupportedProvider
	ErrInvalidAPIKey       = models.ErrInvalidAPIKey
	ErrInvalidModel        = models.ErrInvalidModel
	ErrEmptyMessages       = models.ErrEmptyMessages
	ErrAPIRequest          = models.ErrAPIRequest
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

// UnifiedClient は、複数のプロバイダーを統合したクライアントです。
// モデル名から自動的に適切なプロバイダーを選択します。
type UnifiedClient struct {
	clients              map[Provider]LLMWrapper
	customModelProviders map[Model]Provider // カスタムモデル名とプロバイダーのマッピング
}

// NewUnifiedClient は、複数のプロバイダーを統合した新しいクライアントを作成します。
// APIキーのマップを受け取り、各プロバイダーのクライアントを初期化します。
func NewUnifiedClient(apiKeys map[Provider]string) (*UnifiedClient, error) {
	clients := make(map[Provider]LLMWrapper)

	for provider, apiKey := range apiKeys {
		client, err := NewClient(provider, apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for provider %s: %w", provider, err)
		}
		clients[provider] = client
	}

	return &UnifiedClient{
		clients:              clients,
		customModelProviders: make(map[Model]Provider),
	}, nil
}

// RegisterCustomModel は、カスタムモデル名とプロバイダーのマッピングを登録します。
// これにより、標準のパターンマッチングでは検出できない特殊なモデル名も手動で登録できます。
func (c *UnifiedClient) RegisterCustomModel(model Model, provider Provider) {
	c.customModelProviders[model] = provider
}

// getProviderForModel は、モデル名からプロバイダーを判定します。
func (c *UnifiedClient) getProviderForModel(model Model) Provider {
	// カスタムマッピングを確認
	if provider, ok := c.customModelProviders[model]; ok {
		return provider
	}

	// 標準の判定ロジックを使用
	return model.GetProvider()
}

// GenText は、モデル名から適切なプロバイダーを選択してテキストを生成します。
func (c *UnifiedClient) GenText(params GenTextParams) (string, error, int) {
	provider := c.getProviderForModel(params.Model)

	if provider == "" {
		return "", fmt.Errorf("%w: could not determine provider for model %s", ErrUnsupportedProvider, params.Model), 0
	}

	client, ok := c.clients[provider]
	if !ok {
		return "", fmt.Errorf("%w: no client for provider %s", ErrUnsupportedProvider, provider), 0
	}

	return client.GenText(params)
}
