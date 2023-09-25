// Package llm_test is the test suite for the llm package.
// It uses a mock client to simulate the behavior of the OpenAI API client.
package llm_test

import (
	"context"
	"testing"

	"github.com/apsystole/log"
	"github.com/hlindberg/testutils"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"

	"github.com/cckalen/intellichunk/config"
	"github.com/cckalen/intellichunk/internal/llm"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Errorf("Load env Error: %s", err)

	}
}

// MockClient is a mock implementation of the OpenAI API client.
type MockClient struct {
	mock.Mock
}

// CreateChatCompletion is the mock implementation of the corresponding method in the OpenAI API client.
// This method is called when the ChatCompletion function in the llm package is tested.
// It returns whatever we program it to return.
func (m *MockClient) CreateChatCompletion(ctx context.Context, chatCompletionRequest openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, chatCompletionRequest)
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

func (m *MockClient) CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error) {
	args := m.Called(ctx, conv)
	return args.Get(0).(openai.EmbeddingResponse), args.Error(1)
}

// TestLanguageModelCompletion tests the ChatCompletion function in the llm package.
// It simulates a chat completion with a user message and checks that the function returns the expected response and error.
// It uses the MockClient to simulate the behavior of the OpenAI API client.
func TestLanguageModelCompletion(t *testing.T) {
	ctx := context.Background()
	userMessage := "Hello!"

	expectedResponse := openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "Hi there!",
				},
			},
		},
	}

	// mockClient is an instance of the MockClient.
	mockClient := new(MockClient)
	mockClient.On("CreateChatCompletion", ctx, mock.Anything).Return(expectedResponse, nil)

	// Create a new instance of OpenAI
	lm := llm.NewOpenAI()

	response, err := lm.ChatCompletion(ctx, userMessage)

	testutils.CheckEqual(nil, err, t)
	testutils.CheckEqual(expectedResponse.Choices[0].Message.Content, response, t)
	mockClient.AssertExpectations(t)
}

// TestGenerateEmbeddings tests the GenerateEmbeddings function in the llm package.
// It simulates a generate embeddings request with tokens and checks that the function returns the expected response and error.
// It uses the MockClient to simulate the behavior of the OpenAI API client.
func TestGenerateEmbeddings(t *testing.T) {
	ctx := context.Background()
	tokens := []int{1, 2, 3}
	model := openai.EmbeddingModel(openai.AdaEmbeddingV2)
	user := "system"

	expectedResponse := openai.EmbeddingResponse{
		// Define the expected response here
	}

	// mockClient is an instance of the MockClient.
	mockClient := new(MockClient)
	mockClient.On("CreateEmbeddings", ctx, mock.Anything).Return(expectedResponse, nil)

	// Create a new instance of OpenAI
	openaiClient := llm.NewOpenAI()

	response, err := openaiClient.GenerateEmbeddings(ctx, tokens, model, user)

	testutils.CheckEqual(nil, err, t)
	testutils.CheckEqual(expectedResponse, response, t)
	mockClient.AssertExpectations(t)
}

// TestGenerateMultipleEmbeddings tests the GenerateMultipleEmbeddings function in the llm package.
// It simulates a generate embeddings request with multiple sets of tokens and checks that the function returns the expected response and error.
// It uses the MockClient to simulate the behavior of the OpenAI API client.
func TestGenerateMultipleEmbeddings(t *testing.T) {
	ctx := context.Background()
	multipleTokens := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	model := openai.EmbeddingModel(openai.AdaEmbeddingV2)
	user := "system"

	expectedResponse := openai.EmbeddingResponse{
		// Define the expected response here
	}

	// mockClient is an instance of the MockClient.
	mockClient := new(MockClient)
	mockClient.On("CreateEmbeddings", ctx, mock.Anything).Return(expectedResponse, nil)

	// Create a new instance of OpenAI
	openaiClient := llm.NewOpenAI()

	response, err := openaiClient.GenerateMultipleEmbeddingsFromTokens(ctx, multipleTokens, model, user)

	testutils.CheckEqual(nil, err, t)
	testutils.CheckEqual(expectedResponse, response, t)
	mockClient.AssertExpectations(t)
}
