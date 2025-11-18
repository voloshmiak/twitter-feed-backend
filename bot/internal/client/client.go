package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient[T Sendable] struct {
	endpoint   string
	httpClient *http.Client
}

func NewHTTPClient[T Sendable](endpoint string) *HTTPClient[T] {
	return &HTTPClient[T]{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}
}

func (c *HTTPClient[T]) Send(ctx context.Context, payload T) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
