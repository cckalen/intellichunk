package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cckalen/intellichunk/internal/intellichunk"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [class name] [folder path]",
	Short: "Add container nodes to the vectorstore",
	Long: `The 'add' command takes a class name and the path to a folder containing text files as input. 
	It iterates over the text files within the folder, reads the text from each file, splits it into 
	container nodes, generates embeddings for the nodes, and adds the resulting objects to the vectorstore.
	This function is useful for batch processing of large text files and storing their context in a 
	structured and accessible format.

	path is relative to the project folder.
	-s flag can be used to save everything into .json files.
	For example:
	add "class1" "/files/" -s`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatalf("add command requires exactly 2 arguments: [class name] [folder path]")
		}
		// Retrieve the value of the save flag
		save, err := cmd.Flags().GetBool("save")
		if err != nil {
			log.Fatalf("Error retrieving save flag: %v", err)
		}

		className := args[0]
		relFolderPath := args[1]
		// Getting the current working directory (you should run this where your project root is)
		projectRoot, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting project root: %v", err)
		}

		// Joining the project root with the relative folder path provided
		absFolderPath := filepath.Join(projectRoot, relFolderPath)

		_, err = intellichunk.AddFromFolder(className, absFolderPath, save)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	intellichunkCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("save", "s", false, "Indicate if you want to save the nodes into a file")
}
