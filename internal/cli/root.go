package cli

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yoshi-komoto/gitbucket-cli/internal/config"
	"github.com/yoshi-komoto/gitbucket-cli/internal/gitbucket"
	"github.com/yoshi-komoto/gitbucket-cli/internal/repo"
)

// rootFlags holds parsed global flags.
type rootFlags struct {
	repo   string
	output string
}

var flags rootFlags

var rootCmd = &cobra.Command{
	Use:           "gb",
	Short:         "GitBucket CLI",
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.repo, "repo", "", "Repository in OWNER/REPO format (auto-detected from git remote if omitted)")
	rootCmd.PersistentFlags().StringVar(&flags.output, "output", "table", "Output format: table|json")
}

// Execute is the entrypoint called from main.
func Execute() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	rootCmd.SetContext(ctx)
	err := rootCmd.Execute()
	if err != nil {
		printError(os.Stderr, err)
	}
	return exitCodeFor(err)
}

// session bundles per-command dependencies.
type session struct {
	cfg    *config.Config
	client *gitbucket.Client
	owner  string
	repo   string
}

// newSession resolves config, repo and constructs the API client.
func newSession(ctx context.Context, repoFlag string) (*session, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(home)
	if err != nil {
		return nil, err
	}
	owner, name, err := repo.Resolve(ctx, repoFlag, runGit)
	if err != nil {
		return nil, err
	}
	opts := []gitbucket.Option{gitbucket.WithUserAgent("gb/" + version)}
	if cfg.CACert != "" {
		pool, err := loadCAPool(cfg.CACert)
		if err != nil {
			return nil, err
		}
		opts = append(opts, gitbucket.WithRootCAs(pool))
	}
	c, err := gitbucket.New(cfg.URL, cfg.Token, opts...)
	if err != nil {
		return nil, err
	}
	return &session{cfg: cfg, client: c, owner: owner, repo: name}, nil
}

func loadCAPool(path string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ca_cert %q: %w", path, err)
	}
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("ca_cert %q contains no PEM certificates", path)
	}
	return pool, nil
}

// runGit executes git and returns stdout. Used by repo.Resolve.
func runGit(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}
