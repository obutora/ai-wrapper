package providers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/obutora/ai-wrapper/models"
	"github.com/openai/openai-go/shared"
	"google.golang.org/genai"
)

// このファイルは、プロバイダ横断の ThinkingLevel を各プロバイダのネイティブな
// 推論パラメータへ変換するロジックをまとめたものです。
//
//	Level    | Gemini 3+ / *-latest | Gemini 2.5 (budget) | OpenAI (reasoning_effort) | Anthropic 4.6+ (effort) | Anthropic 3.7〜4.5 (budget)
//	---------+----------------------+---------------------+---------------------------+-------------------------+---------------------------
//	minimal  | LOW (Lite: MINIMAL)*1| 0 (Pro: 128)        | minimal / none(5.1+) (*2) | low                     | 1024
//	low      | LOW                  | 1024                | low                       | low                     | 2048
//	medium   | MEDIUM               | 8192                | medium                    | medium                  | 8192
//	high     | HIGH                 | 24576               | high                      | high                    | 16384
//	max      | HIGH                 | 24576 (Pro: 32768)  | high / xhigh(5.2+)        | max                     | 32768
//
//	*1 MINIMAL のサポートはモデルごとに異なる（2026-09 時点: 3.5 Flash-Lite は対応、3.8 Flash / 3.1 Pro は 400）。
//	   Lite 系のみ MINIMAL を送り、非対応エラーが返った場合は GenText 側で LOW に再試行します。
//	*2 o 系 / codex 系は minimal 非対応のため low に丸めます。
//
// gemini-flash-latest 等のバージョン無し別名は現在 Gemini 3 系に解決されるため、thinkingLevel を使います。
// 推論をサポートしないモデル（Gemini 2.0 以前、GPT-4 系、Claude 3.x など）では
// パラメータを付与しません（指定は無視されます）。

// parseDottedVersion は、正規表現でモデル名から major[.minor] を取り出します。
// 正規表現は 1 番目に major、2 番目（任意）に minor をキャプチャする必要があります。
func parseDottedVersion(re *regexp.Regexp, model string) (major, minor int, ok bool) {
	m := re.FindStringSubmatch(model)
	if m == nil {
		return 0, 0, false
	}
	major, _ = strconv.Atoi(m[1])
	if len(m) > 2 && m[2] != "" {
		minor, _ = strconv.Atoi(m[2])
	}
	return major, minor, true
}

// ---------------------------------------------------------------------------
// Gemini
// ---------------------------------------------------------------------------

var geminiVersionRe = regexp.MustCompile(`^gemini-(\d+)(?:\.(\d+))?`)

// geminiThinkingConfig は、ThinkingLevel を Gemini の ThinkingConfig に変換します。
//   - Gemini 3 以降、および *-latest 等のバージョン無し別名（現在 3 系に解決される）: thinkingLevel
//   - Gemini 2.5: thinkingBudget（トークン数）
//   - Gemini 2.0 以前: thinking 非対応のため nil（指定は無視）
func geminiThinkingConfig(model models.Model, level models.ThinkingLevel) *genai.ThinkingConfig {
	if level == models.ThinkingLevelDefault {
		return nil
	}
	name := strings.ToLower(string(model))
	major, minor, ok := parseDottedVersion(geminiVersionRe, name)
	switch {
	case ok && (major < 2 || (major == 2 && minor < 5)):
		return nil
	case ok && major == 2:
		budget := geminiThinkingBudget(name, level)
		return &genai.ThinkingConfig{ThinkingBudget: &budget}
	default:
		return &genai.ThinkingConfig{ThinkingLevel: geminiThinkingLevel(name, level)}
	}
}

func geminiThinkingLevel(name string, level models.ThinkingLevel) genai.ThinkingLevel {
	switch level {
	case models.ThinkingLevelMinimal:
		// MINIMAL のサポートはモデルごとに異なる（2026-09 時点: 3.5 Flash-Lite は対応、3.8 Flash / 3.1 Pro は非対応）。
		// Lite 系のみ MINIMAL を送り、それ以外は LOW に丸める。
		// 非対応モデルに送ってしまった場合は isGeminiMinimalUnsupported で検知し、GenText 側で LOW に再試行する。
		if strings.Contains(name, "lite") {
			return genai.ThinkingLevelMinimal
		}
		return genai.ThinkingLevelLow
	case models.ThinkingLevelLow:
		return genai.ThinkingLevelLow
	case models.ThinkingLevelMedium:
		return genai.ThinkingLevelMedium
	default: // high, max
		return genai.ThinkingLevelHigh
	}
}

// isGeminiMinimalUnsupported は、thinkingLevel=MINIMAL を非対応モデルに送った際の
// 400 エラー（"Thinking level MINIMAL is not supported for this model"）かを判定します。
func isGeminiMinimalUnsupported(conf *genai.GenerateContentConfig, err error) bool {
	if err == nil || conf == nil || conf.ThinkingConfig == nil ||
		conf.ThinkingConfig.ThinkingLevel != genai.ThinkingLevelMinimal {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "minimal") && strings.Contains(msg, "not supported")
}

// geminiThinkingBudget は、Gemini 2.5 系向けの thinkingBudget を返します。
// 有効範囲: Pro 128〜32768 / Flash 0〜24576 / Flash-Lite 512〜24576（0 で無効化）
func geminiThinkingBudget(name string, level models.ThinkingLevel) int32 {
	isPro := strings.Contains(name, "-pro")
	switch level {
	case models.ThinkingLevelMinimal:
		if isPro {
			// Pro は思考を無効化できないため下限値
			return 128
		}
		return 0
	case models.ThinkingLevelLow:
		return 1024
	case models.ThinkingLevelMedium:
		return 8192
	case models.ThinkingLevelHigh:
		return 24576
	default: // max
		if isPro {
			return 32768
		}
		return 24576
	}
}

// ---------------------------------------------------------------------------
// OpenAI
// ---------------------------------------------------------------------------

var (
	openAIGPTVersionRe = regexp.MustCompile(`^gpt-(\d+)(?:\.(\d+))?`)
	openAIOSeriesRe    = regexp.MustCompile(`^o\d+(?:-|$)`)
)

// openAIReasoningEffort は、ThinkingLevel を OpenAI の reasoning_effort に変換します。
// 推論モデル（o 系 / GPT-5 系 / codex 系）以外では reasoning_effort が拒否されるため、
// 空文字を返します（指定は無視）。
func openAIReasoningEffort(model models.Model, level models.ThinkingLevel) shared.ReasoningEffort {
	if level == models.ThinkingLevelDefault {
		return ""
	}
	name := strings.ToLower(string(model))
	isOSeries := openAIOSeriesRe.MatchString(name)
	isCodex := strings.HasPrefix(name, "codex-")
	major, minor, isGPT := parseDottedVersion(openAIGPTVersionRe, name)
	isGPT5Plus := isGPT && major >= 5

	if !isOSeries && !isCodex && !isGPT5Plus {
		return ""
	}

	gpt51Plus := isGPT && (major > 5 || (major == 5 && minor >= 1))
	gpt52Plus := isGPT && (major > 5 || (major == 5 && minor >= 2))

	switch level {
	case models.ThinkingLevelMinimal:
		switch {
		case gpt51Plus:
			// GPT-5.1 以降は minimal が none に置き換えられた
			return "none"
		case isGPT5Plus:
			return "minimal"
		default:
			// o 系 / codex 系は minimal 非対応
			return shared.ReasoningEffortLow
		}
	case models.ThinkingLevelLow:
		return shared.ReasoningEffortLow
	case models.ThinkingLevelMedium:
		return shared.ReasoningEffortMedium
	case models.ThinkingLevelHigh:
		return shared.ReasoningEffortHigh
	default: // max
		if gpt52Plus || strings.Contains(name, "codex-max") {
			return "xhigh"
		}
		return shared.ReasoningEffortHigh
	}
}

// ---------------------------------------------------------------------------
// Anthropic
// ---------------------------------------------------------------------------

type anthropicThinkingMode int

const (
	// Claude 3.x など thinking 非対応
	anthropicThinkingUnsupported anthropicThinkingMode = iota
	// Claude 3.7 〜 4.5: thinking.type=enabled + budget_tokens
	anthropicThinkingBudget
	// Claude 4.6 以降: thinking.type=adaptive + output_config.effort
	anthropicThinkingAdaptive
)

// anthropicThinkingModeFor は、モデル名から thinking の指定方式を判定します。
func anthropicThinkingModeFor(model string) anthropicThinkingMode {
	major, minor, ok := models.Model(model).AnthropicVersion()
	if !ok {
		return anthropicThinkingUnsupported
	}

	switch {
	case major >= 5, major == 4 && minor >= 6:
		return anthropicThinkingAdaptive
	case major == 4, major == 3 && minor == 7:
		return anthropicThinkingBudget
	default:
		return anthropicThinkingUnsupported
	}
}

// anthropicMinThinkingBudget は、budget_tokens の下限値です。
const anthropicMinThinkingBudget int64 = 1024

// applyAnthropicThinking は、ThinkingLevel を Anthropic のリクエストパラメータに反映します。
//   - Claude 4.6 以降: thinking.type=adaptive + output_config.effort
//   - Claude 3.7 〜 4.5: thinking.type=enabled + budget_tokens（max_tokens 未満に収める）
//   - それ以外: 何もしない（指定は無視）
func applyAnthropicThinking(p *anthropic.MessageNewParams, level models.ThinkingLevel) error {
	if level == models.ThinkingLevelDefault {
		return nil
	}
	switch anthropicThinkingModeFor(string(p.Model)) {
	case anthropicThinkingAdaptive:
		p.Thinking = anthropic.ThinkingConfigParamUnion{OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{}}
		p.OutputConfig = anthropic.OutputConfigParam{Effort: anthropicEffort(level)}
	case anthropicThinkingBudget:
		budget, err := anthropicThinkingBudgetFor(level, p.MaxTokens)
		if err != nil {
			return err
		}
		p.Thinking = anthropic.ThinkingConfigParamOfEnabled(budget)
	}
	return nil
}

func anthropicEffort(level models.ThinkingLevel) anthropic.OutputConfigEffort {
	switch level {
	case models.ThinkingLevelMinimal, models.ThinkingLevelLow:
		return anthropic.OutputConfigEffortLow
	case models.ThinkingLevelMedium:
		return anthropic.OutputConfigEffortMedium
	case models.ThinkingLevelHigh:
		return anthropic.OutputConfigEffortHigh
	default: // max
		return anthropic.OutputConfigEffortMax
	}
}

// anthropicThinkingBudgetFor は、budget_tokens を返します。
// budget_tokens は 1024 以上かつ max_tokens 未満である必要があるため、
// max_tokens が小さい場合はそれに収まるよう丸め、収まらない場合はエラーを返します。
func anthropicThinkingBudgetFor(level models.ThinkingLevel, maxTokens int64) (int64, error) {
	var budget int64
	switch level {
	case models.ThinkingLevelMinimal:
		budget = 1024
	case models.ThinkingLevelLow:
		budget = 2048
	case models.ThinkingLevelMedium:
		budget = 8192
	case models.ThinkingLevelHigh:
		budget = 16384
	default: // max
		budget = 32768
	}
	if maxTokens <= anthropicMinThinkingBudget {
		return 0, fmt.Errorf("%w: MaxToken must be greater than %d to enable thinking on this model (got %d)",
			models.ErrInvalidThinkingConfig, anthropicMinThinkingBudget, maxTokens)
	}
	if budget >= maxTokens {
		budget = maxTokens - 1
	}
	return budget, nil
}
