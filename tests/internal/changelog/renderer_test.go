package changelog_test

import (
	"testing"

	"github.com/lucasbrogni/ai-changelog/internal/changelog"
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
