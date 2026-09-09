package models

import (
	"errors"
	"testing"
)

func TestThinkingLevelValidate(t *testing.T) {
	for _, l := range []ThinkingLevel{
		ThinkingLevelDefault, ThinkingLevelMinimal, ThinkingLevelLow,
		ThinkingLevelMedium, ThinkingLevelHigh, ThinkingLevelMax,
	} {
		if err := l.Validate(); err != nil {
			t.Errorf("%q: unexpected error %v", l, err)
		}
	}
	for _, l := range []ThinkingLevel{"ultra", "HIGH", "xhigh", "none"} {
		err := l.Validate()
		if !errors.Is(err, ErrInvalidThinkingLevel) {
			t.Errorf("%q: expected ErrInvalidThinkingLevel, got %v", l, err)
		}
	}
}

func TestAnthropicVersion(t *testing.T) {
	tests := []struct {
		model        Model
		major, minor int
		ok           bool
	}{
		{ModelClaudeOpus5, 5, 0, true},
		{ModelClaudeSonnet5, 5, 0, true},
		{ModelClaudeFable51, 5, 1, true},
		{ModelClaudeOpus48, 4, 8, true},
		{ModelClaudeOpus46, 4, 6, true},
		{ModelClaudeHaiku45, 4, 5, true},
		{ModelClaudeOpus41, 4, 1, true},
		{"claude-sonnet-4-20250514", 4, 0, true}, // 日付サフィックスは minor ではない
		{"claude-opus-4-1-20250805", 4, 1, true},
		{ModelClaude37Sonnet, 3, 7, true}, // ToAnthropicModel 経由で claude-3-7-sonnet-latest
		{ModelClaude3Haiku, 3, 5, true},   // ToAnthropicModel 経由で claude-3-5-haiku-latest
		{ModelClaude3Opus, 3, 0, true},    // ToAnthropicModel 経由で claude-3-opus-latest
		{ModelGPT4o, 0, 0, false},
		{ModelGeminiFlashLatest, 0, 0, false},
		{"my-custom-model", 0, 0, false},
	}
	for _, tt := range tests {
		major, minor, ok := tt.model.AnthropicVersion()
		if major != tt.major || minor != tt.minor || ok != tt.ok {
			t.Errorf("%s: got (%d, %d, %v), want (%d, %d, %v)", tt.model, major, minor, ok, tt.major, tt.minor, tt.ok)
		}
	}
}

func TestSupportsTemperature(t *testing.T) {
	tests := []struct {
		model Model
		want  bool
	}{
		{ModelClaudeOpus5, false},
		{ModelClaudeSonnet5, false},
		{ModelClaudeFable51, false},
		{ModelClaudeOpus48, false},
		{ModelClaudeOpus47, false},
		{ModelClaudeSonnet46, false},
		{ModelClaudeHaiku45, false},
		{ModelClaudeSonnet45, false},
		{ModelClaudeOpus41, true},
		{"claude-sonnet-4-20250514", true},
		{ModelClaude37Sonnet, true},
		{ModelGPT4o, true},
		{ModelGeminiFlashLatest, true},
	}
	for _, tt := range tests {
		if got := tt.model.SupportsTemperature(); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.model, got, tt.want)
		}
	}
}
