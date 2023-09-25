package cmd

import (
	"github.com/cckalen/intellichunk/api"
	"github.com/spf13/cobra"
)

// runapiCmd represents the runapi command.
// useful for runing the api server on your local machine.
var runapiCmd = &cobra.Command{
	Use:   "runapi",
	Short: "Runs the api server",
	Long:  `Start the api server.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.Run_api()
	},
}

func init() {
	RootCmd.AddCommand(runapiCmd)
}
