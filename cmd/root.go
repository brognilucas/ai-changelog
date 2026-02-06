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

	return rootCmd
}
