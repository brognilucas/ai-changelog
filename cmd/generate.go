package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brognilucas/ai-changelog/internal/changelog"
	"github.com/brognilucas/ai-changelog/internal/git"
	"github.com/brognilucas/ai-changelog/internal/ollama"
)

type CommitReader interface {
	GetCommits(since string) ([]git.Commit, error)
}

type GenerateDeps struct {
	CommitReader CommitReader
	OllamaClient ollama.Client
}

func RunGenerate(deps GenerateDeps, format string, since string, model string, version string, writer io.Writer) error {
	commits, err := deps.CommitReader.GetCommits(since)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Fprintln(writer, "No commits found.")
		return nil
	}

	// Try LLM path first
	if deps.OllamaClient != nil {
		if err := deps.OllamaClient.HealthCheck(); err == nil {
			changelogText, llmErr := deps.OllamaClient.GenerateChangelog(commits, model)
			if llmErr == nil && strings.TrimSpace(changelogText) != "" {
				var output string
				if version != "" {
					output = fmt.Sprintf("# %s\n\n%s", version, changelogText)
				} else {
					output = changelogText
				}
				_, err = fmt.Fprint(writer, output)
				return err
			}
			if llmErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: LLM generation failed (%v), falling back to structured output\n", llmErr)
			}
		}
	}

	// Fallback: structured rendering
	sorted := changelog.SortByDate(commits)
	sections := changelog.GroupByCategory(sorted)

	var renderer changelog.Renderer
	if format == "plain" {
		renderer = &changelog.PlainTextRenderer{}
	} else {
		renderer = &changelog.MarkdownRenderer{}
	}

	output := renderer.Render(sections, version)
	_, err = fmt.Fprint(writer, output)
	return err
}

func WriteToFile(deps GenerateDeps, format string, since string, model string, version string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	return RunGenerate(deps, format, since, model, version, file)
}

func CheckOllamaHealth(client ollama.Client) error {
	if err := client.HealthCheck(); err != nil {
		return fmt.Errorf("Ollama is not running. Start it with: ollama serve")
	}
	return nil
}