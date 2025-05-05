wrapper.goのインターフェースをもっと使いやすいものに洗練させたいです。

現状、README.mdに書いてあるように、各モデルプロバイダごとにクライアントを生成する必要があります。

// OpenAI client
openaiClient, err := wrapper.NewClient(wrapper.ProviderOpenAI, os.Getenv("OPENAI_API_KEY"))

// Anthropic client
anthropicClient, err := wrapper.NewClient(wrapper.ProviderAnthropic, os.Getenv("ANTHROPIC_API_KEY"))

// Gemini client
geminiClient, err := wrapper.NewClient(wrapper.ProviderGemini, os.Getenv("GEMINI_API_KEY"))


しかし、１つのクライアントから、モデル名のみを指定することで内部で自動的に各プロバイダーを切り替えられるようにしたいです。 
それぞれの実装の改善とReadme.mdの改善を要望します。 
