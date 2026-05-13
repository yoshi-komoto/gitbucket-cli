package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrNotConfigured = errors.New("not configured")

type Config struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	CACert string `yaml:"ca_cert,omitempty"`
}

func Load(homeDir string) (*Config, error) {
	envURL := os.Getenv("GITBUCKET_URL")
	envTok := os.Getenv("GITBUCKET_TOKEN")
	envCA := os.Getenv("GITBUCKET_CA_CERT")

	if envURL != "" && envTok != "" {
		return &Config{
			URL:    normalizeURL(envURL),
			Token:  envTok,
			CACert: envCA,
		}, nil
	}

	path := filepath.Join(homeDir, ".config", "gitbucket", "config.yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s does not exist", ErrNotConfigured, path)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.URL == "" || cfg.Token == "" {
		return nil, fmt.Errorf("%w: url or token is missing in %s", ErrNotConfigured, path)
	}
	cfg.URL = normalizeURL(cfg.URL)
	if cfg.CACert != "" && !filepath.IsAbs(cfg.CACert) {
		cfg.CACert = filepath.Join(filepath.Dir(path), cfg.CACert)
	}
	if envCA != "" {
		cfg.CACert = envCA
	}
	return &cfg, nil
}

func normalizeURL(u string) string {
	return strings.TrimRight(strings.TrimSpace(u), "/")
}
