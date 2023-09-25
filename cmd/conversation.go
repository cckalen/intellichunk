package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cckalen/intellichunk/internal/conversation"
	"github.com/cckalen/intellichunk/internal/models"
	"github.com/spf13/cobra"
)

// conversationCmd represents the conversation command
var conversationCmd = &cobra.Command{
	Use:   "conversation",
	Short: "Start a conversation with a specific Class",
	Long:  "This command initiates a conversation with the specified Class using the given query.",
	Run:   runConversation,
}

var chatHistory []string

func runConversation(cmd *cobra.Command, args []string) {
	var ClassID, query string

	if len(args) >= 2 {
		ClassID = args[0]
		query = args[1]
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Please enter your Class ID: ")
		ClassIDInput, _ := reader.ReadString('\n')
		ClassID = strings.TrimSpace(ClassIDInput)

		fmt.Print("Please enter your query: ")
		queryInput, _ := reader.ReadString('\n')
		query = strings.TrimSpace(queryInput)
	}

	for {
		convoReq := models.ConversationRequest{
			ClassID:     ClassID,
			Query:       query,
			ChatHistory: chatHistory,
		}

		convoResp, err := conversation.ClassConversation(convoReq)
		if err != nil {
			log.Fatalf("Error in conversation: %v", err)
		}

		chatHistory = append(chatHistory, query, convoResp.Answer.Answer)
		fmt.Println("------  AI:", convoResp.Answer.Answer)

		fmt.Print("\n ::::::  You: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		query = strings.TrimSpace(input)

		if query == "exit" {
			break
		}
	}

	fmt.Println("Exiting the conversation...")
}

func init() {
	RootCmd.AddCommand(conversationCmd)
}
