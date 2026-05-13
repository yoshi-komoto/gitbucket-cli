package gitbucket

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	token      string
	userAgent  string
	rootCAs    *x509.CertPool
}

type Option func(*Client)

func WithHTTPClient(h *http.Client) Option         { return func(c *Client) { c.httpClient = h } }
func WithUserAgent(ua string) Option               { return func(c *Client) { c.userAgent = ua } }
func WithRootCAs(pool *x509.CertPool) Option       { return func(c *Client) { c.rootCAs = pool } }

func New(baseURL, token string, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is empty")
	}
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	trimmed := strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(trimmed, "/api/v3") {
		trimmed += "/api/v3"
	}
	u, err := url.Parse(trimmed + "/")
	if err != nil {
		return nil, fmt.Errorf("parse baseURL: %w", err)
	}
	c := &Client{
		baseURL:    u,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		userAgent:  "gb/dev",
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.rootCAs != nil {
		tr, ok := c.httpClient.Transport.(*http.Transport)
		if !ok || tr == nil {
			tr = http.DefaultTransport.(*http.Transport).Clone()
		} else {
			tr = tr.Clone()
		}
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		} else {
			tr.TLSClientConfig = tr.TLSClientConfig.Clone()
		}
		tr.TLSClientConfig.RootCAs = c.rootCAs
		c.httpClient.Transport = tr
	}
	return c, nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	rel, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("parse path: %w", err)
	}
	full := c.baseURL.ResolveReference(rel)
	if len(query) > 0 {
		full.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, full.String(), reqBody)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp, full.String())
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func parseAPIError(resp *http.Response, reqURL string) error {
	apiErr := &APIError{StatusCode: resp.StatusCode, URL: reqURL}
	var payload struct {
		Message string `json:"message"`
	}
	body, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(body, &payload) == nil {
		apiErr.Message = payload.Message
	}
	return apiErr
}
