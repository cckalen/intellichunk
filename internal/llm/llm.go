package llm

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

// API interface mainly to decouple from the openai package and easily mock the openai package in tests.
type API interface {
	CreateChatCompletion(ctx context.Context, chatCompletionRequest openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequestConverter) (res openai.EmbeddingResponse, err error)
}

// LanguageModel interface.
type LanguageModel interface {
	ChatCompletion(ctx context.Context, userMessage string, chatHistory []string) (string, error)
	GenerateEmbeddings(ctx context.Context, tokens []int, model openai.EmbeddingModel, user string) (openai.EmbeddingResponse, error)
	GenerateMultipleEmbeddingsFromTokens(ctx context.Context, multipleTokens [][]int, model openai.EmbeddingModel, user string) (openai.EmbeddingResponse, error)
	GenerateMultipleEmbeddingsFromText(ctx context.Context, multipleText []string) ([][]float32, error)
}

type LLMOptions struct {
	APIKey      string
	ModelName   string
	ChatHistory []string
	Temperature float32
	TopP        float32
}

type LLMOption func(*LLMOptions)

func WithAPIKey(APIKey string) LLMOption {
	return func(o *LLMOptions) {
		o.APIKey = APIKey
	}
}

// WithModelName sets the custom model name for the model.
// Suggested models from openai: GPT432K0613 = "gpt-4-32k-0613" / GPT4 = "gpt-4" / GPT3Dot5Turbo16K = "gpt-3.5-turbo-16k" / GPT3Dot5Turbo = "gpt-3.5-turbo" / GPT3Dot5TurboInstruct = "gpt-3.5-turbo-instruct" /
func WithModelName(modelName string) LLMOption {
	return func(o *LLMOptions) {
		o.ModelName = modelName
	}
}

func WithChatHistory(chatHistory []string) LLMOption {
	return func(o *LLMOptions) {
		o.ChatHistory = chatHistory
	}
}

func WithTemperature(temperature float32) LLMOption {
	return func(o *LLMOptions) {
		o.Temperature = temperature
	}
}

func WithTopP(topP float32) LLMOption {
	return func(o *LLMOptions) {
		o.TopP = topP
	}
}
