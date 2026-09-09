package types

import (
	"github.com/obutora/ai-wrapper/models"
)

// Model は、LLMモデルの種類を表す型です。
type Model = models.Model

// 利用可能なモデルの定数
const (
	// OpenAIモデル
	ModelGPT4o      = models.ModelGPT4o
	ModelGPT4       = models.ModelGPT4
	ModelGPT35Turbo = models.ModelGPT35Turbo
	ModelO3Mini     = models.ModelO3Mini
	ModelO4Mini     = models.ModelO4Mini
	Model4_1Nano    = models.Model4_1Nano
	ModelO3         = models.ModelO3

	// Anthropicモデル
	ModelClaudeOpus5    = models.ModelClaudeOpus5
	ModelClaudeSonnet5  = models.ModelClaudeSonnet5
	ModelClaudeHaiku45  = models.ModelClaudeHaiku45
	ModelClaudeFable51  = models.ModelClaudeFable51
	ModelClaudeFable5   = models.ModelClaudeFable5
	ModelClaudeOpus48   = models.ModelClaudeOpus48
	ModelClaudeOpus47   = models.ModelClaudeOpus47
	ModelClaudeOpus46   = models.ModelClaudeOpus46
	ModelClaudeSonnet46 = models.ModelClaudeSonnet46
	ModelClaudeOpus45   = models.ModelClaudeOpus45
	ModelClaudeSonnet45 = models.ModelClaudeSonnet45
	ModelClaudeOpus41   = models.ModelClaudeOpus41   // Deprecated
	ModelClaude3Opus    = models.ModelClaude3Opus    // Deprecated
	ModelClaude37Sonnet = models.ModelClaude37Sonnet // Deprecated
	ModelClaude3Haiku   = models.ModelClaude3Haiku   // Deprecated

	// Geminiモデル
	ModelGeminiFlashLatest     = models.ModelGeminiFlashLatest
	ModelGeminiFlashLiteLatest = models.ModelGeminiFlashLiteLatest
	ModelGeminiProLatest       = models.ModelGeminiProLatest
	ModelGemini38Flash         = models.ModelGemini38Flash
	ModelGemini35FlashLite     = models.ModelGemini35FlashLite
	ModelGemini31ProPreview    = models.ModelGemini31ProPreview
	ModelGemini25Flash         = models.ModelGemini25Flash
	ModelGemini25FlashLite     = models.ModelGemini25FlashLite
	ModelGemini20Flash         = models.ModelGemini20Flash        // Deprecated
	ModelGemini20Pro           = models.ModelGemini20Pro          // Deprecated
	ModelGemini25FlashPreview  = models.ModelGemini25FlashPreview // Deprecated
	ModelGemini25ProPreview    = models.ModelGemini25ProPreview   // Deprecated
	ModelGemini25Pro           = models.ModelGemini25Pro          // Deprecated
)
