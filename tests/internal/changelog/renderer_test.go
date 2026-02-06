package changelog_test

import (
	"strings"
	"testing"
	"time"

	"github.com/lucasbrogni/ai-changelog/internal/changelog"
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type mockRenderer struct{}

func (m *mockRenderer) Render(sections []changelog.ChangelogSection, version string) string {
	return "mock"
}

func TestRendererInterface(t *testing.T) {
	var _ changelog.Renderer = &mockRenderer{}

	mock := &mockRenderer{}
	result := mock.Render(nil, "")
	if result != "mock" {
		t.Errorf("expected 'mock', got %q", result)
	}
}

func TestMarkdownRendererBasic(t *testing.T) {
	var _ changelog.Renderer = &changelog.MarkdownRenderer{}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	sections := []changelog.ChangelogSection{
		{
			Title: "New Features",
			Commits: []git.Commit{
				{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
			},
		},
	}

	renderer := &changelog.MarkdownRenderer{}
	result := renderer.Render(sections, "v1.0.0")

	if result == "" {
		t.Error("expected non-empty output")
	}

	if !strings.Contains(result, "v1.0.0") {
		t.Error("expected output to contain version")
	}

	if !strings.Contains(result, "New Features") {
		t.Error("expected output to contain section title")
	}
}

func TestRenderVersionHeader(t *testing.T) {
	renderer := &changelog.MarkdownRenderer{}

	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"with version", "v1.0.0", "# Changelog v1.0.0"},
		{"without version", "", "# Changelog"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := renderer.Render(nil, tc.version)
			if !strings.HasPrefix(result, tc.expected) {
				t.Errorf("expected output to start with %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestRenderSectionHeaders(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	sections := []changelog.ChangelogSection{
		{
			Title: "New Features",
			Commits: []git.Commit{
				{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
			},
		},
		{
			Title: "Bug Fixes",
			Commits: []git.Commit{
				{Hash: "def4567ghi", Subject: "fix: resolve crash", Author: "Bob", Timestamp: baseTime, Prefix: "fix"},
			},
		},
	}

	renderer := &changelog.MarkdownRenderer{}
	result := renderer.Render(sections, "v1.0.0")

	if !strings.Contains(result, "## New Features") {
		t.Error("expected output to contain '## New Features'")
	}

	if !strings.Contains(result, "## Bug Fixes") {
		t.Error("expected output to contain '## Bug Fixes'")
	}
}

func TestRenderCommitList(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		subject  string
		hash     string
		expected string
	}{
		{
			name:     "strips prefix and truncates hash",
			subject:  "feat: add login",
			hash:     "abc1234def567",
			expected: "- add login (abc1234)",
		},
		{
			name:     "strips scoped prefix",
			subject:  "feat(api): add REST endpoints",
			hash:     "bcd2345efg678",
			expected: "- add REST endpoints (bcd2345)",
		},
		{
			name:     "handles subject without prefix",
			subject:  "random commit message",
			hash:     "cde3456fgh789",
			expected: "- random commit message (cde3456)",
		},
		{
			name:     "handles short hash",
			subject:  "fix: bug",
			hash:     "abc",
			expected: "- bug (abc)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sections := []changelog.ChangelogSection{
				{
					Title: "Test Section",
					Commits: []git.Commit{
						{Hash: tc.hash, Subject: tc.subject, Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
					},
				},
			}

			renderer := &changelog.MarkdownRenderer{}
			result := renderer.Render(sections, "v1.0.0")

			if !strings.Contains(result, tc.expected) {
				t.Errorf("expected output to contain %q, got:\n%s", tc.expected, result)
			}
		})
	}
}

func TestPlainTextRenderer(t *testing.T) {
	var _ changelog.Renderer = &changelog.PlainTextRenderer{}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	sections := []changelog.ChangelogSection{
		{
			Title: "New Features",
			Commits: []git.Commit{
				{Hash: "abc1234def", Subject: "feat: add login", Author: "Alice", Timestamp: baseTime, Prefix: "feat"},
				{Hash: "ghi789jklm", Subject: "feat: add logout", Author: "Alice", Timestamp: baseTime.Add(time.Hour), Prefix: "feat"},
			},
		},
		{
			Title: "Bug Fixes",
			Commits: []git.Commit{
				{Hash: "def456ghij", Subject: "fix: resolve crash", Author: "Bob", Timestamp: baseTime, Prefix: "fix"},
			},
		},
	}

	renderer := &changelog.PlainTextRenderer{}

	t.Run("contains uppercase version header", func(t *testing.T) {
		result := renderer.Render(sections, "v1.0.0")
		if !strings.Contains(result, "CHANGELOG v1.0.0") {
			t.Errorf("expected 'CHANGELOG v1.0.0' in output, got:\n%s", result)
		}
	})

	t.Run("contains underline", func(t *testing.T) {
		result := renderer.Render(sections, "v1.0.0")
		if !strings.Contains(result, "================") {
			t.Errorf("expected underline in output, got:\n%s", result)
		}
	})

	t.Run("contains uppercase section titles", func(t *testing.T) {
		result := renderer.Render(sections, "v1.0.0")
		if !strings.Contains(result, "NEW FEATURES") {
			t.Errorf("expected 'NEW FEATURES' in output, got:\n%s", result)
		}
		if !strings.Contains(result, "BUG FIXES") {
			t.Errorf("expected 'BUG FIXES' in output, got:\n%s", result)
		}
	})

	t.Run("contains indented commit lines", func(t *testing.T) {
		result := renderer.Render(sections, "v1.0.0")
		if !strings.Contains(result, "  * add login (abc1234)") {
			t.Errorf("expected indented commit line in output, got:\n%s", result)
		}
		if !strings.Contains(result, "  * resolve crash (def456g)") {
			t.Errorf("expected indented commit line in output, got:\n%s", result)
		}
	})

	t.Run("version header without version", func(t *testing.T) {
		result := renderer.Render(nil, "")
		if !strings.HasPrefix(result, "CHANGELOG\n") {
			t.Errorf("expected output to start with 'CHANGELOG', got:\n%s", result)
		}
	})
}
