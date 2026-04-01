package main

import (
	"os"

	"github.com/liao-eli/cc-cli-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
