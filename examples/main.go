package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	wrapper "github.com/obutora/ai-wrapper"
	"github.com/obutora/ai-wrapper/models"
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

	config := models.Config{
		MaxToken: 1000, // 最大トークン数
	}

	// OpenAIクライアントの作成
	openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"), config)
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	openaiText1, err, openaiTokens1 := openaiClient.GenText(wrapper.GenTextParams{
		Model: models.Model4_1Nano,
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
	anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"), config)
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	anthropicText1, err, anthropicTokens1 := anthropicClient.GenText(wrapper.GenTextParams{
		Model: models.ModelClaudeSonnet5,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: "田中太郎さんは東京都在住の42歳のエンジニアで、趣味は登山と写真撮影です。彼は先月、富士山に登りました。"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Anthropic Response 1: %s\nTokens used: %d\n\n", anthropicText1, anthropicTokens1)

	// Geminiクライアントの作成
	geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"), config)
	if err != nil {
		panic(err)
	}

	// 1回目のテキスト生成 - 具体的な情報を含む
	geminiText1, err, geminiTokens1 := geminiClient.GenText(wrapper.GenTextParams{
		Model: models.ModelGeminiFlashLatest,
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

	config := models.Config{
		MaxToken: 1000, // 最大トークン数
	}

	// 統合クライアントを作成
	client, err := wrapper.NewUnifiedClient(apiKeys, config)
	if err != nil {
		panic(err)
	}

	// OpenAIモデルを使用（自動的にOpenAIプロバイダが選択される）
	openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
		Model: models.Model4_1Nano,
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
		Model: models.ModelClaudeSonnet5,
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
		Model: models.ModelGeminiFlashLatest,
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

// ThinkingLevel（推論の深さ）を指定した例
// 同じ ThinkingLevel を渡すだけで、各プロバイダのネイティブなパラメータに変換されます。
//   - Gemini 3 以降: thinkingLevel / Gemini 2.5: thinkingBudget
//   - OpenAI 推論モデル: reasoning_effort
//   - Claude 4.6 以降: thinking(adaptive) + output_config.effort / Claude 3.7〜4.5: thinking(enabled) + budget_tokens
func thinkingLevelExample() {
	fmt.Println("=== ThinkingLevel を指定した例 ===")

	apiKeys := map[wrapper.Provider]string{
		wrapper.ProviderOpenAI:    os.Getenv("OPENAI_API_KEY"),
		wrapper.ProviderAnthropic: os.Getenv("ANTHROPIC_API_KEY"),
		wrapper.ProviderGemini:    os.Getenv("GEMINI_API_KEY"),
	}

	// Claude 3.7〜4.5 で thinking を使う場合、budget_tokens < MaxToken である必要があるため余裕を持たせる
	config := models.Config{
		MaxToken: 8000,
	}

	client, err := wrapper.NewUnifiedClient(apiKeys, config)
	if err != nil {
		panic(err)
	}

	question := "3つの箱A,B,Cのうち1つに賞品が入っています。Aを選んだ後、司会者がCを開けて空だと示しました。Bに変更すべきですか？理由を簡潔に。"

	// Gemini: gemini-flash-latest（現在 gemini-3.8-flash）や gemini-3.* なら thinkingLevel、gemini-2.5-* なら thinkingBudget に変換される
	geminiText, err, geminiTokens := client.GenText(wrapper.GenTextParams{
		Model:         models.ModelGeminiFlashLatest,
		ThinkingLevel: wrapper.ThinkingLevelHigh,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: question},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Gemini (high): %s\nTokens used: %d\n\n", geminiText, geminiTokens)

	// OpenAI: 推論モデルでは reasoning_effort に変換される（gpt-4o 等の非推論モデルでは無視される）
	openaiText, err, openaiTokens := client.GenText(wrapper.GenTextParams{
		Model:         models.ModelO4Mini,
		ThinkingLevel: wrapper.ThinkingLevelLow,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: question},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("OpenAI (low): %s\nTokens used: %d\n\n", openaiText, openaiTokens)

	// Anthropic: Claude 4.6 以降（Sonnet 5 / Opus 5 含む）では adaptive thinking + effort に変換される
	anthropicText, err, anthropicTokens := client.GenText(wrapper.GenTextParams{
		Model:         models.ModelClaudeSonnet5,
		ThinkingLevel: wrapper.ThinkingLevelMedium,
		Messages: []wrapper.Message{
			{Role: wrapper.RoleUser, Content: question},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Anthropic (medium): %s\nTokens used: %d\n\n", anthropicText, anthropicTokens)
}

func main() {
	// 従来の方法（個別のクライアント）を使用した例
	traditionalExample()

	fmt.Print("\n-----------------------------------\n\n")

	// 統合クライアントを使用した例
	unifiedClientExample()

	fmt.Print("\n-----------------------------------\n\n")

	// ThinkingLevel を指定した例
	thinkingLevelExample()
}
