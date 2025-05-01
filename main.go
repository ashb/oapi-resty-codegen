package main

import (
	"fmt"
	"os"

	"github.com/ashb/oapi-resty-codegen/cmds"
)

func main() {
	err := cmds.RootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
