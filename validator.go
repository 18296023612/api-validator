package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ============================================================
// Types — standalone copy to avoid importing relay-server
// ============================================================

type ProviderConfig struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type"`
	BaseURL     string   `yaml:"base_url" json:"base_url"`
	APIKey      string   `yaml:"api_key" json:"api_key,omitempty"`
	Models      []string `yaml:"models" json:"models"`
	Weight      int      `yaml:"weight" json:"weight"`
	Priority    int      `yaml:"priority" json:"priority"`
	TargetModel string   `yaml:"target_model,omitempty" json:"target_model,omitempty"`
	Group       string   `yaml:"group,omitempty" json:"group,omitempty"`
	Enabled     bool     `yaml:"enabled" json:"enabled"`
}

type Config struct {
	Server struct {
		Addr    string `yaml:"addr"`
		Timeout string `yaml:"timeout"`
	} `yaml:"server"`
	Admin struct {
		Addr    string `yaml:"addr"`
		Timeout string `yaml:"timeout"`
	} `yaml:"admin"`
	Auth struct {
		Enabled bool     `yaml:"enabled"`
		Keys    []string `yaml:"keys"`
	} `yaml:"auth"`
	RateLimit struct {
		Enabled  bool    `yaml:"enabled"`
		Rate     float64 `yaml:"rate"`
		Capacity int     `yaml:"capacity"`
	} `yaml:"rate_limit"`
	CORS struct {
		AllowOrigin string `yaml:"allow_origin"`
	} `yaml:"cors"`
	Billing struct {
		Enabled     bool   `yaml:"enabled"`
		Currency    string `yaml:"currency"`
		PricingPath string `yaml:"pricing_path"`
		DataDir     string `yaml:"data_dir"`
	} `yaml:"billing"`
	Providers    []ProviderConfig `yaml:"providers"`
	ProvidersDir string           `yaml:"providers_dir"`
}

// ============================================================
// Validation
// ============================================================

type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

func ValidateConfig(cfg *Config) *ValidationResult {
	res := &ValidationResult{Valid: true}

	if cfg.Server.Addr == "" {
		res.Errors = append(res.Errors, "server.addr is required")
	}
	if cfg.Admin.Addr == "" {
		res.Errors = append(res.Errors, "admin.addr is required")
	}
	if cfg.Auth.Enabled && len(cfg.Auth.Keys) == 0 {
		res.Errors = append(res.Errors, "auth.enabled is true but no auth.keys configured")
	}

	for i, p := range cfg.Providers {
		prefix := fmt.Sprintf("providers[%d]", i)
		if p.Name == "" {
			res.Errors = append(res.Errors, prefix+".name is required")
		}
		if p.Type == "" {
			res.Warnings = append(res.Warnings, prefix+": type not set, defaulting to openai")
		}
		if p.BaseURL == "" && p.Type != "mock" {
			res.Errors = append(res.Errors, prefix+"("+p.Name+"): base_url is required for type "+p.Type)
		}
		if p.APIKey == "" && p.Type != "mock" {
			res.Warnings = append(res.Warnings, prefix+"("+p.Name+"): no api_key set")
		}
		if len(p.Models) == 0 {
			res.Warnings = append(res.Warnings, prefix+"("+p.Name+"): no models configured")
		}
	}

	if len(res.Errors) > 0 {
		res.Valid = false
	}
	return res
}

// ============================================================
// File loading
// ============================================================

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	content := os.Expand(string(data), os.Getenv)

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Defaults
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.Admin.Addr == "" {
		cfg.Admin.Addr = ":8081"
	}
	if cfg.CORS.AllowOrigin == "" {
		cfg.CORS.AllowOrigin = "*"
	}
	if cfg.Billing.Currency == "" {
		cfg.Billing.Currency = "CNY"
	}

	return &cfg, nil
}

func LoadProvidersFromDir(dir string) ([]ProviderConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var cfgs []ProviderConfig
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		if strings.HasPrefix(name, "_") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}

		docs := strings.Split(string(data), "\n---\n")
		for _, doc := range docs {
			if strings.TrimSpace(doc) == "" {
				continue
			}
			var pc ProviderConfig
			if err := yaml.Unmarshal([]byte(doc), &pc); err != nil {
				continue
			}
			if pc.Name == "" {
				pc.Name = strings.TrimSuffix(name, filepath.Ext(name))
			}
			pc.APIKey = os.Expand(pc.APIKey, os.Getenv)
			pc.Enabled = true
			cfgs = append(cfgs, pc)
		}
	}
	return cfgs, nil
}

// ============================================================
// Display helpers
// ============================================================

func color(s string, code string) string {
	if noColor {
		return s
	}
	return "\033[" + code + "m" + s + "\033[0m"
}

func green(s string) string  { return color(s, "32") }
func red(s string) string    { return color(s, "31") }
func yellow(s string) string { return color(s, "33") }
func bold(s string) string   { return color(s, "1") }
func dim(s string) string    { return color(s, "2") }

var noColor bool

func printResult(res *ValidationResult, path string) {
	fmt.Printf("\n%s %s\n", bold("🔍 Config:"), path)
	fmt.Printf("  %s %s\n", icon(res.Valid), statusText(res.Valid))
	fmt.Println()

	if len(res.Errors) > 0 {
		fmt.Printf("  %s %s\n", red("✖"), bold("Errors:"))
		for _, e := range res.Errors {
			fmt.Printf("    • %s\n", red(e))
		}
		fmt.Println()
	}

	if len(res.Warnings) > 0 {
		fmt.Printf("  %s %s\n", yellow("⚠"), bold("Warnings:"))
		for _, w := range res.Warnings {
			fmt.Printf("    • %s\n", yellow(w))
		}
		fmt.Println()
	}

	if res.Valid && len(res.Warnings) == 0 {
		fmt.Printf("  %s No issues found. Config looks good!\n", green("✓"))
		fmt.Println()
	}
}

func icon(valid bool) string {
	if valid {
		return green("✓")
	}
	return red("✖")
}

func statusText(valid bool) string {
	if valid {
		return green("PASS")
	}
	return red("FAIL")
}
