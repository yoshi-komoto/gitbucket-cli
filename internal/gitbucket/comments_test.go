package gitbucket

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListIssueComments(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode([]IssueComment{
			{ID: 1, Body: "first", User: User{Login: "alice"}},
			{ID: 2, Body: "second", User: User{Login: "bob"}},
		})
	}))
	defer srv.Close()
	c, _ := New(srv.URL, "T")
	cs, err := c.ListIssueComments(context.Background(), "o", "r", 5)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v3/repos/o/r/issues/5/comments" {
		t.Errorf("path = %q", gotPath)
	}
	if len(cs) != 2 || cs[0].Body != "first" || cs[1].User.Login != "bob" {
		t.Errorf("cs = %+v", cs)
	}
}

func TestCreateIssueComment(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = strings.TrimSpace(string(b))
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"id":99,"body":"hi","user":{"login":"alice"}}`))
	}))
	defer srv.Close()
	c, _ := New(srv.URL, "T")
	got, err := c.CreateIssueComment(context.Background(), "o", "r", 5, "hi")
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v3/repos/o/r/issues/5/comments" {
		t.Errorf("path = %q", gotPath)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q", gotMethod)
	}
	if gotBody != `{"body":"hi"}` {
		t.Errorf("body = %q", gotBody)
	}
	if got.ID != 99 || got.Body != "hi" {
		t.Errorf("got = %+v", got)
	}
}

func TestCreateIssueComment_EmptyBody(t *testing.T) {
	c, _ := New("https://x.example/", "T")
	if _, err := c.CreateIssueComment(context.Background(), "o", "r", 5, ""); err == nil {
		t.Fatal("expected error for empty body")
	}
}
