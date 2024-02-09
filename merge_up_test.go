package main

import (
	"fmt"
	"path"
	"testing"
)

func initialCommit(t *testing.T, tmpdir string) {
	must(t, runCommands(tmpdir, `
		echo content >content
		echo untracked >untracked
		git add content
		git commit -am "Initial commit"
		git branch -M main
	`))
}

func branchAndCommit(t *testing.T, tmpdir string, name string) {
	must(t, runCommands(tmpdir, fmt.Sprintf(`
		git checkout -b %s
		echo %s content >content
		git commit -am "%s commit"
	`, name, name, name)))
}

func TestMergeUp(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	initialCommit(t, tmpdir)
	branchAndCommit(t, tmpdir, "b1")
	branchAndCommit(t, tmpdir, "b2")
	branchAndCommit(t, tmpdir, "b3")
	must(t, runCommands(tmpdir, `
		git checkout main
		echo other >other
		git add other
		git commit -am "other commit"
	`))

	err := mergeUp(runner)
	must(t, err)

	assertFileHasContent(t, path.Join(tmpdir, "content"), "b3 content")
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")
	assertFileHasContent(t, path.Join(tmpdir, "other"), "other")
	snapshotLog(t, tmpdir)
}

func TestMergeUpConflict(t *testing.T) {
	tmpdir, runner := makeTestRepo(t)
	initialCommit(t, tmpdir)
	branchAndCommit(t, tmpdir, "b1")
	branchAndCommit(t, tmpdir, "b2")
	branchAndCommit(t, tmpdir, "b3")
	must(t, runCommands(tmpdir, `
		git checkout main
		echo main >content
		git commit -am "main commit"
	`))

	err := mergeUp(runner)
	if err == nil {
		t.Fatal("should have failed")
	}

	assertFileHasConflict(t, path.Join(tmpdir, "content"))
	assertFileHasContent(t, path.Join(tmpdir, "untracked"), "untracked")
	snapshotLog(t, tmpdir)
}
