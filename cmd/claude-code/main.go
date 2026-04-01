package main

import (
	"os"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
