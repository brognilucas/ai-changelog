package cmd_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lucasbrogni/ai-changelog/cmd"
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type mockCommitReader struct {
	commits []git.Commit
	err     error
}

func (m *mockCommitReader) GetCommits(since string) ([]git.Commit, error) {
	return m.commits, m.err
}

type mockOllamaClient struct {
	summaries []string
	err       error
	healthy   bool
}

func (m *mockOllamaClient) HealthCheck() error {
	if !m.healthy {
		return errOllamaDown
	}
	return nil
}

func (m *mockOllamaClient) SummarizeCommits(commits []git.Commit, model string) ([]string, error) {
	return m.summaries, m.err
}

var errOllamaDown = &ollamaDownError{}

type ollamaDownError struct{}

func (e *ollamaDownError) Error() string { return "ollama not reachable" }

func TestGenerateCommand(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
			{Hash: "def4567ghi", Subject: "fix: resolve crash", Author: "Bob", Timestamp: baseTime.Add(time.Hour), Prefix: "fix"},
		},
	}

	ollamaClient := &mockOllamaClient{
		healthy: true,
		summaries: []string{
			"Added user login feature and fixed application crash",
		},
	}

	var output bytes.Buffer

	deps := cmd.GenerateDeps{
		CommitReader: commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()

	if result == "" {
		t.Error("expected non-empty output")
	}

	if !bytes.Contains(output.Bytes(), []byte("New Features")) {
		t.Errorf("expected output to contain 'New Features', got:\n%s", result)
	}

	if !bytes.Contains(output.Bytes(), []byte("Bug Fixes")) {
		t.Errorf("expected output to contain 'Bug Fixes', got:\n%s", result)
	}
}
