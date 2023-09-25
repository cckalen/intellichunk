// Package main invokes the command line interpreter.
package main

import (
	"os"

	"github.com/cckalen/intellichunk/api"
	"github.com/cckalen/intellichunk/cmd"
	"github.com/cckalen/intellichunk/config"
)

// main runs the command line interpreter
// it runs the api if it's not in local
func main() {
	config.LoadEnv()

	if os.Getenv("RUN_ENV") != "local" {
		api.Run_api()
	} else {
		cmd.Execute()
	}

}
