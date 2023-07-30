package main

import (
	"log"
	"os"

	"github.com/tywil04/slavartdl/cmd"
	"github.com/tywil04/slavartdl/internal/config"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
