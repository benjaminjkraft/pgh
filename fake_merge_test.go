package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/go-git/go-git/v5"
)

var debug = os.Getenv("TEST_DEBUG") != ""

func init() {
	os.Setenv("GIT_AUTHOR_NAME", "Testy McTestFace")
	os.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	os.Setenv("GIT_AUTHOR_DATE", "2023-04-01T12:34:56+00:00")
	os.Setenv("GIT_COMMITTER_NAME", "Testy McTestFace")
	os.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
	os.Setenv("GIT_COMMITTER_DATE", "2023-04-01T12:34:56+00:00")
}

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

func snapshotLog(t *testing.T, tmpdir string) {
	cmd := exec.Command("git", "--no-pager", "log", "--branches", "--graph", "--decorate", "--pretty=fuller", "-p")
	cmd.Dir = tmpdir
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	must(t, err)
	cupaloy.SnapshotT(t, out)
}

func TestFakeMerge(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	commitTwoBranches(t, tmpdir)

	err := fakeMerge(runner, "main")
	must(t, err)

	assertFileHasContent(t, path.Join(tmpdir, "content"), "main content")
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")
	snapshotLog(t, tmpdir)
}

func TestFakeMergeNoArgs(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	commitTwoBranches(t, tmpdir)

	err := fakeMerge(runner)
	must(t, err)

	assertFileHasContent(t, path.Join(tmpdir, "content"), "main content")
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")

	snapshotLog(t, tmpdir)
}
