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
}

// NewOpenAIClient は、OpenAIクライアントの新しいインスタンスを作成します。
func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAIClient{client: client}
}

// GenText は、OpenAI APIを使用してテキストを生成します。
func (c *OpenAIClient) GenText(params models.GenTextParams) (string, error, int) {
	if params.Model == "" {
		return "", models.ErrInvalidModel, 0
	}

	if len(params.Messages) == 0 && params.Prompt == "" {
		return "", models.ErrEmptyMessages, 0
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
