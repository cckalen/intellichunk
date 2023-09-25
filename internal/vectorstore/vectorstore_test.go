package vectorstore_test

import (
	"testing"

	"github.com/apsystole/log"
	"github.com/cckalen/intellichunk/config"
	"github.com/cckalen/intellichunk/internal/vectorstore"
	"github.com/hlindberg/testutils"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Errorf("Load env Error: %s", err)

	}
}
func Test_WeaviateClassExist(t *testing.T) {

	store := vectorstore.NewWeaviateStore()

	store.CheckAndCreateClass("testClassName")
}

func Test_SimilaritySearch(t *testing.T) {
	className := "Class_test5"
	input := "What are the challenges in addressing public health risks associated with climate change?"
	graphFieldNames := []string{"content", "title"}
	withLimit := 2

	store := vectorstore.NewWeaviateStore()

	results, err := store.SimilaritySearch(className, input, graphFieldNames, withLimit)
	if err != nil {
		t.Fatalf("Error in SimilaritySearch: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("No results found")
	}

	testutils.CheckNotError(err, t)
	testutils.CheckNotNil(results, t)

	// for i, result := range results {
	// 	t.Logf("Result %d:", i+1)
	// 	for key, value := range result {
	// 		t.Logf("%s: %v", key, value)
	// 	}
	// }
}

func Test_DeleteObject(t *testing.T) {
	className := "CLassDelete"
	objectId := "ff205228-26e7-430d-a666-71f6a69dc77f" // specific object id from vector store.

	store := vectorstore.NewWeaviateStore()
	err := store.DeleteObjectByID(className, objectId)
	testutils.CheckNotError(err, t)
}
