// Package openrouter implements the OpenRouter provider adapter.
// OpenRouter provides unified access to 200+ models with live pricing,
// capability discovery, and ranking.
package openrouter

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
	defaultBaseURL = "https://openrouter.ai/api/v1"
	defaultTimeout = 120 * time.Second
)

// Adapter implements provider.Provider for OpenRouter.
type Adapter struct {
	id      string
	name    string
	baseURL string
	apiKey  string
	client  *http.Client
}

// New creates an OpenRouter adapter.
func New(apiKey string) *Adapter {
	return &Adapter{
		id:      "openrouter",
		name:    "OpenRouter",
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: defaultTimeout},
	}
}

func (a *Adapter) ID() string   { return a.id }
func (a *Adapter) Name() string { return a.name }

// Authenticate verifies the API key by querying /auth/key or /models.
func (a *Adapter) Authenticate(ctx context.Context) error {
	req, err := a.newRequest(ctx, http.MethodGet, "/auth/key", nil)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("openrouter: authenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("openrouter: authentication failed (401) — check OPENROUTER_API_KEY")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("openrouter: authenticate: status %d", resp.StatusCode)
	}
	return nil
}

// ListModels fetches available models with pricing and context windows.
func (a *Adapter) ListModels(ctx context.Context) ([]provider.Model, error) {
	req, err := a.newRequest(ctx, http.MethodGet, "/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openrouter: list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter: list models: status %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			ContextLength int    `json:"context_length"`
			Pricing       struct {
				Prompt     string `json:"prompt"`
				Completion string `json:"completion"`
			} `json:"pricing"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("openrouter: list models: decode: %w", err)
	}

	models := make([]provider.Model, 0, len(result.Data))
	for _, m := range result.Data {
		models = append(models, provider.Model{
			ProviderID:    a.id,
			ID:            m.ID,
			Name:          m.Name,
			ContextWindow: m.ContextLength,
			Caps: provider.Capabilities{
				Streaming:   true,
				ToolCalling: true,
			},
		})
	}
	return models, nil
}

// Chat sends a non-streaming chat request.
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
		return provider.ChatResponse{}, fmt.Errorf("openrouter: chat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return provider.ChatResponse{}, fmt.Errorf("openrouter: chat: status %d: %s", resp.StatusCode, raw)
	}

	var result struct {
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

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return provider.ChatResponse{}, fmt.Errorf("openrouter: chat: decode: %w", err)
	}

	if len(result.Choices) == 0 {
		return provider.ChatResponse{}, nil
	}

	choice := result.Choices[0]
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
		PromptTokens: result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
	}, nil
}

// Stream sends a streaming chat request.
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
		return nil, fmt.Errorf("openrouter: stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openrouter: stream: status %d: %s", resp.StatusCode, raw)
	}

	return &openrouterStreamReader{body: resp.Body, scanner: bufio.NewScanner(resp.Body)}, nil
}

// TestModel performs model validation.
func (a *Adapter) TestModel(ctx context.Context, modelID string) provider.TestResult {
	result := provider.TestResult{
		ModelID: modelID,
		Stages:  make(map[string]provider.StageResult),
	}
	start := time.Now()

	if err := a.Authenticate(ctx); err != nil {
		result.Stages["authenticate"] = provider.StageResult{Passed: false, Error: err.Error()}
		result.LatencyMs = time.Since(start).Milliseconds()
		return result
	}
	result.Stages["authenticate"] = provider.StageResult{Passed: true}

	resp, err := a.Chat(ctx, provider.ChatRequest{
		Model:    modelID,
		Messages: []provider.Message{{Role: provider.RoleUser, Content: "ping"}},
	})
	if err != nil {
		result.Stages["inference"] = provider.StageResult{Passed: false, Error: err.Error()}
	} else {
		result.Stages["inference"] = provider.StageResult{Passed: resp.Message.Content != ""}
	}

	result.LatencyMs = time.Since(start).Milliseconds()
	result.Overall = result.Stages["authenticate"].Passed && result.Stages["inference"].Passed
	return result
}

func (a *Adapter) Capabilities(modelID string) provider.Capabilities {
	return provider.Capabilities{
		Streaming:   true,
		ToolCalling: true,
	}
}

func (a *Adapter) newRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openrouter: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/sujeevanferos/tercode")
	req.Header.Set("X-Title", "Tercode")
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

type openrouterStreamReader struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
}

func (s *openrouterStreamReader) Read(p []byte) (int, error) { return s.body.Read(p) }
func (s *openrouterStreamReader) Close() error               { return s.body.Close() }

func (s *openrouterStreamReader) Event() (provider.StreamEvent, error) {
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
