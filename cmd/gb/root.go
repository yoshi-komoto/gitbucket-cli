package gb

import (
	"context"
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
	c, err := gitbucket.New(cfg.URL, cfg.Token, gitbucket.WithUserAgent("gb/"+version))
	if err != nil {
		return nil, err
	}
	return &session{cfg: cfg, client: c, owner: owner, repo: name}, nil
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
