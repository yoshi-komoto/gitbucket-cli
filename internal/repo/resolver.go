package repo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrRepoNotDetermined = errors.New("cannot determine repository")

// GitRunner runs a git command and returns its stdout.
type GitRunner func(ctx context.Context, args ...string) (string, error)

var (
	flagPattern  = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
	httpsPattern = regexp.MustCompile(`^https?://[^/]+/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	sshPattern   = regexp.MustCompile(`^[^@]+@[^:]+:([^/]+)/([^/]+?)(?:\.git)?/?$`)
)

func Resolve(ctx context.Context, flag string, run GitRunner) (owner, repo string, err error) {
	if flag != "" {
		m := flagPattern.FindStringSubmatch(flag)
		if m == nil {
			return "", "", fmt.Errorf("invalid --repo value %q (expected OWNER/REPO)", flag)
		}
		return m[1], m[2], nil
	}
	out, err := run(ctx, "remote", "get-url", "origin")
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrRepoNotDetermined, err)
	}
	url := strings.TrimSpace(out)
	if m := httpsPattern.FindStringSubmatch(url); m != nil {
		return m[1], m[2], nil
	}
	if m := sshPattern.FindStringSubmatch(url); m != nil {
		return m[1], m[2], nil
	}
	return "", "", fmt.Errorf("%w: unrecognized remote URL %q", ErrRepoNotDetermined, url)
}
