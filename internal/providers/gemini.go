package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/obutora/ai-wrapper/models"
	"google.golang.org/genai"
)

// GeminiClient は、Geminiプロバイダのクライアントを表す構造体です。
type GeminiClient struct {
	client *genai.Client
	config models.Config
}

// NewGeminiClient は、Geminiクライアントの新しいインスタンスを作成します。
func NewGeminiClient(apiKey string, config models.Config) *GeminiClient {
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
	return &GeminiClient{client: client, config: config}
}

// GenText は、Gemini APIを使用してテキストを生成します。
func (c *GeminiClient) GenText(params models.GenTextParams) (string, error, int) {
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

	// メッセージを変換
	history := []*genai.Content{}
	if len(params.Messages) > 0 {
		for _, msg := range params.Messages {
			var role genai.Role
			switch msg.Role {
			case models.RoleUser:
				role = genai.RoleUser
			case models.RoleAssistant:
				role = genai.RoleModel
			case models.RoleSystem:
				// Geminiでは、システムメッセージはユーザーメッセージとして扱います
				role = genai.RoleUser
			default:
				role = genai.RoleUser
			}

			content := genai.NewContentFromText(msg.Content, role)
			history = append(history, content)
		}
	}

	conf := &genai.GenerateContentConfig{
		MaxOutputTokens: int32(c.config.MaxToken),
		// ThinkingLevel をモデル世代に応じて thinkingLevel / thinkingBudget に変換（非対応モデルでは nil）
		ThinkingConfig: geminiThinkingConfig(params.Model, params.ThinkingLevel),
	}

	// チャットセッションを作成
	chat, err := c.client.Chats.Create(ctx, string(params.Model), conf, history)
	if err != nil {
		return "", fmt.Errorf("%w: %v", models.ErrAPIRequest, err), 0
	}

	// メッセージを送信
	var message string
	if params.Prompt != "" {
		message = params.Prompt
	} else if len(params.Messages) > 0 {
		// 最後のユーザーメッセージを使用
		for i := len(params.Messages) - 1; i >= 0; i-- {
			if params.Messages[i].Role == models.RoleUser {
				message = params.Messages[i].Content
				break
			}
		}
	}

	// メッセージがない場合は、エラーを返します
	if message == "" {
		return "", models.ErrEmptyMessages, 0
	}

	// APIリクエストを実行
	res, err := chat.SendMessage(ctx, genai.Part{Text: message})
	if isGeminiMinimalUnsupported(conf, err) {
		// thinkingLevel=MINIMAL 非対応モデルだった場合は LOW に落として 1 回だけ再試行する
		conf.ThinkingConfig.ThinkingLevel = genai.ThinkingLevelLow
		if chat, err = c.client.Chats.Create(ctx, string(params.Model), conf, history); err == nil {
			res, err = chat.SendMessage(ctx, genai.Part{Text: message})
		}
	}
	if err != nil {
		return "", fmt.Errorf("%w: %v", models.ErrAPIRequest, err), 0
	}

	// レスポンスからテキストを取得
	if len(res.Candidates) == 0 || res.Candidates[0].Content == nil || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned"), 0
	}

	// thinking 有効時は思考パートが含まれ得るため、テキストパートのみを連結する
	var sb strings.Builder
	for _, part := range res.Candidates[0].Content.Parts {
		if part == nil || part.Thought {
			continue
		}
		sb.WriteString(part.Text)
	}
	text := sb.String()

	tokens := int(res.UsageMetadata.TotalTokenCount)

	return text, nil, tokens
}
