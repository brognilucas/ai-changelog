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

func TestGroupByCategory(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commits := []git.Commit{
		{Hash: "abc123", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		{Hash: "def456", Subject: "fix: resolve crash", Author: "Bob", Timestamp: baseTime.Add(time.Hour), Prefix: "fix"},
		{Hash: "ghi789", Subject: "feat: add logout", Author: "Alice", Timestamp: baseTime.Add(2 * time.Hour), Prefix: "feat"},
		{Hash: "jkl012", Subject: "docs: update readme", Author: "Charlie", Timestamp: baseTime.Add(3 * time.Hour), Prefix: "docs"},
	}

	sections := changelog.GroupByCategory(commits)

	expectedOrder := []string{"New Features", "Bug Fixes", "Documentation"}
	if len(sections) != len(expectedOrder) {
		t.Fatalf("expected %d sections, got %d", len(expectedOrder), len(sections))
	}

	for i, expectedTitle := range expectedOrder {
		if sections[i].Title != expectedTitle {
			t.Errorf("section[%d].Title = %q, want %q", i, sections[i].Title, expectedTitle)
		}
	}

	if len(sections[0].Commits) != 2 {
		t.Errorf("expected 2 commits in 'New Features', got %d", len(sections[0].Commits))
	}

	if len(sections[1].Commits) != 1 {
		t.Errorf("expected 1 commit in 'Bug Fixes', got %d", len(sections[1].Commits))
	}

	if len(sections[2].Commits) != 1 {
		t.Errorf("expected 1 commit in 'Documentation', got %d", len(sections[2].Commits))
	}
}

func TestGroupByCategoryEmpty(t *testing.T) {
	sections := changelog.GroupByCategory([]git.Commit{})

	if len(sections) != 0 {
		t.Errorf("expected 0 sections for empty input, got %d", len(sections))
	}

	sections = changelog.GroupByCategory(nil)

	if len(sections) != 0 {
		t.Errorf("expected 0 sections for nil input, got %d", len(sections))
	}
}