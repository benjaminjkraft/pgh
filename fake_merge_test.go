package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
)

var debug = os.Getenv("TEST_DEBUG") != ""

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

func makeTestRepo() (tmpdir string, repo *git.Repository, cleanup func(), err error) {
	tmpdir, err = os.MkdirTemp("", "pgh_test_")
	if err != nil {
		return tmpdir, nil, nil, err
	}
	if debug {
		fmt.Println("repo:", tmpdir)
	}

	cleanup = func() {
		os.RemoveAll(tmpdir)
	}

	err = runCommands(tmpdir, `
		git init
		echo content >content
		git add content
		git commit -am "Initial commit"
		git branch -M main
		echo main content >content
		git commit -am "Main commit"
		git checkout main^ -b mybranch
		echo branch content >content
		git commit -am "Branch commit"
	`)
	if err != nil {
		return tmpdir, nil, cleanup, err
	}

	repo, err = git.PlainOpen(tmpdir)
	return tmpdir, repo, cleanup, err
}

func TestFakeMerge(t *testing.T) {
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	tmpdir, repo, cleanup, err := makeTestRepo()
	if !debug && cleanup != nil {
		defer cleanup()
	}
	must(err)
	var b strings.Builder
	runner := &runner{repo, &b}

	err = fakeMerge(runner, "main")
	must(err)

	content, err := os.ReadFile(path.Join(tmpdir, "content"))
	must(err)
	cleanedContent := strings.TrimSpace(string(content))
	if cleanedContent != "main content" {
		t.Errorf("content wrong, got '%s'", cleanedContent)
	}

	must(runCommands(tmpdir, `
		git --no-pager log --branches --graph --decorate --pretty=fuller
	`))
	// TODO: Assert the commit graph, details, etc. are right.
}
