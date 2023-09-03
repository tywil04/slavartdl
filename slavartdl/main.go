package main

import (
	"os"

	"github.com/tywil04/slavartdl/slavartdl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
