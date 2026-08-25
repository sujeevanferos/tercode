// Package openai implements the generic OpenAI-compatible provider adapter.
// It works with any API that conforms to the OpenAI Chat Completions API
// (OpenAI, Ollama, vLLM, Mistral, Together, local servers, etc.).
package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sujeevanferos/tercode/internal/provider"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultTimeout = 120 * time.Second
)

// Adapter implements provider.Provider for any OpenAI-compatible API.
type Adapter struct {
	id      string
	name    string
	baseURL string
	apiKey  string
	client  *http.Client
	// modelCache caches the last list of models to avoid repeated discovery calls.
	modelCache []provider.Model
}

// New creates a new OpenAI-compatible adapter.
// baseURL defaults to the OpenAI API if empty.
func New(id, name, baseURL, apiKey string) *Adapter {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return &Adapter{
		id:      id,
		name:    name,
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: defaultTimeout},
	}
}

func (a *Adapter) ID() string   { return a.id }
func (a *Adapter) Name() string { return a.name }

// Authenticate calls the models endpoint to verify credentials.
func (a *Adapter) Authenticate(ctx context.Context) error {
	req, err := a.newRequest(ctx, http.MethodGet, "/models", nil)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("openai: authenticate: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("openai: authentication failed (401) — check API key")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("openai: authenticate: status %d", resp.StatusCode)
	}
	return nil
}

// ListModels fetches available models from the API.
func (a *Adapter) ListModels(ctx context.Context) ([]provider.Model, error) {
	req, err := a.newRequest(ctx, http.MethodGet, "/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai: list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: list models: status %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("openai: list models: decode: %w", err)
	}

	models := make([]provider.Model, 0, len(result.Data))
	for _, m := range result.Data {
		models = append(models, provider.Model{
			ProviderID: a.id,
			ID:         m.ID,
			Name:       m.ID,
			Caps:       provider.Capabilities{Streaming: true, ToolCalling: true},
		})
	}
	a.modelCache = models
	return models, nil
}

// Chat sends a non-streaming chat completion request.
func (a *Adapter) Chat(ctx context.Context, req provider.ChatRequest) (provider.ChatResponse, error) {
	body, err := a.buildRequestBody(req, false)
	if err != nil {
		return provider.ChatResponse{}, err
	}

	httpReq, err := a.newRequest(ctx, http.MethodPost, "/chat/completions", body)
	if err != nil {
		return provider.ChatResponse{}, err
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return provider.ChatResponse{}, fmt.Errorf("openai: chat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return provider.ChatResponse{}, fmt.Errorf("openai: chat: status %d: %s", resp.StatusCode, raw)
	}

	var result chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return provider.ChatResponse{}, fmt.Errorf("openai: chat: decode: %w", err)
	}

	return result.toProviderResponse(a.id), nil
}

// Stream sends a streaming chat completion request.
func (a *Adapter) Stream(ctx context.Context, req provider.ChatRequest) (provider.StreamReader, error) {
	body, err := a.buildRequestBody(req, true)
	if err != nil {
		return nil, err
	}

	httpReq, err := a.newRequest(ctx, http.MethodPost, "/chat/completions", body)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openai: stream: status %d: %s", resp.StatusCode, raw)
	}

	return &streamReader{body: resp.Body, scanner: bufio.NewScanner(resp.Body)}, nil
}

// TestModel runs a validation test sequence.
func (a *Adapter) TestModel(ctx context.Context, modelID string) provider.TestResult {
	result := provider.TestResult{
		ModelID: modelID,
		Stages:  make(map[string]provider.StageResult),
	}
	start := time.Now()

	// Stage 1: authentication
	if err := a.Authenticate(ctx); err != nil {
		result.Stages["authenticate"] = provider.StageResult{Passed: false, Error: err.Error()}
		result.LatencyMs = time.Since(start).Milliseconds()
		return result
	}
	result.Stages["authenticate"] = provider.StageResult{Passed: true}

	// Stage 2: simple inference
	resp, err := a.Chat(ctx, provider.ChatRequest{
		Model:    modelID,
		Messages: []provider.Message{{Role: provider.RoleUser, Content: "Reply with 'ok'."}},
	})
	if err != nil {
		result.Stages["inference"] = provider.StageResult{Passed: false, Error: err.Error()}
		result.LatencyMs = time.Since(start).Milliseconds()
		return result
	}
	result.Stages["inference"] = provider.StageResult{Passed: resp.Message.Content != ""}

	// Stage 3: streaming
	sr, err := a.Stream(ctx, provider.ChatRequest{
		Model:    modelID,
		Messages: []provider.Message{{Role: provider.RoleUser, Content: "Reply with 'ok'."}},
	})
	if err != nil {
		result.Stages["streaming"] = provider.StageResult{Passed: false, Error: err.Error()}
	} else {
		defer sr.Close()
		result.Stages["streaming"] = provider.StageResult{Passed: true}
	}

	result.LatencyMs = time.Since(start).Milliseconds()
	result.Overall = allPassed(result.Stages)
	return result
}

// Capabilities returns capability metadata for a model.
func (a *Adapter) Capabilities(modelID string) provider.Capabilities {
	// Generic OpenAI-compatible defaults.
	return provider.Capabilities{
		Streaming:   true,
		ToolCalling: true,
	}
}

// --- internal helpers ---

func (a *Adapter) newRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (a *Adapter) buildRequestBody(req provider.ChatRequest, stream bool) ([]byte, error) {
	msgs := make([]map[string]any, 0, len(req.Messages))
	for _, m := range req.Messages {
		entry := map[string]any{
			"role":    string(m.Role),
			"content": m.Content,
		}
		if m.ToolCallID != "" {
			entry["tool_call_id"] = m.ToolCallID
		}
		if m.ToolName != "" {
			entry["name"] = m.ToolName
		}
		msgs = append(msgs, entry)
	}

	payload := map[string]any{
		"model":    req.Model,
		"messages": msgs,
		"stream":   stream,
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	if req.Temperature >= 0 {
		payload["temperature"] = req.Temperature
	}

	if len(req.Tools) > 0 {
		tools := make([]map[string]any, 0, len(req.Tools))
		for _, t := range req.Tools {
			tools = append(tools, map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        t.Name,
					"description": t.Description,
					"parameters":  t.Parameters,
				},
			})
		}
		payload["tools"] = tools
		payload["tool_choice"] = "auto"
	}

	return json.Marshal(payload)
}

// chatCompletionResponse is the API response structure for non-streaming calls.
type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

func (r *chatCompletionResponse) toProviderResponse(providerID string) provider.ChatResponse {
	if len(r.Choices) == 0 {
		return provider.ChatResponse{}
	}
	choice := r.Choices[0]
	msg := provider.Message{
		Role:    provider.RoleAssistant,
		Content: choice.Message.Content,
	}
	for _, tc := range choice.Message.ToolCalls {
		msg.ToolCalls = append(msg.ToolCalls, provider.ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}
	return provider.ChatResponse{
		Message:      msg,
		FinishReason: choice.FinishReason,
		PromptTokens: r.Usage.PromptTokens,
		OutputTokens: r.Usage.CompletionTokens,
	}
}

// streamReader wraps an HTTP response body and parses SSE lines.
type streamReader struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
}

func (s *streamReader) Read(p []byte) (int, error) { return s.body.Read(p) }
func (s *streamReader) Close() error               { return s.body.Close() }

func (s *streamReader) Event() (provider.StreamEvent, error) {
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return provider.StreamEvent{Done: true}, io.EOF
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   string `json:"content"`
					ToolCalls []struct {
						ID       string `json:"id"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		choice := chunk.Choices[0]
		ev := provider.StreamEvent{
			Delta:        choice.Delta.Content,
			FinishReason: choice.FinishReason,
			Done:         choice.FinishReason != "",
		}
		if len(choice.Delta.ToolCalls) > 0 {
			tc := choice.Delta.ToolCalls[0]
			ev.ToolCall = &provider.ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			}
		}
		return ev, nil
	}
	if err := s.scanner.Err(); err != nil {
		return provider.StreamEvent{}, err
	}
	return provider.StreamEvent{Done: true}, io.EOF
}

func allPassed(stages map[string]provider.StageResult) bool {
	for _, s := range stages {
		if !s.Passed {
			return false
		}
	}
	return true
}
