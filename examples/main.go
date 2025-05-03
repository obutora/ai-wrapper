package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	wrapper "github.com/obutora/ai-wrapper"
)

func init() {
	// プロジェクトのルートディレクトリにある.envファイルを読み込む
	// 現在の作業ディレクトリから親ディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// examplesディレクトリから実行された場合は、親ディレクトリを使用
	envPath := filepath.Join(currentDir, "..", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// 親ディレクトリに.envがない場合は、現在のディレクトリを使用
		envPath = filepath.Join(currentDir, ".env")
	}

	// .envファイルを読み込む
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func main() {
	// OpenAIクライアントの作成
	openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	openaiText1, err, openaiTokens1 := openaiClient.GenText(wrapper.GenTextParams{
		Model: "gpt-4o",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("OpenAI Response 1: %s\nTokens used: %d\n\n", openaiText1, openaiTokens1)

	// 2回目のテキスト生成 - 会話履歴を含む
	openaiText2, err, openaiTokens2 := openaiClient.GenText(wrapper.GenTextParams{
		Model: "gpt-4o",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
			{Role: wrapper.RoleAssistant, Content: openaiText1},
			{Role: wrapper.RoleUser, Content: "田中さんの年齢、職業、趣味、そして先月何をしたか教えてください。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("OpenAI Response 2: %s\nTokens used: %d\n\n", openaiText2, openaiTokens2)

	// Anthropicクライアントの作成
	anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"))
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	anthropicText1, err, anthropicTokens1 := anthropicClient.GenText(wrapper.GenTextParams{
		Model: "claude-3-sonnet",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Anthropic Response 1: %s\nTokens used: %d\n\n", anthropicText1, anthropicTokens1)

	// 2回目のテキスト生成 - 会話履歴を含む
	anthropicText2, err, anthropicTokens2 := anthropicClient.GenText(wrapper.GenTextParams{
		Model: "claude-3-sonnet",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
			{Role: wrapper.RoleAssistant, Content: anthropicText1},
			{Role: wrapper.RoleUser, Content: "田中さんの年齢、職業、趣味、そして先月何をしたか教えてください。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Anthropic Response 2: %s\nTokens used: %d\n\n", anthropicText2, anthropicTokens2)

	// Geminiクライアントの作成
	geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	geminiText1, err, geminiTokens1 := geminiClient.GenText(wrapper.GenTextParams{
		Model: "gemini-2.0-flash",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gemini Response 1: %s\nTokens used: %d\n\n", geminiText1, geminiTokens1)

	// 2回目のテキスト生成 - 会話履歴を含む
	geminiText2, err, geminiTokens2 := geminiClient.GenText(wrapper.GenTextParams{
		Model: "gemini-2.0-flash",
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
			{Role: wrapper.RoleAssistant, Content: geminiText1},
			{Role: wrapper.RoleUser, Content: "田中さんの年齢、職業、趣味、そして先月何をしたか教えてください。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gemini Response 2: %s\nTokens used: %d\n", geminiText2, geminiTokens2)
}
