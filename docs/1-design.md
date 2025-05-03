## go-wrapper

このプロジェクトは、Go言語用のLLMラッパーを提供します。

## 対応するLLMプロバイダ
- OpenAI
- Anthropic
- Gemini

 ## 設計
- `go-wrapper`は、LLMのAPIをラップするためのインターフェースを提供します。
- クライアントの生成時に、Providerを指定します。
- wrapper内でクライアントを保持し、各メソッドではそれを利用するようにしてください
- `client.genText()` メソッドを使用して、テキストを生成させることが出来ます
  - `genText()` メソッドは、プロバイダに応じて適切なAPIを呼び出します。
  - `genText()` メソッドは、引数に model, prompt, cacheEnabled, []messages を含む構造体を受け取ることができます。
    - ただし、Geminiの場合、cacheの管理が煩雑なので初期の実装では省略してかまいません。
  - `[]messages` は、Role, Contentを含む構造体のスライスです。
  - `genText()` メソッドは、返り値として、生成されたテキスト、エラー、利用トークンを返します

## 各プロバイダの実装例
このセクションでは、各プロバイダの実装例を示します。

### OpenAI
```go
package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func main() {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	param := openai.ChatCompletionNewParams{
	Messages: []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("What kind of houseplant is easy to take care of?"),
	},
	Seed:     openai.Int(1),
	Model:    openai.ChatModelGPT4o,
}

  completion, err := client.Chat.Completions.New(ctx, param)

  param.Messages = append(param.Messages, completion.Choices[0].Message.ToParam())
  param.Messages = append(param.Messages, openai.UserMessage("How big are those?"))

  // continue the conversation
  completion, err = client.Chat.Completions.New(ctx, param)
  	if err != nil {
  		panic(err.Error())
  	}
  	println(chatCompletion.Choices[0].Message.Content)
}
```


### Gemini
```go
package main

import (
  "context"
  "fmt"
  "os"
  "google.golang.org/genai"
)

func main() {

  ctx := context.Background()
  client, _ := genai.NewClient(ctx, &genai.ClientConfig{
      APIKey:  os.Getenv("GEMINI_API_KEY"),
      Backend: genai.BackendGeminiAPI,
  })

  history := []*genai.Content{
      genai.NewContentFromText("Hi nice to meet you! I have 2 dogs in my house.", genai.RoleUser),
      genai.NewContentFromText("Great to meet you. What would you like to know?", genai.RoleModel),
  }

  chat, _ := client.Chats.Create(ctx, "gemini-2.0-flash", nil, history)
  res, _ := chat.SendMessage(ctx, genai.Part{Text: "How many paws are in my house?"})

  if len(res.Candidates) > 0 {
      fmt.Println(res.Candidates[0].Content.Parts[0].Text)
  }
}

```
### Anthropic
```go
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: "What is a quaternion?"},
			}},
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}
```
