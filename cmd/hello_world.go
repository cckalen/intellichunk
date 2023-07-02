// Package cmd contains the example hello CLI logic.
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tada/catch"
	"github.com/wyrth-io/goapp-template/internal/example"
)

var helloWorldCmd = &cobra.Command{
	Use:   "world [<name>]",
	Short: "world",
	Long: `This is an example "world" subcommand to the hello command"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var out string
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		if fail {
			err := catch.Error("Failure test")
			nestedFail(err)
		}
		switch {
		case wonderful && UpperCase:
			out = example.GreetWonderfulUpper(name)
		case wonderful:
			out = example.GreetWonderful(name)
		case UpperCase:
			out = example.GreetUpper(name)
		default:
			out = example.Greet(name)
		}

		fmt.Println(out)
	},

	Args: func(cmd *cobra.Command, args []string) error {
		// validate flags/options here return nil if all is fine else an error
		if len(args) > 1 {
			return errors.New("at most one argument accepted")
		}
		return nil
	},
}
var wonderful bool // private
var fail bool      // for test of what a failure looks like

// UpperCase is an example of a global variable (Note: it is required to document all such variables).
var UpperCase bool // global

func init() {
	helloCmd.AddCommand(helloWorldCmd)

	flags := helloWorldCmd.PersistentFlags()
	flags.BoolVarP(&wonderful, "wonderful", "w", false, "for a wonderful world")
	flags.BoolVarP(&UpperCase, "upcase", "u", false, "for upper case hello")
	flags.BoolVarP(&fail, "fail", "f", false, "cause fail to see error handling")
}

func nestedFail(err error) {
	if err != nil {
		panic(catch.Error(err))
	}
}
