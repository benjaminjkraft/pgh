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

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func makeTestRepo(t *testing.T) (tmpdir string, r *runner) {
	tmpdir, err := os.MkdirTemp("", "pgh_test_")
	must(t, err)
	if debug {
		fmt.Println("repo:", tmpdir)
	} else {
		t.Cleanup(func() {
			os.RemoveAll(tmpdir)
		})
	}

	must(t, runCommands(tmpdir, `git init`))
	repo, err := git.PlainOpen(tmpdir)
	must(t, err)

	var b strings.Builder
	return tmpdir, &runner{repo, &b}
}

func commitTwoBranches(t *testing.T, tmpdir string) {
	must(t, runCommands(tmpdir, `
		echo content >content
		echo untracked >untracked
		git add content
		git commit -am "Initial commit"
		git branch -M main
		echo main content >content
		git commit -am "Main commit"
		git checkout main^ -b mybranch
		git branch mybranch -u main
		echo branch content >content
		git commit -am "Branch commit"
	`))
}

func assertFileHasContent(t *testing.T, filename, expectedContent string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
		return
	}
	cleanedContent := strings.TrimSpace(string(content))
	if cleanedContent != expectedContent {
		t.Errorf("content wrong, want '%s' got '%s'", expectedContent, cleanedContent)
	}
}

func assertFileHasConflict(t *testing.T, filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(content), "\n>>>>>>> ") {
		t.Errorf("content wrong, want conflict got '%s'", content)
	}
}

func showLog(t *testing.T, tmpdir string) {
	// TODO: Instead, assert the commit graph, details, etc. are right.
	// Perhaps via stable author-data + snapshots?
	// Crazy parser for ASCII graphs?
	must(t, runCommands(tmpdir, `
		git --no-pager log --branches --graph --decorate --pretty=fuller -p
	`))
}

func TestFakeMerge(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	commitTwoBranches(t, tmpdir)

	err := fakeMerge(runner, "main")
	must(t, err)

	assertFileHasContent(t, path.Join(tmpdir, "content"), "main content")
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")
	showLog(t, tmpdir)
}

func TestFakeMergeNoArgs(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	commitTwoBranches(t, tmpdir)

	err := fakeMerge(runner)
	must(t, err)

	assertFileHasContent(t, path.Join(tmpdir, "content"), "main content")
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")

	must(t, runCommands(tmpdir, `
		git --no-pager log --branches --graph --decorate --pretty=fuller
	`))
}
