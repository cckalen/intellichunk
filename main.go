// Package main invokes the command line interpreter.
package main

import (
	"github.com/wyrth-io/goapp-template/cmd"
)

// main runs the command line interpreter - it has a default function when this application
// is run without any arguments that will check if there is input on stdin and timeout with
// a message if that is not the case. If any subcommand is given it will be executed and
// there will be no check if anything is present on stdin.
func main() {
	cmd.Execute()
}
