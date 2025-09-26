package main

import (
	"fmt"
	"os"

	"github.com/ashb/oapi-resty-codegen/cmds"
)

func main() {
	// We'll print the error ourselves
	cmds.RootCmd.SilenceErrors = true
	err := cmds.RootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
