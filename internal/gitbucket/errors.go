package gitbucket

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("gitbucket API error: %d", e.StatusCode)
	}
	return fmt.Sprintf("gitbucket API error: %d %s", e.StatusCode, e.Message)
}
