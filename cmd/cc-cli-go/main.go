package main

import (
	"os"

	"github.com/user-name/cc-cli-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
