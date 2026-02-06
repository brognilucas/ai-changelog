package changelog

import (
	"fmt"
	"strings"

	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type Renderer interface {
	Render(sections []ChangelogSection, version string) string
}

type MarkdownRenderer struct{}

func (r *MarkdownRenderer) Render(sections []ChangelogSection, version string) string {
	var builder strings.Builder

	builder.WriteString(renderMarkdownVersionHeader(version))

	for _, section := range sections {
		if len(section.Commits) == 0 {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(renderMarkdownSection(section))
	}

	return builder.String()
}

func renderMarkdownVersionHeader(version string) string {
	if version == "" {
		return "# Changelog\n"
	}
	return fmt.Sprintf("# Changelog %s\n", version)
}

func renderMarkdownSection(section ChangelogSection) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("## %s\n\n", section.Title))

	for _, commit := range section.Commits {
		builder.WriteString(renderMarkdownCommitLine(commit))
	}

	return builder.String()
}

func renderMarkdownCommitLine(commit git.Commit) string {
	return fmt.Sprintf("- %s (%s)\n", cleanSubject(commit.Subject), shortHash(commit.Hash))
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

type PlainTextRenderer struct{}

func (r *PlainTextRenderer) Render(sections []ChangelogSection, version string) string {
	var builder strings.Builder

	builder.WriteString(renderPlainTextVersionHeader(version))

	for _, section := range sections {
		if len(section.Commits) == 0 {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(renderPlainTextSection(section))
	}

	return builder.String()
}

func renderPlainTextVersionHeader(version string) string {
	var header string
	if version == "" {
		header = "CHANGELOG"
	} else {
		header = fmt.Sprintf("CHANGELOG %s", version)
	}
	return header + "\n" + strings.Repeat("=", len(header)) + "\n"
}

func renderPlainTextSection(section ChangelogSection) string {
	var builder strings.Builder

	builder.WriteString(strings.ToUpper(section.Title) + "\n\n")

	for _, commit := range section.Commits {
		builder.WriteString(renderPlainTextCommitLine(commit))
	}

	return builder.String()
}

func renderPlainTextCommitLine(commit git.Commit) string {
	return fmt.Sprintf("  * %s (%s)\n", cleanSubject(commit.Subject), shortHash(commit.Hash))
}

func cleanSubject(subject string) string {
	colonIndex := strings.Index(subject, ":")
	if colonIndex == -1 {
		return subject
	}

	cleaned := strings.TrimSpace(subject[colonIndex+1:])
	if cleaned == "" {
		return subject
	}

	return cleaned
}
