package main

import (
	"log"

	"slavartdl/cmd"
	"slavartdl/lib/config"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
