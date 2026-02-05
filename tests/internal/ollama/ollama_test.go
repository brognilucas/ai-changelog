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

func TestSummarizeCommits(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			w.WriteHeader(http.StatusOK)
			return
		}

		callCount++
		response := ollama.GenerateResponse{
			Model:    "llama3",
			Response: "- Improved user experience with new feature",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	commits := []git.Commit{
		{Hash: "abc123", Subject: "feat: add user authentication", Author: "dev", Timestamp: time.Now()},
		{Hash: "def456", Subject: "fix: resolve login timeout", Author: "dev", Timestamp: time.Now()},
		{Hash: "ghi789", Subject: "docs: update README", Author: "dev", Timestamp: time.Now()},
	}

	client := ollama.NewDefaultClient(server.URL)
	summaries, err := client.SummarizeCommits(commits, "llama3")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(summaries) == 0 {
		t.Error("expected at least one summary")
	}

	if callCount == 0 {
		t.Error("expected at least one API call to Ollama")
	}
}

func TestSummarizeCommitsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := ollama.NewDefaultClient(server.URL)
	summaries, err := client.SummarizeCommits([]git.Commit{}, "llama3")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(summaries) != 0 {
		t.Errorf("expected empty summaries for empty commits, got %d", len(summaries))
	}
}

func TestSummarizeFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	commits := []git.Commit{
		{Hash: "abc123", Subject: "feat: add user authentication", Author: "dev", Timestamp: time.Now()},
		{Hash: "def456", Subject: "fix: resolve login timeout", Author: "dev", Timestamp: time.Now()},
	}

	client := ollama.NewDefaultClient(server.URL)
	summaries, err := client.SummarizeCommits(commits, "llama3")

	if err != nil {
		t.Fatalf("expected no error on fallback, got %v", err)
	}

	if len(summaries) != len(commits) {
		t.Errorf("expected %d fallback summaries, got %d", len(commits), len(summaries))
	}

	for i, summary := range summaries {
		if summary != commits[i].Subject {
			t.Errorf("expected fallback summary %q, got %q", commits[i].Subject, summary)
		}
	}
}

func TestOllamaFallbackBehavior(t *testing.T) {
	t.Run("connection refused falls back gracefully", func(t *testing.T) {
		client := ollama.NewDefaultClient("http://localhost:59999")

		commits := []git.Commit{
			{Hash: "abc123", Subject: "feat: add new feature", Author: "dev", Timestamp: time.Now()},
			{Hash: "def456", Subject: "fix: bug fix", Author: "dev", Timestamp: time.Now()},
		}

		summaries, err := client.SummarizeCommits(commits, "llama3")

		if err != nil {
			t.Fatalf("expected graceful fallback without error, got %v", err)
		}

		if len(summaries) != len(commits) {
			t.Errorf("expected %d fallback summaries, got %d", len(commits), len(summaries))
		}

		for i, summary := range summaries {
			if summary != commits[i].Subject {
				t.Errorf("expected fallback to subject %q, got %q", commits[i].Subject, summary)
			}
		}
	})

	t.Run("timeout during generation falls back gracefully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/generate" {
				time.Sleep(200 * time.Millisecond)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := ollama.NewDefaultClientWithTimeout(server.URL, 50*time.Millisecond)

		commits := []git.Commit{
			{Hash: "abc123", Subject: "feat: add timeout test", Author: "dev", Timestamp: time.Now()},
		}

		summaries, err := client.SummarizeCommits(commits, "llama3")

		if err != nil {
			t.Fatalf("expected graceful fallback on timeout, got %v", err)
		}

		if len(summaries) != len(commits) {
			t.Errorf("expected %d fallback summaries, got %d", len(commits), len(summaries))
		}

		if summaries[0] != commits[0].Subject {
			t.Errorf("expected fallback to subject %q, got %q", commits[0].Subject, summaries[0])
		}
	})

	t.Run("partial batch failures preserve successful summaries", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/generate" {
				callCount++
				if callCount == 1 {
					response := ollama.GenerateResponse{
						Model:    "llama3",
						Response: "AI generated summary for batch 1",
						Done:     true,
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		commits := make([]git.Commit, 15)
		for i := 0; i < 15; i++ {
			commits[i] = git.Commit{
				Hash:      "hash" + string(rune('a'+i)),
				Subject:   "commit " + string(rune('A'+i)),
				Author:    "dev",
				Timestamp: time.Now(),
			}
		}

		client := ollama.NewDefaultClient(server.URL)
		summaries, err := client.SummarizeCommits(commits, "llama3")

		if err != nil {
			t.Fatalf("expected no error on partial failure, got %v", err)
		}

		expectedLen := 6
		if len(summaries) != expectedLen {
			t.Fatalf("expected %d summaries (1 AI + 5 fallback), got %d", expectedLen, len(summaries))
		}

		if summaries[0] != "AI generated summary for batch 1" {
			t.Errorf("expected AI summary for first batch, got %q", summaries[0])
		}

		for i := 1; i < len(summaries); i++ {
			commitIndex := 10 + (i - 1)
			expectedSubject := commits[commitIndex].Subject
			if summaries[i] != expectedSubject {
				t.Errorf("expected fallback subject %q for summaries[%d], got %q", expectedSubject, i, summaries[i])
			}
		}
	})
}