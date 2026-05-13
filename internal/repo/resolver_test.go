package repo

import (
	"context"
	"errors"
	"testing"
)

func TestResolve_FlagWins(t *testing.T) {
	owner, name, err := Resolve(context.Background(), "alice/foo", func(ctx context.Context, args ...string) (string, error) {
		t.Fatal("git should not be called when flag is set")
		return "", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if owner != "alice" || name != "foo" {
		t.Errorf("owner=%q name=%q", owner, name)
	}
}

func TestResolve_FlagInvalid(t *testing.T) {
	_, _, err := Resolve(context.Background(), "notslashed", nil)
	if err == nil {
		t.Fatal("want error for invalid flag")
	}
}

func TestResolve_FromHTTPSRemote(t *testing.T) {
	owner, name, err := Resolve(context.Background(), "", func(ctx context.Context, args ...string) (string, error) {
		return "https://gitbucket.example.com/owner/repo.git\n", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if owner != "owner" || name != "repo" {
		t.Errorf("owner=%q name=%q", owner, name)
	}
}

func TestResolve_FromHTTPSRemote_NoSuffix(t *testing.T) {
	owner, name, err := Resolve(context.Background(), "", func(ctx context.Context, args ...string) (string, error) {
		return "https://gitbucket.example.com/owner/repo", nil
	})
	if err != nil || owner != "owner" || name != "repo" {
		t.Errorf("owner=%q name=%q err=%v", owner, name, err)
	}
}

func TestResolve_FromSSHRemote(t *testing.T) {
	owner, name, err := Resolve(context.Background(), "", func(ctx context.Context, args ...string) (string, error) {
		return "git@gitbucket.example.com:owner/repo.git", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if owner != "owner" || name != "repo" {
		t.Errorf("owner=%q name=%q", owner, name)
	}
}

func TestResolve_GitFails(t *testing.T) {
	_, _, err := Resolve(context.Background(), "", func(ctx context.Context, args ...string) (string, error) {
		return "", errors.New("not a git repo")
	})
	if !errors.Is(err, ErrRepoNotDetermined) {
		t.Errorf("err = %v, want ErrRepoNotDetermined", err)
	}
}

func TestResolve_InvalidRemote(t *testing.T) {
	_, _, err := Resolve(context.Background(), "", func(ctx context.Context, args ...string) (string, error) {
		return "ftp://nope", nil
	})
	if !errors.Is(err, ErrRepoNotDetermined) {
		t.Errorf("err = %v, want ErrRepoNotDetermined", err)
	}
}
