package main

import (
	"fmt"
	"os"

	"github.com/lucasbrogni/ai-changelog/cmd"
	"github.com/lucasbrogni/ai-changelog/internal/git"
	"github.com/lucasbrogni/ai-changelog/internal/ollama"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCommand()

	rootCmd.RunE = func(c *cobra.Command, args []string) error {
		since, _ := c.Flags().GetString("since")
		format, _ := c.Flags().GetString("format")
		output, _ := c.Flags().GetString("output")
		model, _ := c.Flags().GetString("model")
		version, _ := c.Flags().GetString("version")

		runner := &git.DefaultRunner{}
		commitReader := git.NewCommitReader(runner)
		ollamaClient := ollama.NewDefaultClient("http://localhost:11434")

		deps := cmd.GenerateDeps{
			CommitReader: commitReader,
			OllamaClient: ollamaClient,
		}

		if err := cmd.CheckOllamaHealth(ollamaClient); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v (using raw commit messages)\n", err)
		}

		if output != "" {
			return cmd.WriteToFile(deps, format, since, model, version, output)
		}

		return cmd.RunGenerate(deps, format, since, model, version, os.Stdout)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
