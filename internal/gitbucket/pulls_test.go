package gitbucket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPulls(t *testing.T) {
	var gotPath, gotState, gotPerPage string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotState = r.URL.Query().Get("state")
		gotPerPage = r.URL.Query().Get("per_page")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]PullRequest{
			{Number: 12, Title: "feat", State: "open", User: User{Login: "alice"}},
			{Number: 11, Title: "fix", State: "open", User: User{Login: "bob"}},
		})
	}))
	defer srv.Close()

	c, _ := New(srv.URL, "T")
	prs, err := c.ListPulls(context.Background(), "owner", "repo", PullsListOptions{State: "open", PerPage: 30})
	if err != nil {
		t.Fatalf("ListPulls: %v", err)
	}
	if gotPath != "/api/v3/repos/owner/repo/pulls" {
		t.Errorf("path = %q", gotPath)
	}
	if gotState != "open" {
		t.Errorf("state = %q", gotState)
	}
	if gotPerPage != "30" {
		t.Errorf("per_page = %q", gotPerPage)
	}
	if len(prs) != 2 || prs[0].Number != 12 || prs[1].User.Login != "bob" {
		t.Errorf("prs = %+v", prs)
	}
}

func TestListPulls_DefaultState(t *testing.T) {
	var gotState string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotState = r.URL.Query().Get("state")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	c, _ := New(srv.URL, "T")
	if _, err := c.ListPulls(context.Background(), "o", "r", PullsListOptions{}); err != nil {
		t.Fatal(err)
	}
	if gotState != "open" {
		t.Errorf("default state = %q, want open", gotState)
	}
}
