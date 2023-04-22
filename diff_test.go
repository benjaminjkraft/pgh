package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	git "github.com/libgit2/git2go/v28"
)

func runCommands(cwd string, commands string) error {
	for _, line := range strings.Split(commands, "\n") {
		cleaned := strings.TrimSpace(line)
		if cleaned == "" {
			continue
		}
		cmd := &exec.Cmd{
			Path:   "/bin/sh",
			Args:   []string{"sh", "-c", cleaned},
			Dir:    cwd,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("cmd `%v` failed: %w", cleaned, err)
		}
	}
	return nil
}

func makeTestRepo() (repo *git.Repository, cleanup func(), err error) {
	tmpdir, err := os.MkdirTemp("", "pgh_test_")
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		if repo != nil {
			repo.Free()
		}
		os.RemoveAll(tmpdir)
	}

	err = runCommands(tmpdir, `
		git init
		echo content >content
		git add content
		git commit -am "Initial commit"
		git branch -M main
		git checkout -b mybranch
		echo updated content >content
		git commit -am "Another commit"
	`)
	if err != nil {
		return nil, cleanup, err
	}

	repo, err = git.OpenRepository(tmpdir)
	return repo, cleanup, err
}

func TestDiffSmoke(t *testing.T) {
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	debug := os.Getenv("TEST_DEBUG") != ""

	repo, cleanup, err := makeTestRepo()
	if debug {
		fmt.Println("repo:", repo.Workdir())
	} else if cleanup != nil {
		defer cleanup()
	}
	must(err)
	var b strings.Builder
	runner := &runner{repo, &b}

	err = diff(runner)
	must(err)
}
