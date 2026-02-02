package git

import (
	"testing"
	"time"
)

func TestCommitStructFields(t *testing.T) {
	timestamp := time.Now()

	commit := Commit{
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
