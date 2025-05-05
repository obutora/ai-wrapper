package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	wrapper "github.com/obutora/ai-wrapper"
	"github.com/obutora/ai-wrapper/internal/types"
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

// 個別のクライアントを使用した例
func traditionalExample() {
	fmt.Println("=== 従来の方法（個別のクライアント） ===")

	// OpenAIクライアントの作成
	openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	openaiText1, err, openaiTokens1 := openaiClient.GenText(wrapper.GenTextParams{
		Model: types.Model4_1Nano,
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
		Model: "gpt-4.1-nano-2025-04-14", // Model4_1Nano
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
		Model: types.ModelClaude3Haiku,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Anthropic Response 1: %s\nTokens used: %d\n\n", anthropicText1, anthropicTokens1)

	// Geminiクライアントの作成
	geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	geminiText1, err, geminiTokens1 := geminiClient.GenText(wrapper.GenTextParams{
		Model: types.ModelGemini25FlashPreview,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gemini Response 1: %s\nTokens used: %d\n\n", geminiText1, geminiTokens1)
}

// 統合クライアントを使用した例
func unifiedClientExample() {
	fmt.Println("=== 統合クライアントを使用した例 ===")

	// 各プロバイダのAPIキーをマップで作成
	apiKeys := map[wrapper.Provider]string{
		wrapper.ProviderOpenAI:    os.Getenv("OPENAI_API_KEY"),
		wrapper.ProviderAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
		wrapper.ProviderGemini:    os.Getenv("GEMINI_API_KEY"),
	}

	// 統合クライアントを作成
	client, err := wrapper.NewUnifiedClient(apiKeys)
	if err != nil {
		panic(err)
	}

	// OpenAIモデルを使用（自動的にOpenAIプロバイダが選択される）
	openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
		Model: types.Model4_1Nano,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "フランスの首都は何ですか？"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("OpenAI Response: %s\nTokens used: %d\n\n", openaiText, openaiTokens)

	// Anthropicモデルを使用（自動的にAnthropicプロバイダが選択される）
	anthropicText, err, anthropicTokens := client.GenText(wrapper.GenTextParams{
		Model: types.ModelClaude3Haiku,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "ドイツの首都は何ですか？"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Anthropic Response: %s\nTokens used: %d\n\n", anthropicText, anthropicTokens)

	// Geminiモデルを使用（自動的にGeminiプロバイダが選択される）
	geminiText, err, geminiTokens := client.GenText(wrapper.GenTextParams{
		Model: types.ModelGemini25FlashPreview,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "日本の首都は何ですか？"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Gemini Response: %s\nTokens used: %d\n\n", geminiText, geminiTokens)

	// カスタムモデルのマッピングを登録
	client.RegisterCustomModel("my-custom-model", wrapper.ProviderOpenAI)
	fmt.Println("カスタムモデル 'my-custom-model' を OpenAI プロバイダに登録しました")
}

func main() {
	// 従来の方法（個別のクライアント）を使用した例
	traditionalExample()

	fmt.Println("\n-----------------------------------\n")

	// 統合クライアントを使用した例
	unifiedClientExample()
}
