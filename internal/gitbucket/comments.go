package gitbucket

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListIssueComments(ctx context.Context, owner, repo string, number int) ([]IssueComment, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	if number <= 0 {
		return nil, fmt.Errorf("number must be > 0")
	}
	path := fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, number)
	var out []IssueComment
	if err := c.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CreateIssueComment(ctx context.Context, owner, repo string, number int, body string) (*IssueComment, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	if number <= 0 {
		return nil, fmt.Errorf("number must be > 0")
	}
	if body == "" {
		return nil, fmt.Errorf("body must not be empty")
	}
	path := fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, number)
	payload := map[string]string{"body": body}
	var out IssueComment
	if err := c.do(ctx, http.MethodPost, path, nil, payload, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
