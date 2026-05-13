package gitbucket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_DoGETSuccess(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotAccept, gotUA string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		gotUA = r.Header.Get("User-Agent")
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
	}))
	defer srv.Close()

	c, err := New(srv.URL, "TOKEN", WithUserAgent("gb/test"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var out map[string]string
	q := url.Values{}
	q.Set("state", "open")
	if err := c.do(context.Background(), http.MethodGet, "repos/o/r/pulls", q, nil, &out); err != nil {
		t.Fatalf("do: %v", err)
	}

	if gotMethod != "GET" {
		t.Errorf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/api/v3/repos/o/r/pulls" {
		t.Errorf("path = %q", gotPath)
	}
	if gotAuth != "token TOKEN" {
		t.Errorf("auth = %q", gotAuth)
	}
	if gotAccept != "application/json" {
		t.Errorf("accept = %q", gotAccept)
	}
	if gotUA != "gb/test" {
		t.Errorf("ua = %q", gotUA)
	}
	if gotQuery.Get("state") != "open" {
		t.Errorf("query state = %q", gotQuery.Get("state"))
	}
	if out["hello"] != "world" {
		t.Errorf("out = %v", out)
	}
}

func TestClient_DoPOSTBody(t *testing.T) {
	var gotBody string
	var gotCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = strings.TrimSpace(string(b))
		gotCT = r.Header.Get("Content-Type")
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"id":42}`))
	}))
	defer srv.Close()

	c, _ := New(srv.URL, "T")
	var out struct {
		ID int `json:"id"`
	}
	in := map[string]string{"body": "hello"}
	if err := c.do(context.Background(), http.MethodPost, "x", nil, in, &out); err != nil {
		t.Fatalf("do: %v", err)
	}
	if gotCT != "application/json" {
		t.Errorf("content-type = %q", gotCT)
	}
	if gotBody != `{"body":"hello"}` {
		t.Errorf("body = %q", gotBody)
	}
	if out.ID != 42 {
		t.Errorf("id = %d", out.ID)
	}
}

func TestClient_DoErrorWithMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer srv.Close()

	c, _ := New(srv.URL, "T")
	err := c.do(context.Background(), http.MethodGet, "x", nil, nil, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v, want *APIError", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("status = %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Not Found" {
		t.Errorf("msg = %q", apiErr.Message)
	}
}

func TestClient_DoErrorWithoutJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`<html>oops</html>`))
	}))
	defer srv.Close()

	c, _ := New(srv.URL, "T")
	err := c.do(context.Background(), http.MethodGet, "x", nil, nil, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v, want *APIError", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("status = %d", apiErr.StatusCode)
	}
	if apiErr.Message != "" {
		t.Errorf("msg should be empty, got %q", apiErr.Message)
	}
}

func TestNew_NormalizesBaseURL(t *testing.T) {
	c, err := New("https://example.com/", "T")
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL.String() != "https://example.com/api/v3/" {
		t.Errorf("baseURL = %q", c.baseURL.String())
	}

	c, err = New("https://example.com/api/v3", "T")
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL.String() != "https://example.com/api/v3/" {
		t.Errorf("baseURL = %q", c.baseURL.String())
	}
}
