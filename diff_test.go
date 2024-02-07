package main

import (
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// TODO(benkraft): replace with the shell setup in fake_merge_test.go
func TestDiffSmoke(t *testing.T) {
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	must(err)
	w, err := r.Worktree()
	must(err)
	var b strings.Builder
	runner := &runner{r, &b}

	// make a commit
	err = util.WriteFile(fs, "a", []byte("content\n"), 0o644)
	must(err)
	_, err = w.Add(".")
	must(err)
	h, err := w.Commit("initial commit", &git.CommitOptions{All: true})
	must(err)

	// make main point to it
	err = r.CreateBranch(&config.Branch{Name: "main", Remote: "origin", Merge: "refs/heads/main"})
	must(err)
	err = r.Storer.SetReference(plumbing.NewHashReference("refs/heads/main", h))
	must(err)

	// make another commit
	err = util.WriteFile(fs, "a", []byte("updated content\n"), 0o644)
	must(err)
	_, err = w.Commit("initial commit", &git.CommitOptions{All: true})
	must(err)

	err = diff(runner)
	must(err)
}
