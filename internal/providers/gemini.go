package providers

import (
	"context"
	"fmt"

	"github.com/obutora/ai-wrapper/internal/types"
	"google.golang.org/genai"
)

// GeminiClient は、Geminiプロバイダのクライアントを表す構造体です。
type GeminiClient struct {
	client *genai.Client
}

// NewGeminiClient は、Geminiクライアントの新しいインスタンスを作成します。
func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		// エラーが発生した場合は、nilを返します
		// 実際のアプリケーションでは、エラーハンドリングを改善する必要があります
		return nil
	}
	return &GeminiClient{client: client}
}

// GenText は、Gemini APIを使用してテキストを生成します。
func (c *GeminiClient) GenText(params types.GenTextParams) (string, error, int) {
	if params.Model == "" {
		return "", types.ErrInvalidModel, 0
	}

	if len(params.Messages) == 0 && params.Prompt == "" {
		return "", types.ErrEmptyMessages, 0
	}

	ctx := context.Background()

	// メッセージを変換
	history := []*genai.Content{}
	if len(params.Messages) > 0 {
		for _, msg := range params.Messages {
			var role genai.Role
			switch msg.Role {
			case types.RoleUser:
				role = genai.RoleUser
			case types.RoleAssistant:
				role = genai.RoleModel
			case types.RoleSystem:
				// Geminiでは、システムメッセージはユーザーメッセージとして扱います
				role = genai.RoleUser
			default:
				role = genai.RoleUser
			}

			content := genai.NewContentFromText(msg.Content, role)
			history = append(history, content)
		}
	}

	// チャットセッションを作成
	chat, err := c.client.Chats.Create(ctx, string(params.Model), nil, history)
	if err != nil {
		return "", fmt.Errorf("%w: %v", types.ErrAPIRequest, err), 0
	}

	// メッセージを送信
	var message string
	if params.Prompt != "" {
		message = params.Prompt
	} else if len(params.Messages) > 0 {
		// 最後のユーザーメッセージを使用
		for i := len(params.Messages) - 1; i >= 0; i-- {
			if params.Messages[i].Role == types.RoleUser {
				message = params.Messages[i].Content
				break
			}
		}
	}

	// メッセージがない場合は、エラーを返します
	if message == "" {
		return "", types.ErrEmptyMessages, 0
	}

	// APIリクエストを実行
	res, err := chat.SendMessage(ctx, genai.Part{Text: message})
	if err != nil {
		return "", fmt.Errorf("%w: %v", types.ErrAPIRequest, err), 0
	}

	// レスポンスからテキストを取得
	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned"), 0
	}

	text := res.Candidates[0].Content.Parts[0].Text

	tokens := int(res.UsageMetadata.TotalTokenCount)

	return text, nil, tokens
}
