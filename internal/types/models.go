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
	ModelClaude3Opus    = models.ModelClaude3Opus
	ModelClaude37Sonnet = models.ModelClaude37Sonnet
	ModelClaude3Haiku   = models.ModelClaude3Haiku

	// Geminiモデル
	ModelGemini20Flash        = models.ModelGemini20Flash
	ModelGemini20Pro          = models.ModelGemini20Pro
	ModelGemini25FlashPreview = models.ModelGemini25FlashPreview
	ModelGemini25ProPreview   = models.ModelGemini25ProPreview
	ModelGemini25Pro          = models.ModelGemini25Pro
)
