package gb

import (
	"errors"
	"fmt"
	"io"

	"github.com/yoshi-komoto/gitbucket-cli/internal/config"
	"github.com/yoshi-komoto/gitbucket-cli/internal/gitbucket"
	"github.com/yoshi-komoto/gitbucket-cli/internal/repo"
)

func exitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, config.ErrNotConfigured) || errors.Is(err, repo.ErrRepoNotDetermined) {
		return 2
	}
	var apiErr *gitbucket.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case 401, 403:
			return 4
		case 404:
			return 5
		default:
			return 6
		}
	}
	return 1
}

func printError(w io.Writer, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, config.ErrNotConfigured) {
		fmt.Fprintln(w, "gb: not configured. Create ~/.config/gitbucket/config.yaml with url and token (see README).")
		return
	}
	if errors.Is(err, repo.ErrRepoNotDetermined) {
		fmt.Fprintln(w, "gb: cannot determine repository. Pass --repo OWNER/REPO or run inside a git repo with origin set.")
		return
	}
	var apiErr *gitbucket.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case 401:
			fmt.Fprintln(w, "gb: authentication failed. Check token in ~/.config/gitbucket/config.yaml")
			return
		case 403:
			fmt.Fprintln(w, "gb: forbidden (403). Token may lack required scopes.")
			return
		case 404:
			fmt.Fprintln(w, "gb: not found (404). Repo, PR, or comment does not exist.")
			return
		}
	}
	fmt.Fprintf(w, "gb: %v\n", err)
}
