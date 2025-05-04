package infra

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient はOpenAIのAPIを呼び出すためのクライアントです
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient は新しいOpenAIClientを作成します
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(apiKey)
	return &OpenAIClient{
		client: client,
	}, nil
}

func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, maxTokens int) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:               openai.O3Mini,
		Messages:            messages,
		MaxCompletionTokens: maxTokens,
		// Temperature:         0.7,
		// TopP:                1.0,
		// N:                   1,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
