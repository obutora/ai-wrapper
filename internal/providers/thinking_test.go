package providers

import (
	"errors"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/obutora/ai-wrapper/models"
	"github.com/openai/openai-go/shared"
	"google.golang.org/genai"
)

func TestGeminiThinkingConfig(t *testing.T) {
	i32 := func(v int32) *int32 { return &v }
	tests := []struct {
		model models.Model
		level models.ThinkingLevel
		want  *genai.ThinkingConfig
	}{
		// 未指定
		{"gemini-2.5-flash", models.ThinkingLevelDefault, nil},
		// Gemini 3 系: thinkingLevel
		{"gemini-3.8-flash", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelLow}}, // 3.8 Flash は MINIMAL 非対応
		{"gemini-3.8-flash", models.ThinkingLevelLow, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelLow}},
		{"gemini-3.8-flash", models.ThinkingLevelMedium, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMedium}},
		{"gemini-3.8-flash", models.ThinkingLevelHigh, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelHigh}},
		{"gemini-3.8-flash", models.ThinkingLevelMax, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelHigh}},
		{"gemini-3.5-flash-lite", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMinimal}}, // Lite は MINIMAL 対応
		{"gemini-3-pro-preview", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelLow}},      // Pro は MINIMAL 非対応
		{"gemini-3.1-pro-preview", models.ThinkingLevelMedium, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMedium}},
		// バージョン無しの別名（現在 3 系に解決される）: thinkingLevel
		{"gemini-flash-latest", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelLow}},
		{"gemini-flash-latest", models.ThinkingLevelHigh, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelHigh}},
		{"gemini-flash-lite-latest", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMinimal}},
		{"gemini-pro-latest", models.ThinkingLevelMedium, &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMedium}},
		// Gemini 2.5 系: thinkingBudget
		{"gemini-2.5-flash", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingBudget: i32(0)}},
		{"gemini-2.5-pro", models.ThinkingLevelMinimal, &genai.ThinkingConfig{ThinkingBudget: i32(128)}},
		{"gemini-2.5-flash", models.ThinkingLevelLow, &genai.ThinkingConfig{ThinkingBudget: i32(1024)}},
		{"gemini-2.5-flash", models.ThinkingLevelMedium, &genai.ThinkingConfig{ThinkingBudget: i32(8192)}},
		{"gemini-2.5-flash", models.ThinkingLevelHigh, &genai.ThinkingConfig{ThinkingBudget: i32(24576)}},
		{"gemini-2.5-flash-lite", models.ThinkingLevelMax, &genai.ThinkingConfig{ThinkingBudget: i32(24576)}},
		{"gemini-2.5-pro", models.ThinkingLevelMax, &genai.ThinkingConfig{ThinkingBudget: i32(32768)}},
		{"gemini-2.5-pro-preview-03-25", models.ThinkingLevelHigh, &genai.ThinkingConfig{ThinkingBudget: i32(24576)}},
		// 非対応モデル: 無視
		{"gemini-2.0-flash", models.ThinkingLevelHigh, nil},
		{"gemini-1.5-pro", models.ThinkingLevelHigh, nil},
	}
	for _, tt := range tests {
		got := geminiThinkingConfig(tt.model, tt.level)
		if (got == nil) != (tt.want == nil) {
			t.Errorf("%s/%s: got %+v, want %+v", tt.model, tt.level, got, tt.want)
			continue
		}
		if got == nil {
			continue
		}
		if got.ThinkingLevel != tt.want.ThinkingLevel {
			t.Errorf("%s/%s: ThinkingLevel got %q, want %q", tt.model, tt.level, got.ThinkingLevel, tt.want.ThinkingLevel)
		}
		if (got.ThinkingBudget == nil) != (tt.want.ThinkingBudget == nil) ||
			(got.ThinkingBudget != nil && *got.ThinkingBudget != *tt.want.ThinkingBudget) {
			t.Errorf("%s/%s: ThinkingBudget got %v, want %v", tt.model, tt.level, deref(got.ThinkingBudget), deref(tt.want.ThinkingBudget))
		}
		if got.IncludeThoughts {
			t.Errorf("%s/%s: IncludeThoughts must stay false", tt.model, tt.level)
		}
	}
}

func TestIsGeminiMinimalUnsupported(t *testing.T) {
	minimalConf := &genai.GenerateContentConfig{ThinkingConfig: &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMinimal}}
	lowConf := &genai.GenerateContentConfig{ThinkingConfig: &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelLow}}
	apiErr := errors.New(`Error 400, Message: Thinking level MINIMAL is not supported for this model. Please retry with other thinking level.`)

	if !isGeminiMinimalUnsupported(minimalConf, apiErr) {
		t.Error("expected true for MINIMAL + unsupported error")
	}
	if isGeminiMinimalUnsupported(lowConf, apiErr) {
		t.Error("expected false when MINIMAL was not sent")
	}
	if isGeminiMinimalUnsupported(minimalConf, errors.New("Error 429, rate limited")) {
		t.Error("expected false for unrelated error")
	}
	if isGeminiMinimalUnsupported(minimalConf, nil) || isGeminiMinimalUnsupported(nil, apiErr) {
		t.Error("expected false for nil inputs")
	}
}

func deref(p *int32) any {
	if p == nil {
		return nil
	}
	return *p
}

func TestOpenAIReasoningEffort(t *testing.T) {
	tests := []struct {
		model models.Model
		level models.ThinkingLevel
		want  shared.ReasoningEffort
	}{
		// 未指定
		{"o3-mini-2025-01-31k", models.ThinkingLevelDefault, ""},
		// 非推論モデル: 無視
		{"gpt-4o", models.ThinkingLevelHigh, ""},
		{"gpt-4.1-nano-2025-04-14", models.ThinkingLevelHigh, ""},
		{"gpt-3.5-turbo", models.ThinkingLevelLow, ""},
		// o 系
		{"o3-2025-04-16", models.ThinkingLevelHigh, shared.ReasoningEffortHigh},
		{"o4-mini-2025-04-16", models.ThinkingLevelMinimal, shared.ReasoningEffortLow},
		{"o4-mini-2025-04-16", models.ThinkingLevelMax, shared.ReasoningEffortHigh},
		{"o1", models.ThinkingLevelMedium, shared.ReasoningEffortMedium},
		// GPT-5 系
		{"gpt-5", models.ThinkingLevelMinimal, "minimal"},
		{"gpt-5-mini", models.ThinkingLevelLow, shared.ReasoningEffortLow},
		{"gpt-5", models.ThinkingLevelMax, shared.ReasoningEffortHigh},
		{"gpt-5.1", models.ThinkingLevelMinimal, "none"},
		{"gpt-5.1", models.ThinkingLevelMax, shared.ReasoningEffortHigh},
		{"gpt-5.2", models.ThinkingLevelMax, "xhigh"},
		{"gpt-5.2-pro", models.ThinkingLevelMedium, shared.ReasoningEffortMedium},
		{"gpt-5.1-codex-max", models.ThinkingLevelMax, "xhigh"},
		{"gpt-6", models.ThinkingLevelMax, "xhigh"},
		// codex 系
		{"codex-mini-latest", models.ThinkingLevelMedium, shared.ReasoningEffortMedium},
		{"codex-mini-latest", models.ThinkingLevelMinimal, shared.ReasoningEffortLow},
	}
	for _, tt := range tests {
		if got := openAIReasoningEffort(tt.model, tt.level); got != tt.want {
			t.Errorf("%s/%s: got %q, want %q", tt.model, tt.level, got, tt.want)
		}
	}
}

func TestAnthropicThinkingModeFor(t *testing.T) {
	tests := []struct {
		model string
		want  anthropicThinkingMode
	}{
		// 4.6 以降 / 5 系: adaptive
		{"claude-opus-4-6", anthropicThinkingAdaptive},
		{"claude-sonnet-4-6", anthropicThinkingAdaptive},
		{"claude-opus-4-7", anthropicThinkingAdaptive},
		{"claude-opus-4-8", anthropicThinkingAdaptive},
		{"claude-opus-5", anthropicThinkingAdaptive},
		{"claude-sonnet-5", anthropicThinkingAdaptive},
		{"claude-fable-5-1", anthropicThinkingAdaptive},
		// 3.7 〜 4.5: budget
		{"claude-haiku-4-5", anthropicThinkingBudget},
		{"claude-sonnet-4-5", anthropicThinkingBudget},
		{"claude-opus-4-5", anthropicThinkingBudget},
		{"claude-opus-4-1", anthropicThinkingBudget},
		{"claude-opus-4-1-20250805", anthropicThinkingBudget},
		{"claude-sonnet-4-20250514", anthropicThinkingBudget},
		{"claude-opus-4-20250514", anthropicThinkingBudget},
		{"claude-3-7-sonnet-latest", anthropicThinkingBudget},
		// 非対応
		{"claude-3-5-haiku-latest", anthropicThinkingUnsupported},
		{"claude-3-5-sonnet-20241022", anthropicThinkingUnsupported},
		{"claude-3-opus-latest", anthropicThinkingUnsupported},
		{"claude-3-haiku-20240307", anthropicThinkingUnsupported},
		{"my-custom-model", anthropicThinkingUnsupported},
	}
	for _, tt := range tests {
		if got := anthropicThinkingModeFor(tt.model); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.model, got, tt.want)
		}
	}
}

func TestApplyAnthropicThinking(t *testing.T) {
	newParams := func(model string, maxTokens int64) anthropic.MessageNewParams {
		return anthropic.MessageNewParams{Model: anthropic.Model(model), MaxTokens: maxTokens}
	}

	t.Run("default level leaves params untouched", func(t *testing.T) {
		p := newParams("claude-opus-4-6", 16000)
		if err := applyAnthropicThinking(&p, models.ThinkingLevelDefault); err != nil {
			t.Fatal(err)
		}
		if p.Thinking.OfAdaptive != nil || p.Thinking.OfEnabled != nil || p.OutputConfig.Effort != "" {
			t.Errorf("expected no thinking params, got %+v / %+v", p.Thinking, p.OutputConfig)
		}
	})

	t.Run("adaptive + effort on 4.6+", func(t *testing.T) {
		cases := map[models.ThinkingLevel]anthropic.OutputConfigEffort{
			models.ThinkingLevelMinimal: anthropic.OutputConfigEffortLow,
			models.ThinkingLevelLow:     anthropic.OutputConfigEffortLow,
			models.ThinkingLevelMedium:  anthropic.OutputConfigEffortMedium,
			models.ThinkingLevelHigh:    anthropic.OutputConfigEffortHigh,
			models.ThinkingLevelMax:     anthropic.OutputConfigEffortMax,
		}
		for level, wantEffort := range cases {
			p := newParams("claude-opus-4-7", 16000)
			if err := applyAnthropicThinking(&p, level); err != nil {
				t.Fatal(err)
			}
			if p.Thinking.OfAdaptive == nil || p.Thinking.OfEnabled != nil {
				t.Errorf("%s: expected adaptive thinking, got %+v", level, p.Thinking)
			}
			if p.OutputConfig.Effort != wantEffort {
				t.Errorf("%s: effort got %q, want %q", level, p.OutputConfig.Effort, wantEffort)
			}
		}
	})

	t.Run("enabled + budget_tokens on 3.7-4.5", func(t *testing.T) {
		cases := []struct {
			level      models.ThinkingLevel
			maxTokens  int64
			wantBudget int64
		}{
			{models.ThinkingLevelMinimal, 16000, 1024},
			{models.ThinkingLevelLow, 16000, 2048},
			{models.ThinkingLevelMedium, 16000, 8192},
			{models.ThinkingLevelHigh, 32000, 16384},
			{models.ThinkingLevelMax, 64000, 32768},
			// max_tokens 未満に丸める
			{models.ThinkingLevelHigh, 4000, 3999},
			{models.ThinkingLevelMinimal, 1025, 1024},
		}
		for _, c := range cases {
			p := newParams("claude-sonnet-4-5", c.maxTokens)
			if err := applyAnthropicThinking(&p, c.level); err != nil {
				t.Fatalf("%s/%d: %v", c.level, c.maxTokens, err)
			}
			if p.Thinking.OfEnabled == nil || p.Thinking.OfAdaptive != nil {
				t.Errorf("%s/%d: expected enabled thinking, got %+v", c.level, c.maxTokens, p.Thinking)
				continue
			}
			if p.Thinking.OfEnabled.BudgetTokens != c.wantBudget {
				t.Errorf("%s/%d: budget got %d, want %d", c.level, c.maxTokens, p.Thinking.OfEnabled.BudgetTokens, c.wantBudget)
			}
			if p.OutputConfig.Effort != "" {
				t.Errorf("%s/%d: effort must not be set on budget models, got %q", c.level, c.maxTokens, p.OutputConfig.Effort)
			}
		}
	})

	t.Run("budget models reject too small MaxToken", func(t *testing.T) {
		p := newParams("claude-haiku-4-5", 1000)
		err := applyAnthropicThinking(&p, models.ThinkingLevelLow)
		if !errors.Is(err, models.ErrInvalidThinkingConfig) {
			t.Errorf("expected ErrInvalidThinkingConfig, got %v", err)
		}
	})

	t.Run("unsupported models ignore the level", func(t *testing.T) {
		p := newParams("claude-3-5-haiku-latest", 1000)
		if err := applyAnthropicThinking(&p, models.ThinkingLevelHigh); err != nil {
			t.Fatal(err)
		}
		if p.Thinking.OfAdaptive != nil || p.Thinking.OfEnabled != nil || p.OutputConfig.Effort != "" {
			t.Errorf("expected no thinking params, got %+v / %+v", p.Thinking, p.OutputConfig)
		}
	})
}

// 各プロバイダの GenText が、API を呼ぶ前に不正な ThinkingLevel を弾くことを確認する。
func TestGenTextRejectsInvalidThinkingLevel(t *testing.T) {
	cfg := models.Config{MaxToken: 100}
	params := models.GenTextParams{Model: "dummy-model", Prompt: "hi", ThinkingLevel: "ultra"}

	clients := map[string]models.LLMWrapper{
		"openai":    NewOpenAIClient("dummy-key", cfg),
		"anthropic": NewAnthropicClient("dummy-key", cfg),
	}
	if g := NewGeminiClient("dummy-key", cfg); g != nil {
		clients["gemini"] = g
	}
	for name, c := range clients {
		_, err, tokens := c.GenText(params)
		if !errors.Is(err, models.ErrInvalidThinkingLevel) {
			t.Errorf("%s: expected ErrInvalidThinkingLevel, got %v", name, err)
		}
		if tokens != 0 {
			t.Errorf("%s: expected 0 tokens, got %d", name, tokens)
		}
	}
}
