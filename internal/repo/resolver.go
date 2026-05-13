package repo

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var ErrRepoNotDetermined = errors.New("cannot determine repository")

// GitRunner runs a git command and returns its stdout.
type GitRunner func(ctx context.Context, args ...string) (string, error)

var flagPattern = regexp.MustCompile(`^([^/]+)/([^/]+)$`)

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
	o, r, ok := parseRemoteURL(out)
	if !ok {
		return "", "", fmt.Errorf("%w: unrecognized remote URL %q", ErrRepoNotDetermined, strings.TrimSpace(out))
	}
	return o, r, nil
}

// parseRemoteURL extracts owner and repo from a git remote URL.
// It takes the last two path segments after stripping a trailing ".git",
// so it tolerates GitBucket-style context paths like
// "https://host/gitbucket/git/owner/repo.git".
func parseRemoteURL(raw string) (owner, repo string, ok bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", "", false
	}

	var path string
	switch {
	case strings.HasPrefix(s, "http://"), strings.HasPrefix(s, "https://"), strings.HasPrefix(s, "ssh://"), strings.HasPrefix(s, "git://"):
		u, err := url.Parse(s)
		if err != nil {
			return "", "", false
		}
		path = u.Path
	default:
		// SCP-like: user@host:path (must contain "@" before the ":")
		colon := strings.Index(s, ":")
		if colon < 0 || !strings.Contains(s[:colon], "@") {
			return "", "", false
		}
		path = s[colon+1:]
	}

	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", false
	}
	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]
	if owner == "" || repo == "" {
		return "", "", false
	}
	return owner, repo, true
}
