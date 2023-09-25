package intellichunk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/apsystole/log"
	"github.com/cckalen/intellichunk/internal/llm"
	"github.com/cckalen/intellichunk/internal/models"
	"github.com/cckalen/intellichunk/internal/util"
	"github.com/cckalen/intellichunk/internal/vectorstore"
)

// SplitTextIntoContainerNodes function takes a long text as input and splits it into smaller sections/nodes,
// each containing around 300 characters. It also adds relevant metadata to each section,
// including 5 keywords and 2 specific questions that the section can answer.
// The function returns the resulting text in a stringfied JSON format representing the nodes.
func SplitTextIntoContainerNodes(longText string) (chunkedResp string, err error) {
	// Define the JSON schema for Sections
	nodesSchema := &models.Definition{
		Type: models.Object,
		Properties: map[string]models.Definition{
			"content": {
				Type:        models.String,
				Description: "The text content of the section.",
			},
			"keywords": {
				Type:        models.Array,
				Items:       &models.Definition{Type: models.String},
				Description: "The array of 5 keywords related to this section",
			},
			"questions": {
				Type:        models.Array,
				Items:       &models.Definition{Type: models.String},
				Description: "The array of 2 questions this section can answer",
			},
			"sectionNumber": {
				Type:        models.Integer,
				Description: "Section number so if there are 3 sections and this is the last one, it is 3",
			},
		},
		Required: []string{"content", "keywords", "questions", "sectionNumber"},
	}

	funcDef := []models.FunctionDefinition{
		{
			Name:        "split_into_sections",
			Description: "Splits input into smaller sections/nodes",
			Parameters: models.Definition{
				Type: models.Object,
				Properties: map[string]models.Definition{
					"title": {
						Type:        models.String,
						Description: "Title of the whole input",
					},
					"summary": {
						Type:        models.String,
						Description: "Summary of the whole input.",
					},
					"abstract_description": {
						Type:        models.String,
						Description: "Very abstract description of the input",
					},
					"nodes": {
						Type:        models.Array,
						Items:       nodesSchema,
						Description: "Smaller sections/nodes that compose the whole input",
					},
				},
				Required: []string{"title", "summary", "abstract_description", "nodes"},
			},
		},
	}

	var container models.DataContainer
	// Create a new OpenAI language model instance.
	languageModel := llm.NewOpenAI()
	llmOptions := []llm.LLMOption{
		llm.WithTemperature(0.3),
	}

	jsonString := `{"title": "Title", "summary": "Summary", "abstract_description": "Description", "nodes": [{"content": "Content1", "keywords": ["k1", "k2", "k3", "k4", "k5"], "questions": ["Q1?", "Q2?"], "sectionNumber": 1}, {"content": "Content2", "keywords": ["k6", "k7", "k8", "k9", "k10"], "questions": ["Q3?", "Q4?"], "sectionNumber": 2}]}`
	promptToSplit := "Create a single paragraph summary, an abstract description that describes the input and a title considering unique entities found in the following input.  Also, Split the input into smaller sections called nodes, each around 200 words(this is important), for each section add 3 relevant keywords from that section/chunk and 2 questions this section can provide specific answers to which are unlikely to be found elsewhere. Example Output: " + jsonString + "\n\n Input:" + longText

	for retries := 0; retries < 3; retries++ {
		chunkedResp, err = languageModel.ChatCompletionFunctionsOptions(context.Background(), promptToSplit, funcDef, llmOptions...)
		if err != nil {
			log.Errorf("error : %s", err)
			continue // Retry if there's an error
		}

		// Check if valid input before returning
		err = json.Unmarshal([]byte(chunkedResp), &container)
		if err != nil {
			log.Println("Problematic response from LLM, retrying...", err)
			continue // Retry if there's an unmarshalling error
		}

		break
	}

	if err != nil {
		// tried 3 times and still haven't got a valid response
		log.Errorf("Failed to get a valid response after 3 attempts: %s", err)
		return chunkedResp, err
	}

	return chunkedResp, nil
}

// GenerateContainerNodes function processes the string JSON response from SplitTextIntoContainerNodes and generates container nodes
// with additional embeddings. It returns a slice of models.ContainerNodeVector, which contains information about each node
// along with its associated embeddings.
func GenerateContainerNodes(chunkedResp, reftitle, refUrl string) (nodes []models.ContainerNodeVector, err error) {

	var container models.DataContainer
	err = json.Unmarshal([]byte(chunkedResp), &container)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return nodes, err
	}

	embedTextSlice := make([]string, 0, len(container.Nodes))
	// Range through each node and build our string output
	for _, node := range container.Nodes {
		var embedBuilder strings.Builder
		// Start building the embedding text
		questions := strings.Join(node.Questions, ", ")
		embedBuilder.WriteString(fmt.Sprintf("Q: %s\n", questions))
		// Add node content
		embedBuilder.WriteString(fmt.Sprintf("A: %s\n\n", node.Content))
		embedBuilder.WriteString(fmt.Sprintf("Taken from '%s': \n%s", container.Title, container.AbstactSum))
		// Print the formatted string
		embedTextSlice = append(embedTextSlice, embedBuilder.String())
		//fmt.Println(embedBuilder.String())
		// Add space between nodes
		//fmt.Print("\n\n")
	}

	// Create a new OpenAI language model instance.
	languageModel := llm.NewOpenAI()

	//embedBatch holds [][]float32 of embeddings
	embedBatch, err := languageModel.GenerateMultipleEmbeddingsFromText(context.Background(), embedTextSlice)
	if err != nil {
		log.Println("Error GenerateMultipleEmbeddingsFromText  :", err)
		return nodes, err
	}

	// Assigning each node with relavant embeddings returned
	for i, node := range container.Nodes {
		var cNV models.ContainerNodeVector
		cNV.Title = container.Title
		cNV.Summary = container.Summary
		cNV.AbstactSum = container.AbstactSum
		cNV.Content = node.Content
		cNV.Keywords = node.Keywords
		cNV.Questions = node.Questions
		cNV.NodeNumber = node.NodeNumber
		cNV.RefTitle = reftitle
		cNV.ReferenceURL = refUrl
		cNV.Embedding = embedBatch[i]
		nodes = append(nodes, cNV)
	}

	return nodes, nil
}

// Add function is responsible for adding the generated container nodes to the vectorstore.
// It takes a class name and long text as input, splits the text into container nodes using SplitTextIntoContainerNodes,
// generates embeddings for the nodes using GenerateContainerNodes, and finally adds the resulting objects to the vectorstore.
// The function returns a slice of object IDs returned from the vectorstore.
func Add(className, longText string) (objIDs []string, err error) {

	nodesInString, err := SplitTextIntoContainerNodes(longText)
	if err != nil {
		log.Println("Error Add SplitTextIntoContainerNodes  :", err)
		return objIDs, err
	}

	dataContainerNodes, err := GenerateContainerNodes(nodesInString, "", "")
	if err != nil {
		log.Println("Error Add GenerateContainerNodes  :", err)
		return objIDs, err
	}

	store := vectorstore.NewWeaviateStore()
	objIDs, err = store.AddNodeObjects(className, dataContainerNodes)
	if err != nil {
		log.Println("Error Add AddObjects  :", err)
		return objIDs, err
	}
	return objIDs, nil
}

// AddFromFolder function is Similar to Add function but works from a given folder,
// reads the .txt files from the folder and looks for a particular format as following then process them similar to Add()
//
// Title: Another article title
// RefURL:https://....
// Content: Long content
func AddFromFolder(className, folderPath string, saveToFile bool) (objIDs []string, err error) {
	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		util.Red("Error getting absolute path: %v\n", err)
		return nil, err
	}

	err = filepath.WalkDir(absFolderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			util.Red("Error walking the path: %v\n", err)
			return err
		}

		if !d.IsDir() && filepath.Ext(d.Name()) == ".txt" {
			util.Yellow(":: Analysing... > %s\n", d.Name())

			file, err := os.ReadFile(path)
			if err != nil {
				util.Red("Error walking the directory: %v\n", err)
				return nil // skip this file and move to the next one
			}

			contentWithoutBOM := strings.TrimPrefix(string(file), "\ufeff")
			articles := strings.Split(contentWithoutBOM, "\nTitle: ")
			firstArticle := true
			for _, article := range articles {
				if len(article) == 0 {
					continue
				}

				// Adding "Title: " back to the start of the article string
				if !firstArticle {
					article = "Title: " + article
				}
				firstArticle = false

				refTitleStart := strings.Index(article, "Title:") + len("Title:")
				refURLStart := strings.Index(article, "\nRefURL:") + len("\nRefURL:")
				longTextStart := strings.Index(article, "\nContent:") + len("\nContent:")

				if refTitleStart < 0 || refURLStart < 0 || longTextStart < 0 {
					continue
				}

				refTitle := strings.TrimSpace(article[refTitleStart : refURLStart-len("\nRefURL:")])
				util.Green("::::: Processing Article... > %s\n", refTitle)

				refURL := strings.TrimSpace(article[refURLStart : longTextStart-len("\nContent:")])
				longText := strings.TrimSpace(article[longTextStart:])

				nodesInString, err := SplitTextIntoContainerNodes(longText)
				if err != nil {
					util.Red("--------> Skipping this article! --%s-- \n Err:  %v\n", refTitle, err)
					continue // skip this article and move to the next one
				}

				dataContainerNodes, err := GenerateContainerNodes(nodesInString, refTitle, refURL)
				if err != nil {
					util.Red("--------> Error:GenerateContainerNodes. -- Skipping this article! \n %v\n", err)
					continue // skip this article and move to the next one
				}
				fmt.Println("--------> Intellichunked into ", len(dataContainerNodes), " nodes successfully.")

				// User have the option to save every node into a json file.
				if saveToFile {
					// Create the JSON file with the reftitle as the name
					sanitizedTitle := util.SanitizeFileName(refTitle)
					fileName := filepath.Join(filepath.Dir(path), sanitizedTitle+".json")
					file, err := os.Create(fileName)
					if err != nil {
						log.Println("Error creating JSON file:", err)
						continue
					}
					// Create a copy of dataContainerNodes
					dataWithoutEmbedding := make([]models.ContainerNodeVector, len(dataContainerNodes))
					copy(dataWithoutEmbedding, dataContainerNodes)

					// Iterate over the copy and set the Embedding field to nil
					for i := range dataWithoutEmbedding {
						dataWithoutEmbedding[i].Embedding = nil
					}

					// Create a new JSON encoder and write the nodes to the file
					encoder := json.NewEncoder(file)
					err = encoder.Encode(dataWithoutEmbedding)
					if err != nil {
						log.Println("Error encoding JSON to file:", err)
						// log the error and move to the next iteration
						continue
					}
				}

				store := vectorstore.NewWeaviateStore()
				fileObjIDs, err := store.AddNodeObjects(className, dataContainerNodes)
				if err != nil {
					util.Red("--------> Skipping this article! --%s-- \n Err:  %v\n", refTitle, err)
					continue // skip this article and move to the next one
				}

				util.Green("--------> Successful. Added to the vectorstore.")
				fmt.Println("--------> Vector Object IDs:", strings.Join(fileObjIDs, ", "))
				objIDs = append(objIDs, fileObjIDs...)
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error walking the directory: ", err)
		return nil, err
	}

	return objIDs, nil
}

func AddWithoutNodes(className string, dataObjects []models.GeneralDataHolder) (objIDs []string, err error) {

	var embedTextSlice []string
	for _, dObj := range dataObjects {
		mergedText := dObj.Title + " \n " + dObj.Content
		embedTextSlice = append(embedTextSlice, mergedText)
	}

	// Create a new OpenAI language model instance.
	languageModel := llm.NewOpenAI()

	//embedBatch holds [][]float32 of embeddings
	embedBatch, err := languageModel.GenerateMultipleEmbeddingsFromText(context.Background(), embedTextSlice)
	if err != nil {
		log.Println("Error GenerateMultipleEmbeddingsFromText  :", err)
		return objIDs, err
	}

	for i := range dataObjects {
		dataObjects[i].Embedding = embedBatch[i]
	}

	store := vectorstore.NewWeaviateStore()
	objIDs, err = store.AddGenericObjects(className, dataObjects)
	if err != nil {
		log.Println("Error Intellichunk AddGenericObjects  :", err)
		return objIDs, err
	}

	return objIDs, nil
}
