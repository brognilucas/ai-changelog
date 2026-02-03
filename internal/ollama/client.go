package ollama

import (
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type Client interface {
	HealthCheck() error
	SummarizeCommits(commits []git.Commit, model string) ([]string, error)
}