package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HP-jh/TarsDify/pkg/upstream"
)

func TestHealthHandler(t *testing.T) {
	client := upstream.NewClient("", "", 5)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req, client)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if resp["status"] != "ok" && resp["status"] != "degraded" { t.Fatalf("expected ok, got %v", resp["status"]) }
}

func TestToolsHandler(t *testing.T) {
	client := upstream.NewClient("", "", 5)
	req := httptest.NewRequest(http.MethodGet, "/tools", nil)
	w := httptest.NewRecorder()
	toolsHandler(w, req, client)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	tools, ok := resp["tools"].([]interface{})
	if !ok || len(tools) == 0 { t.Fatal("expected tools") }
}

func TestInvokeHandler_MethodNotAllowed(t *testing.T) {
	client := upstream.NewClient("", "", 5)
	req := httptest.NewRequest(http.MethodGet, "/invoke", nil)
	w := httptest.NewRecorder()
	invokeHandler(w, req, client)
	if w.Code != http.StatusMethodNotAllowed { t.Fatalf("expected 405, got %d", w.Code) }
}

func TestInvokeHandler_UnknownTool(t *testing.T) {
	client := upstream.NewClient("", "", 5)
	body, _ := json.Marshal(map[string]interface{}{ "tool": "xxx", "params": map[string]interface{}{} })
	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewReader(body))
	w := httptest.NewRecorder()
	invokeHandler(w, req, client)
	if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestMetricsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	metricsHandler(w, req)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if resp["name"] != RepoName { t.Fatalf("expected name %s, got %v", RepoName, resp["name"]) }
}

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	versionHandler(w, req)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if resp["version"] != Version { t.Fatalf("expected version %s, got %v", Version, resp["version"]) }
}

func TestUpdateHandler(t *testing.T) {
	client := upstream.NewClient("", "", 5)
	req := httptest.NewRequest(http.MethodGet, "/update", nil)
	w := httptest.NewRecorder()
	updateHandler(w, req, client)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if resp["current"] != Version { t.Fatalf("expected current version") }
}
