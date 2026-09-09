package wrapper_test

import (
	"os"
	"strings"
	"testing"

	wrapper "github.com/obutora/ai-wrapper"
	"github.com/obutora/ai-wrapper/models"
)

// TestIntegrationThinkingLevel は、実 API に対して ThinkingLevel 付きのリクエストが受理されることを確認します。
// 課金が発生するため、AI_WRAPPER_INTEGRATION=1 が設定されている場合のみ実行され、
// 各プロバイダの API キーが未設定のケースはスキップされます。
//
//	AI_WRAPPER_INTEGRATION=1 GEMINI_API_KEY=... go test -run TestIntegrationThinkingLevel -v .
func TestIntegrationThinkingLevel(t *testing.T) {
	if os.Getenv("AI_WRAPPER_INTEGRATION") == "" {
		t.Skip("set AI_WRAPPER_INTEGRATION=1 to run integration tests (real API calls)")
	}

	cases := []struct {
		provider wrapper.Provider
		envKey   string
		model    wrapper.Model
		level    wrapper.ThinkingLevel
		note     string
	}{
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGeminiFlashLatest, wrapper.ThinkingLevelLow, "alias -> thinkingLevel"},
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGeminiFlashLatest, wrapper.ThinkingLevelMinimal, "alias minimal -> LOW"},
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGeminiFlashLiteLatest, wrapper.ThinkingLevelMinimal, "lite minimal -> MINIMAL"},
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGeminiProLatest, wrapper.ThinkingLevelMinimal, "pro minimal -> LOW"},
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGemini38Flash, wrapper.ThinkingLevelHigh, "thinkingLevel"},
		{wrapper.ProviderGemini, "GEMINI_API_KEY", models.ModelGemini25Flash, wrapper.ThinkingLevelLow, "thinkingBudget"},
		{wrapper.ProviderOpenAI, "OPENAI_API_KEY", models.ModelO4Mini, wrapper.ThinkingLevelLow, "reasoning_effort"},
		{wrapper.ProviderOpenAI, "OPENAI_API_KEY", models.Model4_1Nano, wrapper.ThinkingLevelHigh, "non-reasoning model: level ignored"},
		{wrapper.ProviderAnthropic, "ANTHROPIC_API_KEY", models.ModelClaudeHaiku45, wrapper.ThinkingLevelLow, "enabled + budget_tokens"},
		{wrapper.ProviderAnthropic, "ANTHROPIC_API_KEY", models.ModelClaudeSonnet5, wrapper.ThinkingLevelLow, "adaptive + effort"},
	}

	for _, c := range cases {
		t.Run(string(c.model)+"/"+string(c.level), func(t *testing.T) {
			key := os.Getenv(c.envKey)
			if key == "" {
				t.Skipf("%s not set", c.envKey)
			}
			client, err := wrapper.NewClient(c.provider, key, models.Config{MaxToken: 4000})
			if err != nil {
				t.Fatal(err)
			}
			text, err, tokens := client.GenText(wrapper.GenTextParams{
				Model:         c.model,
				ThinkingLevel: c.level,
				Prompt:        "Answer with a single word: what is the capital of Japan?",
			})
			if err != nil {
				t.Fatalf("[%s] %v", c.note, err)
			}
			if strings.TrimSpace(text) == "" {
				t.Fatalf("[%s] empty text returned", c.note)
			}
			if tokens <= 0 {
				t.Errorf("[%s] tokens = %d, want > 0", c.note, tokens)
			}
			t.Logf("[%s] text=%q tokens=%d", c.note, strings.TrimSpace(text), tokens)
		})
	}
}
