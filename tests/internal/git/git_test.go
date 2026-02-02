package git_test

import (
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
