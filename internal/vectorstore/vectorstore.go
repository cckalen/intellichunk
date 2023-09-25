// Package vectorstore provides an interface and implementation for a vector database.
// It includes a Weaviate implementation of the VectorStore interface.
package vectorstore

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/apsystole/log"
	"github.com/cckalen/intellichunk/internal/models"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	wmodels "github.com/weaviate/weaviate/entities/models"
)

// VectorStore is an abstraction of a vector database.
// This interface makes it easier to test the code and swap the underlying implementation.
type VectorStore interface {
	CheckOrCreateClass(className string) error
	AddObjects(className string, objects []models.ContainerNodeVector, weaviateKey string) (objIDs []string, err error)
	AddGenericObjects(className string, objects []models.GeneralDataHolder) (objIDs []string, err error)
	DeleteObjectByID(className, objectID string) (err error)
	GetObjects(weaviateKey string, openAIKey string, className string, graphFieldNames []string, withLimit int) (interface{}, error)
	SimilaritySearch(weaviateKey string, openAIKey string, className string, input string, graphFieldNames []string, withLimit int) ([]map[string]interface{}, error)
}

// WeaviateStore is a Weaviate implementation of the VectorStore interface.
type WeaviateStore struct {
	Host        string
	Scheme      string
	WeaviateKey string
	OpenAIKey   string
}

// Option is a function that can modify the WeaviateStore configuration.
type Option func(*WeaviateStore)

// WithHost sets the custom host for the WeaviateStore.
func WithHost(host string) Option {
	return func(store *WeaviateStore) {
		store.Host = host
	}
}

// WithScheme sets the custom scheme for the WeaviateStore.
func WithScheme(scheme string) Option {
	return func(store *WeaviateStore) {
		store.Scheme = scheme
	}
}

// WithWeaviateKey sets the Weaviate API key for the WeaviateStore.
func WithWeaviateKey(key string) Option {
	return func(store *WeaviateStore) {
		store.WeaviateKey = key
	}
}

// WithOpenAIKey sets the OpenAI API key for the WeaviateStore.
func WithOpenAIKey(key string) Option {
	return func(store *WeaviateStore) {
		store.OpenAIKey = key
	}
}

// NewWeaviateStore creates a new instance of WeaviateStore with optional configurations.
func NewWeaviateStore(options ...Option) *WeaviateStore {
	store := &WeaviateStore{
		Host:        os.Getenv("WEAVIATE_URL"),
		Scheme:      "https",
		WeaviateKey: os.Getenv("WEAVIATE_API_KEY"),
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
	}

	for _, option := range options {
		option(store)
	}

	return store
}

// convertNamesToFields converts field names to graphql.Field.
func convertNamesToFields(names []string) []graphql.Field {
	fields := make([]graphql.Field, 0, len(names))
	for _, name := range names {
		fields = append(fields, graphql.Field{Name: name})
	}
	return fields
}

// GetObjects retrieves objects from Weaviate based on the specified parameters.
func (store WeaviateStore) GetObjects(className string, graphFieldNames []string, withLimit int) (interface{}, error) {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
		Headers:    map[string]string{"X-OpenAI-Api-Key": store.OpenAIKey},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	fields := convertNamesToFields(graphFieldNames)

	result, err := client.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithLimit(withLimit).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SimilaritySearch performs a similarity search in Weaviate based on the specified parameters.
// It utilizes the OpenAI Ada embeddings model to vectorize the given query and then searches the Weaviate database for similar results.
//
// Parameters:
//   - className: The classname of the objects to search in Weaviate.
//   - input: The query or concept to search for similarities.
//   - graphFieldNames: The field names in the Weaviate objects to consider for the search.
//   - withLimit: The maximum number of similar results to retrieve.
//
// Returns:
//   - An array of maps, where each map represents a relevant object found in Weaviate.
//     Each map contains the selected graph field names as keys and their corresponding values.
//   - An error if the search or retrieval process fails.
//
// Note: The function uses the Weaviate client and the GraphQL Get method to perform the similarity search.
// It constructs a nearText argument with the provided input and a distance threshold of 0.8 to find similar objects.
// The resulting objects are extracted and transformed into a map-based structure for easier retrieval of desired fields.
func (store WeaviateStore) SimilaritySearch(className string, input string, graphFieldNames []string, withLimit int) ([]map[string]interface{}, error) {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
		Headers:    map[string]string{"X-OpenAI-Api-Key": store.OpenAIKey},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Errorf("Failed to create new Weaviate client: %v", err)
		return nil, err
	}

	fields := convertNamesToFields(graphFieldNames)

	concepts := []string{input}
	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts(concepts).
		WithDistance(0.8)

	result, err := client.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(withLimit).
		Do(context.Background())

	if err != nil {
		log.Errorf("Failed to perform GraphQL Get: %v", err)
		return nil, err
	} else if len(result.Errors) > 0 {
		log.Errorf("Failed to perform GraphQL: %v", result.Errors[0].Message)
		return nil, errors.New(result.Errors[0].Message)
	}

	rawObjects, ok := result.Data["Get"].(map[string]interface{})[className]
	if !ok {
		log.Errorf("ClassName not found in result")
		return nil, errors.New("ClassName not found in result")
	}

	objectList, ok := rawObjects.([]interface{})
	if !ok {
		// we silently continue if this type assertion fails
		return nil, nil
	}

	objects := make([]map[string]interface{}, len(objectList))

	for i, rawObj := range objectList {
		rawMap, ok := rawObj.(map[string]interface{})
		if !ok {
			// we silently continue if this type assertion fails
			continue
		}

		// Create a map to store the desired fields
		object := make(map[string]interface{})

		// Loop over the desired fields and extract them from the rawMap
		for _, fieldName := range graphFieldNames {
			fieldValue, ok := rawMap[fieldName]
			if !ok {
				// we silently continue if this type assertion fails
				continue
			}

			object[fieldName] = fieldValue
		}

		objects[i] = object
	}

	return objects, nil
}

// AddObjects adds the provided ContainerNodeVector objects to the specified Weaviate class in batch mode. The objects are represented
// by a slice of models.ContainerNodeVector. The function utilizes the Weaviate client's batch mode to efficiently
// add multiple objects at once. It first checks if the class exists, and if not, creates the class using CheckAndCreateClass.
// The function returns the IDs of the added objects and an error if any issues occur during the process.
func (store WeaviateStore) AddNodeObjects(className string, objects []models.ContainerNodeVector) (objIDs []string, err error) {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Errorf("Failed to create new Weaviate client: %v", err)
		return objIDs, err
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, obj := range objects {
		properties := make(map[string]interface{})
		objValue := reflect.ValueOf(obj)
		objType := reflect.TypeOf(obj)
		for i := 0; i < objValue.NumField(); i++ {
			field := objValue.Field(i)
			fieldInfo := objType.Field(i)
			jsonTag := fieldInfo.Tag.Get("json")
			// Check if there's a valid JSON tag, and use only the part before any comma
			if jsonTag != "" {
				jsonTag = strings.Split(jsonTag, ",")[0]
			}
			// Exclude the "embedding" field from properties
			if jsonTag != "embedding" {
				properties[jsonTag] = field.Interface()
			}
		}

		weaviateObject := &wmodels.Object{
			Class:      className,
			Properties: properties,
			Vector:     obj.Embedding,
		}

		batcher.WithObjects(weaviateObject)
	}

	//Check if the class exist already and create otherwise
	err = store.CheckAndCreateClass(className)
	if err != nil {
		log.Errorf("Failed to add objects to Weaviate: %v", err)
		return objIDs, err
	}

	result, err := batcher.Do(context.Background())
	if err != nil {
		log.Errorf("Failed to add objects to Weaviate: %v", err)
		return objIDs, err
	}

	for _, res := range result {
		objIDs = append(objIDs, res.ID.String())
	}
	return objIDs, nil
}

// CheckAndCreateClass checks if the given class exists in the Weaviate database. If the class does not exist,
// it creates the class with the specified className and the required schema for text-based vectorization.
// It utilizes the Weaviate client and the GraphQL Get method to perform the existence check and creation.
func (store WeaviateStore) CheckAndCreateClass(className string) error {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Errorf("Failed to create new Weaviate client: %v", err)

	}

	if className[0] >= 'a' && className[0] <= 'z' {
		// Convert the first letter to uppercase
		className = strings.ToUpper(string(className[0])) + className[1:]
	}

	// Check if class exist
	classExist, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(context.Background())

	if err != nil {
		log.Print(err)
		return err
	}

	if classExist {
		return nil
	}

	// If Class doesn't exist, create one
	classObj := &wmodels.Class{
		Class:      className,
		Vectorizer: "text2vec-openai", // If set to "none" you must always provide vectors yourself. Could be any other "text2vec-*" also.
		ModuleConfig: map[string]interface{}{
			"text2vec-openai": map[string]interface{}{
				"model":        "ada",
				"modelVersion": "002",
				"type":         "text",
			},
		},
	}
	// add the schema
	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (store WeaviateStore) AddGenericObjects(className string, objects []models.GeneralDataHolder) (objIDs []string, err error) {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Errorf("Failed to create new Weaviate client: %v", err)
		return objIDs, err
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, obj := range objects {
		weaviateObject := &wmodels.Object{
			Class: className,
			Properties: map[string]interface{}{
				"title":   obj.Title,
				"content": obj.Content,
			},
			Vector: obj.Embedding,
		}

		batcher.WithObjects(weaviateObject)
	}

	//Check if the class exist already and create otherwise
	err = store.CheckAndCreateClass(className)
	if err != nil {
		log.Errorf("Failed CheckAndCreateClass to Weaviate: %v", err)
		return objIDs, err
	}

	result, err := batcher.Do(context.Background())
	if err != nil {
		log.Errorf("Failed batcher. Do to Weaviate: %v", err)
		return objIDs, err
	}

	for _, res := range result {
		objIDs = append(objIDs, res.ID.String())
	}
	return objIDs, nil
}

func (store WeaviateStore) DeleteObjectByID(className, objectID string) (err error) {

	cfg := weaviate.Config{
		Host:       store.Host,
		Scheme:     store.Scheme,
		AuthConfig: auth.ApiKey{Value: store.WeaviateKey},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Errorf("Failed to create new Weaviate client: %v", err)
		return err
	}

	err = client.Data().Deleter().
		WithClassName(className).
		WithID(objectID).
		//WithConsistencyLevel(replication.ConsistencyLevel.ALL).  // default QUORUM
		Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}
