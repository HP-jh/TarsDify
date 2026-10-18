package config

import (
	"os"
	"testing"
)

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("TARS_MCP_ENDPOINT", "https://dify.example.com")
	os.Setenv("TARS_MCP_API_KEY", "key123")
	defer os.Unsetenv("TARS_MCP_ENDPOINT")
	defer os.Unsetenv("TARS_MCP_API_KEY")
	cfg, err := Load("")
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if cfg.Endpoint != "https://dify.example.com" { t.Fatalf("endpoint mismatch") }
	if cfg.APIKey != "key123" { t.Fatalf("apikey mismatch") }
}

func TestLoad_FromFile(t *testing.T) {
	f, _ := os.CreateTemp("", "config*.json")
	defer os.Remove(f.Name())
	f.WriteString(`{"endpoint":"https://file.example.com","timeout":60}`)
	f.Close()
	cfg, err := Load(f.Name())
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if cfg.Endpoint != "https://file.example.com" { t.Fatalf("endpoint mismatch: %s", cfg.Endpoint) }
	if cfg.Timeout != 60 { t.Fatalf("timeout mismatch: %d", cfg.Timeout) }
}

func TestLoad_MissingEndpoint(t *testing.T) {
	os.Unsetenv("TARS_MCP_ENDPOINT")
	os.Unsetenv("DIFY_ENDPOINT")
	_, err := Load("")
	if err == nil { t.Fatal("expected error for missing endpoint") }
}

func TestLoadForEnv(t *testing.T) {
	os.Setenv("TARS_MCP_ENDPOINT", "https://env.example.com")
	defer os.Unsetenv("TARS_MCP_ENDPOINT")
	cfg, err := LoadForEnv("dev")
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if cfg.Endpoint != "https://env.example.com" { t.Fatalf("endpoint mismatch") }
}
