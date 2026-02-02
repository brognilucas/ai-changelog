package git_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lucasbrogni/ai-changelog/internal/git"
)

func TestCommitStructFields(t *testing.T) {
	timestamp := time.Now()

	commit := git.Commit{
		Hash:      "abc123",
		Subject:   "feat: add new feature",
		Author:    "John Doe",
		Timestamp: timestamp,
		Prefix:    "feat",
	}

	if commit.Hash != "abc123" {
		t.Errorf("expected Hash to be 'abc123', got '%s'", commit.Hash)
	}

	if commit.Subject != "feat: add new feature" {
		t.Errorf("expected Subject to be 'feat: add new feature', got '%s'", commit.Subject)
	}

	if commit.Author != "John Doe" {
		t.Errorf("expected Author to be 'John Doe', got '%s'", commit.Author)
	}

	if !commit.Timestamp.Equal(timestamp) {
		t.Errorf("expected Timestamp to be '%v', got '%v'", timestamp, commit.Timestamp)
	}

	if commit.Prefix != "feat" {
		t.Errorf("expected Prefix to be 'feat', got '%s'", commit.Prefix)
	}
}

func TestParseCommitLine(t *testing.T) {
	line := "abc123def456|feat: add new feature|John Doe|1706745600"

	commit, err := git.ParseCommitLine(line)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if commit.Hash != "abc123def456" {
		t.Errorf("expected Hash to be 'abc123def456', got '%s'", commit.Hash)
	}

	if commit.Subject != "feat: add new feature" {
		t.Errorf("expected Subject to be 'feat: add new feature', got '%s'", commit.Subject)
	}

	if commit.Author != "John Doe" {
		t.Errorf("expected Author to be 'John Doe', got '%s'", commit.Author)
	}

	expectedTime := time.Unix(1706745600, 0)
	if !commit.Timestamp.Equal(expectedTime) {
		t.Errorf("expected Timestamp to be '%v', got '%v'", expectedTime, commit.Timestamp)
	}
}

func TestParseCommitLineInvalidFormat(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"empty line", ""},
		{"missing fields", "abc123|feat: something"},
		{"too few pipes", "abc123|subject|author"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := git.ParseCommitLine(tt.line)
			if err == nil {
				t.Errorf("expected error for line '%s', got nil", tt.line)
			}
		})
	}
}

func TestParseCommitLineInvalidTimestamp(t *testing.T) {
	line := "abc123|feat: something|John Doe|not-a-number"

	_, err := git.ParseCommitLine(line)

	if err == nil {
		t.Error("expected error for invalid timestamp, got nil")
	}
}

func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		subject  string
		expected string
	}{
		{"feat: add new feature", "feat"},
		{"fix: resolve bug", "fix"},
		{"docs: update readme", "docs"},
		{"chore: update dependencies", "chore"},
		{"refactor: simplify code", "refactor"},
		{"test: add unit tests", "test"},
		{"style: fix formatting", "style"},
		{"perf: improve performance", "perf"},
		{"feat(scope): scoped feature", "feat"},
		{"fix(auth): fix login", "fix"},
	}

	for _, tt := range tests {
		t.Run(tt.subject, func(t *testing.T) {
			result := git.ExtractPrefix(tt.subject)
			if result != tt.expected {
				t.Errorf("ExtractPrefix(%q) = %q, want %q", tt.subject, result, tt.expected)
			}
		})
	}
}

func TestExtractPrefixUnknown(t *testing.T) {
	tests := []struct {
		name    string
		subject string
	}{
		{"no colon", "update readme"},
		{"unknown prefix", "unknown: something"},
		{"empty string", ""},
		{"just colon", ":"},
		{"space before colon", "feat : something"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := git.ExtractPrefix(tt.subject)
			if result != "other" {
				t.Errorf("ExtractPrefix(%q) = %q, want 'other'", tt.subject, result)
			}
		})
	}
}

type mockRunner struct {
	output string
	err    error
}

func (m *mockRunner) Run(args ...string) (string, error) {
	return m.output, m.err
}

func TestRunnerInterface(t *testing.T) {
	var runner git.Runner = &mockRunner{output: "test output", err: nil}

	output, err := runner.Run("log")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "test output" {
		t.Errorf("expected output 'test output', got '%s'", output)
	}
}

func TestDefaultRunnerImplementsInterface(t *testing.T) {
	var _ git.Runner = &git.DefaultRunner{}
}

func TestGetCommits(t *testing.T) {
	mockOutput := "abc123|feat: add feature|John Doe|1706745600\ndef456|fix: fix bug|Jane Doe|1706746600"
	runner := &mockRunner{output: mockOutput, err: nil}

	reader := git.NewCommitReader(runner)
	commits, err := reader.GetCommits("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	if commits[0].Hash != "abc123" {
		t.Errorf("expected first commit hash 'abc123', got '%s'", commits[0].Hash)
	}

	if commits[0].Prefix != "feat" {
		t.Errorf("expected first commit prefix 'feat', got '%s'", commits[0].Prefix)
	}

	if commits[1].Hash != "def456" {
		t.Errorf("expected second commit hash 'def456', got '%s'", commits[1].Hash)
	}

	if commits[1].Prefix != "fix" {
		t.Errorf("expected second commit prefix 'fix', got '%s'", commits[1].Prefix)
	}
}

func TestGetCommitsEmpty(t *testing.T) {
	runner := &mockRunner{output: "", err: nil}

	reader := git.NewCommitReader(runner)
	commits, err := reader.GetCommits("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(commits) != 0 {
		t.Errorf("expected 0 commits, got %d", len(commits))
	}
}

func TestGetCommitsSinceTag(t *testing.T) {
	mockOutput := "abc123|feat: add feature|John Doe|1706745600"
	argsReceived := []string{}

	customRunner := &mockRunnerWithArgs{
		output: mockOutput,
		err:    nil,
		onRun: func(args ...string) {
			argsReceived = args
		},
	}

	reader := git.NewCommitReader(customRunner)
	_, err := reader.GetCommits("v1.0.0")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, arg := range argsReceived {
		if arg == "v1.0.0..HEAD" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected 'v1.0.0..HEAD' in git args, got %v", argsReceived)
	}
}

type mockRunnerWithArgs struct {
	output string
	err    error
	onRun  func(args ...string)
}

func (m *mockRunnerWithArgs) Run(args ...string) (string, error) {
	if m.onRun != nil {
		m.onRun(args...)
	}
	return m.output, m.err
}

func TestGetCommitsGitError(t *testing.T) {
	gitError := errors.New("fatal: not a git repository")
	runner := &mockRunner{output: "", err: gitError}

	reader := git.NewCommitReader(runner)
	_, err := reader.GetCommits("")

	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, gitError) {
		t.Errorf("expected git error, got %v", err)
	}
}

func TestGetCommitsNotARepo(t *testing.T) {
	gitError := errors.New("fatal: not a git repository (or any of the parent directories): .git")
	runner := &mockRunner{output: "", err: gitError}

	reader := git.NewCommitReader(runner)
	_, err := reader.GetCommits("")

	if err == nil {
		t.Error("expected error for not a repo, got nil")
	}
}

func TestGetCommitsSkipsInvalidLines(t *testing.T) {
	mockOutput := "abc123|feat: valid commit|John Doe|1706745600\ninvalid line without pipes\ndef456|fix: another valid|Jane Doe|1706746600"
	runner := &mockRunner{output: mockOutput, err: nil}

	reader := git.NewCommitReader(runner)
	commits, err := reader.GetCommits("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(commits) != 2 {
		t.Errorf("expected 2 valid commits, got %d", len(commits))
	}
}
