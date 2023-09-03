package main

import (
	"os"

	"github.com/tywil04/slavartdl/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
