package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const defaultBaseURL = "http://localhost:11435"

func BaseURL() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return defaultBaseURL
}

func TestModel() string {
	if v := os.Getenv("TEST_MODEL"); v != "" {
		return v
	}
	return "Qwen/Qwen3-0.6B-GGUF"
}

func urlFor(path string) string {
	return fmt.Sprintf("%s%s", BaseURL(), path)
}

func Get(path string) (*http.Response, error) {
	return http.Get(urlFor(path))
}

func Post(path string, body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return http.Post(urlFor(path), "application/json", strings.NewReader(string(b)))
}

func Delete(path string, body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, urlFor(path), strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ParseJSON[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// ParseSSEEvents splits a body into SSE data events (strips "data: " prefix).
func ParseSSEEvents(body []byte) []string {
	var events []string
	for line := range strings.SplitSeq(string(body), "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			events = append(events, after)
		}
	}
	return events
}

// --- Response types ---

type HealthResponse struct {
	Status string `json:"status"`
}

type ModelInfo struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	Size       int64  `json:"size"`
	Format     string `json:"format"`
	ModifiedAt string `json:"modified_at,omitempty"`
	ExpiresAt  string `json:"expires_at,omitempty"`
}

type TagsResponse struct {
	Models []ModelInfo `json:"models"`
}

type ShowRequest struct {
	Model string `json:"model"`
}

type ShowResponse struct {
	Modelfile string       `json:"modelfile"`
	Details   *ModelDetail `json:"details"`
}

type ModelDetail struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	Size       int64  `json:"size"`
	Format     string `json:"format"`
	ModifiedAt string `json:"modified_at"`
}

type PullRequest struct {
	Model string `json:"model"`
}

type PullEvent struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

type DeleteRequest struct {
	Model string `json:"model"`
}

type DeleteResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []ChatMessage  `json:"messages"`
	Stream   bool           `json:"stream,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

type ChatResponse struct {
	Model     string       `json:"model"`
	Message   *ChatMessage `json:"message"`
	Done      bool         `json:"done"`
	CreatedAt string       `json:"created_at"`
}

type GenerateRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

type GenerateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

type PSResponse struct {
	Models []ModelInfo `json:"models"`
}

type StopRequest struct {
	Model string `json:"model"`
}

type StopResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
