package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/brognilucas/ai-changelog/cmd"
	"github.com/brognilucas/ai-changelog/internal/git"
)

type mockCommitReader struct {
	commits []git.Commit
	err     error
}

func (m *mockCommitReader) GetCommits(since string) ([]git.Commit, error) {
	return m.commits, m.err
}

type mockOllamaClient struct {
	summaries       []string
	err             error
	healthy         bool
	changelogOutput string
	changelogErr    error
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

func (m *mockOllamaClient) GenerateChangelog(commits []git.Commit, model string) (string, error) {
	return m.changelogOutput, m.changelogErr
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
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", "tinyllama", "", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()

	if result == "" {
		t.Error("expected non-empty output")
	}

	// With empty changelogOutput, LLM path fails → falls back to structured
	if !bytes.Contains(output.Bytes(), []byte("New Features")) {
		t.Errorf("expected output to contain 'New Features', got:\n%s", result)
	}

	if !bytes.Contains(output.Bytes(), []byte("Bug Fixes")) {
		t.Errorf("expected output to contain 'Bug Fixes', got:\n%s", result)
	}
}

func TestWriteToFile(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		},
	}

	ollamaClient := &mockOllamaClient{healthy: true}

	deps := cmd.GenerateDeps{
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "CHANGELOG.md")

	err := cmd.WriteToFile(deps, "markdown", "", "tinyllama", "", outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected non-empty file content")
	}

	if !bytes.Contains(content, []byte("New Features")) {
		t.Errorf("expected file to contain 'New Features', got:\n%s", string(content))
	}
}

func TestOllamaStartupCheck(t *testing.T) {
	t.Run("returns user-friendly error when ollama is down", func(t *testing.T) {
		ollamaClient := &mockOllamaClient{healthy: false}

		err := cmd.CheckOllamaHealth(ollamaClient)
		if err == nil {
			t.Fatal("expected error when ollama is not reachable")
		}

		errMsg := err.Error()
		if !bytes.Contains([]byte(errMsg), []byte("Ollama is not running")) {
			t.Errorf("expected user-friendly error message, got: %s", errMsg)
		}

		if !bytes.Contains([]byte(errMsg), []byte("ollama serve")) {
			t.Errorf("expected error to suggest 'ollama serve', got: %s", errMsg)
		}
	})

	t.Run("returns nil when ollama is healthy", func(t *testing.T) {
		ollamaClient := &mockOllamaClient{healthy: true}

		err := cmd.CheckOllamaHealth(ollamaClient)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGenerateWithLLM(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
			{Hash: "def4567ghi", Subject: "test: add login tests", Author: "Alice", Timestamp: baseTime.Add(time.Hour), Prefix: "test"},
		},
	}

	llmOutput := "_This release adds user authentication._\n\n## Highlights\n\n- User login support\n"

	ollamaClient := &mockOllamaClient{
		healthy:         true,
		changelogOutput: llmOutput,
	}

	var output bytes.Buffer
	deps := cmd.GenerateDeps{
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", "tinyllama", "", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()

	if !strings.Contains(result, "Highlights") {
		t.Errorf("expected LLM output with 'Highlights', got:\n%s", result)
	}

	if strings.Contains(result, "New Features") {
		t.Errorf("expected LLM path, but got structured output with 'New Features':\n%s", result)
	}

	if strings.Contains(result, "Testing") {
		t.Errorf("expected LLM to omit testing section, got:\n%s", result)
	}
}

func TestGenerateLLMFallback(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		},
	}

	ollamaClient := &mockOllamaClient{
		healthy:      true,
		changelogErr: fmt.Errorf("model not found"),
	}

	var output bytes.Buffer
	deps := cmd.GenerateDeps{
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", "tinyllama", "", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()

	if !strings.Contains(result, "New Features") {
		t.Errorf("expected fallback to structured output with 'New Features', got:\n%s", result)
	}
}

func TestGenerateWithVersion(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		},
	}

	llmOutput := "_Summary._\n\n## Highlights\n\n- Login\n"

	ollamaClient := &mockOllamaClient{
		healthy:         true,
		changelogOutput: llmOutput,
	}

	var output bytes.Buffer
	deps := cmd.GenerateDeps{
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", "tinyllama", "v1.0.0", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "# v1.0.0") {
		t.Errorf("expected version header in output, got:\n%s", result)
	}
}

func TestGenerateOllamaUnhealthyFallback(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commitReader := &mockCommitReader{
		commits: []git.Commit{
			{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		},
	}

	ollamaClient := &mockOllamaClient{
		healthy: false,
	}

	var output bytes.Buffer
	deps := cmd.GenerateDeps{
		CommitReader:  commitReader,
		OllamaClient: ollamaClient,
	}

	err := cmd.RunGenerate(deps, "markdown", "", "tinyllama", "", &output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "New Features") {
		t.Errorf("expected fallback to structured output, got:\n%s", result)
	}
}
