package cmd

import (
	"fmt"
	"io"

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