package cli

import (
	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "claude-code",
	Short: "Claude Code CLI - AI-powered coding assistant",
	Long:  `A Go implementation of Claude Code CLI for learning purposes.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Version = Version
}
