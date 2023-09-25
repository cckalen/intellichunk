/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// intellichunkCmd represents the intellichunk command
var intellichunkCmd = &cobra.Command{
	Use:   "intellichunk",
	Short: "Intelligent text chunking and metadata enhancement",
	Long: `The 'intellichunk' command leverages AI to divide large texts into manageable chunks 
    and enriches them with relevant metadata. It's designed for analysis of long factual 
    documents, overcoming token limitations, and it supports batch embeddings for efficient 
    vector generation. Although it doesn't vectorize keywords due to LLM performance, they 
    can be used for other purposes and kept in the db.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("intellichunk called")
	},
}

func init() {
	RootCmd.AddCommand(intellichunkCmd)

}
