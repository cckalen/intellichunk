package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/cckalen/intellichunk/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

// OpenAI struct.
type OpenAI struct {
	client     API
	llmOptions *LLMOptions
}

// NewOpenAI creates a new OpenAI instance with an optional API key.
func NewOpenAI(opts ...LLMOption) *OpenAI {
	defaultAPIKey := os.Getenv("OPENAI_API_KEY")
	llmOptions := &LLMOptions{
		APIKey:    defaultAPIKey,
		ModelName: "gpt-3.5-turbo",
	}

	// Apply any specified options
	for _, opt := range opts {
		opt(llmOptions)
	}

	client := openai.NewClient(llmOptions.APIKey)

	return &OpenAI{
		client:     client,
		llmOptions: llmOptions,
	}
}

func ConvertToOpenAIFunctionDefinition(funcDefs []models.FunctionDefinition) []openai.FunctionDefinition {
	var openaiFuncDefs []openai.FunctionDefinition
	for _, fd := range funcDefs {
		// Check and assert the type of Parameters
		param, ok := fd.Parameters.(models.Definition)
		if !ok {
			// Handle the error: Parameters isn't of type models.Definition
			continue // or return an error, based on your needs
		}
		parameters := convertDefinitionToInterface(param)

		openaiFuncDef := openai.FunctionDefinition{
			Name:        fd.Name,
			Description: fd.Description,
			Parameters:  parameters,
		}
		openaiFuncDefs = append(openaiFuncDefs, openaiFuncDef)
	}
	return openaiFuncDefs
}

func convertDefinitionToInterface(def models.Definition) interface{} {
	// Convert the custom models.Definition to a map[string]interface{}
	// Assuming that the openai library expects it in this format
	result := make(map[string]interface{})
	result["type"] = def.Type
	if def.Properties != nil {
		properties := make(map[string]interface{})
		for key, propDef := range def.Properties {
			properties[key] = convertDefinitionToInterface(propDef)
		}
		result["properties"] = properties
	}
	if def.Items != nil {
		result["items"] = convertDefinitionToInterface(*def.Items)
	}
	if def.Required != nil {
		result["required"] = def.Required
	}

	return result
}

// ChatCompletion sends a chat completion request to the OpenAI API.
func (o *OpenAI) ChatCompletion(ctx context.Context, userMessage string) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: o.llmOptions.ModelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in chat completion response")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatCompletionWithInstructions sends a chat completion request with two roles.
// Even indexed strings are the user’s input;  and the odd index strings are the llm response
func (o *OpenAI) ChatCompletionWithInstructions(ctx context.Context, systemMessage, userMessage string, chatHistory []string) (string, error) {
	// Start with system message
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
	}

	// Add chat history
	// Append an additional message if chatHistory is not empty
	if len(chatHistory) > 0 {
		histInfo := "Chat History: "
		// Assuming the role for the extra message, you can change as needed
		extraMessageRole := openai.ChatMessageRoleSystem
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    extraMessageRole,
			Content: histInfo,
		})
	}
	// loop through chatHistory[], assigning the role based on whether the index is even or odd.
	for i, message := range chatHistory {
		role := openai.ChatMessageRoleUser
		if i%2 != 0 {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: message,
		})
	}

	// Add current user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    o.llmOptions.ModelName,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}
	// fmt.Println("=============================================================")
	// fmt.Println(messages)
	// fmt.Println("=============================================================")
	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) ChatCompletionFunctionsOptions(ctx context.Context, systemMessage string, funcDetails []models.FunctionDefinition, opts ...LLMOption) (string, error) {
	options := &LLMOptions{}
	for _, opt := range opts {
		opt(options)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
	}

	// Add chat history similar to ChatCompletionWithInstructions
	for i, message := range options.ChatHistory {
		role := openai.ChatMessageRoleUser
		if i%2 != 0 {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: message,
		})
	}

	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Temperature: options.Temperature,
			TopP:        options.TopP,
			Model:       o.llmOptions.ModelName,
			Functions:   ConvertToOpenAIFunctionDefinition(funcDetails),
			Messages:    messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in chat completion response")
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content == "" {
		if resp.Choices[0].Message.FunctionCall != nil {
			return resp.Choices[0].Message.FunctionCall.Arguments, nil
		}
		return "", fmt.Errorf("function call is nil")
	}
	return resp.Choices[0].Message.Content, nil

}

// CreateEmbeddings sends a create embeddings request to the OpenAI API.
func (o *OpenAI) GenerateEmbeddings(ctx context.Context, tokens []int, model openai.EmbeddingModel, user string) (openai.EmbeddingResponse, error) {

	response, err := o.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestTokens{
			Input: [][]int{tokens, tokens},
			Model: model,
			User:  user,
		},
	)

	if err != nil {
		return openai.EmbeddingResponse{}, fmt.Errorf("failed to create embeddings: %w", err)
	}

	return response, nil
}

// GenerateMultipleEmbeddings creates embeddings request to the OpenAI API with multiple sets of tokens, an embedding model, and a user.
func (o *OpenAI) GenerateMultipleEmbeddingsFromTokens(ctx context.Context, multipleTokens [][]int, model openai.EmbeddingModel, user string) (openai.EmbeddingResponse, error) {

	response, err := o.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestTokens{
			Input: multipleTokens,
			Model: model,
			User:  user,
		},
	)

	if err != nil {
		return openai.EmbeddingResponse{}, fmt.Errorf("failed to create multiple embeddings: %w", err)
	}

	return response, nil
}

// GenerateMultipleEmbeddingsFromText creates embeddings request to the OpenAI API with multiple sets of tokens, an embedding model, and a user.
func (o *OpenAI) GenerateMultipleEmbeddingsFromText(ctx context.Context, multipleText []string) ([][]float32, error) {

	response, err := o.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequestStrings{
			Input: multipleText,
			Model: openai.AdaEmbeddingV2,
			User:  "system",
		},
	)
	embedBatch := make([][]float32, 0, len(response.Data))

	if err != nil {
		return embedBatch, fmt.Errorf("failed to create multiple embeddings: %w", err)
	}

	for _, em := range response.Data {
		embedBatch = append(embedBatch, em.Embedding)
	}
	return embedBatch, nil
}
