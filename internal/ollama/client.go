package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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

const defaultTimeout = 30 * time.Second

func NewDefaultClient(baseURL string) *DefaultClient {
	return &DefaultClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

func NewDefaultClientWithTimeout(baseURL string, timeout time.Duration) *DefaultClient {
	return &DefaultClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
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

const defaultBatchSize = 10

func (c *DefaultClient) SummarizeCommits(commits []git.Commit, model string) ([]string, error) {
	if len(commits) == 0 {
		return []string{}, nil
	}

	var summaries []string

	for i := 0; i < len(commits); i += defaultBatchSize {
		end := i + defaultBatchSize
		if end > len(commits) {
			end = len(commits)
		}

		batch := commits[i:end]
		prompt := BuildPrompt(batch)

		response, err := c.Generate(model, prompt)
		if err != nil {
			for _, commit := range batch {
				summaries = append(summaries, commit.Subject)
			}
			continue
		}

		summaries = append(summaries, response)
	}

	return summaries, nil
}

func (c *DefaultClient) Generate(model string, prompt string) (string, error) {
	request := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("generate request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
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