package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// TODO(benkraft): unit tests
func upstream(repo *git.Repository, refName plumbing.ReferenceName, includeRemote bool) plumbing.ReferenceName {
	if !refName.IsBranch() {
		return ""
	}
	b, err := repo.Branch(refName.Short())
	if err != nil {
		return ""
	}

	if b.Remote == "." {
		return b.Merge
	} else if !includeRemote {
		return ""
	}
	return plumbing.NewRemoteReferenceName(b.Remote, b.Merge.Short())
}

func fakeMerge(runner *runner, args ...string) error {
	head, err := runner.repo.Head()
	if err != nil {
		return err
	}

	var otherRefName plumbing.ReferenceName
	switch len(args) {
	case 0:
		otherRefName = upstream(runner.repo, head.Name(), true)
		if otherRefName == "" {
			return fmt.Errorf("no upstream for %v, so must pass branch-name", head.Name())
		}
	case 1:
		otherRefName = plumbing.NewBranchReferenceName(args[0])
	default:
		return fmt.Errorf("usage: fake-merge [branch-name]")
	}

	otherRef, err := runner.repo.Reference(otherRefName, false)
	if err != nil {
		return err
	}

	// TODO: test this case (main is a symbolic-ref)
	for otherRef.Type() == plumbing.SymbolicReference {
		otherRef, err = runner.repo.Reference(otherRef.Target(), false)
		if err != nil {
			return err
		}
	}
	otherCommit, err := runner.repo.CommitObject(otherRef.Hash())
	if err != nil {
		return err
	}

	author, committer, err := getAuthor(runner.repo)
	if err != nil {
		return err
	}

	// Modifed from worktree.buildCommitObject
	mergeCommit := &object.Commit{
		Author:       author,
		Committer:    committer,
		Message:      mergeCommitMessage(otherRef, head),
		TreeHash:     otherCommit.TreeHash,
		ParentHashes: []plumbing.Hash{head.Hash(), otherRef.Hash()},
	}

	obj := runner.repo.Storer.NewEncodedObject()
	err = mergeCommit.Encode(obj)
	if err != nil {
		return err
	}
	mergeHash, err := runner.repo.Storer.SetEncodedObject(obj)
	if err != nil {
		return err
	}

	// Copied from Worktree.updateHEAD
	headRef, err := runner.repo.Storer.Reference(plumbing.HEAD)
	if err != nil {
		return err
	}

	name := plumbing.HEAD
	if headRef.Type() != plumbing.HashReference {
		name = headRef.Target()
	}

	ref := plumbing.NewHashReference(name, mergeHash)
	err = runner.repo.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	// Shell out to work around incorrect behavior of Worktree.Reset
	// https://github.com/src-d/go-git/issues/1026#issue-382413262
	err = callGit(runner, "reset", "--hard", mergeHash.String())
	if err != nil {
		return err
	}

	fmt.Fprintf(runner.out, "fake-merged %s\n", otherRefName)

	return nil
}

func getAuthor(repo *git.Repository) (author, committer object.Signature, err error) {
	var empty object.Signature
	// Validate() is what sets the author/committer defaults.
	var opts git.CommitOptions
	err = opts.Validate(repo)
	if err != nil {
		return empty, empty, err
	}

	author = *opts.Author
	committer = *opts.Committer

	// But it doesn't respect the envvars, which we care about for tests. (And
	// for, you know, actual usage.)
	if name, ok := os.LookupEnv("GIT_AUTHOR_NAME"); ok {
		author.Name = name
	}
	if email, ok := os.LookupEnv("GIT_AUTHOR_EMAIL"); ok {
		author.Email = email
	}
	if date, ok := os.LookupEnv("GIT_AUTHOR_DATE"); ok {
		// TODO: match git's date formats.
		dt, err := time.Parse(time.RFC3339, date)
		if err != nil {
			// errors should not pass silently...
			return empty, empty, err
		}
		author.When = dt
	}
	if name, ok := os.LookupEnv("GIT_COMMITTER_NAME"); ok {
		committer.Name = name
	}
	if email, ok := os.LookupEnv("GIT_COMMITTER_EMAIL"); ok {
		committer.Email = email
	}
	if date, ok := os.LookupEnv("GIT_COMMITTER_DATE"); ok {
		// TODO: match git's date formats.
		dt, err := time.Parse(time.RFC3339, date)
		if err != nil {
			// errors should not pass silently...
			return empty, empty, err
		}
		committer.When = dt
	}

	return author, committer, nil
}

func callGit(runner *runner, args ...string) error {
	wt, err := runner.repo.Worktree()
	if err != nil {
		return err
	}

	ch, ok := wt.Filesystem.(*chroot.ChrootHelper)
	if !ok {
		return fmt.Errorf("not implemented: %T", wt.Filesystem)
	}
	gitdir := ch.Root()

	cmd := exec.Command("git", append([]string{"-C", gitdir}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func mergeCommitMessage(other, head *plumbing.Reference) string {
	return fmt.Sprintf("Merge branch '%s' into '%s', clobbering all changes", other.Name().Short(), head.Name().Short())
}
