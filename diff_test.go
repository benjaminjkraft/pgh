package main

import (
	"strings"
	"testing"
)

func TestDiffSmoke(t *testing.T) {
	tmpdir, repo := makeTestRepo(t)

	must(t, runCommands(tmpdir, `
		echo content >content
		echo untracked >untracked
		git add content
		git commit -am "Initial commit"
		git remote add origin .
		git fetch origin
		git branch -M main
		git branch -u origin/main
		git checkout --detach
		echo updated content >content
		git commit -am "New commit"
	`))

	var b strings.Builder
	runner := &runner{repo, &b}

	must(t, diff(runner))
}
