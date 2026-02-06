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

func TestSortByDate(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commits := []git.Commit{
		{Hash: "oldest", Subject: "feat: oldest", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		{Hash: "newest", Subject: "feat: newest", Author: "Alice", Timestamp: baseTime.Add(2 * time.Hour), Prefix: "feat"},
		{Hash: "middle", Subject: "feat: middle", Author: "Alice", Timestamp: baseTime.Add(time.Hour), Prefix: "feat"},
	}

	sorted := changelog.SortByDate(commits)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 commits, got %d", len(sorted))
	}

	if sorted[0].Hash != "newest" {
		t.Errorf("expected first commit to be 'newest', got %q", sorted[0].Hash)
	}

	if sorted[1].Hash != "middle" {
		t.Errorf("expected second commit to be 'middle', got %q", sorted[1].Hash)
	}

	if sorted[2].Hash != "oldest" {
		t.Errorf("expected third commit to be 'oldest', got %q", sorted[2].Hash)
	}
}

func TestSortByDateEmpty(t *testing.T) {
	sorted := changelog.SortByDate([]git.Commit{})
	if len(sorted) != 0 {
		t.Errorf("expected 0 commits, got %d", len(sorted))
	}

	sorted = changelog.SortByDate(nil)
	if sorted != nil {
		t.Errorf("expected nil for nil input, got %v", sorted)
	}
}

func TestSortByDateDoesNotMutateOriginal(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commits := []git.Commit{
		{Hash: "oldest", Timestamp: baseTime},
		{Hash: "newest", Timestamp: baseTime.Add(time.Hour)},
	}

	changelog.SortByDate(commits)

	if commits[0].Hash != "oldest" {
		t.Errorf("original slice was mutated: expected first commit to be 'oldest', got %q", commits[0].Hash)
	}
}

func TestGroupUnknownPrefix(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commits := []git.Commit{
		{Hash: "abc123", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
		{Hash: "def456", Subject: "random commit message", Author: "Bob", Timestamp: baseTime.Add(time.Hour), Prefix: "other"},
		{Hash: "ghi789", Subject: "WIP: work in progress", Author: "Charlie", Timestamp: baseTime.Add(2 * time.Hour), Prefix: "other"},
	}

	sections := changelog.GroupByCategory(commits)

	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}

	if sections[0].Title != "New Features" {
		t.Errorf("expected first section to be 'New Features', got %q", sections[0].Title)
	}

	if sections[1].Title != "Other" {
		t.Errorf("expected last section to be 'Other', got %q", sections[1].Title)
	}

	if len(sections[1].Commits) != 2 {
		t.Errorf("expected 2 commits in 'Other' section, got %d", len(sections[1].Commits))
	}
}

func TestGroupUnknownPrefixAppearsLast(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	commits := []git.Commit{
		{Hash: "abc123", Subject: "random commit", Author: "Alice", Timestamp: baseTime, Prefix: "other"},
		{Hash: "def456", Subject: "fix: bug fix", Author: "Bob", Timestamp: baseTime.Add(time.Hour), Prefix: "fix"},
		{Hash: "ghi789", Subject: "chore: cleanup", Author: "Charlie", Timestamp: baseTime.Add(2 * time.Hour), Prefix: "chore"},
	}

	sections := changelog.GroupByCategory(commits)

	if len(sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(sections))
	}

	lastSection := sections[len(sections)-1]
	if lastSection.Title != "Other" {
		t.Errorf("expected 'Other' section to appear last, but got %q", lastSection.Title)
	}
}