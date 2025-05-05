package providers

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/obutora/ai-wrapper/internal/types"
	"github.com/obutora/ai-wrapper/models"
)

// AnthropicClient は、Anthropicプロバイダのクライアントを表す構造体です。
type AnthropicClient struct {
	client anthropic.Client
}

// NewAnthropicClient は、Anthropicクライアントの新しいインスタンスを作成します。
func NewAnthropicClient(apiKey string) *AnthropicClient {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &AnthropicClient{client: client}
}

// GenText は、Anthropic APIを使用してテキストを生成します。
func (c *AnthropicClient) GenText(params types.GenTextParams) (string, error, int) {
	if params.Model == "" {
		return "", types.ErrInvalidModel, 0
	}

	if len(params.Messages) == 0 && params.Prompt == "" {
		return "", types.ErrEmptyMessages, 0
	}

	ctx := context.Background()
	messages := []anthropic.MessageParam{}

	// メッセージがある場合は、それらを変換して使用します
	if len(params.Messages) > 0 {
		for _, msg := range params.Messages {
			var role anthropic.MessageParamRole
			switch msg.Role {
			case types.RoleUser:
				role = anthropic.MessageParamRoleUser
			case types.RoleAssistant:
				role = anthropic.MessageParamRoleAssistant
			case types.RoleSystem:
				// Anthropicでは、システムメッセージは特別な処理が必要
				// システムメッセージはシステムプロンプトとして扱います
				continue
			default:
				role = anthropic.MessageParamRoleUser
			}

			content := []anthropic.ContentBlockParamUnion{
				{
					OfRequestTextBlock: &anthropic.TextBlockParam{
						Text: msg.Content,
						// cacheを有効化
						CacheControl: anthropic.CacheControlEphemeralParam{},
					},
				},
			}

			messages = append(messages, anthropic.MessageParam{
				Role:    role,
				Content: content,
			})
		}
	} else if params.Prompt != "" {
		// プロンプトがある場合は、ユーザーメッセージとして追加します
		content := []anthropic.ContentBlockParamUnion{
			{
				OfRequestTextBlock: &anthropic.TextBlockParam{
					Text: params.Prompt,
				},
			},
		}

		messages = append(messages, anthropic.MessageParam{
			Role:    anthropic.MessageParamRoleUser,
			Content: content,
		})
	}

	// モデル名を取得
	model := models.Model(params.Model).ToAnthropicModel()

	// APIリクエストパラメータを作成
	messageParams := anthropic.MessageNewParams{
		Model:     model,
		Messages:  messages,
		MaxTokens: 1024 * 100, // NOTE: トークン数は適宜調整してください
	}

	// APIリクエストを実行
	response, err := c.client.Messages.New(ctx, messageParams)
	if err != nil {
		return "", fmt.Errorf("%w: %v", types.ErrAPIRequest, err), 0
	}

	// レスポンスからテキストを取得
	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content returned"), 0
	}

	// レスポンスからテキストを取得
	text := response.Content[0].Text

	// トークン数を取得
	tokens := int(response.Usage.OutputTokens + response.Usage.InputTokens)

	return text, nil, tokens
}
