package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, homeDir, content string) {
	t.Helper()
	dir := filepath.Join(homeDir, ".config", "gitbucket")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLoad_File(t *testing.T) {
	home := t.TempDir()
	writeConfig(t, home, "url: https://gitbucket.example.com/\ntoken: TKN\n")
	t.Setenv("GITBUCKET_URL", "")
	t.Setenv("GITBUCKET_TOKEN", "")

	cfg, err := Load(home)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.URL != "https://gitbucket.example.com" {
		t.Errorf("URL = %q (should be trimmed)", cfg.URL)
	}
	if cfg.Token != "TKN" {
		t.Errorf("Token = %q", cfg.Token)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	home := t.TempDir()
	writeConfig(t, home, "url: https://from-file.example/\ntoken: FILE\n")
	t.Setenv("GITBUCKET_URL", "https://from-env.example")
	t.Setenv("GITBUCKET_TOKEN", "ENV")

	cfg, err := Load(home)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.URL != "https://from-env.example" || cfg.Token != "ENV" {
		t.Errorf("cfg = %+v", cfg)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GITBUCKET_URL", "")
	t.Setenv("GITBUCKET_TOKEN", "")
	_, err := Load(home)
	if !errors.Is(err, ErrNotConfigured) {
		t.Errorf("err = %v, want ErrNotConfigured", err)
	}
}

func TestLoad_PartialFile(t *testing.T) {
	home := t.TempDir()
	writeConfig(t, home, "url: https://gitbucket.example.com\n")
	t.Setenv("GITBUCKET_URL", "")
	t.Setenv("GITBUCKET_TOKEN", "")
	_, err := Load(home)
	if !errors.Is(err, ErrNotConfigured) {
		t.Errorf("err = %v, want ErrNotConfigured", err)
	}
}

func TestLoad_EnvOnlyPartial(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GITBUCKET_URL", "https://x.example")
	t.Setenv("GITBUCKET_TOKEN", "")
	_, err := Load(home)
	if !errors.Is(err, ErrNotConfigured) {
		t.Errorf("err = %v", err)
	}
}
