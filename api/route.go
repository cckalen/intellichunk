package api

import (
	"encoding/json"
	"net/http"

	"github.com/cckalen/intellichunk/internal/conversation"
	"github.com/cckalen/intellichunk/internal/intellichunk"
	"github.com/cckalen/intellichunk/internal/models"
)

func ConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var convReq models.ConversationRequest
	err := json.NewDecoder(r.Body).Decode(&convReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error Marshalling the request, provide valid ConversationRequest"}`))
		return
	}

	convResp, err := conversation.ClassConversation(convReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error calling ClassConversation"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(convResp)
}

func IntellichunkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intReq models.IntellichunkRequest
	err := json.NewDecoder(r.Body).Decode(&intReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error Marshalling the request, provide valid IntellichunkRequest"}`))
		return
	}

	objIDs, err := intellichunk.Add(intReq.ClassName, intReq.LongText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error calling ClassConversation"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objIDs)
}
