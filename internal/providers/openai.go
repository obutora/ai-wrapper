package providers

import (
	"context"
	"fmt"

	"github.com/obutora/ai-wrapper/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIClient は、OpenAIプロバイダのクライアントを表す構造体です。
type OpenAIClient struct {
	client openai.Client
	config models.Config
}

// NewOpenAIClient は、OpenAIクライアントの新しいインスタンスを作成します。
func NewOpenAIClient(apiKey string, config models.Config) *OpenAIClient {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAIClient{client: client, config: config}
}

// GenText は、OpenAI APIを使用してテキストを生成します。
func (c *OpenAIClient) GenText(params models.GenTextParams) (string, error, int) {
	if params.Model == "" {
		return "", models.ErrInvalidModel, 0
	}

	if len(params.Messages) == 0 && params.Prompt == "" {
		return "", models.ErrEmptyMessages, 0
	}

	if err := params.ThinkingLevel.Validate(); err != nil {
		return "", err, 0
	}

	ctx := context.Background()
	messages := []openai.ChatCompletionMessageParamUnion{}

	// メッセージがある場合は、それらを変換して使用します
	if len(params.Messages) > 0 {
		for _, msg := range params.Messages {
			switch msg.Role {
			case models.RoleUser:
				messages = append(messages, openai.UserMessage(msg.Content))
			case models.RoleAssistant:
				messages = append(messages, openai.AssistantMessage(msg.Content))
			case models.RoleSystem:
				messages = append(messages, openai.SystemMessage(msg.Content))
			default:
				messages = append(messages, openai.UserMessage(msg.Content))
			}
		}
	} else if params.Prompt != "" {
		// プロンプトがある場合は、ユーザーメッセージとして追加します
		messages = append(messages, openai.UserMessage(params.Prompt))
	}

	// モデル名を取得
	model := models.Model(params.Model).ToOpenAIModel()

	// APIリクエストパラメータを作成
	chatParams := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
		// max_tokens は推論モデル（o 系 / GPT-5 系）で拒否されるため、全モデルで受理される
		// max_completion_tokens のみを送信する
		MaxCompletionTokens: openai.Int(int64(c.config.MaxToken)),
		// ThinkingLevel を reasoning_effort に変換（推論モデル以外では空 = 未送信）
		ReasoningEffort: openAIReasoningEffort(params.Model, params.ThinkingLevel),
	}

	// APIリクエストを実行
	completion, err := c.client.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		return "", fmt.Errorf("%w: %v", models.ErrAPIRequest, err), 0
	}

	// レスポンスからテキストとトークン数を取得
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned"), 0
	}

	text := completion.Choices[0].Message.Content
	tokens := int(completion.Usage.TotalTokens)

	return text, nil, tokens
}
