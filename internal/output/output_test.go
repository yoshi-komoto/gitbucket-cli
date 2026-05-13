package output

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yoshi-komoto/gitbucket-cli/internal/gitbucket"
)

var update = flag.Bool("update", false, "update golden files")

func mustReadGolden(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	p := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(p, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	if want := mustReadGolden(t, name); want != got {
		t.Errorf("%s mismatch\nWANT:\n%s\nGOT:\n%s", name, want, got)
	}
}

func ts(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestRenderPullList_Table(t *testing.T) {
	prs := []gitbucket.PullRequest{
		{Number: 12, Title: "Add new feature", State: "open", User: gitbucket.User{Login: "yoshi"}, UpdatedAt: ts("2026-05-12T10:00:00Z")},
		{Number: 11, Title: "Fix typo in README", State: "closed", User: gitbucket.User{Login: "alice"}, UpdatedAt: ts("2026-05-10T01:23:00Z")},
	}
	var buf bytes.Buffer
	if err := RenderPullList(&buf, "table", prs); err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "pr_list.golden", buf.String())
}

func TestRenderPullList_JSON(t *testing.T) {
	prs := []gitbucket.PullRequest{{Number: 1, Title: "x", State: "open"}}
	var buf bytes.Buffer
	if err := RenderPullList(&buf, "json", prs); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"number": 1`)) {
		t.Errorf("json output missing number: %s", buf.String())
	}
}

func TestRenderPullView_Table(t *testing.T) {
	pr := gitbucket.PullRequest{
		Number:    12,
		Title:     "Add new feature",
		Body:      "This adds a brand new feature.\nSecond line.",
		State:     "open",
		User:      gitbucket.User{Login: "yoshi"},
		Base:      gitbucket.Branch{Ref: "main"},
		Head:      gitbucket.Branch{Ref: "topic"},
		CreatedAt: ts("2026-05-10T00:00:00Z"),
		UpdatedAt: ts("2026-05-12T10:00:00Z"),
	}
	var buf bytes.Buffer
	if err := RenderPullView(&buf, "table", &pr); err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "pr_view.golden", buf.String())
}

func TestRenderCommentList_Table(t *testing.T) {
	cs := []gitbucket.IssueComment{
		{ID: 101, Body: "LGTM", User: gitbucket.User{Login: "alice"}, CreatedAt: ts("2026-05-11T14:23:00Z")},
		{ID: 102, Body: "please squash before merge\nthanks", User: gitbucket.User{Login: "bob"}, CreatedAt: ts("2026-05-12T09:01:00Z")},
	}
	var buf bytes.Buffer
	if err := RenderCommentList(&buf, "table", cs); err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "comment_list.golden", buf.String())
}

func TestRenderComment_TableSingle(t *testing.T) {
	c := gitbucket.IssueComment{ID: 200, Body: "ok", User: gitbucket.User{Login: "yoshi"}, CreatedAt: ts("2026-05-13T01:02:03Z")}
	var buf bytes.Buffer
	if err := RenderComment(&buf, "table", &c); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("200")) || !bytes.Contains(buf.Bytes(), []byte("yoshi")) {
		t.Errorf("output = %s", buf.String())
	}
}
