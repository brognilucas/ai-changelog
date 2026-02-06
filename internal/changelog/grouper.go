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

type ChangelogSection struct {
	Title   string
	Commits []git.Commit
}