package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	git "github.com/libgit2/git2go/v28"
)

func commit(repo *git.Repository, message string, filenames ...string) error {
	loc, err := time.LoadLocation("America/Buenos_Aires")
	if err != nil {
		return err
	}

	sig := &git.Signature{
		Name:  "Jorge Luis Borges",
		Email: "borges@example.com",
		When:  time.Date(1986, time.June, 14, 0, 0, 0, 0, loc),
	}

	index, err := repo.Index()
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		err = index.AddByPath(filename)
		if err != nil {
			return err
		}
	}
	err = index.Write()
	if err != nil {
		return err
	}
	treeID, err := index.WriteTree()
	if err != nil {
		return err
	}
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return err
	}
	var parents []*git.Commit
	headID, err := referenceNameToID(repo, "HEAD")
	if err == nil {
		headCommit, err := repo.LookupCommit(headID)
		if err != nil {
			return err
		}
		parents = append(parents, headCommit)
	} else if !git.IsErrorCode(err, git.ErrorCodeNotFound) {
		return err
	}

	_, err = repo.CreateCommit("HEAD", sig, sig, message, tree, parents...)
	return err
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
	// TODO(benkraft): Might actually be better to rewrite all this to just
	// shell out to make it easy to understand what we are making.
	repo, err = git.InitRepository(tmpdir, false)
	if err != nil {
		return nil, cleanup, err
	}

	filename := "content"
	err = ioutil.WriteFile(filepath.Join(tmpdir, filename), []byte("content\n"), 0o644)
	if err != nil {
		return nil, cleanup, err
	}
	err = commit(repo, "Initial commit\n", filename)
	if err != nil {
		return nil, cleanup, err
	}

	mainRef, err := repo.References.Lookup("refs/heads/master")
	if err != nil {
		return nil, cleanup, err
	}
	_, err = mainRef.Rename("refs/heads/main", false, "rename master -> main")
	if err != nil {
		return nil, cleanup, err
	}
	mainCommit, err := repo.LookupCommit(mainRef.Target())
	if err != nil {
		return nil, cleanup, err
	}

	branch, err := repo.CreateBranch("mybranch", mainCommit, false)
	if err != nil {
		return nil, cleanup, err
	}

	headRef, err := repo.References.Lookup("HEAD")
	if err != nil {
		return nil, cleanup, err
	}
	_, err = headRef.SetSymbolicTarget(branch.Reference.Name(), "checkout mybranch")
	if err != nil {
		return nil, cleanup, err
	}

	err = ioutil.WriteFile(filepath.Join(tmpdir, filename), []byte("updated content\n"), 0o644)
	if err != nil {
		return nil, cleanup, err
	}
	err = commit(repo, "Another commit\n", filename)
	if err != nil {
		return nil, cleanup, err
	}
	return repo, cleanup, nil
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
