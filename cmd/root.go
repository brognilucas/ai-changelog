package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ai-changelog",
		Short: "Generate changelogs from git commits using AI",
	}

	rootCmd.PersistentFlags().StringP("output", "o", "", "write changelog to file instead of stdout")
	rootCmd.PersistentFlags().StringP("since", "s", "", "generate changelog since tag or date")
	rootCmd.PersistentFlags().StringP("model", "m", "tinyllama", "ollama model to use for summarization")
	rootCmd.PersistentFlags().StringP("format", "f", "markdown", "output format: markdown or plain")

	return rootCmd
}
