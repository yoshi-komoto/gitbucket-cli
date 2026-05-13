package gitbucket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type PullsListOptions struct {
	State   string
	PerPage int
}

func (c *Client) ListPulls(ctx context.Context, owner, repo string, opt PullsListOptions) ([]PullRequest, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	state := opt.State
	if state == "" {
		state = "open"
	}
	q := url.Values{}
	q.Set("state", state)
	if opt.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opt.PerPage))
	}
	var out []PullRequest
	path := fmt.Sprintf("repos/%s/%s/pulls", owner, repo)
	if err := c.do(ctx, http.MethodGet, path, q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetPull(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	if number <= 0 {
		return nil, fmt.Errorf("number must be > 0")
	}
	path := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, number)
	var out PullRequest
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
