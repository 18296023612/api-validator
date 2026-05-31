package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Server.Addr = ":8080"
	cfg.Admin.Addr = ":8081"

	res := ValidateConfig(cfg)
	if !res.Valid {
		t.Fatalf("expected valid, got errors: %v", res.Errors)
	}
}

func TestValidateConfig_MissingFields(t *testing.T) {
	cfg := &Config{}
	res := ValidateConfig(cfg)
	if res.Valid {
		t.Fatal("expected invalid")
	}
	if len(res.Errors) == 0 {
		t.Fatal("expected errors")
	}
}

func TestValidateConfig_AuthEnabledNoKeys(t *testing.T) {
	cfg := &Config{}
	cfg.Server.Addr = ":8080"
	cfg.Admin.Addr = ":8081"
	cfg.Auth.Enabled = true

	res := ValidateConfig(cfg)
	if res.Valid {
		t.Fatal("expected invalid when auth enabled but no keys")
	}
}

func TestValidateConfig_Providers(t *testing.T) {
	cfg := &Config{}
	cfg.Server.Addr = ":8080"
	cfg.Admin.Addr = ":8081"

	// Valid provider
	cfg.Providers = append(cfg.Providers, ProviderConfig{
		Name:    "valid", Type: "openai",
		BaseURL: "https://api.example.com", APIKey: "sk-xxx",
		Models: []string{"gpt-4"},
	})
	// Provider missing base_url
	cfg.Providers = append(cfg.Providers, ProviderConfig{
		Name: "nourl", Type: "openai",
	})

	res := ValidateConfig(cfg)
	if res.Valid {
		t.Fatal("expected invalid due to provider without base_url")
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	content := `server:
  addr: ":9090"
admin:
  addr: ":9091"
`
	os.WriteFile(path, []byte(content), 0644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Server.Addr != ":9090" {
		t.Fatalf("expected :9090, got %s", cfg.Server.Addr)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	content := `server:
  addr: ":8080"
admin:
  addr: ":8081"
`
	os.WriteFile(path, []byte(content), 0644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.CORS.AllowOrigin != "*" {
		t.Fatalf("expected default CORS *, got %s", cfg.CORS.AllowOrigin)
	}
	if cfg.Billing.Currency != "CNY" {
		t.Fatalf("expected default currency CNY, got %s", cfg.Billing.Currency)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadProvidersFromDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "provider.yaml"), []byte("name: test\ntype: mock\nmodels: [\"m\"]"), 0644)

	cfgs, err := LoadProvidersFromDir(dir)
	if err != nil {
		t.Fatalf("LoadProvidersFromDir failed: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(cfgs))
	}
}

func TestLoadProvidersFromDir_SkipHidden(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "_hidden.yaml"), []byte("name: hidden"), 0644)
	os.WriteFile(filepath.Join(dir, "visible.yaml"), []byte("name: visible\ntype: mock\nmodels: [\"m\"]"), 0644)

	cfgs, err := LoadProvidersFromDir(dir)
	if err != nil {
		t.Fatalf("LoadProvidersFromDir failed: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("expected 1, got %d", len(cfgs))
	}
}

func TestLoadProvidersFromDir_MultiDoc(t *testing.T) {
	dir := t.TempDir()
	content := `name: provider-a
type: openai
base_url: https://a.com
models: ["a"]
---
name: provider-b
type: openai
base_url: https://b.com
models: ["b"]
`
	os.WriteFile(filepath.Join(dir, "multi.yaml"), []byte(content), 0644)

	cfgs, err := LoadProvidersFromDir(dir)
	if err != nil {
		t.Fatalf("LoadProvidersFromDir failed: %v", err)
	}
	if len(cfgs) != 2 {
		t.Fatalf("expected 2, got %d", len(cfgs))
	}
}

func TestEnvVarExpansion(t *testing.T) {
	os.Setenv("TEST_KEY", "sk-test-val")
	defer os.Unsetenv("TEST_KEY")

	dir := t.TempDir()
	content := `name: env-test
type: openai
base_url: https://api.example.com
api_key: ${TEST_KEY}
models: ["m"]
`
	os.WriteFile(filepath.Join(dir, "env.yaml"), []byte(content), 0644)

	cfgs, err := LoadProvidersFromDir(dir)
	if err != nil {
		t.Fatalf("LoadProvidersFromDir failed: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("expected 1, got %d", len(cfgs))
	}
	if cfgs[0].APIKey != "sk-test-val" {
		t.Fatalf("expected sk-test-val, got %s", cfgs[0].APIKey)
	}
}
