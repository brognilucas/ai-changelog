package changelog

type Renderer interface {
	Render(sections []ChangelogSection, version string) string
}
