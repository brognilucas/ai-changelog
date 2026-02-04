package ollama

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lucasbrogni/ai-changelog/internal/git"
)

type Client interface {
	HealthCheck() error
	SummarizeCommits(commits []git.Commit, model string) ([]string, error)
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type DefaultClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDefaultClient(baseURL string) *DefaultClient {
	return &DefaultClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *DefaultClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL)
	if err != nil {
		return fmt.Errorf("ollama not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *DefaultClient) SummarizeCommits(commits []git.Commit, model string) ([]string, error) {
	return nil, nil
}

func BuildPrompt(commits []git.Commit) string {
	if len(commits) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("You are a changelog generator. Summarize the following git commits into clear, user-friendly changelog entries.\n\n")
	builder.WriteString("Commits:\n")

	for _, commit := range commits {
		builder.WriteString(fmt.Sprintf("- %s\n", commit.Subject))
	}

	builder.WriteString("\nGenerate a concise changelog summary grouped by type (features, fixes, etc.).")

	return builder.String()
}