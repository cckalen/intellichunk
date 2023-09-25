package conversation_test

import (
	"testing"

	"github.com/apsystole/log"

	"github.com/cckalen/intellichunk/config"
	"github.com/cckalen/intellichunk/internal/conversation"
	"github.com/cckalen/intellichunk/internal/models"
	"github.com/hlindberg/testutils"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Errorf("Load env Error: %s", err)

	}
}

// Test_GetVectors tests searching vector database and bringing relevant vectors/content.
func Test_GetVectors(t *testing.T) {
	content, refurls, err := conversation.GetRelevantContent("Class_demo1_kansas", "How much of Fossil-fuel combustion attributed to residential and commercial buildings?")
	testutils.CheckError(err, t)
	testutils.CheckNotNil(content, t)
	testutils.CheckNotNil(refurls, t)
}

// Test_ClassConversationInvalidID tests the ClassConversation function with invalid Classid.
func Test_ClassConversationInvalidID(t *testing.T) {
	// Define the conversation request with an invalid ClassID.
	convoReq := models.ConversationRequest{
		ConversationID: "TestID",
		ClassID:        "veryWrongClassID",
		ChatHistory:    []string{""},
		Query:          "How far is the moon?",
	}

	// Call the ClassConversation function with the invalid request.
	convoResp, err := conversation.ClassConversation(convoReq)

	// Assert that an error was returned.
	testutils.CheckNotNil(err, t)

	// Assert that the returned response is the zero value.
	expectedConvoResp := models.ConversationResponse{}
	testutils.CheckEqual(expectedConvoResp, convoResp, t)
}

// Test_ClassConvoSuccess tests the whole function with chat history
func Test_ClassConvoSuccess(t *testing.T) {
	convoReq := models.ConversationRequest{
		ConversationID: "TestID",
		ClassID:        "Class_testRun",
		ChatHistory:    []string{"What is the location of Ladakh?", "Ladakh is located in the eastern part of the larger Kashmir region and is administered as a union territory by India. It is bordered by the Tibet Autonomous Region to the east and the Indian state of Himachal Pradesh to the south.", "Does Tibet have authority over it?", "No, Tibet does not have authority over Ladakh. Ladakh is administered as a union territory by the Government of India. While Ladakh shares a border with the Tibet Autonomous Region, it is governed by the Indian authorities."},
		Query:          "Can you print out each question and answer as bullet points so far",
	}

	convoResp, err := conversation.ClassConversation(convoReq)

	log.Println(convoResp)
	testutils.CheckNotError(err, t)
}
