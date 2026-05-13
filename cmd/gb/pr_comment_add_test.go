package gb

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveCommentBody_FlagWins(t *testing.T) {
	got, err := resolveCommentBody("hello", "", nil)
	if err != nil || got != "hello" {
		t.Errorf("got=%q err=%v", got, err)
	}
}

func TestResolveCommentBody_FileAndFlagConflict(t *testing.T) {
	_, err := resolveCommentBody("a", "b", nil)
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("err = %v", err)
	}
}

func TestResolveCommentBody_FromFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "body.md")
	if err := os.WriteFile(p, []byte("from file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := resolveCommentBody("", p, nil)
	if err != nil || got != "from file" {
		t.Errorf("got=%q err=%v", got, err)
	}
}

func TestResolveCommentBody_EmptyFlag(t *testing.T) {
	_, err := resolveCommentBody("", "", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveCommentBody_WhitespaceOnly(t *testing.T) {
	_, err := resolveCommentBody("   \n\t  ", "", nil)
	if !errors.Is(err, errEmptyBody) {
		t.Errorf("err = %v, want errEmptyBody", err)
	}
}
