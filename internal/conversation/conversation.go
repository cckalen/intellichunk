package conversation

import (
	"context"
	"strings"

	"github.com/apsystole/log"
	"github.com/cckalen/intellichunk/internal/llm"
	"github.com/cckalen/intellichunk/internal/models"
	"github.com/cckalen/intellichunk/internal/templateprompt"
	"github.com/cckalen/intellichunk/internal/vectorstore"
)

// GetRelevantContent retrieves relevant content based on a given classname and question.
// It uses a vectorstore to perform similarity search and returns the relevant objects/documents in the form of content and reference URLs.
//
// Parameters:
// - classname: The classname of the objects to search in Weaviate.
// - question: The question or concept to search for similarities.
//
// Returns:
// - content: The merged content of the relevant contents into a string.
// - refurls: The reference URLs associated with the relevant objects as array.
//
// Note: To use custom parameters, look into the vectorstore package or uncomment the 'vectorstore.WithHost("custom-host")' option.
func GetRelevantContent(classname string, question string) (string, []string, error) {
	store := vectorstore.NewWeaviateStore(
	// vectorstore.WithHost("custom-host").
	// You can add or remove options as needed.
	)

	className := classname
	input := question
	graphFieldNames := []string{"content", "title"}
	withLimit := 3

	// Define the arrays for contents and refurls.
	var contents []string
	var refurls []string

	result, err := store.SimilaritySearch(className, input, graphFieldNames, withLimit)
	if err != nil {
		log.Errorf("Failed to perform similarity search: %v", err)
		return "", refurls, err
	}

	// slice of map[string]interface{}.
	for _, object := range result {
		for key, value := range object {
			if key == "content" {
				contents = append(contents, value.(string))
			} else if key == "title" {
				refurls = append(refurls, value.(string))
			}
		}
	}

	// Merge contents.
	content := strings.Join(contents, "/n/n")
	return content, refurls, nil
}

// ClassConversation main function dealing with incoming api calls.
func ClassConversation(convoReq models.ConversationRequest) (convoResp models.ConversationResponse, err error) {
	var content string

	// Creating a new TemplateRenderer using our prompt.
	tr, err := templateprompt.NewTemplateRenderer(templateprompt.HelperAgentPrompt)
	if err != nil {
		return
	}

	// Get Relevant Content from this Classs vector database.
	content, convoResp.Answer.Sources, err = GetRelevantContent(convoReq.ClassID, convoReq.Query)
	if err != nil {
		return
	}

	// Defining parameters for the prompt.
	params := map[string]string{
		"Details":   "some dynamic instructive text text",
		"SSContent": content,
	}

	// Rendering the template with the provided parameters.
	promptSystem, err := tr.Render(params)
	if err != nil {
		log.Error(err)
		return
	}

	languageModel := llm.NewOpenAI()

	// Use the language model to generate a chat completion.
	convoResp.Answer.Answer, err = languageModel.ChatCompletionWithInstructions(context.Background(), promptSystem, convoReq.Query, convoReq.ChatHistory)
	if err != nil {
		log.Errorf("Failed to generate chat completion: %v", err)
		return
	}

	convoResp.ClassID = convoReq.ClassID
	convoResp.ConversationID = convoReq.ConversationID
	convoResp.Query = convoReq.Query

	return
}
