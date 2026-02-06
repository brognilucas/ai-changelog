package changelog

import (
	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type ChangelogSection struct {
	Title   string
	Commits []git.Commit
}