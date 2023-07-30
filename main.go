package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/tywil04/slavartdl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}
}
