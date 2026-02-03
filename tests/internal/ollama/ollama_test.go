package ollama_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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