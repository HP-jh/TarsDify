package upstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/HP-jh/TarsDify/pkg/retry"
)

type Client struct {
	endpoint string
	apiKey   string
	client   *http.Client
	retry    *retry.Strategy
}

func NewClient(endpoint, apiKey string, timeout int) *Client {
	if endpoint == "" { endpoint = "https://api.example.com" }
	if timeout <= 0 { timeout = 30 }
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
	}
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		apiKey:   apiKey,
		client:   &http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: transport},
		retry:    retry.Default(),
	}
}

func (c *Client) Health(ctx context.Context) map[string]interface{} {
	res, err := c.InvokeRaw(ctx, "GET", "/health", nil)
	if err != nil { return map[string]interface{}{ "healthy": false, "error": err.Error() } }
	res["healthy"] = true
	return res
}

func (c *Client) ListModels(ctx context.Context) ([]interface{}, error) {
	res, err := c.InvokeRaw(ctx, "GET", "/v1/models", nil)
	if err != nil { return nil, err }
	models, _ := res["data"].([]interface{})
	if models == nil { models = []interface{}{} }
	return models, nil
}

func (c *Client) InvokeRaw(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.retry.Do(ctx, func() error {
		var bodyR io.Reader
		if body != nil {
			b, _ := json.Marshal(body)
			bodyR = bytes.NewReader(b)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, bodyR)
		if err != nil { return err }
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}
		if body != nil { req.Header.Set("Content-Type", "application/json") }
		resp, err := c.client.Do(req)
		if err != nil { return err }
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(b))
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			result = map[string]interface{}{ "raw": "non-JSON response" }
		}
		return nil
	})
	if err != nil { return nil, err }
	return result, nil
}
