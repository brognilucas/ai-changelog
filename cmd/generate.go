package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/lucasbrogni/ai-changelog/internal/changelog"
	"github.com/lucasbrogni/ai-changelog/internal/git"
	"github.com/lucasbrogni/ai-changelog/internal/ollama"
)

type CommitReader interface {
	GetCommits(since string) ([]git.Commit, error)
}

type GenerateDeps struct {
	CommitReader CommitReader
	OllamaClient ollama.Client
}

func RunGenerate(deps GenerateDeps, format string, since string, writer io.Writer) error {
	commits, err := deps.CommitReader.GetCommits(since)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Fprintln(writer, "No commits found.")
		return nil
	}

	sorted := changelog.SortByDate(commits)
	sections := changelog.GroupByCategory(sorted)

	var renderer changelog.Renderer
	if format == "plain" {
		renderer = &changelog.PlainTextRenderer{}
	} else {
		renderer = &changelog.MarkdownRenderer{}
	}

	output := renderer.Render(sections, "")
	_, err = fmt.Fprint(writer, output)
	return err
}

func WriteToFile(deps GenerateDeps, format string, since string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	return RunGenerate(deps, format, since, file)
}

func CheckOllamaHealth(client ollama.Client) error {
	if err := client.HealthCheck(); err != nil {
		return fmt.Errorf("Ollama is not running. Start it with: ollama serve")
	}
	return nil
}