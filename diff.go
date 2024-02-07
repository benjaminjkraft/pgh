package main

import (
	"fmt"
)

func diff(runner *runner, args ...string) error {
	head, err := runner.repo.Head()
	if err != nil {
		return err
	}
	headCommit, err := runner.repo.CommitObject(head.Hash())
	if err != nil {
		return err
	}
	switch headCommit.NumParents() {
	case 0:
		return fmt.Errorf("can't diff on initial commit")
	case 1:
	default:
		return fmt.Errorf("can't diff on merge commit")
	}

	parent, err := headCommit.Parent(0)
	if err != nil {
		return err
	}

	mainBranch, err := runner.repo.Branch("main")
	if err != nil {
		return err
	}
	main, err := runner.repo.Reference(mainBranch.Merge, true)
	if err != nil {
		return err
	}
	if parent.Hash != main.Hash() {
		return fmt.Errorf("TODO: must be one commit ahead of main for now, but parent was %v and main is %v", parent.Hash, main.Hash())
	}
	fmt.Println(parent, main)

	return nil
}
