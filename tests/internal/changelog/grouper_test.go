package changelog_test

import (
	"testing"
	"time"

	"github.com/lucasbrogni/ai-changelog/internal/changelog"
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

func TestCategoryDisplayNames(t *testing.T) {
	tests := []struct {
		prefix      string
		displayName string
	}{
		{"feat", "New Features"},
		{"fix", "Bug Fixes"},
		{"perf", "Performance"},
		{"docs", "Documentation"},
		{"refactor", "Internal Changes"},
		{"chore", "Maintenance"},
		{"test", "Testing"},
		{"style", "Style"},
		{"other", "Other"},
	}

	for _, tc := range tests {
		t.Run(tc.prefix, func(t *testing.T) {
			displayName := changelog.GetDisplayName(tc.prefix)
			if displayName != tc.displayName {
				t.Errorf("GetDisplayName(%q) = %q, want %q", tc.prefix, displayName, tc.displayName)
			}
		})
	}
}

func TestChangelogSectionStruct(t *testing.T) {
	commit := git.Commit{
		Hash:      "abc123",
		Subject:   "feat: add new feature",
		Author:    "Test Author",
		Timestamp: time.Now(),
		Prefix:    "feat",
	}

	section := changelog.ChangelogSection{
		Title:   "New Features",
		Commits: []git.Commit{commit},
	}

	if section.Title != "New Features" {
		t.Errorf("expected Title to be 'New Features', got '%s'", section.Title)
	}

	if len(section.Commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(section.Commits))
	}

	if section.Commits[0].Hash != "abc123" {
		t.Errorf("expected commit hash 'abc123', got '%s'", section.Commits[0].Hash)
	}
}