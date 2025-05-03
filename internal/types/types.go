package types

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
	// Model は、使用するLLMモデルの名前です。
	Model string `json:"model"`
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
