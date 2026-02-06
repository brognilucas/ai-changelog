package changelog

import (
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

const (
	CategoryFeat     = "feat"
	CategoryFix      = "fix"
	CategoryPerf     = "perf"
	CategoryDocs     = "docs"
	CategoryRefactor = "refactor"
	CategoryChore    = "chore"
	CategoryTest     = "test"
	CategoryStyle    = "style"
	CategoryOther    = "other"
)

var categoryDisplayNames = map[string]string{
	CategoryFeat:     "New Features",
	CategoryFix:      "Bug Fixes",
	CategoryPerf:     "Performance",
	CategoryDocs:     "Documentation",
	CategoryRefactor: "Internal Changes",
	CategoryChore:    "Maintenance",
	CategoryTest:     "Testing",
	CategoryStyle:    "Style",
	CategoryOther:    "Other",
}

func GetDisplayName(prefix string) string {
	if name, ok := categoryDisplayNames[prefix]; ok {
		return name
	}
	return "Other"
}

var categoryOrder = []string{
	CategoryFeat,
	CategoryFix,
	CategoryPerf,
	CategoryDocs,
	CategoryRefactor,
	CategoryChore,
	CategoryTest,
	CategoryStyle,
	CategoryOther,
}

type ChangelogSection struct {
	Title   string
	Commits []git.Commit
}

func GroupByCategory(commits []git.Commit) []ChangelogSection {
	if len(commits) == 0 {
		return []ChangelogSection{}
	}

	grouped := make(map[string][]git.Commit)
	for _, commit := range commits {
		grouped[commit.Prefix] = append(grouped[commit.Prefix], commit)
	}

	var sections []ChangelogSection
	for _, category := range categoryOrder {
		if commitList, ok := grouped[category]; ok && len(commitList) > 0 {
			sections = append(sections, ChangelogSection{
				Title:   GetDisplayName(category),
				Commits: commitList,
			})
		}
	}

	return sections
}