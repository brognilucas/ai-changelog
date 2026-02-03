package ollama_test

import (
	"testing"

	"github.com/lucasbrogni/ai-changelog/internal/git"
	"github.com/lucasbrogni/ai-changelog/internal/ollama"
)

type mockOllamaClient struct{}

func (m *mockOllamaClient) HealthCheck() error {
	return nil
}

func (m *mockOllamaClient) SummarizeCommits(commits []git.Commit, model string) ([]string, error) {
	return nil, nil
}

func TestOllamaClientInterface(t *testing.T) {
	var client ollama.Client = &mockOllamaClient{}

	err := client.HealthCheck()
	if err != nil {
		t.Errorf("unexpected error from HealthCheck: %v", err)
	}

	_, err = client.SummarizeCommits([]git.Commit{}, "llama3")
	if err != nil {
		t.Errorf("unexpected error from SummarizeCommits: %v", err)
	}
}