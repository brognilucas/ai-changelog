package ollama_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lucasbrogni/ai-changelog/internal/git"
	"github.com/lucasbrogni/ai-changelog/internal/ollama"
)

type mockOllamaClient struct{}

func (m *mockOllamaClient) HealthCheck() error {
	return nil
}

func (m *mockOllamaClient) SummarizeCommits(commits []git.Commit, model string) ([]string, error) {
	return nil, nil
}

func TestOllamaClientInterface(t *testing.T) {
	var client ollama.Client = &mockOllamaClient{}

	err := client.HealthCheck()
	if err != nil {
		t.Errorf("unexpected error from HealthCheck: %v", err)
	}

	_, err = client.SummarizeCommits([]git.Commit{}, "llama3")
	if err != nil {
		t.Errorf("unexpected error from SummarizeCommits: %v", err)
	}
}

func TestRequestSerialization(t *testing.T) {
	request := ollama.GenerateRequest{
		Model:  "llama3",
		Prompt: "Summarize these commits",
		Stream: false,
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var decoded ollama.GenerateRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if decoded.Model != request.Model {
		t.Errorf("expected Model %q, got %q", request.Model, decoded.Model)
	}

	if decoded.Prompt != request.Prompt {
		t.Errorf("expected Prompt %q, got %q", request.Prompt, decoded.Prompt)
	}

	if decoded.Stream != request.Stream {
		t.Errorf("expected Stream %v, got %v", request.Stream, decoded.Stream)
	}
}

func TestResponseSerialization(t *testing.T) {
	jsonData := `{"model":"llama3","response":"Summary of changes","done":true}`

	var response ollama.GenerateResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Model != "llama3" {
		t.Errorf("expected Model 'llama3', got %q", response.Model)
	}

	if response.Response != "Summary of changes" {
		t.Errorf("expected Response 'Summary of changes', got %q", response.Response)
	}

	if response.Done != true {
		t.Errorf("expected Done true, got %v", response.Done)
	}
}

func TestHealthCheckSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := ollama.NewDefaultClient(server.URL)
	err := client.HealthCheck()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHealthCheckFail(t *testing.T) {
	client := ollama.NewDefaultClient("http://localhost:99999")
	err := client.HealthCheck()

	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}

func TestBuildPrompt(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc123", Subject: "feat: add user authentication", Author: "dev"},
		{Hash: "def456", Subject: "fix: resolve login timeout", Author: "dev"},
	}

	prompt := ollama.BuildPrompt(commits)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
	if !strings.Contains(prompt, "feat: add user authentication") {
		t.Error("prompt should contain commit subjects")
	}
	if !strings.Contains(prompt, "fix: resolve login timeout") {
		t.Error("prompt should contain all commit subjects")
	}
}

func TestBuildPromptEmpty(t *testing.T) {
	commits := []git.Commit{}

	prompt := ollama.BuildPrompt(commits)

	if prompt != "" {
		t.Errorf("expected empty string for empty commits, got %q", prompt)
	}
}

func TestGenerateSuccess(t *testing.T) {
	expectedResponse := "## Changelog\n- Added user authentication"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("expected path /api/generate, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		var req ollama.GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Model != "llama3" {
			t.Errorf("expected model llama3, got %s", req.Model)
		}

		if req.Stream != false {
			t.Errorf("expected stream false, got %v", req.Stream)
		}

		response := ollama.GenerateResponse{
			Model:    "llama3",
			Response: expectedResponse,
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := ollama.NewDefaultClient(server.URL)
	result, err := client.Generate("llama3", "test prompt")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expectedResponse {
		t.Errorf("expected response %q, got %q", expectedResponse, result)
	}
}

func TestGenerateTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	client := ollama.NewDefaultClientWithTimeout(server.URL, 50*time.Millisecond)
	_, err := client.Generate("llama3", "test prompt")

	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}